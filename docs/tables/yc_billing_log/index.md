---
title: "Tailpipe Table: yc_billing_log - Query YC Billing Logs"
description: "Analyze Yandex Cloud spend by service, folder, cloud, SKU, currency, and time."
---

# Table: yc_billing_log - Query YC Billing Logs

The `yc_billing_log` table allows you to query Yandex Cloud billing detail exports. It provides a structured view of service consumption, costs, credits, currency and dates, enabling cost analysis, trend monitoring, and optimization across clouds, folders, services, and SKUs.

This table is built on Yandex Cloud’s periodic export of billing data as CSV files to an Object Storage bucket. We configure an integration so that Yandex Query can read the exported files directly from the bucket, and Tailpipe uses that connection as the data source.

## Configure

Create a [partition](https://tailpipe.io/docs/manage/partition) for `yc_billing_log` ([examples](https://hub.tailpipe.io/plugins/turbot/yc/tables/yc_billing_log#example-configurations)):

```sh
vi ~/.tailpipe/config/yc.tpc
```

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

## Collect

[Collect](https://tailpipe.io/docs/manage/collection) logs for all `yc_billing_log` partitions:

```sh
tailpipe collect yc_billing_log
```

Or collect for a single partition:

```sh
tailpipe collect yc_billing_log.my_analytics
```

## Query

Explore example queries for this table in the queries page.

**[Explore 10+ example queries for this table →](https://hub.tailpipe.io/plugins/turbot/yc/queries/yc_billing_log)**

### Top SKUs by Cost

Find the SKUs contributing most to total spend.

```sql
select
  sku_name,
  sku_id,
  service_name,
  sum(cost) as total_cost
from
  yc_billing_log
group by
  sku_name,
  sku_id,
  service_name
order by
  total_cost desc
limit 20;
```

### Daily Cost Trends by Service

Analyze daily spend per service to understand usage and trends.

```sql
select
  date,
  service_name,
  sum(cost) as total_cost
from
  yc_billing_log
group by
  date,
  service_name
order by
  date asc,
  total_cost desc;
```

### High Daily Spend (Threshold)

Flag days where spend exceeds a given threshold.

```sql
select
  date,
  sum(cost) as total_cost
from
  yc_billing_log
group by
  date
having
  sum(cost) > 1000 -- adjust threshold to your currency
order by
  total_cost desc;
```

### Latest Exported Data by Day

Verify export freshness using the `exported_at` timestamp.

```sql
select
  date,
  max(exported_at) as latest_exported_at
from
  yc_billing_log
group by
  date
order by
  date desc;
```
