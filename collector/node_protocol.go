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
	"encoding/json"
	"fmt"
	"strings"

	"github.com/adobe/prometheus-emcisilon-exporter/isiclient"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

type nodeProtoCollector struct {
	nodeProtocolInMax        *prometheus.Desc
	nodeProtocolInMin        *prometheus.Desc
	nodeProtocolInRate       *prometheus.Desc
	nodeProtocolOpCount      *prometheus.Desc
	nodeProtocolOpRate       *prometheus.Desc
	nodeProtocolOutMax       *prometheus.Desc
	nodeProtocolOutMin       *prometheus.Desc
	nodeProtocolOutRate      *prometheus.Desc
	nodeProtocolTimeAvg      *prometheus.Desc
	nodeProtocolTimeMax      *prometheus.Desc
	nodeProtocolTimeMin      *prometheus.Desc
	nodeProtocolTotalInMax   *prometheus.Desc
	nodeProtocolTotalInMin   *prometheus.Desc
	nodeProtocolTotalInRate  *prometheus.Desc
	nodeProtocolTotalOpCount *prometheus.Desc
	nodeProtocolTotalOpRate  *prometheus.Desc
	nodeProtocolTotalOutMax  *prometheus.Desc
	nodeProtocolTotalOutMin  *prometheus.Desc
	nodeProtocolTotalOutRate *prometheus.Desc
	nodeProtocolTotalTimeAvg *prometheus.Desc
	nodeProtocolTotalTimeMax *prometheus.Desc
	nodeProtocolTotalTimeMin *prometheus.Desc
	nodeClientsConnected     *prometheus.Desc
	nodeClientsActive        *prometheus.Desc
}

func init() {
	registerCollector("node_protocol", defaultEnabled, NewNodeProtoCollector)
	if !protosUpdated {
		GetProtos()
	}
}

//NewNodeProtoCollector returns a new Collector exposing Node protocol statistics.
func NewNodeProtoCollector() (Collector, error) {
	return &nodeProtoCollector{
		nodeProtocolInMax: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "protostats_in_max"),
			"Node protocol operation in max.",
			[]string{"node", "proto", "op"}, ConstLabels,
		),
		nodeProtocolInMin: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "protostats_in_min"),
			"Node protocol operation in min.",
			[]string{"node", "proto", "op"}, ConstLabels,
		),
		nodeProtocolInRate: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "protostats_in_rate"),
			"Node protocol operation in rate.",
			[]string{"node", "proto", "op"}, ConstLabels,
		),
		nodeProtocolOpCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "protostats_op_count"),
			"Node protocol operation count.",
			[]string{"node", "proto", "op"}, ConstLabels,
		),
		nodeProtocolOpRate: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "protostats_op_rate"),
			"Node protocol operation rate.",
			[]string{"node", "proto", "op"}, ConstLabels,
		),
		nodeProtocolOutMax: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "protostats_out_max"),
			"Node protocol operation out max.",
			[]string{"node", "proto", "op"}, ConstLabels,
		),
		nodeProtocolOutMin: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "protostats_out_min"),
			"Node protocol operation out min.",
			[]string{"node", "proto", "op"}, ConstLabels,
		),
		nodeProtocolOutRate: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "protostats_out_rate"),
			"Node protocol operation out rate.",
			[]string{"node", "proto", "op"}, ConstLabels,
		),
		nodeProtocolTimeAvg: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "protostats_time_avg"),
			"Node protocol operation time average.",
			[]string{"node", "proto", "op"}, ConstLabels,
		),
		nodeProtocolTimeMax: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "protostats_time_max"),
			"Node protocol operation time max.",
			[]string{"node", "proto", "op"}, ConstLabels,
		),
		nodeProtocolTimeMin: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "protostats_time_min"),
			"Node protocol operation in rate.",
			[]string{"node", "proto", "op"}, ConstLabels,
		),
		nodeProtocolTotalInMax: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "protostats_in_max_total"),
			"Total node protocol operation in max.",
			[]string{"node", "proto"}, ConstLabels,
		),
		nodeProtocolTotalInMin: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "protostats_in_min_total"),
			"Total node protocol operation in min.",
			[]string{"node", "proto"}, ConstLabels,
		),
		nodeProtocolTotalInRate: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "protostats_in_rate_total"),
			"Total node protocol operation in rate.",
			[]string{"node", "proto"}, ConstLabels,
		),
		nodeProtocolTotalOpCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "protostats_op_count_total"),
			"Total node protocol operation count.",
			[]string{"node", "proto"}, ConstLabels,
		),
		nodeProtocolTotalOpRate: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "protostats_op_rate_total"),
			"Total node protocol operation rate.",
			[]string{"node", "proto"}, ConstLabels,
		),
		nodeProtocolTotalOutMax: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "protostats_out_max_total"),
			"Total node protocol operation out max.",
			[]string{"node", "proto"}, ConstLabels,
		),
		nodeProtocolTotalOutMin: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "protostats_out_min_total"),
			"Node protocol operation out min.",
			[]string{"node", "proto"}, ConstLabels,
		),
		nodeProtocolTotalOutRate: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "protostats_out_rate_total"),
			"Total node protocol operation out rate.",
			[]string{"node", "proto"}, ConstLabels,
		),
		nodeProtocolTotalTimeAvg: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "protostats_time_avg_total"),
			"Total node protocol operation time average.",
			[]string{"node", "proto"}, ConstLabels,
		),
		nodeProtocolTotalTimeMax: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "protostats_time_max_total"),
			"Total node protocol operation time max.",
			[]string{"node", "proto"}, ConstLabels,
		),
		nodeProtocolTotalTimeMin: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "protostats_time_min_total"),
			"Total node protocol operation in rate.",
			[]string{"node", "proto"}, ConstLabels,
		),
		nodeClientsConnected: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "clientstats_connected"),
			"Total node protocol operation in rate.",
			[]string{"node", "proto"}, ConstLabels,
		),
		nodeClientsActive: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeCollectorSubsystem, "clientstats_active"),
			"Total node protocol operation in rate.",
			[]string{"node", "proto"}, ConstLabels,
		),
	}, nil
}

func (c *nodeProtoCollector) Update(ch chan<- prometheus.Metric) error {
	//Set client stats for nfs and smb to ungathered
	nfsClientstatGathered = false
	smbClientstatsGathered = false

	var err error
	//Attemp to update over the list of all protocol
	for proto, state := range protocolState {
		// Only execute if state is true
		if *state {
			err = c.updateProtoOpStats(ch, proto)
			if err != nil {
				log.Warnf("Unabled to collect protocol operation stats for %s", proto)
			}
			err = c.updateProtoStats(ch, proto)
			if err != nil {
				log.Warnf("Unable to collect protocol stats for %s", proto)
			}
			err = c.updateProtoClientstatsActive(ch, proto)
			if err != nil {
				log.Warnf("Unable to collect protocol stats for %s", proto)
			}
			err = c.updateProtoClientstatsConnected(ch, proto)
			if err != nil {
				log.Warnf("Unable to collect protocol stats for %s", proto)
			}
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func (c *nodeProtoCollector) updateProtoOpStats(ch chan<- prometheus.Metric, protocol string) error {
	key := fmt.Sprintf("node.protostats.%v", protocol)
	resp, err := isiclient.GetProtoStat(IsiCluster.Client, key)
	if err != nil {
		log.Warnf("Unable to collect node protocol stats for protocol %s.", protocol)
		return err
	}
	//Get stats for each node
	for _, stat := range resp.Stats {
		if stat.Value != nil {
			values := stat.Value.([]interface{})
			if len(values) > 0 {
				for _, value := range values {
					//Now that we know there is data in this interface, marshal into json and back out into a struct
					var protoStat isiclient.IsiProtoStatOp
					j, err := json.Marshal(value)
					if err != nil {
						log.Warnf("Could not marshal back into json: %v", err)
						return err
					}
					err = json.Unmarshal(j, &protoStat)
					if err != nil {
						log.Warnf("Could not unmarshl into stuct: %v", err)
						return err
					}
					node := fmt.Sprintf("%v", stat.Devid)
					//Add metrics for each item
					ch <- prometheus.MustNewConstMetric(c.nodeProtocolInMax, prometheus.GaugeValue, protoStat.InMax, node, protocol, protoStat.OpName)
					ch <- prometheus.MustNewConstMetric(c.nodeProtocolInMin, prometheus.GaugeValue, protoStat.InMin, node, protocol, protoStat.OpName)
					ch <- prometheus.MustNewConstMetric(c.nodeProtocolInRate, prometheus.GaugeValue, protoStat.InRate, node, protocol, protoStat.OpName)
					ch <- prometheus.MustNewConstMetric(c.nodeProtocolOpCount, prometheus.GaugeValue, protoStat.OpCount, node, protocol, protoStat.OpName)
					ch <- prometheus.MustNewConstMetric(c.nodeProtocolOpRate, prometheus.GaugeValue, protoStat.OpRate, node, protocol, protoStat.OpName)
					ch <- prometheus.MustNewConstMetric(c.nodeProtocolOutMax, prometheus.GaugeValue, protoStat.OutMax, node, protocol, protoStat.OpName)
					ch <- prometheus.MustNewConstMetric(c.nodeProtocolOutMin, prometheus.GaugeValue, protoStat.OutMin, node, protocol, protoStat.OpName)
					ch <- prometheus.MustNewConstMetric(c.nodeProtocolOutRate, prometheus.GaugeValue, protoStat.OutRate, node, protocol, protoStat.OpName)
					ch <- prometheus.MustNewConstMetric(c.nodeProtocolTimeAvg, prometheus.GaugeValue, protoStat.TimeAvg, node, protocol, protoStat.OpName)
					ch <- prometheus.MustNewConstMetric(c.nodeProtocolTimeMax, prometheus.GaugeValue, protoStat.TimeMax, node, protocol, protoStat.OpName)
					ch <- prometheus.MustNewConstMetric(c.nodeProtocolTimeMin, prometheus.GaugeValue, protoStat.TimeMin, node, protocol, protoStat.OpName)
				}
			}
		}
	}
	return nil
}

func (c *nodeProtoCollector) updateProtoStats(ch chan<- prometheus.Metric, protocol string) error {
	key := fmt.Sprintf("node.protostats.%s.total", protocol)
	resp, err := isiclient.GetProtoStat(IsiCluster.Client, key)
	if err != nil {
		log.Warnf("Unable to collect node protocol stats for protocol %s.", protocol)
		return err
	}
	for _, stat := range resp.Stats {
		if stat.Value != nil {
			values := stat.Value.([]interface{})

			if len(values) > 0 {
				for _, value := range values {
					//Now that we know there is data in this interface, marshal into json and back out into a struct
					var protoStat isiclient.IsiProtoStatTotal
					j, err := json.Marshal(value)
					if err != nil {
						log.Warnf("Could not marshal back into json: %v", err)
						return err
					}
					err = json.Unmarshal(j, &protoStat)
					if err != nil {
						log.Warnf("Could not unmarshl into stuct: %v", err)
						return err
					}
					node := fmt.Sprintf("%v", stat.Devid)
					ch <- prometheus.MustNewConstMetric(c.nodeProtocolTotalInMax, prometheus.GaugeValue, protoStat.InMax, node, protocol)
					ch <- prometheus.MustNewConstMetric(c.nodeProtocolTotalInMin, prometheus.GaugeValue, protoStat.InMin, node, protocol)
					ch <- prometheus.MustNewConstMetric(c.nodeProtocolTotalInRate, prometheus.GaugeValue, protoStat.InRate, node, protocol)
					ch <- prometheus.MustNewConstMetric(c.nodeProtocolTotalOpCount, prometheus.GaugeValue, protoStat.OpCount, node, protocol)
					ch <- prometheus.MustNewConstMetric(c.nodeProtocolTotalOpRate, prometheus.GaugeValue, protoStat.OpRate, node, protocol)
					ch <- prometheus.MustNewConstMetric(c.nodeProtocolTotalOutMax, prometheus.GaugeValue, protoStat.OutMax, node, protocol)
					ch <- prometheus.MustNewConstMetric(c.nodeProtocolTotalOutMin, prometheus.GaugeValue, protoStat.OutMin, node, protocol)
					ch <- prometheus.MustNewConstMetric(c.nodeProtocolTotalOutRate, prometheus.GaugeValue, protoStat.OutRate, node, protocol)
					ch <- prometheus.MustNewConstMetric(c.nodeProtocolTotalTimeAvg, prometheus.GaugeValue, protoStat.TimeAvg, node, protocol)
					ch <- prometheus.MustNewConstMetric(c.nodeProtocolTotalTimeMax, prometheus.GaugeValue, protoStat.TimeMax, node, protocol)
					ch <- prometheus.MustNewConstMetric(c.nodeProtocolTotalTimeMin, prometheus.GaugeValue, protoStat.TimeMin, node, protocol)
				}
			}
		}
	}
	return nil
}

func (c *nodeProtoCollector) updateProtoClientstatsActive(ch chan<- prometheus.Metric, protocol string) error {
	// There are not client stats for lsass_in or nfs4
	if protocol == "nfs4" || protocol == "lsass_in" {
		return nil
	}

	activeKey := fmt.Sprintf("node.clientstats.active.%s", protocol)

	// Both stats are single stat values in the normal format
	resp, err := isiclient.QueryStatsEngineSingleVal(IsiCluster.Client, activeKey)
	if err != nil {
		log.Warnf("Unable to collect node protocol client stats for protocol %s.", protocol)
		return err
	}

	for _, stat := range resp.Stats {
		node := fmt.Sprintf("%v", stat.Devid)
		ch <- prometheus.MustNewConstMetric(c.nodeClientsActive, prometheus.GaugeValue, stat.Value, node, protocol)
	}
	return nil
}

func (c *nodeProtoCollector) updateProtoClientstatsConnected(ch chan<- prometheus.Metric, protocol string) error {
	// There are not client stats for lsass
	if protocol == "jobd" || strings.Contains(protocol, "lsass") {
		return nil
	}

	// All versions of NFS and SMB fall under the same key name.
	if strings.Contains(protocol, "nfs") {
		if !nfsClientstatGathered {
			protocol = "nfs"
			nfsClientstatGathered = true
		} else {
			return nil
		}
	}

	if strings.Contains(protocol, "smb") {
		if !smbClientstatsGathered {
			protocol = "smb"
			smbClientstatsGathered = true
		} else {
			return nil
		}
	}

	connectedKey := fmt.Sprintf("node.clientstats.connected.%s", protocol)

	resp, err := isiclient.QueryStatsEngineSingleVal(IsiCluster.Client, connectedKey)
	if err != nil {
		log.Warnf("Unable to collect node protocol client stats for protocol %s.", protocol)
		return err
	}
	for _, stat := range resp.Stats {
		node := fmt.Sprintf("%v", stat.Devid)
		ch <- prometheus.MustNewConstMetric(c.nodeClientsConnected, prometheus.GaugeValue, stat.Value, node, protocol)
	}
	return nil
}
