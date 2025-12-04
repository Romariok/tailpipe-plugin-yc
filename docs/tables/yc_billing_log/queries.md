## Activity Examples

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

### Top Services by Cost (Last 30 Days)

Identify the top services by total cost in the last 30 days.

```sql
select
  service_name,
  sum(cost) as total_cost
from
  yc_billing_log
where
  date >= now() - interval 30 day
group by
  service_name
order by
  total_cost desc
limit 10;
```

### Cost by Folder

Compare spend across folders to see where costs are concentrated.

```sql
select
  folder_name,
  folder_id,
  sum(cost) as total_cost
from
  yc_billing_log
group by
  folder_name,
  folder_id
order by
  total_cost desc;
```

### Cost by Cloud

Break down spend by cloud.

```sql
select
  cloud_name,
  cloud_id,
  sum(cost) as total_cost
from
  yc_billing_log
group by
  cloud_name,
  cloud_id
order by
  total_cost desc;
```

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

### Credits Breakdown by Type

See how different credits reduce spend.

```sql
select
  date,
  sum(credit) as total_credit,
  sum(monetary_grant_credit) as monetary_grant_credit,
  sum(volume_incentive_credit) as volume_incentive_credit,
  sum(cud_credit) as cud_credit,
  sum(misc_credit) as misc_credit
from
  yc_billing_log
group by
  date
order by
  date asc;
```

### Currency Distribution

Understand the currency mix in your billing data.

```sql
select
  currency,
  count(*) as row_count,
  sum(cost) as total_cost,
  sum(credit) as total_credit
from
  yc_billing_log
group by
  currency
order by
  total_cost desc;
```

### Cost per Unit by SKU (Sample Efficiency View)

Approximate efficiency by dividing cost by usage quantity for rows where quantity is present.

```sql
select
  sku_name,
  service_name,
  pricing_unit,
  sum(pricing_quantity) as total_quantity,
  sum(cost) as total_cost,
  case when sum(pricing_quantity) > 0 then sum(cost) / sum(pricing_quantity) end as cost_per_unit
from
  yc_billing_log
where
  pricing_quantity is not null
group by
  sku_name,
  service_name,
  pricing_unit
order by
  total_cost desc
limit 50;
```

## Detection Examples

### Rows with Zero Cost but Non‑Zero Usage

Potential misconfiguration or grant coverage where usage is reported but cost is zero.

```sql
select
  date,
  service_name,
  sku_name,
  pricing_quantity,
  pricing_unit,
  cost,
  credit
from
  yc_billing_log
where
  coalesce(pricing_quantity, 0) > 0
  and coalesce(cost, 0) = 0
order by
  date desc;
```

### Negative or Net‑Negative Rows (High Credits)

Identify items where credits fully offset or exceed costs.

```sql
select
  date,
  service_name,
  sku_name,
  cost,
  credit,
  (cost + credit) as net_amount
from
  yc_billing_log
where
  (cost + credit) <= 0
order by
  date desc,
  net_amount asc;
```

### Missing Folder Names (Deleted or Renamed)

Rows where `folder_name` is empty can indicate deleted folders at export time.

```sql
select
  date,
  folder_id,
  folder_name,
  service_name,
  sum(cost) as total_cost
from
  yc_billing_log
where
  (folder_name is null or trim(folder_name) = '')
group by
  date,
  folder_id,
  folder_name,
  service_name
order by
  date desc,
  total_cost desc;
```

## Operational Examples

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

### Most Recently Updated Rows

Check the most recently updated billing rows to confirm ingestion recency.

```sql
select
  date,
  service_name,
  sku_name,
  updated_at
from
  yc_billing_log
order by
  updated_at desc
limit 50;
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


