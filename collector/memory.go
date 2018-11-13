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

type memoryCollector struct {
	memoryUsed  *prometheus.Desc
	memoryFree  *prometheus.Desc
	memoryCache *prometheus.Desc
}

func init() {
	registerCollector("memory", defaultEnabled, NewMemoryCollector)
}

//NewMemoryCollector returns a new Collector exposing node memory statistics.
func NewMemoryCollector() (Collector, error) {
	return &memoryCollector{
		memoryUsed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "memory_used"),
			"RAM memory currently in use in bytes.",
			[]string{"node"}, ConstLabels,
		),
		memoryFree: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "memory_free"),
			"RAM memory currently free in bytes.",
			[]string{"node"}, ConstLabels,
		),
		memoryCache: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "memory_cache"),
			"RAM memory currently used for cache in bytes.",
			[]string{"node"}, ConstLabels,
		),
	}, nil
}

func (c *memoryCollector) Update(ch chan<- prometheus.Metric) error {
	keyMap := make(map[string]string)

	keyMap["memoryUsed"] = "node.memory.used"
	keyMap["memoryFree"] = "node.memory.free"
	keyMap["memoryCache"] = "node.memory.cache"

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
			case "memoryUsed":
				ch <- prometheus.MustNewConstMetric(c.memoryUsed, prometheus.GaugeValue, val, node)
			case "memoryFree":
				ch <- prometheus.MustNewConstMetric(c.memoryFree, prometheus.GaugeValue, val, node)
			case "memoryCache":
				ch <- prometheus.MustNewConstMetric(c.memoryCache, prometheus.GaugeValue, val, node)
			}
		}
	}
	return nil
}
