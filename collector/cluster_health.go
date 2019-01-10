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

type clusterHealthCollector struct {
	clusterHealth *prometheus.Desc
	onefsVersion  *prometheus.Desc
}

func init() {
	registerCollector("cluster_health", defaultEnabled, NewClusterHealthCollector)
}

//NewClusterHealthCollector returns a new Collector exposing cluster health information.
func NewClusterHealthCollector() (Collector, error) {
	return &clusterHealthCollector{
		clusterHealth: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, clusterCollectorSubsystem, "health"),
			"Current health of the cluster. Int of 1 2 or 3",
			nil, ConstLabels,
		),
		onefsVersion: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, clusterCollectorSubsystem, "onefs_version"),
			"Current OneFS version. This returns a 1 always, the version is a label to the metric.",
			[]string{"version"}, ConstLabels,
		),
	}, nil
}

func (c *clusterHealthCollector) Update(ch chan<- prometheus.Metric) error {
	var errCount int64
	keyMap := make(map[*prometheus.Desc]string)

	keyMap[c.clusterHealth] = "cluster.health"

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
				ch <- prometheus.MustNewConstMetric(promStat, prometheus.GaugeValue, stat.Value)
			}
		}
	}

	version, err := isiclient.GetOneFsVersion(IsiCluster.Client)
	if err != nil {
		log.Warnf("Unable to update the Onefs version stat.")
		errCount++
	}
	ch <- prometheus.MustNewConstMetric(c.onefsVersion, prometheus.GaugeValue, 1, version)

	if errCount != 0 {
		err := fmt.Errorf("There where %v errors", errCount)
		return err
	}
	return nil
}
