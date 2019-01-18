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
	"strconv"
	"strings"

	"github.com/adobe/prometheus-emcisilon-exporter/isiclient"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

type nodePartitionCollector struct {
	nodePartitionUsedSpacePercentage  *prometheus.Desc
	nodePartitionCount                *prometheus.Desc
	nodePartitionFileNodesFree        *prometheus.Desc
	nodePartitionFileNodesTotal       *prometheus.Desc
	nodePartitionFileNodesFreePercent *prometheus.Desc
}

func init() {
	registerCollector("node_partition", defaultEnabled, NewNodePartitionCollector)
}

//NewNodePartitionCollector exposed various metrics and information about nodes.
func NewNodePartitionCollector() (Collector, error) {
	return &nodePartitionCollector{
		nodePartitionUsedSpacePercentage: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "partition_used_space_percentage"),
			"Percentage of space used on a partition.",
			[]string{"node", "node_id", "mount_point"}, ConstLabels,
		),
		nodePartitionCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "partition_count"),
			"Count of the total number of partitions on a node.",
			[]string{"node", "node_id"}, ConstLabels,
		),
		nodePartitionFileNodesFree: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "partition_filenodes_free"),
			"Number of filenodes free on a partition.",
			[]string{"node", "node_id", "mount_point"}, ConstLabels,
		),
		nodePartitionFileNodesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "partition_filenodes_total"),
			"Total number of filenodes on a partition.",
			[]string{"node", "node_id", "mount_point"}, ConstLabels,
		),
		nodePartitionFileNodesFreePercent: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "partition_filenodes_free_percent"),
			"Percentage of filenodes free on a partition.",
			[]string{"node", "node_id", "mount_point"}, ConstLabels,
		),
	}, nil
}

func (c *nodePartitionCollector) Update(ch chan<- prometheus.Metric) error {
	err := c.updatePartitionStats(ch)
	return err
}

func (c *nodePartitionCollector) updatePartitionStats(ch chan<- prometheus.Metric) error {
	resp, err := isiclient.GetNodesPartitions(IsiCluster.Client)
	if err != nil {
		return err
	}

	for _, node := range resp.Nodes {
		nodeID := fmt.Sprintf("%v", node.ID)
		nodeLNN := fmt.Sprintf("%v", node.Lnn)
		ch <- prometheus.MustNewConstMetric(c.nodePartitionCount, prometheus.GaugeValue, node.Count, nodeLNN, nodeID)
		for _, partition := range node.Partitions {
			if strings.Contains(partition.MountPoint, "Unknown") {
				continue
			}
			used, err := strconv.Atoi(strings.Replace(partition.PercentUsed, "%", "", -1))
			if err != nil {
				log.Infof("Error converting used percentage: %s", err)
			}
			usedPercentage := float64(used) / 100.0
			ch <- prometheus.MustNewConstMetric(c.nodePartitionUsedSpacePercentage, prometheus.GaugeValue, usedPercentage, nodeLNN, nodeID, partition.MountPoint)

			percentFree := partition.Statfs.FFfree / partition.Statfs.FFiles
			ch <- prometheus.MustNewConstMetric(c.nodePartitionFileNodesFreePercent, prometheus.GaugeValue, percentFree, nodeLNN, nodeID, partition.MountPoint)
			ch <- prometheus.MustNewConstMetric(c.nodePartitionFileNodesFree, prometheus.GaugeValue, partition.Statfs.FFfree, nodeLNN, nodeID, partition.MountPoint)
			ch <- prometheus.MustNewConstMetric(c.nodePartitionFileNodesTotal, prometheus.GaugeValue, partition.Statfs.FFiles, nodeLNN, nodeID, partition.MountPoint)
		}
	}
	return err
}
