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

type nfsExportsCollector struct {
	exportCount *prometheus.Desc
}

func init() {
	registerCollector("nfs_exports", defaultEnabled, NewNfsExportsCollector)
}

//NewNfsExportsCollector exposed various metrics and information about nodes.
func NewNfsExportsCollector() (Collector, error) {
	return &nfsExportsCollector{
		exportCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "nfs", "export_total"),
			"Total number of NFS exports on a cluster.",
			nil, ConstLabels,
		),
	}, nil
}

func (c *nfsExportsCollector) Update(ch chan<- prometheus.Metric) error {
	resp, err := isiclient.GetExportSummary(IsiCluster.Client)
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(c.exportCount, prometheus.GaugeValue, resp.Summary.Count)
	return err
}
