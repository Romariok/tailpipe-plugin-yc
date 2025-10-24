package yc

import (
	"github.com/turbot/go-kit/helpers"
	"github.com/Romariok/tailpipe-plugin-yc/config"
	"github.com/Romariok/tailpipe-plugin-yc/sources/billing_log_api"
	"github.com/Romariok/tailpipe-plugin-yc/tables/billing_log"
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

type Plugin struct {
	plugin.PluginImpl
}

func init() {
	// Register tables, with type parameters:
	// 1. row struct
	// 2. table implementation
	table.RegisterTable[*billing_log.BillingLog, *billing_log.BillingLogTable]()
	// register sources
	row_source.RegisterRowSource[*billing_log_api.BillingLogAPISource]()
}

func NewPlugin() (_ plugin.TailpipePlugin, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = helpers.ToError(r)
		}
	}()

	p := &Plugin{
		PluginImpl: plugin.NewPluginImpl(config.PluginName),
	}

	return p, nil
}