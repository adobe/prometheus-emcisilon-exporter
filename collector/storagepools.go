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
	"strconv"

	"github.com/adobe/prometheus-emcisilon-exporter/isiclient"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

type storagePoolsCollector struct {
	storagePoolTotal               *prometheus.Desc
	storagePoolManual              *prometheus.Desc
	storagePoolAvailBytes          *prometheus.Desc
	storagePoolAvailSSDBytes       *prometheus.Desc
	storagePoolBalaced             *prometheus.Desc
	storagePoolFreeBytes           *prometheus.Desc
	storagePoolFreeSSDBytes        *prometheus.Desc
	storagePoolTotalBytes          *prometheus.Desc
	storagePoolTotalSSDBytes       *prometheus.Desc
	storagePoolVirtalHotSpareBytes *prometheus.Desc
}

func init() {
	registerCollector("storage_pools", defaultEnabled, NewStoragePoolsCollector)
}

//NewStoragePoolsCollector exposed various metrics and information about storage pools.
func NewStoragePoolsCollector() (Collector, error) {
	return &storagePoolsCollector{
		storagePoolTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "storage_pool", "total"),
			"Total number of storage pools on a cluster.",
			nil, ConstLabels,
		),
		storagePoolManual: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "storage_pool", "manual"),
			"0 of storage pool is not manually managed, 1 is it is.",
			[]string{"name"}, ConstLabels,
		),
		storagePoolAvailBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "storage_pool", "bytes_avail"),
			"Number of bytes available on the storage pool.",
			[]string{"name"}, ConstLabels,
		),
		storagePoolAvailSSDBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "storage_pool", "bytes_avail_ssd"),
			"Number of bytes available on ssd for the storage pool.",
			[]string{"name"}, ConstLabels,
		),
		storagePoolBalaced: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "storage_pool", "balanced"),
			"0 if the storage pool is balanced, 1 if it is not.",
			[]string{"name"}, ConstLabels,
		),
		storagePoolFreeBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "storage_pool", "bytes_free"),
			"Number of bytes available on the storage pool.",
			[]string{"name"}, ConstLabels,
		),
		storagePoolFreeSSDBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "storage_pool", "bytes_free_ssd"),
			"Number of bytes free on ssd for the storage pool.",
			[]string{"name"}, ConstLabels,
		),
		storagePoolTotalBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "storage_pool", "bytes_total"),
			"Total number of bytes on the storage pool.",
			[]string{"name"}, ConstLabels,
		),
		storagePoolTotalSSDBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "storage_pool", "bytes_total_ssd"),
			"Total number of bytes on ssd for the storage pool.",
			[]string{"name"}, ConstLabels,
		),
		storagePoolVirtalHotSpareBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "storage_pool", "bytes_virtual_hot_spare"),
			"Number of bytes in vhs for the storage pool.",
			[]string{"name"}, ConstLabels,
		),
	}, nil
}

func (c *storagePoolsCollector) Update(ch chan<- prometheus.Metric) error {
	resp, err := isiclient.GetStoragePools(IsiCluster.Client)
	if err != nil {
		return err
	}
	ch <- prometheus.MustNewConstMetric(c.storagePoolTotal, prometheus.GaugeValue, resp.Total)

	for _, pool := range resp.Storagepools {
		var (
			manual   float64
			balanced float64
		)
		if pool.Manual {
			manual = 1
		} else {
			manual = 0
		}
		ch <- prometheus.MustNewConstMetric(c.storagePoolManual, prometheus.GaugeValue, manual, pool.Name)

		if pool.Usage.Balanced {
			balanced = 0
		} else {
			balanced = 1
		}
		ch <- prometheus.MustNewConstMetric(c.storagePoolBalaced, prometheus.GaugeValue, balanced, pool.Name)

		bAvail, err := strconv.Atoi(pool.Usage.AvailBytes)
		if err != nil {
			log.Warn("Unable to convert availbytes to int.")
		} else {
			ch <- prometheus.MustNewConstMetric(c.storagePoolAvailBytes, prometheus.GaugeValue, float64(bAvail), pool.Name)
		}

		bAvailSSD, err := strconv.Atoi(pool.Usage.AvailSsdBytes)
		if err != nil {
			log.Warn("Unable to convert availbytesssd to int.")
		} else {
			ch <- prometheus.MustNewConstMetric(c.storagePoolAvailSSDBytes, prometheus.GaugeValue, float64(bAvailSSD), pool.Name)
		}

		bFree, err := strconv.Atoi(pool.Usage.FreeBytes)
		if err != nil {
			log.Warn("Unable to convert freebytes to int.")
		} else {
			ch <- prometheus.MustNewConstMetric(c.storagePoolFreeBytes, prometheus.GaugeValue, float64(bFree), pool.Name)
		}

		bFreeSSD, err := strconv.Atoi(pool.Usage.FreeSsdBytes)
		if err != nil {
			log.Warn("Unable to convert freebytesssd to int.")
		} else {
			ch <- prometheus.MustNewConstMetric(c.storagePoolFreeSSDBytes, prometheus.GaugeValue, float64(bFreeSSD), pool.Name)
		}

		bTotal, err := strconv.Atoi(pool.Usage.TotalBytes)
		if err != nil {
			log.Warn("Unable to convert totalbytes to int.")
		} else {
			ch <- prometheus.MustNewConstMetric(c.storagePoolTotalBytes, prometheus.GaugeValue, float64(bTotal), pool.Name)
		}

		bTotalSsd, err := strconv.Atoi(pool.Usage.TotalSsdBytes)
		if err != nil {
			log.Warn("Unable to convert totalssdbytes to int.")
		} else {
			ch <- prometheus.MustNewConstMetric(c.storagePoolTotalSSDBytes, prometheus.GaugeValue, float64(bTotalSsd), pool.Name)
		}

		bVHS, err := strconv.Atoi(pool.Usage.VirtualHotSpareBytes)
		if err != nil {
			log.Warn("Unable to convert virtualhotsparebytes to int.")
		} else {
			ch <- prometheus.MustNewConstMetric(c.storagePoolVirtalHotSpareBytes, prometheus.GaugeValue, float64(bVHS), pool.Name)
		}

	}

	return err
}
