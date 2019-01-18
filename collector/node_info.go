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
	nodeDriveState  *prometheus.Desc
	nodeInfo        *prometheus.Desc
}

var state float64

func init() {
	registerCollector("node_info", defaultEnabled, NewNodeStatusCollector)
}

//NewNodeStatusCollector exposed various metrics and information about nodes.
func NewNodeStatusCollector() (Collector, error) {
	return &nodeStatusCollector{
		nodeInfo: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "info"),
			"Contains information about each node in labels. Always returns a 1.",
			[]string{"id", "infiniband", "motherboard", "generation_code", "chassis_code", "lnn", "hwgen", "nvram", "chassis_count", "serial_number", "disk_expander", "disk_collector", "family_code", "product", "class", "cpu", "chassis", "proc_count", "proc_type", "name"}, ConstLabels,
		),
		nodeBattery: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "status_battery"),
			"Status for batteries.",
			[]string{"node", "node_id", "result1", "result2"}, ConstLabels,
		),
		nodePowerSupply: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "status_power_supply"),
			"Status for power supplies.",
			[]string{"node", "node_id", "power_supply", "status"}, ConstLabels,
		),
		nodeDriveState: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "drive_state"),
			"Current state of the drive in a bay. 0 = HEALTHY/L3, 1 = STALLED, 2 = FW_UPDATE, 3 = SMARTFAILED, 4 = USED, 5 = PREPARING, 10 = NEW, 11 = EMPTY, 12 = REPLACE, 99 = UNKNOWN.",
			[]string{"node", "node_id", "bay_num", "media_type", "model", "interaface_type", "dev_name", "state"}, ConstLabels,
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
	err = c.updateDriveStatus(ch)
	if err != nil {
		log.Warnf("Unable to update drive states. %s", err)
	}
	err = c.updateNodeInfo(ch)
	if err != nil {
		log.Warnf("Unable to update node info: %s", err)
	}
	return err
}

func (c *nodeStatusCollector) updateBatteryStatus(ch chan<- prometheus.Metric, nodeStatus isiclient.IsiNodesStatus) error {
	for _, node := range nodeStatus.Nodes {
		var status float64
		nodeID := fmt.Sprintf("%v", node.ID)
		nodeLNN := fmt.Sprintf("%v", node.Lnn)

		b1, b2 := c.checkStatus(node.Batterystatus.Result1), c.checkStatus(node.Batterystatus.Result2)

		if v := b1 + b2; v == 0 {
			status = 0
		} else {
			status = 1
		}
		ch <- prometheus.MustNewConstMetric(c.nodeBattery, prometheus.GaugeValue, status, nodeLNN, nodeID, node.Batterystatus.Result1, node.Batterystatus.Result2)
	}
	return nil
}

func (c *nodeStatusCollector) updatePowerSupplyStatus(ch chan<- prometheus.Metric, nodeStatus isiclient.IsiNodesStatus) error {
	for _, node := range nodeStatus.Nodes {
		nodeID := fmt.Sprintf("%v", node.ID)
		nodeLNN := fmt.Sprintf("%v", node.Lnn)
		for _, powerSupply := range node.Powersupplies.Supplies {
			supplyID := fmt.Sprintf("%v", powerSupply.ID)
			var status float64
			if powerSupply.Good == "Good" {
				status = 0
			} else {
				status = 1
			}
			ch <- prometheus.MustNewConstMetric(c.nodePowerSupply, prometheus.GaugeValue, status, nodeLNN, nodeID, supplyID, powerSupply.Status)
		}
	}
	return nil
}

func (c *nodeStatusCollector) updateDriveStatus(ch chan<- prometheus.Metric) error {
	resp, err := isiclient.GetDriveInfo(IsiCluster.Client)
	if err != nil {
		log.Warnf("Unabled to collect drive status. %s", err)
		return err
	}
	for _, node := range resp.Nodes {
		nodeID := fmt.Sprintf("%v", node.ID)
		nodeLNN := fmt.Sprintf("%v", node.Lnn)
		for _, drive := range node.Drives {
			bayID := fmt.Sprintf("%v", drive.Baynum)
			devID := fmt.Sprintf("%v", drive.Devname)

			switch drive.UIState {
			case "HEALTHY", "L3":
				state = 0
			case "STALLED":
				state = 1
			case "FW_UPDATE":
				state = 2
			case "SMARTFAIL":
				state = 3
			case "USED":
				state = 4
			case "PREPARING":
				state = 5
			case "NEW":
				state = 10
			case "EMPTY":
				state = 11
			case "REPLACE":
				state = 12
			default:
				state = 99
			}
			ch <- prometheus.MustNewConstMetric(c.nodeDriveState, prometheus.GaugeValue, state, nodeLNN, nodeID, bayID, drive.MediaType, drive.Model, drive.InterfaceType, devID, drive.UIState)
		}
	}
	return nil
}

func (c *nodeStatusCollector) updateNodeInfo(ch chan<- prometheus.Metric) error {
	var na = "n/a"
	resp, err := isiclient.GetNodesHardware(IsiCluster.Client)
	if err != nil {
		return fmt.Errorf("Unable to collect hardware info. %s", err)
	}
	for _, node := range resp.Nodes {
		nodeID := fmt.Sprintf("%v", node.ID)
		infini, err := c.labelTrimmer(node.Infiniband)
		if err != nil {
			log.Warnf("Unable to trim label: %s", node.Infiniband)
			infini = na
		}
		mobo, err := c.labelTrimmer(node.Motherboard)
		if err != nil {
			log.Warnf("Unable to trim label: %s", node.Motherboard)
			mobo = na
		}
		lnnID := fmt.Sprintf("%v", node.Lnn)
		hwgen, err := c.labelTrimmer(node.Hwgen)
		if err != nil {
			log.Warnf("Unable to trim label: %s", node.Hwgen)
			hwgen = na
		}
		nvram, err := c.labelTrimmer(node.Nvram)
		if err != nil {
			log.Warnf("Unable to trim label: %s", node.Nvram)
			nvram = na
		}
		chassisCount, err := c.labelTrimmer(node.ChassisCount)
		if err != nil {
			log.Warnf("Unable to trim label: %s", node.ChassisCount)
			chassisCount = na
		}
		diskExp, err := c.labelTrimmer(node.DiskExpander)
		if err != nil {
			log.Warnf("Unable to trim label: %s", node.DiskExpander)
			diskExp = na
		}
		diskCtl, err := c.labelTrimmer(node.DiskController)
		if err != nil {
			log.Warnf("Unable to trim label: %s", node.DiskController)
			diskCtl = na
		}
		product, err := c.labelTrimmer(node.Product)
		if err != nil {
			log.Warnf("Unable to trim label: %s", node.Product)
			product = na
		}
		cpu, err := c.labelTrimmer(node.CPU)
		if err != nil {
			log.Warnf("Unable to trim label: %s", node.CPU)
			cpu = na
		}
		chassis, err := c.labelTrimmer(node.Chassis)
		if err != nil {
			log.Warnf("Unable to trim label: %s", node.Chassis)
			chassis = na
		}
		procCount := na
		procType := na
		subProcs := strings.Split(node.Processor, ",")
		if len(subProcs) == 2 {
			procCount = subProcs[0]
			procType = strings.TrimSpace(subProcs[1])
		}
		name := fmt.Sprintf("%v-%v", IsiCluster.Name, lnnID)
		ch <- prometheus.MustNewConstMetric(c.nodeInfo, prometheus.GaugeValue, float64(1), nodeID, infini, mobo, node.GenerationCode, node.ChassisCode, lnnID, hwgen, nvram, chassisCount, node.SerialNumber, diskExp, diskCtl, node.FamilyCode, product, node.Class, cpu, chassis, procCount, procType, name)
	}
	return nil
}

func (c *nodeStatusCollector) labelTrimmer(label string) (string, error) {
	substrings := strings.Split(label, " ")
	if len(substrings) == 1 {
		return label, nil
	}

	if len(substrings) < 2 {
		return label, fmt.Errorf("Unable to split label string")
	}
	return substrings[0], nil
}

func (c *nodeStatusCollector) checkStatus(status string) float64 {
	switch status {
	case "passed", "N/A":
		return 0
	default:
		return 1
	}
}
