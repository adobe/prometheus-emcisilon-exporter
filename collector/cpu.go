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
	"strings"

	"github.com/adobe/prometheus-emcisilon-exporter/isiclient"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

type cpuCollector struct {
	cpuCount  *prometheus.Desc
	cpuIdle   *prometheus.Desc
	cpuUser   *prometheus.Desc
	cpuSys    *prometheus.Desc
	load1min  *prometheus.Desc
	load5min  *prometheus.Desc
	load15min *prometheus.Desc
}

func init() {
	registerCollector("cpu", defaultEnabled, NewCPUCollector)
}

//NewCPUCollector returns a new Collector exposing node cpu statistics.
func NewCPUCollector() (Collector, error) {
	return &cpuCollector{
		cpuCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "cpu_count"),
			"Count of number of cpu a node contains.",
			[]string{"node"}, ConstLabels,
		),
		cpuIdle: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "cpu_idle_avg"),
			"Current cpu idle percentage for the node.",
			[]string{"node"}, ConstLabels,
		),
		cpuUser: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "cpu_user_avg"),
			"Current cpu busy percentage for user mode represented in 0.0-1.0.",
			[]string{"node"}, ConstLabels,
		),
		cpuSys: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "cpu_sys_avg"),
			"Current cpu busy percentage for sys mode represented in 0.0-1.0.",
			[]string{"node"}, ConstLabels,
		),
		load1min: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "load_1min"),
			"Current 1min node load.",
			[]string{"node"}, ConstLabels,
		),
		load5min: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "load_5min"),
			"Current 5min node load.",
			[]string{"node"}, ConstLabels,
		),
		load15min: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "load_15min"),
			"Current 15min node load.",
			[]string{"node"}, ConstLabels,
		),
	}, nil
}

func (c *cpuCollector) Update(ch chan<- prometheus.Metric) error {
	keyMap := make(map[*prometheus.Desc]string)

	keyMap[c.cpuCount] = "node.cpu.count"
	keyMap[c.cpuIdle] = "node.cpu.idle.avg"
	keyMap[c.cpuUser] = "node.cpu.user.avg"
	keyMap[c.cpuSys] = "node.cpu.sys.avg"
	keyMap[c.load1min] = "node.load.1min"
	keyMap[c.load5min] = "node.load.5min"
	keyMap[c.load15min] = "node.load.15min"

	for promStat, statKey := range keyMap {
		resp, err := isiclient.QueryStatsEngineSingleVal(IsiCluster.Client, statKey)
		if err != nil {
			log.Warnf("Error attempting to query stats engine with key %s: %s", statKey, err)
			return err
		}
		for _, stat := range resp.Stats {
			var val float64
			node := fmt.Sprintf("%v", stat.Devid)
			if strings.Contains(statKey, "cpu") {
				if !(strings.Contains(statKey, "count")) {
					val = stat.Value / 10
				}
			}
			if strings.Contains(statKey, "load") {
				val = stat.Value / 100
			}
			ch <- prometheus.MustNewConstMetric(promStat, prometheus.GaugeValue, val, node)
		}
	}
	return nil
}
