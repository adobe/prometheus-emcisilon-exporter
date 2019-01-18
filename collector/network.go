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

type networkCollector struct {
	netBytesInRate   *prometheus.Desc
	netBytesOutRate  *prometheus.Desc
	netErrorsInRate  *prometheus.Desc
	netErrorsOutRate *prometheus.Desc
}

func init() {
	registerCollector("network", defaultEnabled, NewNetworkCollector)
}

//NewNetworkCollector returns a new Collector exposing node network statistics.
func NewNetworkCollector() (Collector, error) {
	return &networkCollector{
		netBytesInRate: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "net_ext_bytes_in_rate"),
			"Current network bytes in rate from external interfaces.",
			[]string{"node"}, ConstLabels,
		),
		netBytesOutRate: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "net_ext_bytes_out_rate"),
			"Current network bytes out rate from external interfaces.",
			[]string{"node"}, ConstLabels,
		),
		netErrorsInRate: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "net_ext_errors_in_rate"),
			"Input errors per second for a node's external interfaces.",
			[]string{"node"}, ConstLabels,
		),
		netErrorsOutRate: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "net_ext_errors_out_rate"),
			"Output errors per seccond for a node's external interfaces.",
			[]string{"node"}, ConstLabels,
		),
	}, nil
}

func (c *networkCollector) Update(ch chan<- prometheus.Metric) error {
	keyMap := make(map[*prometheus.Desc]string)

	keyMap[c.netBytesInRate] = "node.net.ext.bytes.in.rate"
	keyMap[c.netBytesOutRate] = "node.net.ext.bytes.out.rate"
	keyMap[c.netErrorsInRate] = "node.net.ext.errors.in.rate"
	keyMap[c.netErrorsOutRate] = "node.net.ext.errors.out.rate"

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
