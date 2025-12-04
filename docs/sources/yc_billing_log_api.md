---
title: "Source: yc_billing_log_api - Collect YC billing logs"
description: "Allows users to collect YC billing logs."
---

# Source: yc_billing_log_api - Collect YC billing logs

Yandex Cloud provides a scheduled export that periodically writes billing data as CSV files to an Object Storage bucket. This source is built on top of that export: once the export is enabled, we configure the integration so that Yandex Query can read the exported files from the bucket. Tailpipe then uses the Yandex Query connection to collect and query the billing records.

## Example Configurations

### Collect Billing logs

Collect Billing logs.

```hcl
connection "yc" "yc_analytics" {
   key_file = "/path/to/json"
   connection_name = "billing-billing-data-storage__"
   folder_id = "folder_id"
   endpoint_url = "grpcs://grpc.yandex-query.cloud.yandex.net:2135"
   max_error_retry_attempts = 3
   # milliseconds
   min_error_retry_delay = 200
}

partition "yc_billing_log" "my_analytics" {
  source "yc_billing_log_api" {
    connection = connection.yc.yc_analytics
  }
}
```


## Arguments

| Argument                  | Type   | Required | Default                                   | Description                                                                                 |
|---------------------------|--------|----------|-------------------------------------------|-------------------------------------------------------------------------------------------------|
| connection_name           | String | Yes      |                                           | Name of the connection in Yandex Query to read billing data from.                                                                                            |
| key_file                  | String | Yes      |                                           | Path to the service account JSON key (`authorized_key.json`).                                                                                          |
| folder_id                 | String | Yes      |                                           | Yandex Cloud Folder ID.                                                                                              |
| endpoint_url              | String | No       | `grpc.yandex-query.cloud.yandex.net:2135` | Custom endpoint for Yandex Query. Either `host:port` or a full prefix like `grpcs://host:port`.                                                                                      |
| max_error_retry_attempts  | Number | No       |                                           | Maximum number of retries for YQ requests.                                                                                        |
| min_error_retry_delay     | Number | No       |                                           | Minimum delay between retries (milliseconds).                                                                                  |

### Table Defaults

The following tables define their own default values for certain source arguments:

- **[yc_billing_log](https://hub.tailpipe.io/plugins/Romariok/yc/tables/yc_billing_log)**