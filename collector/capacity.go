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
	var errCount int64
	keyMap := make(map[*prometheus.Desc]string)

	keyMap[c.bytesTotal] = "ifs.bytes.total"
	keyMap[c.bytesUsed] = "ifs.bytes.used"
	keyMap[c.bytesAvail] = "ifs.bytes.avail"
	keyMap[c.bytesFree] = "ifs.bytes.free"
	keyMap[c.percentUsed] = "ifs.percent.used"
	keyMap[c.percentAvail] = "ifs.percent.avail"
	keyMap[c.percentFree] = "ifs.percent.free"

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
				val := stat.Value
				ch <- prometheus.MustNewConstMetric(promStat, prometheus.GaugeValue, val)
			}
		}
	}
	if errCount != 0 {
		err := fmt.Errorf("There where %v errors", errCount)
		return err
	}
	return nil
}
