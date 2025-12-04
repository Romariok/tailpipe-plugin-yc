---
organization: Turbot
category: ["public cloud"]
icon_url: "some path"
brand_color: "#5282FF"
display_name: "Yandex Cloud"
description: "Tailpipe plugin for collecting and querying various logs from Yandex Cloud."
og_description: "Collect Yandex Cloud logs and query them instantly with SQL! Open source CLI. No DB required."
og_image: "some path"
---

# Yandex Cloud + Tailpipe

[Tailpipe](https://tailpipe.io) is an open-source CLI tool that allows you to collect logs and query them with SQL.

[Yandex Cloud](https://yandex.cloud/en) providing scalable infrastructure, storage, machine learning and development tools to build and enhance digital services and applications.

The [YC Plugin for Tailpipe](https://hub.tailpipe.io/plugins/turbot/yc) allows you to collect and query YC billing data using SQL to track activity, monitor trends, and more!

- Documentation: [Table definitions & examples](https://hub.tailpipe.io/plugins/turbot/yc/tables)
- Community: [Join #tailpipe on Slack →](https://turbot.com/community/join)
- Get involved: [Issues](https://github.com/Romariok/tailpipe-plugin-yc/issues)

![image](https://raw.githubusercontent.com/Romariok/tailpipe-plugin-yc/main/docs/images/yc_billing_log_terminal.png?type=thumbnail)
![image](https://raw.githubusercontent.com/Romariok/tailpipe-plugin-yc/main/docs/images/yc_billing_log_dashboard.png?type=thumbnail)

## Getting Started

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

Dashboards and detections are [open source](https://github.com/topics/tailpipe-mod), allowing easy customization and collaboration.

To get started, choose a mod from the [Powerpipe Hub](https://hub.powerpipe.io/?engines=tailpipe&q=yc).

## Connection Credentials

### Arguments

| Name                        | Type   | Required | Description                                                                                         |
|-----------------------------|--------|----------|-----------------------------------------------------------------------------------------------------|
| `connection_name`           | String | Yes      | Name of the connection in Yandex Query to read billing data from.                                   |
| `key_file`                  | String | Yes      | Path to the service account JSON key (`authorized_key.json`).                                       |
| `folder_id`                 | String | Yes      | Yandex Cloud Folder ID.                                                                             |
| `endpoint_url`              | String | No       | Custom endpoint for Yandex Query. Either `host:port` or a full prefix like `grpcs://host:port`.     |
| `max_error_retry_attempts`  | Number | No       | Maximum number of retries for YQ requests.                                                          |
| `min_error_retry_delay`     | Number | No       | Minimum delay between retries (milliseconds).                                                       |

