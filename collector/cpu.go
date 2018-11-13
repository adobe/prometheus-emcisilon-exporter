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

type cpuCollector struct {
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
	keyMap := make(map[string]string)

	keyMap["cpuUser"] = "node.cpu.user.avg"
	keyMap["cpuSys"] = "node.cpu.sys.avg"
	keyMap["load1min"] = "node.load.1min"
	keyMap["load5min"] = "node.load.5min"
	keyMap["load15min"] = "node.load.15min"

	for promStat, statKey := range keyMap {
		resp, err := isiclient.QueryStatsEngineSingleVal(IsiCluster.Client, statKey)
		if err != nil {
			log.Warnf("Error attempting to query stats engine with key %s: %s", statKey, err)
			return err
		}
		for _, stat := range resp.Stats {
			node := fmt.Sprintf("%v", stat.Devid)
			val := stat.Value
			err = c.updateCPUStats(promStat, node, val, ch)
			if err != nil {
				log.Warnf("Error updating metric for CPU stat %s", promStat)
			}
		}
	}
	return nil
}

func (c *cpuCollector) updateCPUStats(promStat string, node string, val float64, ch chan<- prometheus.Metric) error {
	switch promStat {
	case "cpuUser":
		mVal := val / 10
		ch <- prometheus.MustNewConstMetric(c.cpuUser, prometheus.GaugeValue, mVal, node)
	case "cpuSys":
		mVal := val / 10
		ch <- prometheus.MustNewConstMetric(c.cpuSys, prometheus.GaugeValue, mVal, node)
	case "load1min":
		mVal := val / 100
		ch <- prometheus.MustNewConstMetric(c.load1min, prometheus.GaugeValue, mVal, node)
	case "load5min":
		mVal := val / 100
		ch <- prometheus.MustNewConstMetric(c.load5min, prometheus.GaugeValue, mVal, node)
	case "load15min":
		mVal := val / 100
		ch <- prometheus.MustNewConstMetric(c.load15min, prometheus.GaugeValue, mVal, node)
	}
	return nil
}
