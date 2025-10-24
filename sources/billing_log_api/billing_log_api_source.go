package billing_log_api

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"log/slog"

	"github.com/Romariok/tailpipe-plugin-yc/config"
	"github.com/turbot/tailpipe-plugin-sdk/collection_state"
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
	s.NewCollectionStateFunc = collection_state.NewTimeRangeCollectionState
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

	dsn := fmt.Sprintf("grpcs://%s/%s", "grpc.yandex-query.cloud.yandex.net:2135", s.Connection.FolderID)

	db, err := ydb.Open(ctx, dsn, ydbOpts...)
	if err != nil {
		return fmt.Errorf("failed to open YDB connection: %w", err)
	}
	s.db = db
	return nil
}

func (s *BillingLogAPISource) Collect(ctx context.Context) error {

	qCSV := `
SELECT
    billing_account_id,
    billing_account_name,
    cloud_id,
    cloud_name,
    folder_id,
    folder_name,
    resource_id,
    service_id,
    service_name,
    sku_id,
    sku_name,
    CAST(updated_at AS String) AS updated_at,
    CAST(exported_at AS String) AS exported_at,
    CAST(date AS String) AS date,
    currency,
    CAST(pricing_quantity AS Double) AS pricing_quantity,
    pricing_unit,
    CAST(cost AS Double) AS cost,
    CAST(credit AS Double) AS credit,
    CAST(monetary_grant_credit AS Double) AS monetary_grant_credit,
    CAST(volume_incentive_credit AS Double) AS volume_incentive_credit,
    CAST(cud_credit AS Double) AS cud_credit,
    CAST(misc_credit AS Double) AS misc_credit,
    locale
FROM ` + "`" + s.Connection.ConnectionName + "`" + `
LIMIT 1000;`

	enrichment := schema.NewSourceEnrichment(nil)

	if s.db != nil {
		ctxQ, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()
		err := s.db.Table().Do(ctxQ, func(ctx context.Context, sess table.Session) error {
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
					row := map[string]interface{}{
						"billing_account_id":      vString(billingAccountID),
						"billing_account_name":    vString(billingAccountName),
						"cloud_id":                vString(cloudID),
						"cloud_name":              vString(cloudName),
						"folder_id":               vString(folderID),
						"folder_name":             vString(folderName),
						"resource_id":             vString(resourceID),
						"service_id":              vString(serviceID),
						"service_name":            vString(serviceName),
						"sku_id":                  vString(skuID),
						"sku_name":                vString(skuName),
						"date":                    vTimeString(date),
						"currency":                vString(currency),
						"pricing_quantity":        vFloat(pricingQuantity),
						"pricing_unit":            vString(pricingUnit),
						"cost":                    vFloat(cost),
						"credit":                  vFloat(credit),
						"monetary_grant_credit":   vFloat(monetaryGrantCredit),
						"volume_incentive_credit": vFloat(volumeIncentiveCredit),
						"cud_credit":              vFloat(cudCredit),
						"misc_credit":             vFloat(miscCredit),
						"locale":                  vString(locale),
						"updated_at":              vTimeString(updatedAt),
						"exported_at":             vTimeString(exportedAt),
					}
					if err := s.RowSourceImpl.OnRow(ctx, &types.RowData{Data: row, SourceEnrichment: enrichment}); err != nil {
						return err
					}
				}
			}
			return res.Err()
		})
		if err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}
	return s.RowSourceImpl.OnCollectionComplete()
}

func vString(v sql.NullString) string {
	if v.Valid {
		return v.String
	}
	return ""
}

func vTimeString(v sql.NullString) string {
	if !v.Valid || v.String == "" {
		return ""
	}
	s := v.String
	if ts, err := time.Parse(time.RFC3339Nano, s); err == nil {
		return ts.UTC().Format(time.RFC3339Nano)
	}
	if ts, err := time.Parse(time.RFC3339, s); err == nil {
		return ts.UTC().Format(time.RFC3339Nano)
	}
	if secs, err := strconv.ParseInt(s, 10, 64); err == nil {
		return time.Unix(secs, 0).UTC().Format(time.RFC3339Nano)
	}
	if ts, err := time.Parse("2006-01-02", s); err == nil {
		return ts.UTC().Format(time.RFC3339Nano)
	}
	return s
}
func vFloat(v sql.NullFloat64) float64 {
	if v.Valid {
		return v.Float64
	}
	return 0
}
