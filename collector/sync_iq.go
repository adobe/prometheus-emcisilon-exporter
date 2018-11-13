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

type syncIQPoliciesCollector struct {
	syncPolicyState          *prometheus.Desc
	syncPolicyLastSuccess    *prometheus.Desc
	syncPolicyLastStart      *prometheus.Desc
	syncPolicyPriority       *prometheus.Desc
	syncPolicyWorkersPerNode *prometheus.Desc
	syncPolicyEnabled        *prometheus.Desc
	syncPolicyTotalCount     *prometheus.Desc
}

func init() {
	registerCollector("sync_iq", defaultEnabled, NewSyncIQCollector)
}

//NewSyncIQCollector returns a new Collector exposing sync IQ policy information.
func NewSyncIQCollector() (Collector, error) {
	return &syncIQPoliciesCollector{
		syncPolicyState: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "sync", "policy_state"),
			"Last state from run of sync policy.",
			[]string{"name"}, ConstLabels,
		),
		syncPolicyLastSuccess: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "sync", "policy_last_success"),
			"Epoch timestamp of the last successful sync for a policy.",
			[]string{"name"}, ConstLabels,
		),
		syncPolicyLastStart: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "sync", "policy_last_start"),
			"Epoch timestame for last sync start for a policy.",
			[]string{"name"}, ConstLabels,
		),
		syncPolicyPriority: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "sync", "policy_priority"),
			"Current priority for the policy.",
			[]string{"name"}, ConstLabels,
		),
		syncPolicyWorkersPerNode: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "sync", "policy_workers_per_node"),
			"Number of worker threads per node for a policy.",
			[]string{"name"}, ConstLabels,
		),
		syncPolicyEnabled: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "sync", "policy_enabled"),
			"1 = Enabled, 0 = Disabled for the specified policy",
			[]string{"name"}, ConstLabels,
		),
		syncPolicyTotalCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "sync", "policies_total_count"),
			"Total number of sync policies on the cluster.", nil, ConstLabels,
		),
	}, nil
}

func (c *syncIQPoliciesCollector) Update(ch chan<- prometheus.Metric) error {
	resp, err := isiclient.GetSyncPolicies(IsiCluster.Client)
	if err != nil {
		log.Warnf("Error attempting to view sync policies.")
	}
	ch <- prometheus.MustNewConstMetric(c.syncPolicyTotalCount, prometheus.GaugeValue, float64(len(resp.Policies)))
	for _, policy := range resp.Policies {
		var enabled float64
		if policy.Enabled {
			enabled = 1
		} else {
			enabled = 0
		}
		ch <- prometheus.MustNewConstMetric(c.syncPolicyEnabled, prometheus.GaugeValue, enabled, policy.Name)
		ch <- prometheus.MustNewConstMetric(c.syncPolicyLastStart, prometheus.GaugeValue, policy.LastStarted, policy.Name)
		ch <- prometheus.MustNewConstMetric(c.syncPolicyLastSuccess, prometheus.GaugeValue, policy.LastSuccess, policy.Name)
		ch <- prometheus.MustNewConstMetric(c.syncPolicyPriority, prometheus.GaugeValue, policy.Priority, policy.Name)

		var state float64
		if policy.LastJobState == "finished" {
			state = 0
		} else {
			state = 1
		}
		ch <- prometheus.MustNewConstMetric(c.syncPolicyState, prometheus.GaugeValue, state, policy.Name)
		ch <- prometheus.MustNewConstMetric(c.syncPolicyWorkersPerNode, prometheus.GaugeValue, policy.WorkersPerNode, policy.Name)
	}
	return nil
}
