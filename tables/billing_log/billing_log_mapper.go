package billing_log

import (
	"context"
	"time"

	sdkerrors "github.com/turbot/tailpipe-plugin-sdk/error_types"
	"github.com/turbot/tailpipe-plugin-sdk/mappers"
	sdktypes "github.com/turbot/tailpipe-plugin-sdk/types"
)

type BillingLogMapper struct{}

func (m *BillingLogMapper) Map(_ context.Context, data any, _ ...mappers.MapOption[*BillingLog]) (*BillingLog, error) {
	var src map[string]interface{}

	switch v := data.(type) {
	case *sdktypes.RowData:
		if mp, ok := v.Data.(map[string]interface{}); ok {
			src = mp
		}
	case map[string]interface{}:
		src = v
	}

	if src == nil {
		return nil, sdkerrors.NewRowErrorWithMessage("unable to map row, invalid type received")
	}

	missing := []string{}
	invalid := []string{}

	billingAccountID, _ := toString(src["billing_account_id"])
	billingAccountName, _ := toString(src["billing_account_name"])
	cloudID, _ := toString(src["cloud_id"])
	cloudName, _ := toString(src["cloud_name"])
	folderID, _ := toString(src["folder_id"])
	folderName, _ := toString(src["folder_name"])
	resourceID, _ := toString(src["resource_id"])
	serviceID, _ := toString(src["service_id"])
	serviceName, _ := toString(src["service_name"])
	skuID, _ := toString(src["sku_id"])
	skuName, _ := toString(src["sku_name"])
	currency, ok := toString(src["currency"])
	if !ok {
		invalid = append(invalid, "currency")
	}
	pricingQuantity, _ := toFloat64(src["pricing_quantity"])
	pricingUnit, _ := toString(src["pricing_unit"])
	cost, ok := toFloat64(src["cost"])
	if !ok {
		invalid = append(invalid, "cost")
	}
	credit, _ := toFloat64(src["credit"])
	monetaryGrantCredit, _ := toFloat64(src["monetary_grant_credit"])
	volumeIncentiveCredit, _ := toFloat64(src["volume_incentive_credit"])
	cudCredit, _ := toFloat64(src["cud_credit"])
	miscCredit, _ := toFloat64(src["misc_credit"])
	locale, _ := toString(src["locale"])
	date, ok := toTime(src["date"])
	if !ok {
		invalid = append(invalid, "date")
	}
	updatedAt, _ := toTime(src["updated_at"])
	exportedAt, _ := toTime(src["exported_at"])

	if len(missing) > 0 || len(invalid) > 0 {
		return nil, sdkerrors.NewRowErrorWithFields(missing, invalid)
	}

	return &BillingLog{
		BillingAccountID:      billingAccountID,
		BillingAccountName:    billingAccountName,
		CloudID:               cloudID,
		CloudName:             cloudName,
		FolderID:              folderID,
		FolderName:            folderName,
		ResourceID:            resourceID,
		ServiceID:             serviceID,
		ServiceName:           serviceName,
		SkuID:                 skuID,
		SkuName:               skuName,
		Date:                  date,
		Currency:              currency,
		PricingQuantity:       pricingQuantity,
		PricingUnit:           pricingUnit,
		Cost:                  cost,
		Credit:                credit,
		MonetaryGrantCredit:   monetaryGrantCredit,
		VolumeIncentiveCredit: volumeIncentiveCredit,
		CUDCredit:             cudCredit,
		MiscCredit:            miscCredit,
		Locale:                locale,
		UpdatedAt:             updatedAt,
		ExportedAt:            exportedAt,
	}, nil
}

func (m *BillingLogMapper) Identifier() string {
	return "yc_billing_log_mapper"
}

func toString(v interface{}) (string, bool) {
	switch t := v.(type) {
	case string:
		return t, true
	default:
		return "", false
	}
}

func toFloat64(v interface{}) (float64, bool) {
	switch t := v.(type) {
	case float64:
		return t, true
	case float32:
		return float64(t), true
	case int:
		return float64(t), true
	case int64:
		return float64(t), true
	case int32:
		return float64(t), true
	default:
		return 0, false
	}
}

func toTime(v interface{}) (time.Time, bool) {
	switch t := v.(type) {
	case time.Time:
		return t, true
	case string:
		if ts, err := time.Parse(time.RFC3339Nano, t); err == nil {
			return ts, true
		}
		if ts, err := time.Parse(time.RFC3339, t); err == nil {
			return ts, true
		}
		if ts, err := time.Parse("2006-01-02", t); err == nil {
			return ts, true
		}
		return time.Time{}, false
	default:
		return time.Time{}, false
	}
}
