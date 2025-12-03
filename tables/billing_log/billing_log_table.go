package billing_log

import (
	"time"

	"github.com/Romariok/tailpipe-plugin-yc/sources/billing_log_api"
	"github.com/rs/xid"
	"github.com/turbot/tailpipe-plugin-sdk/schema"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

const BillingLogTableIdentifier = "yc_billing_log"

type BillingLogTable struct {
}

func (c *BillingLogTable) Identifier() string {
	return BillingLogTableIdentifier
}

func (c *BillingLogTable) GetSourceMetadata() ([]*table.SourceMetadata[*BillingLog], error) {
	return []*table.SourceMetadata[*BillingLog]{
		{
			SourceName: billing_log_api.BillingLogAPISourceIdentifier,
			Mapper:     &BillingLogMapper{},
		},
	}, nil
}

func (c *BillingLogTable) EnrichRow(row *BillingLog, sourceEnrichmentFields schema.SourceEnrichment) (*BillingLog, error) {
	row.CommonFields = sourceEnrichmentFields.CommonFields
	row.TpID = xid.New().String()
	row.TpDate = row.Date.Truncate(24 * time.Hour)
	// Ensure TpTimestamp falls within the collection window for the billing day.
	// Use end of the UTC day for the billing date to avoid filtering out the boundary day.
	row.TpTimestamp = row.Date.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	row.TpIngestTimestamp = time.Now()

	return row, nil
}

func (c *BillingLogTable) GetDescription() string {
	return "Yandex Cloud billing logs (daily, Moscow time): service/SKU, costs and credits, currency."
}
