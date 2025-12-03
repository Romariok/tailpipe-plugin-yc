package billing_log_api

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"log/slog"

	"github.com/Romariok/tailpipe-plugin-yc/config"
	"github.com/Romariok/tailpipe-plugin-yc/tables"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/schema"
	"github.com/turbot/tailpipe-plugin-sdk/types"
	ycsdk "github.com/yandex-cloud/go-sdk"
	"github.com/yandex-cloud/go-sdk/iamkey"
	ydb "github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
)

const BillingLogAPISourceIdentifier = "yc_billing_log_api"

type BillingLogAPISource struct {
	row_source.RowSourceImpl[BillingLogAPISourceConfig, *config.YandexCloudConnection]
	db *ydb.Driver
}

func (s *BillingLogAPISource) Identifier() string { return BillingLogAPISourceIdentifier }

func (s *BillingLogAPISource) Init(ctx context.Context, params *row_source.RowSourceParams, opts ...row_source.RowSourceOption) error {
	s.RegisterSource(s)
	s.NewCollectionStateFunc = NewBillingLogCollectionState
	s.GetGranularityFunc = func() time.Duration { return time.Hour }
	if err := s.RowSourceImpl.Init(ctx, params, opts...); err != nil {
		return err
	}

	var ydbOpts []ydb.Option
	if s.Connection.KeyFile != "" {
		key, err := iamkey.ReadFromJSONFile(s.Connection.KeyFile)
		if err != nil {
			slog.Error("Failed to read service account key file", "error", err)
			return err
		} else {
			creds, err := ycsdk.ServiceAccountKey(key)
			if err != nil {
				slog.Error("Failed to build YC credentials from key", "error", err)
				return err
			} else {
				sdk, err := ycsdk.Build(ctx, ycsdk.Config{Credentials: creds})
				if err != nil {
					slog.Error("Failed to init YC SDK for IAM exchange", "error", err)
					return err
				} else {
					resp, err := sdk.CreateIAMTokenForServiceAccount(ctx, key.GetServiceAccountId())
					if err != nil {
						slog.Error("Failed to create IAM token from key", "error", err)
						return err
					} else if resp != nil {
						ydbOpts = append(ydbOpts, ydb.WithAccessTokenCredentials(resp.IamToken))
					}
				}
			}
		}
	}

	// Build DSN with optional custom endpoint
	var dsn string
	const defaultEndpoint = "grpc.yandex-query.cloud.yandex.net:2135"
	if s.Connection.EndpointUrl != nil {
		endpoint := strings.TrimSpace(*s.Connection.EndpointUrl)
		if endpoint != "" {
			// If a scheme is provided, assume it's a full DSN base; otherwise default to grpcs://
			if strings.Contains(endpoint, "://") {
				dsn = fmt.Sprintf("%s/%s", strings.TrimRight(endpoint, "/"), s.Connection.FolderID)
			} else {
				dsn = fmt.Sprintf("grpcs://%s/%s", endpoint, s.Connection.FolderID)
			}
		}
	}
	if dsn == "" {
		dsn = fmt.Sprintf("grpcs://%s/%s", defaultEndpoint, s.Connection.FolderID)
	}

	db, err := ydb.Open(ctx, dsn, ydbOpts...)
	if err != nil {
		return fmt.Errorf("failed to open YDB connection: %w", err)
	}
	s.db = db
	return nil
}

func (s *BillingLogAPISource) Collect(ctx context.Context) error {

	enrichment := schema.NewSourceEnrichment(nil)

	if s.db != nil {
		// Moscow timezone handling (UTC+3):
		// - compute 'today' and 'yesterday' by Moscow calendar date
		// - collection 'to' is Moscow today 00:00 converted to UTC (exclusive upper boundary)
		utcNow := time.Now().UTC()
		mskNow := utcNow.Add(3 * time.Hour)
		mskTodayStart := time.Date(mskNow.Year(), mskNow.Month(), mskNow.Day(), 0, 0, 0, 0, time.UTC)
		_ = mskTodayStart.Add(-3 * time.Hour)

		mskYesterday := mskTodayStart.Add(-24 * time.Hour)
		toDateStr := mskYesterday.Format("2006-01-02")
		// derive lower bound from saved collection state end time (only if state is not empty)
		var fromDateStr string
		if s.RowSourceImpl.CollectionState != nil && !s.RowSourceImpl.CollectionState.IsEmpty() {
			lastTo := s.RowSourceImpl.CollectionState.GetToTime()
			if !lastTo.IsZero() {
				// Continue from the Moscow calendar day corresponding to the saved 'to' boundary (exclusive).
				// Example: if 'to' was 2025-11-27 00:00 MSK (stored as 2025-11-26 21:00 UTC),
				// we should start from 2025-11-27 (no re-fetch of 26th).
				fromDateStr = lastTo.Add(3 * time.Hour).Format("2006-01-02")
			}
		}

		var qCSV string
		selectClause := fmt.Sprintf(
			"SELECT\n"+
				"    billing_account_id,\n"+
				"    billing_account_name,\n"+
				"    cloud_id,\n"+
				"    cloud_name,\n"+
				"    folder_id,\n"+
				"    folder_name,\n"+
				"    resource_id,\n"+
				"    service_id,\n"+
				"    service_name,\n"+
				"    sku_id,\n"+
				"    sku_name,\n"+
				"    CAST(updated_at AS String) AS updated_at,\n"+
				"    CAST(exported_at AS String) AS exported_at,\n"+
				"    CAST(date AS String) AS date,\n"+
				"    currency,\n"+
				"    CAST(pricing_quantity AS Double) AS pricing_quantity,\n"+
				"    pricing_unit,\n"+
				"    CAST(cost AS Double) AS cost,\n"+
				"    CAST(credit AS Double) AS credit,\n"+
				"    CAST(monetary_grant_credit AS Double) AS monetary_grant_credit,\n"+
				"    CAST(volume_incentive_credit AS Double) AS volume_incentive_credit,\n"+
				"    CAST(cud_credit AS Double) AS cud_credit,\n"+
				"    CAST(misc_credit AS Double) AS misc_credit,\n"+
				"    locale\n"+
				"FROM `%s`\n",
			s.Connection.ConnectionName,
		)
		if fromDateStr != "" {
			qCSV = fmt.Sprintf("%sWHERE CAST(date AS Utf8) >= '%s' AND CAST(date AS Utf8) <= '%s';", selectClause, fromDateStr, toDateStr)
		} else {
			qCSV = fmt.Sprintf("%sWHERE CAST(date AS Utf8) <= '%s';", selectClause, toDateStr)
		}

		ctxQ, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()
		// Retry configuration
		attempts := 1
		if s.Connection.MaxErrorRetryAttempts != nil && *s.Connection.MaxErrorRetryAttempts > 0 {
			attempts = *s.Connection.MaxErrorRetryAttempts
		}
		var baseDelay time.Duration
		if s.Connection.MinErrorRetryDelay != nil && *s.Connection.MinErrorRetryDelay > 0 {
			baseDelay = time.Duration(*s.Connection.MinErrorRetryDelay) * time.Millisecond
		}

		var execErr error
		for try := 0; try < attempts; try++ {
			execErr = s.db.Table().Do(ctxQ, func(ctx context.Context, sess table.Session) error {
				tx := table.TxControl(table.BeginTx(table.WithOnlineReadOnly()), table.CommitTx())
				_, res, err := sess.Execute(ctx, tx, qCSV, nil)
				if err != nil {
					return err
				}
				defer res.Close()
				for res.NextResultSet(ctx) {
					for res.NextRow() {
						var (
							billingAccountID      sql.NullString
							billingAccountName    sql.NullString
							cloudID               sql.NullString
							cloudName             sql.NullString
							folderID              sql.NullString
							folderName            sql.NullString
							resourceID            sql.NullString
							serviceID             sql.NullString
							serviceName           sql.NullString
							skuID                 sql.NullString
							skuName               sql.NullString
							updatedAt             sql.NullString
							exportedAt            sql.NullString
							date                  sql.NullString
							currency              sql.NullString
							pricingQuantity       sql.NullFloat64
							pricingUnit           sql.NullString
							cost                  sql.NullFloat64
							credit                sql.NullFloat64
							monetaryGrantCredit   sql.NullFloat64
							volumeIncentiveCredit sql.NullFloat64
							cudCredit             sql.NullFloat64
							miscCredit            sql.NullFloat64
							locale                sql.NullString
						)
						if err := res.Scan(
							&billingAccountID,
							&billingAccountName,
							&cloudID,
							&cloudName,
							&folderID,
							&folderName,
							&resourceID,
							&serviceID,
							&serviceName,
							&skuID,
							&skuName,
							&updatedAt,
							&exportedAt,
							&date,
							&currency,
							&pricingQuantity,
							&pricingUnit,
							&cost,
							&credit,
							&monetaryGrantCredit,
							&volumeIncentiveCredit,
							&cudCredit,
							&miscCredit,
							&locale,
						); err != nil {
							return err
						}
						dateStr := tables.VString(date)
						dateTime := tables.MoscowDateMidnightUTC(dateStr)
						if dateTime.IsZero() {
							continue
						}
						row := map[string]interface{}{
							"billing_account_id":      tables.VString(billingAccountID),
							"billing_account_name":    tables.VString(billingAccountName),
							"cloud_id":                tables.VString(cloudID),
							"cloud_name":              tables.VString(cloudName),
							"folder_id":               tables.VString(folderID),
							"folder_name":             tables.VString(folderName),
							"resource_id":             tables.VString(resourceID),
							"service_id":              tables.VString(serviceID),
							"service_name":            tables.VString(serviceName),
							"sku_id":                  tables.VString(skuID),
							"sku_name":                tables.VString(skuName),
							"date":                    tables.VTimeString(date),
							"currency":                tables.VString(currency),
							"pricing_quantity":        tables.VFloat(pricingQuantity),
							"pricing_unit":            tables.VString(pricingUnit),
							"cost":                    tables.VFloat(cost),
							"credit":                  tables.VFloat(credit),
							"monetary_grant_credit":   tables.VFloat(monetaryGrantCredit),
							"volume_incentive_credit": tables.VFloat(volumeIncentiveCredit),
							"cud_credit":              tables.VFloat(cudCredit),
							"misc_credit":             tables.VFloat(miscCredit),
							"locale":                  tables.VString(locale),
							"updated_at":              tables.VTimeString(updatedAt),
							"exported_at":             tables.VTimeString(exportedAt),
						}
						if err := s.RowSourceImpl.OnRow(ctx, &types.RowData{Data: row, SourceEnrichment: enrichment}); err != nil {
							return err
						}
					}
				}
				return res.Err()
			})
			if execErr == nil {
				break
			}
			// backoff
			if baseDelay > 0 && try+1 < attempts {
				time.Sleep(baseDelay * time.Duration(1<<try))
			}
		}
		if execErr != nil {
			return fmt.Errorf("failed to execute query: %w", execErr)
		}
	}
	utcNow := time.Now().UTC()
	mskNow := utcNow.Add(3 * time.Hour)
	mskTodayStart := time.Date(mskNow.Year(), mskNow.Month(), mskNow.Day(), 0, 0, 0, 0, time.UTC)
	s.CollectionTimeRange.UpperBoundary = mskTodayStart.Add(-3 * time.Hour)
	return s.RowSourceImpl.OnCollectionComplete()
}
