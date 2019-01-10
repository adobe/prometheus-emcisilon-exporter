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

type diskCollector struct {
	diskBusyAll         *prometheus.Desc
	diskIoschedQueueAll *prometheus.Desc
	diskXfersInRateAll  *prometheus.Desc
	diskXfersOutRateAll *prometheus.Desc
	diskLatencyAll      *prometheus.Desc
}

func init() {
	registerCollector("disk", defaultEnabled, NewDiskCollector)
}

//NewDiskCollector returns a new Collector exposing node disk statistics.
func NewDiskCollector() (Collector, error) {
	return &diskCollector{
		diskBusyAll: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "disk_busy_all"),
			"Current disk busy percentage represented in 0.0-1.0.",
			[]string{"node", "disk"}, ConstLabels,
		),
		diskIoschedQueueAll: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "disk_iosched_queued_all"),
			"Current queue depth for IO sceduler.",
			[]string{"node", "disk"}, ConstLabels,
		),
		diskXfersInRateAll: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "disk_xfers_in_rate_all"),
			"Current disk ingest transfer rate.",
			[]string{"node", "disk"}, ConstLabels,
		),
		diskXfersOutRateAll: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "disk_xfers_out_rate_all"),
			"Current disk egress transfer rate.",
			[]string{"node", "disk"}, ConstLabels,
		),
		diskLatencyAll: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "disk_latency_all"),
			"Current disk latency.",
			[]string{"node", "disk"}, ConstLabels,
		),
	}, nil
}

func (c *diskCollector) Update(ch chan<- prometheus.Metric) error {
	var errCount int64
	keyMap := make(map[*prometheus.Desc]string)

	keyMap[c.diskBusyAll] = "node.disk.busy.all"
	keyMap[c.diskIoschedQueueAll] = "node.disk.iosched.queue.all"
	keyMap[c.diskXfersInRateAll] = "node.disk.xfers.in.rate.all"
	keyMap[c.diskXfersOutRateAll] = "node.disk.xfers.out.rate.all"
	keyMap[c.diskLatencyAll] = "node.disk.access.latency.all"

	for promStat, statKey := range keyMap {
		begin := time.Now()
		resp, err := isiclient.QueryStatsEngineMultiVal(IsiCluster.Client, statKey)
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
				for _, valset := range stat.ValueSet {
					for disk, val := range valset {
						if statKey == "node.disk.busy.all" {
							val = val / 10
						}
						ch <- prometheus.MustNewConstMetric(promStat, prometheus.GaugeValue, val, node, disk)
					}
				}
			}
		}
	}
	if errCount != 0 {
		err := fmt.Errorf("There where %v errors", errCount)
		return err
	}
	return nil
}
