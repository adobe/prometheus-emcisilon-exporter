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
	keyMap := make(map[string]string)

	keyMap["clusterHealth"] = "cluster.health"

	for promStat, statKey := range keyMap {
		resp, err := isiclient.QueryStatsEngineSingleVal(IsiCluster.Client, statKey)
		if err != nil {
			log.Warnf("Error attempting to query stats engine with key %s: %s", statKey, err)
		}
		for _, stat := range resp.Stats {
			val := stat.Value
			switch promStat {
			case "clusterHealth":
				ch <- prometheus.MustNewConstMetric(c.clusterHealth, prometheus.GaugeValue, val)
			}
		}
	}

	version, err := isiclient.GetOneFsVersion(IsiCluster.Client)
	if err != nil {
		log.Warnf("Unable to update the Onefs version stat.")
	}
	ch <- prometheus.MustNewConstMetric(c.onefsVersion, prometheus.GaugeValue, 1, version)
	return nil
}
