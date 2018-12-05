/*
Copyright 2018 Adobe
All Rights Reserved.

NOTICE: Adobe permits you to use, modify, and distribute this file in
accordance with the terms of the Adobe license agreement accompanying
it. If you have received this file from a source other than Adobe,
then your use, modification, or distribution of it requires the prior
written permission of Adobe.
*/
package collector

import (
	"fmt"

	"github.com/adobe/prometheus-emcisilon-exporter/isiclient"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

type nodeHealthCollector struct {
	nodeNvramBatteryStatus *prometheus.Desc
	nodeProcessCount       *prometheus.Desc
	nodeFilesOpen          *prometheus.Desc
	nodeDiskUnhealthyCount *prometheus.Desc
	nodeHealth             *prometheus.Desc
	nodeDiskCount          *prometheus.Desc
	nodeBootTime           *prometheus.Desc
	nodeUptime             *prometheus.Desc
}

func init() {
	registerCollector("node_health", defaultEnabled, NewNodeHealthCollector)
}

//NewNodeHealthCollector returns a new Collector exposing node health information.
func NewNodeHealthCollector() (Collector, error) {
	return &nodeHealthCollector{
		nodeProcessCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "process_count"),
			"Number of processess on the node.",
			[]string{"node"}, ConstLabels,
		),
		nodeNvramBatteryStatus: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "nvram_battery_status"),
			"Combined charge status for all batteries. 0 = Not available, 1 = Good, 2 = Caution, 3 = Error.",
			[]string{"node"}, ConstLabels,
		),
		nodeFilesOpen: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "open_files"),
			"Number of open files on the node.",
			[]string{"node"}, ConstLabels,
		),
		nodeDiskUnhealthyCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "disk_unhealthy_count"),
			"Number of unhealthy disk per node as an int.",
			[]string{"node"}, ConstLabels,
		),
		nodeHealth: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "health"),
			"Current health of a node from the view of the onefs cluster.",
			[]string{"node"}, ConstLabels,
		),
		nodeDiskCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "disk_count"),
			"Number of disk per node as seen by the onefs system.",
			[]string{"node"}, ConstLabels,
		),
		nodeBootTime: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "boottime"),
			"Unix timestamp of when a load booted.",
			[]string{"node"}, ConstLabels,
		),
		nodeUptime: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "uptime"),
			"Current uptime of a node in seconds.",
			[]string{"node"}, ConstLabels,
		),
	}, nil
}

func (c *nodeHealthCollector) Update(ch chan<- prometheus.Metric) error {
	keyMap := make(map[*prometheus.Desc]string)

	keyMap[c.nodeNvramBatteryStatus] = "node.nvram.charge.status"
	keyMap[c.nodeDiskUnhealthyCount] = "node.disk.unhealthy.count"
	keyMap[c.nodeHealth] = "node.health"
	keyMap[c.nodeDiskCount] = "node.disk.count"
	keyMap[c.nodeBootTime] = "node.boottime"
	keyMap[c.nodeUptime] = "node.uptime"
	keyMap[c.nodeFilesOpen] = "node.open.files"
	keyMap[c.nodeProcessCount] = "node.process.count"

	for promStat, statKey := range keyMap {
		resp, err := isiclient.QueryStatsEngineSingleVal(IsiCluster.Client, statKey)
		if err != nil {
			log.Warnf("Error attempting to query stats engine with key %s: %s", statKey, err)
			return err
		}
		for _, stat := range resp.Stats {
			node := fmt.Sprintf("%v", stat.Devid)
			ch <- prometheus.MustNewConstMetric(promStat, prometheus.GaugeValue, stat.Value, node)
		}
	}
	return nil
}
