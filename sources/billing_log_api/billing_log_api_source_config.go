package billing_log_api

import "github.com/hashicorp/hcl/v2"


type BillingLogAPISourceConfig struct {
	Remain hcl.Body `hcl:",remain" json:"-"`
}

func (c BillingLogAPISourceConfig) Validate() error    { return nil }
func (c BillingLogAPISourceConfig) Identifier() string { return BillingLogAPISourceIdentifier }
