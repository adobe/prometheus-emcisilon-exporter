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
	keyMap := make(map[string]string)

	keyMap["nodeDiskUnhealthyCount"] = "node.disk.unhealthy.count"
	keyMap["nodeHealth"] = "node.health"
	keyMap["nodeDiskCount"] = "node.disk.count"
	keyMap["nodeBootTime"] = "node.boottime"
	keyMap["nodeUptime"] = "node.uptime"

	for promStat, statKey := range keyMap {
		resp, err := isiclient.QueryStatsEngineSingleVal(IsiCluster.Client, statKey)
		if err != nil {
			log.Warnf("Error attempting to query stats engine with key %s: %s", statKey, err)
			return err
		}
		for _, stat := range resp.Stats {
			node := fmt.Sprintf("%v", stat.Devid)
			val := stat.Value
			switch promStat {
			case "nodeDiskUnhealthyCount":
				ch <- prometheus.MustNewConstMetric(c.nodeDiskUnhealthyCount, prometheus.GaugeValue, val, node)
			case "nodeHealth":
				ch <- prometheus.MustNewConstMetric(c.nodeHealth, prometheus.GaugeValue, val, node)
			case "nodeDiskCount":
				ch <- prometheus.MustNewConstMetric(c.nodeDiskCount, prometheus.GaugeValue, val, node)
			case "nodeBootTime":
				ch <- prometheus.MustNewConstMetric(c.nodeBootTime, prometheus.GaugeValue, val, node)
			case "nodeUptime":
				ch <- prometheus.MustNewConstMetric(c.nodeUptime, prometheus.GaugeValue, val, node)
			}
		}
	}
	return nil
}
