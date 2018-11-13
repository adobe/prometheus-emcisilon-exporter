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

type capacityCollector struct {
	bytesTotal   *prometheus.Desc
	bytesUsed    *prometheus.Desc
	bytesAvail   *prometheus.Desc
	bytesFree    *prometheus.Desc
	percentUsed  *prometheus.Desc
	percentAvail *prometheus.Desc
	percentFree  *prometheus.Desc
}

const (
	ifsSubSystem = "ifs"
)

func init() {
	registerCollector("capacity", defaultEnabled, NewCapacityCollector)
}

//NewCapacityCollector returns a new Collector exposing cluster capacity/disk space statistics.
func NewCapacityCollector() (Collector, error) {
	return &capacityCollector{
		bytesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, ifsSubSystem, "bytes_total"),
			"Current ifs filesystem capacity total in bytes.",
			nil, ConstLabels,
		),
		bytesUsed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, ifsSubSystem, "bytes_used"),
			"Current ifs filesystem capacity used in bytes.",
			nil, ConstLabels,
		),
		bytesAvail: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, ifsSubSystem, "bytes_avail"),
			"Current ifs filesystem capacity available in bytes.",
			nil, ConstLabels,
		),
		bytesFree: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, ifsSubSystem, "bytes_free"),
			"Current ifs filesystem capacity free in bytes.",
			nil, ConstLabels,
		),
		percentUsed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, ifsSubSystem, "percent_used"),
			"Current ifs filesystem capacity used in as a percentage from 0.0 - 1.0.",
			nil, ConstLabels,
		),
		percentAvail: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, ifsSubSystem, "percent_avail"),
			"Current ifs filesystem capacity available as a percentage from 0.0 - 1.0.",
			nil, ConstLabels,
		),
		percentFree: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, ifsSubSystem, "percent_free"),
			"Current ifs filesystem capacity free as a percentage from 0.0 - 1.0.",
			nil, ConstLabels,
		),
	}, nil
}

func (c *capacityCollector) Update(ch chan<- prometheus.Metric) error {
	keyMap := make(map[string]string)

	keyMap["bytesTotal"] = "ifs.bytes.total"
	keyMap["bytesUsed"] = "ifs.bytes.used"
	keyMap["bytesAvail"] = "ifs.bytes.avail"
	keyMap["bytesFree"] = "ifs.bytes.free"
	keyMap["percentUsed"] = "ifs.percent.used"
	keyMap["percentAvail"] = "ifs.percent.avail"
	keyMap["percentFree"] = "ifs.percent.free"

	for promStat, statKey := range keyMap {
		resp, err := isiclient.QueryStatsEngineSingleVal(IsiCluster.Client, statKey)
		if err != nil {
			log.Warnf("Error attempting to query stats engine with key %s: %s", statKey, err)
		}
		for _, stat := range resp.Stats {
			val := stat.Value
			err = c.updateCapacity(promStat, val, ch)
			if err != nil {
				log.Infof("Unable to update capacity metric for %s", statKey)
			}
		}
	}
	return nil
}

func (c *capacityCollector) updateCapacity(promStat string, val float64, ch chan<- prometheus.Metric) error {
	switch promStat {
	case "bytesTotal":
		ch <- prometheus.MustNewConstMetric(c.bytesTotal, prometheus.GaugeValue, val)
	case "bytesUsed":
		ch <- prometheus.MustNewConstMetric(c.bytesUsed, prometheus.GaugeValue, val)
	case "bytesAvail":
		ch <- prometheus.MustNewConstMetric(c.bytesAvail, prometheus.GaugeValue, val)
	case "bytesFree":
		ch <- prometheus.MustNewConstMetric(c.bytesFree, prometheus.GaugeValue, val)
	case "percentUsed":
		ch <- prometheus.MustNewConstMetric(c.percentUsed, prometheus.GaugeValue, val)
	case "percentAvail":
		ch <- prometheus.MustNewConstMetric(c.percentAvail, prometheus.GaugeValue, val)
	case "percentFree":
		ch <- prometheus.MustNewConstMetric(c.percentFree, prometheus.GaugeValue, val)
	}
	return nil
}
