/*
Copyright 2019 Adobe
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

type nodeCapacityCollector struct {
	nodeIfsBytesFree  *prometheus.Desc
	nodeIfsBytesUsed  *prometheus.Desc
	nodeIfsBytesTotal *prometheus.Desc
}

func init() {
	registerCollector("node_capacity", defaultEnabled, NodeCapacityCollector)
}

//NodeCapacityCollector returns a new Collector exposing node ifs capacity information.
func NodeCapacityCollector() (Collector, error) {
	return &nodeCapacityCollector{
		nodeIfsBytesFree: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "ifs_bytes_free"),
			"Number of ifs bytes free on the node.",
			[]string{"node", "node_id"}, ConstLabels,
		),
		nodeIfsBytesUsed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "ifs_bytes_used"),
			"Number of ifs bytes used on the node.",
			[]string{"node", "node_id"}, ConstLabels,
		),
		nodeIfsBytesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "ifs_bytes_total"),
			"Number of ifs bytes total on the node.",
			[]string{"node", "node_id"}, ConstLabels,
		),
	}, nil
}

func (c *nodeCapacityCollector) Update(ch chan<- prometheus.Metric) error {
	keyMap := make(map[*prometheus.Desc]string)

	keyMap[c.nodeIfsBytesFree] = "node.ifs.bytes.free"
	keyMap[c.nodeIfsBytesUsed] = "node.ifs.bytes.used"
	keyMap[c.nodeIfsBytesTotal] = "node.ifs.bytes.total"

	for promStat, statKey := range keyMap {
		begin := time.Now()
		resp, err := isiclient.QueryStatsEngineSingleVal(IsiCluster.Client, statKey)
		duration := time.Since(begin)
		ch <- prometheus.MustNewConstMetric(statsEngineCallDuration, prometheus.GaugeValue, duration.Seconds(), statKey)
		if err != nil {
			log.Warnf("Error attempting to query stats engine with key %s: %s", statKey, err)
			ch <- prometheus.MustNewConstMetric(statsEngineCallFailure, prometheus.GaugeValue, 1, statKey)
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
