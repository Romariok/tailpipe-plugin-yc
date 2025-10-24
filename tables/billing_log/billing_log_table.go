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
	if !row.ExportedAt.IsZero() {
		row.TpDate = row.ExportedAt.Truncate(24 * time.Hour)
	} else {
		row.TpDate = row.Date.Truncate(24 * time.Hour)
	}

	if !row.ExportedAt.IsZero() {
		row.TpTimestamp = row.ExportedAt
	} else {
		row.TpTimestamp = row.Date
	}
	row.TpIngestTimestamp = time.Now()

	return row, nil
}

func (c *BillingLogTable) GetDescription() string {
	return "Yandex Cloud billing logs with service, cost, currency and date."
}
