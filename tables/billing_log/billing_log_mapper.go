package billing_log

import (
	"context"

	"github.com/Romariok/tailpipe-plugin-yc/tables"
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

	billingAccountID, _ := tables.ToString(src["billing_account_id"])
	billingAccountName, _ := tables.ToString(src["billing_account_name"])
	cloudID, _ := tables.ToString(src["cloud_id"])
	cloudName, _ := tables.ToString(src["cloud_name"])
	folderID, _ := tables.ToString(src["folder_id"])
	folderName, _ := tables.ToString(src["folder_name"])
	resourceID, _ := tables.ToString(src["resource_id"])
	serviceID, _ := tables.ToString(src["service_id"])
	serviceName, _ := tables.ToString(src["service_name"])
	skuID, _ := tables.ToString(src["sku_id"])
	skuName, _ := tables.ToString(src["sku_name"])
	currency, ok := tables.ToString(src["currency"])
	if !ok {
		invalid = append(invalid, "currency")
	}
	pricingQuantity, _ := tables.ToFloat64(src["pricing_quantity"])
	pricingUnit, _ := tables.ToString(src["pricing_unit"])
	cost, ok := tables.ToFloat64(src["cost"])
	if !ok {
		invalid = append(invalid, "cost")
	}
	credit, _ := tables.ToFloat64(src["credit"])
	monetaryGrantCredit, _ := tables.ToFloat64(src["monetary_grant_credit"])
	volumeIncentiveCredit, _ := tables.ToFloat64(src["volume_incentive_credit"])
	cudCredit, _ := tables.ToFloat64(src["cud_credit"])
	miscCredit, _ := tables.ToFloat64(src["misc_credit"])
	locale, _ := tables.ToString(src["locale"])
	date, ok := tables.ToTime(src["date"])
	if !ok {
		invalid = append(invalid, "date")
	}
	updatedAt, _ := tables.ToTime(src["updated_at"])
	exportedAt, _ := tables.ToTime(src["exported_at"])

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
