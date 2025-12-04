# Yandex Cloud Plugin for Tailpipe

[Tailpipe](https://tailpipe.io) is an open-source CLI tool that allows you to collect logs and query them with SQL.

[Yandex Cloud](https://yandex.cloud/en) providing scalable infrastructure, storage, machine learning and development tools to build and enhance digital services and applications.

The [YC Plugin for Tailpipe](https://hub.tailpipe.io/plugins/turbot/yc) allows you to collect and query YC billing data using SQL to track activity, monitor trends, and more!

- **[Get started →](https://hub.tailpipe.io/plugins/turbot/yc)**
- Documentation: [Table definitions & examples](https://hub.tailpipe.io/plugins/turbot/yc/tables)
- Community: [Join #tailpipe on Slack →](https://turbot.com/community/join)
- Get involved: [Issues](https://github.com/Romariok/tailpipe-plugin-yc/issues)

Collect and query logs:
![image](https://raw.githubusercontent.com/Romariok/tailpipe-plugin-yc/main/docs/images/yc_billing_log_terminal.png)

## Getting Started

### Yandex Cloud 

Sign up for Yandex Cloud and create a [billing account](https://yandex.cloud/en/docs/billing/concepts/billing-account):

1. Navigate to the [management console](https://console.yandex.cloud/) and log in to Yandex Cloud or create a new account.
2. On the [Yandex Cloud Billing](https://center.yandex.cloud/billing/accounts) page, make sure you have a billing account linked and it has the `ACTIVE` or `TRIAL_ACTIVE` [status](https://yandex.cloud/en/docs/billing/concepts/billing-account-statuses). If you do not have a billing account, [create one](https://yandex.cloud/en/docs/billing/quickstart/) and [link](https://yandex.cloud/en/docs/billing/operations/pin-cloud) a cloud to it.


Create service account and generate new authorized key for the service account and after that download json file with keys.

Set roles in the system:
- `billing.accounts.owner` for your account in cloud
- `yq.editor` and `storage.viewer` for service account in the folder

Set up regular export, instruction is available in the [Yandex Cloud docs](https://yandex.cloud/en/docs/billing/operations/get-folder-report#set-up-regular-download). In the `Create periodic export` window you must enter the value `/` ​​in the `Path within the bucket` field.

Then you must set up integration between Yandex Cloud Billing and Query:

1. Open the list of expense detail exports in the Yandex Cloud console.
2. Select the required detail and click **Process in YQ**:
![image](https://raw.githubusercontent.com/Romariok/tailpipe-plugin-yc/main/docs/images/yc_billing_log_yq_button_location.png)
1. When switching from Yandex Cloud Billing to Yandex Query for the first time, set up integration: 
   1. In the Query interface, select the service account for reading data from Object Storage in the connection creation dialog box and click **Create**.
   2. In the Query interface, choose `Billing` option in the **Automatically fill settings for** dropdown list. Next, click **Create** to complete the integration.

Current instructions of setting up integration between Yandex Cloud Billing and Query can be found at [Yandex Cloud docs](https://yandex.cloud/en/docs/billing/operations/query-integration#integration).

Configure lifecycle in the bucket:

1. Go to **Settings** page of the bucket, which you created previously
2. Select **Lifecycle** category and click **Configure** button
3. In the Rule Configuration interface add description and in the **Type** section set `Number of Days` trigger and enter value `1` to `Trigger time` field 

### Tailpipe

Install Tailpipe from the [downloads](https://tailpipe.io/downloads) page:

```sh
# MacOS
brew install turbot/tap/tailpipe
```

```sh
# Linux or Windows (WSL)
sudo /bin/sh -c "$(curl -fsSL https://tailpipe.io/install/tailpipe.sh)"
```

Install the plugin:

```sh
tailpipe plugin install yc
```

Configure your [connection credentials](https://hub.tailpipe.io/plugins/Romariok/yc#connection-credentials), table partition, and data source ([examples](https://hub.tailpipe.io/plugins/Romariok/yc/tables/yc_billing_log#example-configurations)):

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

Download, enrich, and save logs from your source ([examples](https://tailpipe.io/docs/reference/cli/collect)):

```sh
tailpipe collect yc_billing_log
```

Enter interactive query mode:

```sh
tailpipe query
```

Run a query:

```sql
select 
  date,
  service_name,
  count(*) as service_count
from 
  yc_billing_log 
group by 
  date,
  service_name 
order by
  date asc
```

```
+---------------------+-----------------------+---------------+
| date                | service_name          | service_count |
+---------------------+-----------------------+---------------+
| 2025-11-25 00:00:00 | Virtual Private Cloud | 3             |
| 2025-11-25 00:00:00 | Object Storage        | 7             |
| 2025-11-25 00:00:00 | Yandex Query          | 1             |
| 2025-11-26 00:00:00 | Virtual Private Cloud | 1             |
| 2025-11-26 00:00:00 | Object Storage        | 3             |
| 2025-11-27 00:00:00 | Object Storage        | 6             |
| 2025-11-27 00:00:00 | Yandex Query          | 1             |
| 2025-11-27 00:00:00 | Virtual Private Cloud | 3             |
| 2025-11-28 00:00:00 | Object Storage        | 6             |
| 2025-11-28 00:00:00 | Yandex Query          | 1             |
| 2025-11-28 00:00:00 | Virtual Private Cloud | 3             |
| 2025-11-29 00:00:00 | Virtual Private Cloud | 4             |
| 2025-11-29 00:00:00 | Object Storage        | 6             |
| 2025-11-29 00:00:00 | Yandex Query          | 1             |
| 2025-11-30 00:00:00 | Object Storage        | 1             |
| 2025-12-01 00:00:00 | Object Storage        | 5             |
| 2025-12-01 00:00:00 | Virtual Private Cloud | 3             |
+---------------------+-----------------------+---------------+
```

## Detections as Code with Powerpipe

Pre-built dashboard for the YC plugin is available in [GitHub](https://github.com/Romariok/tailpipe-mod-yc-cost-and-usage-insights), helping you monitor and analyze activity across your YC folders.

![image](https://raw.githubusercontent.com/Romariok/tailpipe-plugin-yc/main/docs/images/yc_billing_log_dashboard.png)

## Developing

Prerequisites:

- [Tailpipe](https://tailpipe.io/downloads)
- [Golang](https://golang.org/doc/install)

Clone:

```sh
git clone https://github.com/Romariok/tailpipe-plugin-yc.git
cd tailpipe-plugin-yc
```

After making your local changes, build the plugin, which automatically installs the new version to your `~/.tailpipe/plugins` directory:

```sh
make
```

Re-collect your data:

```sh
tailpipe collect yc_billing_log
```

Try it!

```sh
tailpipe query
> .inspect yc_billing_log
```

## Open Source & Contributing

This repository is published under the [Apache 2.0](https://www.apache.org/licenses/LICENSE-2.0) (source code) and [CC BY-NC-ND](https://creativecommons.org/licenses/by-nc-nd/2.0/) (docs) licenses.
