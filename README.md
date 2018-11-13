Prometheus-emcisilon-exporter
======

### Introduction

`Prometheus-emcisilon-exporter` is a golang prometheus exporter that facilities the monitoring of EMC Isilon cluster with prometheus.  It uses the OneFS API to expose metrics and configuration items.

The prometheus-emcisilon-exporter and be build on any go compatibile system with the required dependencies.

Similar projects are e.g.:

 - [Prometheus Node Exporter](https://github.com/prometheus/node_exporter)
 - [Prometheus Isilon Exporter](https://github.com/paychex/prometheus-isilon-exporter)
 - [Prometheus ECS Exporter](https://github.com/paychex/prometheus-emcecs-exporter)

### Usage

#### Building the package

This exporter will run on any supported go platform. To build run: `go build`

#### Dependencies

    "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
    "gopkg.in/alecthomas/kingpin.v2"
    "github.com/thecodeteam/goisilon"

#### Run

The build go binary can be executed directly on the command-line with the desired configration flags.

Make sure the password environment variable has been properly exported.

`./prometheus-emcisilon-exporter --isilon.cluster.fqdn=<FQDN> --isilon.cluster.port=<PORT> --isilon.cluster.username=<USERNAME> --isilon.cluster.password.env=<PASSWD_ENV>

Collectors can be enabled by using `--collector.<name>` or disabled with `--no-collector.<name>`

##### Configuration

###### System Flags

| Program Flags   | Description | Default Value | Required |
|-----------------------------|-----------------------------------------------------------------------------------------------|---------------------------------|--------|
| --isilon.cluster.fqdn    | The DNS FQDN for reaching the cluster. HTTPS is forced as auth is made with basic auth.       | "localhost" | Yes |
| --isilon.cluster.port     | The port for the API on the Isilon cluster.                                  | 8080 | Yes |
| --isilon.cluster.username | The username for access the API                                                               | | Yes |
| --isilon.cluster.password.env | The password environment variabled that contains the password.                                               | "ISILON_CLUSTER_PASSWORD" | Yes |
| --isilon.cluster.site | The site the cluster resides in. Added as a label. | | No |
| --web.listen-address | The port that the exporter is bound to. | ":9300" | Yes |
| --web.telemtry-path | HTTP path for access metrics. | "/metrics" | Yes |
| --log.level | Log level of the exporter. | "info" | Yes |
| --log.format | Sets log target and format | "logger:stderr" | Yes |

###### Collector Flags
| Flag | Collector | Description | Default |
|-------------------------------|------------|----------------------------------------------------------------------|---------------|
| --collector.capacity | capacity | Exposes system /ifs capacity information. | enabled |
| --collector.cluster_health | cluster_health | Exposes cluster health information | enabled |
| --collector.cluster_protocol | cluster_protocol | Exposes protocol statistics at the cluster level | enabled |
| --collector.common_cifs | cluster_protocol & node_protocol | Enables the collection of cifs protocol statistics | enabled |
| --collector.common_ftp | cluster_protocol & node_protocol | Enables the collection of ftp protocol statistics | enabled |
| --collector.common_hdfs | cluster_protocol & node_protocol | Enables the collection of hdfs protocol statistics | enabled |
| --collector.common_http | cluster_protocol & node_protocol | Enables the collection of http protocol statistics | enabled |
| --collector.common_irp | cluster_protocol & node_protocol | Enables the collection of irp protocol statistics | enabled |
| --collector.common_jobd | cluster_protocol & node_protocol | Enables the collection of jobd protocol statistics | enabled |
| --collector.common_lsass_in | cluster_protocol & node_protocol | Enables the collection of lsass inbound protocol statistics | enabled |
| --collector.common_lsass_out | cluster_protocol & node_protocol | Enables the collection of lsass outbound protocol statistics | enabled |
| --collector.common_nfs | cluster_protocol & node_protocol | Enables the collection of nfs protocol statistics | enabled |
| --collector.common_nfs3 | cluster_protocol & node_protocol | Enables the collection of nfs3 protocol statistics | enabled |
| --collector.common_nfs4 | cluster_protocol & node_protocol | Enables the collection of nfs4 protocol statistics | enabled |
| --collector.common_nlm | cluster_protocol & node_protocol | Enables the collection of nlm protocol statistics | enabled |
| --collector.common_papi | cluster_protocol & node_protocol | Enables the collection of papi protocol statistics | enabled |
| --collector.common_siq | cluster_protocol & node_protocol | Enables the collection of siq protocol statistics | enabled |
| --collector.common_smb1 | cluster_protocol & node_protocol | Enables the collection of smb1 protocol statistics | enabled |
| --collector.common_smb2 | cluster_protocol & node_protocol | Enables the collection of smb2 protocol statistics | enabled |
| --collector.cpu | cpu | Enables the collection of CPU statistics | enabled |
| --collector.disk | disk | Enables the collection of disk statistics |
| --collector.memory | memory | Enables the collection of memory statistics | enabled |
| --collector.network | network | Enables the collection of network statistics | enabled |
| --collector.nfs_exports | nfs_exports | Enables the collection of summary information about nfs exports | enabled |
| --collector.node_health | node_health | Enables the collection of node health information | enabled |
| --collector.node_partition | node_partition | Enables the collection of node partition information (\, \var, \var\crash, etc.) | enabled |
| --collector.node_protocol | node_protocol | Enables the collection of node level procotol statistics | enabled |
| --collector.node_status | node_status | Enables the collection of node status information (Power supplies & batteries) | enabled |
| --collector.quota | quota | Enables the collection of quota information (thresholds and status) | disabled |
| --collector.quota.type | quota | Sets the type of quotas to be collected (directory, user, group, etc.) | all |
| --collector.quota.exceeded | quota | Sets the quota collector to return only exceeded quotas | disabled | 
| --collector.quota_summary | quota_summary | Enables the collection of summary information about all quotas | enabled |
| --collector.smb_shares | smb_share | Enables the colleciton of summary information about smb share. |
| --collector.snapshots | snapshots | Enables the collection of summary information about snapshots | enabled |
| --collector.statfs | statfs | Enables the collection of statfs statistics about the general /ifs system | enabled |
| --collector.storage_pools | storage_pools | Enables the collection of information about storage pools (virtual hot spare size, etc.) | enabled |
| --collector.sync_iq | sync_iq | Enables the collection of sync iq policies | enabled |

#### Provided Metrics
```# HELP isilon_cluster_health Current health of the cluster. Int of 1 2 or 3
# TYPE isilon_cluster_health gauge
 
# HELP isilon_cluster_onefs_version Current OneFS version. This returns a 1 always, the version is a label to the metric.
# TYPE isilon_cluster_onefs_version gauge
 
# HELP isilon_ifs_bytes_avail Current ifs filesystem capacity available in bytes.
# TYPE isilon_ifs_bytes_avail gauge
 
# HELP isilon_ifs_bytes_free Current ifs filesystem capacity free in bytes.
# TYPE isilon_ifs_bytes_free gauge
 
# HELP isilon_ifs_bytes_total Current ifs filesystem capacity total in bytes.
# TYPE isilon_ifs_bytes_total gauge
 
# HELP isilon_ifs_bytes_used Current ifs filesystem capacity used in bytes.
# TYPE isilon_ifs_bytes_used gauge
 
# HELP isilon_ifs_percent_avail Current ifs filesystem capacity available as a percentage from 0.0 - 1.0.
# TYPE isilon_ifs_percent_avail gauge
 
# HELP isilon_ifs_percent_free Current ifs filesystem capacity free as a percentage from 0.0 - 1.0.
# TYPE isilon_ifs_percent_free gauge
 
# HELP isilon_ifs_percent_used Current ifs filesystem capacity used in as a percentage from 0.0 - 1.0.
# TYPE isilon_ifs_percent_used gauge
 
# HELP isilon_nfs_export_total Total number of NFS exports on a cluster.
# TYPE isilon_nfs_export_total gauge
 
# HELP isilon_node_boottime Unix timestamp of when a load booted.
# TYPE isilon_node_boottime gauge
 
# HELP isilon_node_clientstats_active Total node protocol operation in rate.
# TYPE isilon_node_clientstats_active gauge
 
# HELP isilon_node_clientstats_connected Total node protocol operation in rate.
# TYPE isilon_node_clientstats_connected gauge
 
# HELP isilon_node_cpu_sys_avg Current cpu busy percentage for sys mode represented in 0.0-1.0.
# TYPE isilon_node_cpu_sys_avg gauge
 
# HELP isilon_node_cpu_user_avg Current cpu busy percentage for user mode represented in 0.0-1.0.
# TYPE isilon_node_cpu_user_avg gauge
 
# HELP isilon_node_disk_busy_all Current disk busy percentage represented in 0.0-1.0.
# TYPE isilon_node_disk_busy_all gauge
 
# HELP isilon_node_disk_count Number of disk per node as seen by the onefs system.
# TYPE isilon_node_disk_count gauge
 
# HELP isilon_node_disk_iosched_queued_all Current queue depth for IO sceduler.
# TYPE isilon_node_disk_iosched_queued_all gauge

# HELP isilon_node_disk_unhealthy_count Number of unhealthy disk per node as an int.
# TYPE isilon_node_disk_unhealthy_count gauge
 
# HELP isilon_node_disk_xfers_in_rate_all Current disk ingest transfer rate.
# TYPE isilon_node_disk_xfers_in_rate_all gauge
 
# HELP isilon_node_disk_xfers_out_rate_all Current disk egress transfer rate.
# TYPE isilon_node_disk_xfers_out_rate_all gauge
 
# HELP isilon_node_health Current health of a node from the view of the onefs cluster.
# TYPE isilon_node_health gauge
 
# HELP isilon_node_load_15min Current 15min node load.
# TYPE isilon_node_load_15min gauge
 
# HELP isilon_node_load_1min Current 1min node load.
# TYPE isilon_node_load_1min gauge
 
# HELP isilon_node_load_5min Current 5min node load.
# TYPE isilon_node_load_5min gauge
 
# HELP isilon_node_memory_cache RAM memory currently used for cache in bytes.
# TYPE isilon_node_memory_cache gauge
 
# HELP isilon_node_memory_free RAM memory currently free in bytes.
# TYPE isilon_node_memory_free gauge
 
# HELP isilon_node_memory_used RAM memory currently in use in bytes.
# TYPE isilon_node_memory_used gauge
 
# HELP isilon_node_net_ext_bytes_in_rate Current network bytes in rate from external interfaces.
# TYPE isilon_node_net_ext_bytes_in_rate gauge
 
# HELP isilon_node_net_ext_bytes_out_rate Current network bytes out rate from external interfaces.
# TYPE isilon_node_net_ext_bytes_out_rate gauge
 
# HELP isilon_node_partition_count Count of the total number of partitions on a node.
# TYPE isilon_node_partition_count gauge
 
# HELP isilon_node_partition_filenodes_free Number of filenodes free on a partition.
# TYPE isilon_node_partition_filenodes_free gauge
 
# HELP isilon_node_partition_filenodes_free_percent Percentage of filenodes free on a partition.
# TYPE isilon_node_partition_filenodes_free_percent gauge
 
# HELP isilon_node_partition_filenodes_total Total number of filenodes on a partition.
# TYPE isilon_node_partition_filenodes_total gauge
 
# HELP isilon_node_partition_used_space_percentage Percentage of space used on a partition.
# TYPE isilon_node_partition_used_space_percentage gauge
 
# HELP isilon_node_status_battery Status for batteries.
# TYPE isilon_node_status_battery gauge
 
# HELP isilon_node_status_power_supply Status for power supplies.
# TYPE isilon_node_status_power_supply gauge
 
# HELP isilon_node_uptime Current uptime of a node in seconds.
# TYPE isilon_node_uptime gauge
 
# HELP isilon_quota_api_collection_duration Returns the amount of time it took to collect an iteration of quotas from the api.
# TYPE isilon_quota_api_collection_duration gauge
 
# HELP isilon_quota_container 1 if quota is a container quota, 0 if not.
# TYPE isilon_quota_container gauge
 
# HELP isilon_quota_enforced 1 if quota is enforced, 2 if quota is an advisory quota.
# TYPE isilon_quota_enforced gauge
 
# HELP isilon_quota_include_snapshots 1 if quota includes snapshots in usage, 0 if not.
# TYPE isilon_quota_include_snapshots gauge
 
# HELP isilon_quota_summary_default_group_quotas_count Number of default group quotas.
# TYPE isilon_quota_summary_default_group_quotas_count gauge
 
# HELP isilon_quota_summary_default_user_quotas_count Number of default user quotas.
# TYPE isilon_quota_summary_default_user_quotas_count gauge
 
# HELP isilon_quota_summary_directory_quotas_count Number of directory quotas.
# TYPE isilon_quota_summary_directory_quotas_count gauge
 
# HELP isilon_quota_summary_group_quotas_count Number of group quotas.
# TYPE isilon_quota_summary_group_quotas_count gauge
 
# HELP isilon_quota_summary_linked_quotas_count Number of linked quotas.
# TYPE isilon_quota_summary_linked_quotas_count gauge
 
# HELP isilon_quota_summary_quotas_user Number of user quotas.
# TYPE isilon_quota_summary_quotas_user gauge
 
# HELP isilon_quota_summary_total_quotas_count Total number of quotas on a cluster.
# TYPE isilon_quota_summary_total_quotas_count gauge
 
# HELP isilon_quota_threshold_advisory Usage bytes at which notifications will be sent but writes will not be denied.
# TYPE isilon_quota_threshold_advisory gauge
 
# HELP isilon_quota_threshold_advisory_exceeded 1 if the advisory threshold has been hit.
# TYPE isilon_quota_threshold_advisory_exceeded gauge
 
# HELP isilon_quota_threshold_hard Usage bytes at which further writes will be denied.
# TYPE isilon_quota_threshold_hard gauge
 
# HELP isilon_quota_threshold_hard_exceeded True if the hard threshold has been hit.
# TYPE isilon_quota_threshold_hard_exceeded gauge
 
# HELP isilon_quota_threshold_soft Usage bytes at which notifications will be sent and soft grace time will be started.
# TYPE isilon_quota_threshold_soft gauge
 
# HELP isilon_quota_threshold_soft_exceeded 1 if the soft threshold has been hit.
# TYPE isilon_quota_threshold_soft_exceeded gauge
 
# HELP isilon_quota_threshold_soft_grace Time in seconds after which the soft threshold has been hit before writes will be denied.
# TYPE isilon_quota_threshold_soft_grace gauge
 
# HELP isilon_quota_usage_inodes Number of inodes (filesystem entities) used by governed data.
# TYPE isilon_quota_usage_inodes gauge
 
# HELP isilon_quota_usage_logical Apparent bytes used by governed data.
# TYPE isilon_quota_usage_logical gauge
 
# HELP isilon_quota_usage_physical Bytes used for governed data and filesystem overhead.
# TYPE isilon_quota_usage_physical gauge
 
# HELP isilon_scrape_collector_duration_seconds isilon_exporter: Duration of a collector scrape,
# TYPE isilon_scrape_collector_duration_seconds gauge
 
# HELP isilon_scrape_collector_success isilon_exporter: Whether a collector succeeded.
# TYPE isilon_scrape_collector_success gauge
  
# HELP isilon_smb_share_total Total number of SMB shares on a cluster.
# TYPE isilon_smb_share_total gauge
 
# HELP isilon_snapshots_15_day_count Number of snapshots older than 15 days.
# TYPE isilon_snapshots_15_day_count gauge
 
# HELP isilon_snapshots_30_day_count Number of snapshots older than 30 days
# TYPE isilon_snapshots_30_day_count gauge
 
# HELP isilon_snapshots_60_day_count Number of snapshots older than 60 days.
# TYPE isilon_snapshots_60_day_count gauge
 
# HELP isilon_snapshots_7_day_count Number of snapshots older than 7 days.
# TYPE isilon_snapshots_7_day_count gauge
 
# HELP isilon_snapshots_90_day_count Number of snapshots older than 90 days.
# TYPE isilon_snapshots_90_day_count gauge
 
# HELP isilon_snapshots_active_count Number of snapshots that are active on the system.
# TYPE isilon_snapshots_active_count gauge
 
# HELP isilon_snapshots_active_size Size in bytes of space occupied by active snapshots.
# TYPE isilon_snapshots_active_size gauge
 
# HELP isilon_snapshots_deleting_count Number of snapshots that are being deleted from the system.
# TYPE isilon_snapshots_deleting_count gauge
 
# HELP isilon_snapshots_deleting_size Size in bytes of space occupied by snapshots being deleted. 
# TYPE isilon_snapshots_deleting_size gauge
 
# HELP isilon_snapshots_total_count Total number of snapshots (both active and deleting) on a cluster.
# TYPE isilon_snapshots_total_count gauge
 
# HELP isilon_snapshots_total_size Size in bytes of space occupides by all snapshots.
# TYPE isilon_snapshots_total_size gauge
 
# HELP isilon_stafs_file_block_avail The filesystem fragment size.
# TYPE isilon_stafs_file_block_avail gauge
 
# HELP isilon_statfs_file_block_free The number of free blocks in the filesystem.
# TYPE isilon_statfs_file_block_free gauge
 
# HELP isilon_statfs_file_block_size The filesystem fragment size.
# TYPE isilon_statfs_file_block_size gauge
 
# HELP isilon_statfs_file_block_used The total number of data blocks in the filesystem.
# TYPE isilon_statfs_file_block_used gauge
 
# HELP isilon_statfs_file_io_size The optimal transfer block size.
# TYPE isilon_statfs_file_io_size gauge
 
# HELP isilon_statfs_file_name_max The maximum length of a file name.
# TYPE isilon_statfs_file_name_max gauge
 
# HELP isilon_statfs_file_node_free The number of free blocks in the filesystem.
# TYPE isilon_statfs_file_node_free gauge
 
# HELP isilon_statfs_file_node_free_percent The percentage of free file nodes in the filesystem.
# TYPE isilon_statfs_file_node_free_percent gauge
 
# HELP isilon_statfs_file_node_total The total number of file nodes in the filesystem.
# TYPE isilon_statfs_file_node_total gauge
 
# HELP isilon_storage_pool_balanced 0 if the storage pool is balanced, 1 if it is not.
# TYPE isilon_storage_pool_balanced gauge
 
# HELP isilon_storage_pool_bytes_avail Number of bytes available on the storage pool.
# TYPE isilon_storage_pool_bytes_avail gauge
 
# HELP isilon_storage_pool_bytes_avail_ssd Number of bytes available on ssd for the storage pool.
# TYPE isilon_storage_pool_bytes_avail_ssd gauge
 
# HELP isilon_storage_pool_bytes_free Number of bytes available on the storage pool.
# TYPE isilon_storage_pool_bytes_free gauge
 
# HELP isilon_storage_pool_bytes_free_ssd Number of bytes free on ssd for the storage pool.
# TYPE isilon_storage_pool_bytes_free_ssd gauge
 
# HELP isilon_storage_pool_bytes_total Total number of bytes on the storage pool.
# TYPE isilon_storage_pool_bytes_total gauge
 
# HELP isilon_storage_pool_bytes_total_ssd Total number of bytes on ssd for the storage pool.
# TYPE isilon_storage_pool_bytes_total_ssd gauge
 
# HELP isilon_storage_pool_bytes_virtual_hot_spare Number of bytes in vhs for the storage pool.
# TYPE isilon_storage_pool_bytes_virtual_hot_spare gauge
 
# HELP isilon_storage_pool_manual 0 of storage pool is not manually managed, 1 is it is.
# TYPE isilon_storage_pool_manual gauge
 
# HELP isilon_storage_pool_total Total number of storage pools on a cluster.
# TYPE isilon_storage_pool_total gauge
 
# HELP isilon_sync_policies_total_count Total number of sync policies on the cluster.
# TYPE isilon_sync_policies_total_count gauge
 
# HELP isilon_sync_policy_enabled 1 = Enabled, 0 = Disabled for the specified policy
# TYPE isilon_sync_policy_enabled gauge
 
# HELP isilon_sync_policy_last_start Epoch timestame for last sync start for a policy.
# TYPE isilon_sync_policy_last_start gauge
 
# HELP isilon_sync_policy_last_success Epoch timestamp of the last successful sync for a policy.
# TYPE isilon_sync_policy_last_success gauge
 
# HELP isilon_sync_policy_priority Current priority for the policy.
# TYPE isilon_sync_policy_priority gauge
 
# HELP isilon_sync_policy_state Last state from run of sync policy.
# TYPE isilon_sync_policy_state gauge
 
# HELP isilon_sync_policy_workers_per_node Number of worker threads per node for a policy.
# TYPE isilon_sync_policy_workers_per_node gauge
```

### Contributing

Contributions are welcomed! Read the [Contributing Guide](./.github/CONTRIBUTING.md) for more information.

### Licensing

This project is licensed under the MIT license. See [LICENSE](LICENSE) for more information.
