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
		"cloud_id":                "Cloud ID for which the details are collected",
		"cloud_name":              "Cloud name",
		"folder_id":               "Folder ID",
		"folder_name":             "Folder name at the time of export. The value can be empty if the folder was deleted before the export",
		"resource_id":             "Resource ID, resource name, or subscription ID. The value can be empty if the service usage applies to the entire folder or contains no resources",
		"service_id":              "Service ID that owns the consumed product",
		"service_name":            "Service name that owns the consumed product",
		"sku_id":                  "Consumed product ID",
		"sku_name":                "Product name",
		"date":                    "Date for which the consumption cost is charged. The date is defined as a range from 0:00 to 23:59 Moscow time (UTC +3)",
		"currency":                "Billing account currency. Possible values: RUB, USD, and KZT",
		"pricing_quantity":        "Number of consumed product units. Decimal separator is a period",
		"pricing_unit":            "Unit of measurement for product consumption",
		"cost":                    "Total consumption cost. Decimal separator is a period",
		"credit":                  "Total discount amount. Always negative. Decimal separator is a period",
		"monetary_grant_credit":   "Grant discount, including the platform trial grant. Decimal separator is a period",
		"volume_incentive_credit": "Discount for product consumption volume. Decimal separator is a period",
		"cud_credit":              "Discount for committed use of resources. The cost of consumption exceeding the committed amount equals the sum of the cost and credit columns. Decimal separator is a period",
		"misc_credit":             "Other types of discounts, including discounts for resource consumption after the trial grant expires but before switching to the paid version. Decimal separator is a period",
		"locale":                  "Language of each row in the export. The value determines the language of the sku_name column. Possible values: en and ru",
		"updated_at":              "Date and time of the last row update in Unix Timestamp format",
		"exported_at":             "Date and time when the row was added to the details file",
	}
}
