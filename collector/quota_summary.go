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

type quotaSummaryCollector struct {
	quotaSummaryTotalCount        *prometheus.Desc
	quotaSummaryDefaultGroupCount *prometheus.Desc
	quotaSummaryDefaultUserCount  *prometheus.Desc
	quotaSummaryDirectoryCount    *prometheus.Desc
	quotaSummaryGroupCount        *prometheus.Desc
	quotaSummaryLinkedCount       *prometheus.Desc
	quotaSummaryUserCount         *prometheus.Desc
}

func init() {
	registerCollector("quota_summary", defaultEnabled, NewQuotaSummaryCollector)
}

//NewQuotaCollector returns a new Collector exposing node health information.
func NewQuotaSummaryCollector() (Collector, error) {
	return &quotaSummaryCollector{
		quotaSummaryTotalCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, quotaCollectorSubsystem, "summary_total_quotas_count"),
			"Total number of quotas on a cluster.",
			nil, ConstLabels,
		),
		quotaSummaryDefaultGroupCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, quotaCollectorSubsystem, "summary_default_group_quotas_count"),
			"Number of default group quotas.",
			nil, ConstLabels,
		),
		quotaSummaryDefaultUserCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, quotaCollectorSubsystem, "summary_default_user_quotas_count"),
			"Number of default user quotas.",
			nil, ConstLabels,
		),
		quotaSummaryDirectoryCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, quotaCollectorSubsystem, "summary_directory_quotas_count"),
			"Number of directory quotas.",
			nil, ConstLabels,
		),
		quotaSummaryGroupCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, quotaCollectorSubsystem, "summary_group_quotas_count"),
			"Number of group quotas.",
			nil, ConstLabels,
		),
		quotaSummaryLinkedCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, quotaCollectorSubsystem, "summary_linked_quotas_count"),
			"Number of linked quotas.",
			nil, ConstLabels,
		),
		quotaSummaryUserCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, quotaCollectorSubsystem, "summary_quotas_user"),
			"Number of user quotas.",
			nil, ConstLabels,
		),
	}, nil
}

func (c *quotaSummaryCollector) Update(ch chan<- prometheus.Metric) error {
	//Get quota summary statistics
	err := c.updateQuotaSummary(ch)
	if err != nil {
		log.Warn("Unable to collect quota summary information.")
		return err
	}
	return nil
}

func (c *quotaSummaryCollector) updateQuotaSummary(ch chan<- prometheus.Metric) error {
	summary, err := isiclient.GetQuotaSummary(IsiCluster.Client)
	if err != nil {
		log.Warn("Unabled to update quota summary information.")
		return err
	}

	ch <- prometheus.MustNewConstMetric(c.quotaSummaryTotalCount, prometheus.GaugeValue, summary.Count)
	ch <- prometheus.MustNewConstMetric(c.quotaSummaryDefaultGroupCount, prometheus.GaugeValue, summary.DefaultGroupQuotasCount)
	ch <- prometheus.MustNewConstMetric(c.quotaSummaryDefaultUserCount, prometheus.GaugeValue, summary.DefaultUserQuotasCount)
	ch <- prometheus.MustNewConstMetric(c.quotaSummaryDirectoryCount, prometheus.GaugeValue, summary.DirectoryQuotasCount)
	ch <- prometheus.MustNewConstMetric(c.quotaSummaryGroupCount, prometheus.GaugeValue, summary.GroupQuotasCount)
	ch <- prometheus.MustNewConstMetric(c.quotaSummaryLinkedCount, prometheus.GaugeValue, summary.LinkedQuotasCount)
	ch <- prometheus.MustNewConstMetric(c.quotaSummaryUserCount, prometheus.GaugeValue, summary.UserQuotasCount)
	return nil
}
