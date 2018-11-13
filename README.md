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

### Contributing

Contributions are welcomed! Read the [Contributing Guide](./.github/CONTRIBUTING.md) for more information.

### Licensing

This project is licensed under the MIT license. See [LICENSE](LICENSE) for more information.