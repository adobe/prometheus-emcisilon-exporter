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
)

type statfsCollector struct {
	statfsFileBlockAvail      *prometheus.Desc
	statfsFileBlockFree       *prometheus.Desc
	statfsFileBlockTotal      *prometheus.Desc
	statfsFileBlockSize       *prometheus.Desc
	statfsFileNodeFree        *prometheus.Desc
	statfsFileNodeTotal       *prometheus.Desc
	statfsFileNodeFreePercent *prometheus.Desc
	statfsFileIOSize          *prometheus.Desc
	statfsFileNameMax         *prometheus.Desc
}

func init() {
	registerCollector("statfs", defaultEnabled, NewStatfsCollector)
}

//NewStatfsCollector exposed various metrics and information about nodes.
func NewStatfsCollector() (Collector, error) {
	return &statfsCollector{
		statfsFileBlockAvail: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "stafs", "file_block_avail"),
			"The filesystem fragment size.",
			[]string{"mount_point"}, ConstLabels,
		),
		statfsFileBlockFree: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "statfs", "file_block_free"),
			"The number of free blocks in the filesystem.",
			[]string{"mount_point"}, ConstLabels,
		),
		statfsFileBlockSize: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "statfs", "file_block_size"),
			"The filesystem fragment size.",
			[]string{"mount_point"}, ConstLabels,
		),
		statfsFileBlockTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "statfs", "file_block_used"),
			"The total number of data blocks in the filesystem.",
			[]string{"mount_point"}, ConstLabels,
		),
		statfsFileIOSize: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "statfs", "file_io_size"),
			"The optimal transfer block size.",
			[]string{"mount_point"}, ConstLabels,
		),
		statfsFileNameMax: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "statfs", "file_name_max"),
			"The maximum length of a file name.",
			[]string{"mount_point"}, ConstLabels,
		),
		statfsFileNodeTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "statfs", "file_node_total"),
			"The total number of file nodes in the filesystem.",
			[]string{"mount_point"}, ConstLabels,
		),
		statfsFileNodeFree: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "statfs", "file_node_free"),
			"The number of free blocks in the filesystem.",
			[]string{"mount_point"}, ConstLabels,
		),
		statfsFileNodeFreePercent: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "statfs", "file_node_free_percent"),
			"The percentage of free file nodes in the filesystem.",
			[]string{"mount_point"}, ConstLabels,
		),
	}, nil
}

func (c *statfsCollector) Update(ch chan<- prometheus.Metric) error {
	resp, err := isiclient.GetStatfs(IsiCluster.Client)
	if err != nil {
		return err
	}

	percentFree := resp.FFfree / resp.FFiles

	ch <- prometheus.MustNewConstMetric(c.statfsFileBlockAvail, prometheus.GaugeValue, resp.FBavail, resp.FMntonname)
	ch <- prometheus.MustNewConstMetric(c.statfsFileBlockFree, prometheus.GaugeValue, resp.FBfree, resp.FMntonname)
	ch <- prometheus.MustNewConstMetric(c.statfsFileBlockSize, prometheus.GaugeValue, resp.FBsize, resp.FMntonname)
	ch <- prometheus.MustNewConstMetric(c.statfsFileBlockTotal, prometheus.GaugeValue, resp.FBlocks, resp.FMntonname)
	ch <- prometheus.MustNewConstMetric(c.statfsFileIOSize, prometheus.GaugeValue, resp.FIosize, resp.FMntonname)
	ch <- prometheus.MustNewConstMetric(c.statfsFileNameMax, prometheus.GaugeValue, resp.FNamemax, resp.FMntonname)
	ch <- prometheus.MustNewConstMetric(c.statfsFileNodeFree, prometheus.GaugeValue, resp.FFfree, resp.FMntonname)
	ch <- prometheus.MustNewConstMetric(c.statfsFileNodeFreePercent, prometheus.GaugeValue, float64(percentFree), resp.FMntonname)
	ch <- prometheus.MustNewConstMetric(c.statfsFileNodeTotal, prometheus.GaugeValue, resp.FFiles, resp.FMntonname)
	return err
}
