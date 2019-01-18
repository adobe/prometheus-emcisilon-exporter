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
	"time"

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
	var errCount int64
	keyMap := make(map[*prometheus.Desc]string)

	keyMap[c.memoryUsed] = "node.memory.used"
	keyMap[c.memoryFree] = "node.memory.free"
	keyMap[c.memoryCache] = "node.memory.cache"

	for promStat, statKey := range keyMap {
		begin := time.Now()
		resp, err := isiclient.QueryStatsEngineSingleVal(IsiCluster.Client, statKey)
		duration := time.Since(begin)
		ch <- prometheus.MustNewConstMetric(statsEngineCallDuration, prometheus.GaugeValue, duration.Seconds(), statKey)
		if err != nil {
			log.Warnf("Error attempting to query stats engine with key %s: %s", statKey, err)
			ch <- prometheus.MustNewConstMetric(statsEngineCallFailure, prometheus.GaugeValue, 1, statKey)
			errCount++
		} else {
			ch <- prometheus.MustNewConstMetric(statsEngineCallFailure, prometheus.GaugeValue, 0, statKey)
			for _, stat := range resp.Stats {
				node := fmt.Sprintf("%v", stat.Devid)
				ch <- prometheus.MustNewConstMetric(promStat, prometheus.GaugeValue, stat.Value, node)
			}
		}
	}
	return nil
}
