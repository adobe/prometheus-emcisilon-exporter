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
	"strings"

	"github.com/adobe/prometheus-emcisilon-exporter/isiclient"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

type nodeStatusCollector struct {
	nodeBattery     *prometheus.Desc
	nodePowerSupply *prometheus.Desc
}

func init() {
	registerCollector("node_status", defaultEnabled, NewNodeStatusCollector)
}

//NewNodeStatusCollector exposed various metrics and information about nodes.
func NewNodeStatusCollector() (Collector, error) {
	return &nodeStatusCollector{
		nodeBattery: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "status_battery"),
			"Status for batteries.",
			[]string{"node", "battery"}, ConstLabels,
		),
		nodePowerSupply: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "status_power_supply"),
			"Status for power supplies.",
			[]string{"node", "power_supply"}, ConstLabels,
		),
	}, nil
}

func (c *nodeStatusCollector) Update(ch chan<- prometheus.Metric) error {
	resp, err := isiclient.GetNodesStatus(IsiCluster.Client)
	if err != nil {
		log.Warnf("Unable to get node status from API. %s", err)
	}
	err = c.updateBatteryStatus(ch, resp)
	if err != nil {
		log.Warnf("Unable to update battery status. %s", err)
	}
	err = c.updatePowerSupplyStatus(ch, resp)
	if err != nil {
		log.Warnf("Unable to update power supply status. %s", err)
	}
	return err
}

func (c *nodeStatusCollector) updateBatteryStatus(ch chan<- prometheus.Metric, nodeStatus isiclient.IsiNodesStatus) error {
	for _, node := range nodeStatus.Nodes {
		nodeID := fmt.Sprintf("%v", node.ID)

		bs1 := strings.ToLower(node.Batterystatus.Status1)
		bs2 := strings.ToLower(node.Batterystatus.Status2)

		if (strings.Contains(bs1, "good") || strings.Contains(bs1, "ready")) && !(strings.Contains(bs1, "N/A")) {
			ch <- prometheus.MustNewConstMetric(c.nodeBattery, prometheus.GaugeValue, float64(0), nodeID, "1")
		} else if !(strings.Contains(bs1, "n/a")) {
			ch <- prometheus.MustNewConstMetric(c.nodeBattery, prometheus.GaugeValue, float64(1), nodeID, "1")
		}

		if (strings.Contains(bs2, "good") || strings.Contains(bs2, "ready")) && !(strings.Contains(bs2, "N/A")) {
			ch <- prometheus.MustNewConstMetric(c.nodeBattery, prometheus.GaugeValue, float64(0), nodeID, "2")
		} else if !(strings.Contains(bs2, "n/a")) {
			ch <- prometheus.MustNewConstMetric(c.nodeBattery, prometheus.GaugeValue, float64(1), nodeID, "2")
		}
	}
	return nil
}

func (c *nodeStatusCollector) updatePowerSupplyStatus(ch chan<- prometheus.Metric, nodeStatus isiclient.IsiNodesStatus) error {
	for _, node := range nodeStatus.Nodes {
		nodeID := fmt.Sprintf("%v", node.ID)
		for _, powerSupply := range node.Powersupplies.Supplies {
			supplyID := fmt.Sprintf("%v", powerSupply.ID)
			var status float64
			if powerSupply.Good == "Good" {
				status = 0
			} else {
				status = 1
			}
			ch <- prometheus.MustNewConstMetric(c.nodePowerSupply, prometheus.GaugeValue, status, nodeID, supplyID)
		}
	}
	return nil
}
