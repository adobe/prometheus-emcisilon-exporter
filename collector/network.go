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

type networkCollector struct {
	netBytesInRate  *prometheus.Desc
	netBytesOutRate *prometheus.Desc
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
	}, nil
}

func (c *networkCollector) Update(ch chan<- prometheus.Metric) error {
	keyMap := make(map[string]string)

	keyMap["netBytesInRate"] = "node.net.ext.bytes.in.rate"
	keyMap["netBytesOutRate"] = "node.net.ext.bytes.out.rate"

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
			case "netBytesInRate":
				ch <- prometheus.MustNewConstMetric(c.netBytesInRate, prometheus.GaugeValue, val, node)
			case "netBytesOutRate":
				ch <- prometheus.MustNewConstMetric(c.netBytesOutRate, prometheus.GaugeValue, val, node)
			}
		}
	}
	return nil
}
