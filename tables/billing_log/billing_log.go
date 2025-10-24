package billing_log

import (
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/schema"
)

type BillingLog struct {
	schema.CommonFields

	BillingAccountID      string    `json:"billing_account_id"`
	BillingAccountName    string    `json:"billing_account_name"`
	CloudID               string    `json:"cloud_id"`
	CloudName             string    `json:"cloud_name"`
	FolderID              string    `json:"folder_id"`
	FolderName            string    `json:"folder_name"`
	ResourceID            string    `json:"resource_id"`
	ServiceID             string    `json:"service_id"`
	ServiceName           string    `json:"service_name"`
	SkuID                 string    `json:"sku_id"`
	SkuName               string    `json:"sku_name"`
	Date                  time.Time `json:"date"`
	Currency              string    `json:"currency"`
	PricingQuantity       float64   `json:"pricing_quantity"`
	PricingUnit           string    `json:"pricing_unit"`
	Cost                  float64   `json:"cost"`
	Credit                float64   `json:"credit"`
	MonetaryGrantCredit   float64   `json:"monetary_grant_credit"`
	VolumeIncentiveCredit float64   `json:"volume_incentive_credit"`
	CUDCredit             float64   `json:"cud_credit"`
	MiscCredit            float64   `json:"misc_credit"`
	Locale                string    `json:"locale"`
	UpdatedAt             time.Time `json:"updated_at"`
	ExportedAt            time.Time `json:"exported_at"`
}

func (a *BillingLog) GetColumnDescriptions() map[string]string {
	return map[string]string{
		"billing_account_id":      "Billing account ID",
		"billing_account_name":    "Billing account name",
		"cloud_id":                "Cloud ID",
		"cloud_name":              "Cloud name",
		"folder_id":               "Folder ID",
		"folder_name":             "Folder name",
		"resource_id":             "Resource ID",
		"service_id":              "Service ID",
		"service_name":            "Service name",
		"sku_id":                  "SKU ID",
		"sku_name":                "SKU name",
		"date":                    "Date",
		"currency":                "Currency",
		"pricing_quantity":        "Pricing quantity",
		"pricing_unit":            "Pricing unit",
		"cost":                    "Cost",
		"credit":                  "Credit (negative)",
		"monetary_grant_credit":   "Monetary grant credit",
		"volume_incentive_credit": "Volume incentive credit",
		"cud_credit":              "CUDCredit",
		"misc_credit":             "Misc credit",
		"locale":                  "Locale (ru/en)",
		"updated_at":              "Updated at (UTC)",
		"exported_at":             "Exported at (UTC)",
	}
}
