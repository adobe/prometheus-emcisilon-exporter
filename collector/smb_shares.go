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

type smbSharesCollector struct {
	sharesCount *prometheus.Desc
}

func init() {
	registerCollector("smb_shares", defaultEnabled, NewSmbSharesCollector)
}

//NewSmbSharesCollector exposed various metrics and information about nodes.
func NewSmbSharesCollector() (Collector, error) {
	return &smbSharesCollector{
		sharesCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "smb", "share_total"),
			"Total number of SMB shares on a cluster.",
			nil, ConstLabels,
		),
	}, nil
}

func (c *smbSharesCollector) Update(ch chan<- prometheus.Metric) error {
	resp, err := isiclient.GetSharesSummary(IsiCluster.Client)
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(c.sharesCount, prometheus.GaugeValue, resp.Summary.Count)
	return err
}
