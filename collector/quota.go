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
	"errors"
	"fmt"
	"time"

	"github.com/adobe/prometheus-emcisilon-exporter/isiclient"
	"github.com/thecodeteam/goisilon"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"gopkg.in/alecthomas/kingpin.v2"
)

type quotaCollector struct {
	quotaIterationCollectionTime   *prometheus.Desc
	quotaContainer                 *prometheus.Desc
	quotaEnforced                  *prometheus.Desc
	quotaIncludeSnapshots          *prometheus.Desc
	quotaUsageLogical              *prometheus.Desc
	quotaUsageInodes               *prometheus.Desc
	quotaUsagePhysical             *prometheus.Desc
	quotaThresholdAdvisory         *prometheus.Desc
	quotaThresholdAdvisoryExceeded *prometheus.Desc
	quotaThresholdHard             *prometheus.Desc
	quotaThresholdHardExceeded     *prometheus.Desc
	quotaThresholdSoft             *prometheus.Desc
	quotaThresholdSoftGrace        *prometheus.Desc
	quotaThresholdSoftExceeded     *prometheus.Desc
}

var (
	typeFlag      *string
	exceededFlag  *string
	numQuotas     string
	rtoken        = "unset"
	collectNumber int
)

func init() {
	registerCollector("quota", defaultDisabled, NewQuotaCollector)

	//Quota type flag.
	typeFlagName := "collector.quota.type"
	typeFlagHelp := "Quota type to collect stats for (default: all). One of type (directory, user, group, default-user, default-group, all)"
	typeFlag = kingpin.Flag(typeFlagName, typeFlagHelp).Default("all").String()

	//Quota exceeded flag.
	exceededFlagName := "collector.quota.exceeded"
	exceededFlagHelp := "Only turn quotas that have exceeded one of more thresholds (default: false). Boolean of type (false, true)."
	exceededFlag = kingpin.Flag(exceededFlagName, exceededFlagHelp).Default("false").String()

}

//NewQuotaCollector returns a new Collector exposing node health information.
func NewQuotaCollector() (Collector, error) {
	return &quotaCollector{
		quotaIterationCollectionTime: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, quotaCollectorSubsystem, "api_collection_duration"),
			"Returns the amount of time it took to collect an iteration of quotas from the api.",
			[]string{"iteration"}, ConstLabels,
		),
		quotaContainer: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, quotaCollectorSubsystem, "container"),
			"1 if quota is a container quota, 0 if not.",
			[]string{"id", "path", "name", "type"}, ConstLabels,
		),
		quotaEnforced: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, quotaCollectorSubsystem, "enforced"),
			"1 if quota is enforced, 2 if quota is an advisory quota.",
			[]string{"id", "path", "name", "type"}, ConstLabels,
		),
		quotaIncludeSnapshots: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, quotaCollectorSubsystem, "include_snapshots"),
			"1 if quota includes snapshots in usage, 0 if not.",
			[]string{"id", "path", "name", "type"}, ConstLabels,
		),
		quotaUsageLogical: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, quotaCollectorSubsystem, "usage_logical"),
			"Apparent bytes used by governed data.",
			[]string{"id", "path", "name", "type"}, ConstLabels,
		),
		quotaUsageInodes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, quotaCollectorSubsystem, "usage_inodes"),
			"Number of inodes (filesystem entities) used by governed data.",
			[]string{"id", "path", "name", "type"}, ConstLabels,
		),
		quotaUsagePhysical: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, quotaCollectorSubsystem, "usage_physical"),
			"Bytes used for governed data and filesystem overhead.",
			[]string{"id", "path", "name", "type"}, ConstLabels,
		),
		quotaThresholdAdvisory: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, quotaCollectorSubsystem, "threshold_advisory"),
			"Usage bytes at which notifications will be sent but writes will not be denied.",
			[]string{"id", "path", "name", "type"}, ConstLabels,
		),
		quotaThresholdAdvisoryExceeded: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, quotaCollectorSubsystem, "threshold_advisory_exceeded"),
			"1 if the advisory threshold has been hit.",
			[]string{"id", "path", "name", "type"}, ConstLabels,
		),
		quotaThresholdSoft: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, quotaCollectorSubsystem, "threshold_soft"),
			"Usage bytes at which notifications will be sent and soft grace time will be started.",
			[]string{"id", "path", "name", "type"}, ConstLabels,
		),
		quotaThresholdSoftExceeded: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, quotaCollectorSubsystem, "threshold_soft_exceeded"),
			"1 if the soft threshold has been hit.",
			[]string{"id", "path", "name", "type"}, ConstLabels,
		),
		quotaThresholdSoftGrace: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, quotaCollectorSubsystem, "threshold_soft_grace"),
			"Time in seconds after which the soft threshold has been hit before writes will be denied.",
			[]string{"id", "path", "name", "type"}, ConstLabels,
		),
		quotaThresholdHard: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, quotaCollectorSubsystem, "threshold_hard"),
			"Usage bytes at which further writes will be denied.",
			[]string{"id", "path", "name", "type"}, ConstLabels,
		),
		quotaThresholdHardExceeded: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, quotaCollectorSubsystem, "threshold_hard_exceeded"),
			"True if the hard threshold has been hit.",
			[]string{"id", "path", "name", "type"}, ConstLabels,
		),
	}, nil
}

func (c *quotaCollector) Update(ch chan<- prometheus.Metric) error {
	log.Debugf("Collecting quota type(s): %s", *typeFlag)
	log.Debugf("Collected only exceeded quotas: %v", *exceededFlag)

	// Keep going until there is no resume token
	var err error
	for rtoken != "" {
		// Collect a time and counter for each iteration of quotas
		collectNumber++
		begin := time.Now()

		//Ask for first set of quotas
		quotas, qerr := c.getQuotas(IsiCluster.Client)
		if qerr != nil {
			log.Warnf("Unable to collect quotas for type: %s", *typeFlag)
			err = qerr
			return err
		}

		//Calculate the time it took to gather this iteration of quotas
		duration := time.Since(begin)

		//Grab the resume token if there is one. If there is not one it will be a empty string.
		rtoken = quotas.Resume

		//Create a new metric of the amount of time it took to collect this iteration of quotas.
		ch <- prometheus.MustNewConstMetric(c.quotaIterationCollectionTime, prometheus.GaugeValue, duration.Seconds(), fmt.Sprintf("%v", collectNumber))

		//Range over all quotas in this iteration.

		for _, quota := range quotas.Quotas {
			// Get username for the quota
			var name string
			if quota.Type != "directory" {
				name, err = c.getQuotaUserName(quota)
				if err != nil {
					log.Infof("Unabled to get a name for quota: %s", err)
				}
			} else {
				name = quota.Path
			}

			//Gather meta-data metrics
			err = c.updateMetaData(ch, quota, name)
			if err != nil {
				log.Warnf("Unable to update meta data for quota: %s", name)
			}

			//Gather usage metrics
			err = c.updateUsage(ch, quota, name)
			if err != nil {
				log.Warnf("Unable to update usage for quota: %s", name)
			}

			//Gather threshold metrics
			err = c.updateThresholds(ch, quota, name)
			if err != nil {
				log.Warnf("Unable to update usage for quota: %s", name)
			}
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func (c *quotaCollector) getQuotas(con *goisilon.Client) (isiclient.IsiQuotas, error) {
	//Get count of quotas
	var quotas isiclient.IsiQuotas
	sum, err := isiclient.GetQuotaSummary(con)
	if err != nil {
		log.Warn("Unable to collect quota summary information.")
		return quotas, err
	}

	//Check to see what type of quota is being collected.
	if rtoken != "" && rtoken != "unset" {
		quotas, err = isiclient.GetQuotasWithResume(con, rtoken)
	} else {
		switch *typeFlag {
		case "directory":
			numQuotas = fmt.Sprintf("%v", sum.DirectoryQuotasCount)
			quotas, err = isiclient.GetQuotasOfType(con, *exceededFlag, numQuotas, *typeFlag)
		case "user":
			numQuotas = fmt.Sprintf("%v", sum.UserQuotasCount)
			quotas, err = isiclient.GetQuotasOfType(con, *exceededFlag, numQuotas, *typeFlag)
		case "group":
			numQuotas = fmt.Sprintf("%v", sum.GroupQuotasCount)
			quotas, err = isiclient.GetQuotasOfType(con, *exceededFlag, numQuotas, *typeFlag)
		case "default-user":
			numQuotas = fmt.Sprintf("%v", sum.DefaultUserQuotasCount)
			quotas, err = isiclient.GetQuotasOfType(con, *exceededFlag, numQuotas, *typeFlag)
		case "default-group":
			numQuotas = fmt.Sprintf("%v", sum.DefaultGroupQuotasCount)
			quotas, err = isiclient.GetQuotasOfType(con, *exceededFlag, numQuotas, *typeFlag)
		case "all":
			numQuotas = fmt.Sprintf("%v", sum.Count)
			quotas, err = isiclient.GetAllQuotas(con, *exceededFlag, numQuotas)
		default:
			mesg := fmt.Sprintf("Unknown quota type: %s", *typeFlag)
			err = errors.New(mesg)
			return quotas, err
		}
	}

	if err != nil {
		mesg := fmt.Sprintf("Unable to collect quota information: %s", err)
		err := errors.New(mesg)
		return quotas, err
	}
	return quotas, nil
}

func (c *quotaCollector) getQuotaUserName(q isiclient.IsiQuota) (string, error) {
	var (
		name string
	)

	//Try and unmarshal the persona interface.
	if q.Persona != nil {
		per, ok := q.Persona.(map[string]interface{})
		if !ok {
			err := errors.New("unable to unmarshal the persona interface{}")
			return "", err
		}
		val, ok := per["name"]
		if !ok {
			name = q.Path
		} else {
			name = val.(string)
		}
	} else {
		log.Debugf("Persona is nil: %v", q.Path)
		name = q.Path
	}
	return name, nil
}

func (c *quotaCollector) updateMetaData(ch chan<- prometheus.Metric, q isiclient.IsiQuota, n string) error {
	var (
		container       float64
		enforced        float64
		includeSnapshot float64
	)
	//Gather if quota is a container quota
	if q.Container {
		container = 1
	} else {
		container = 0
	}
	ch <- prometheus.MustNewConstMetric(c.quotaContainer, prometheus.GaugeValue, container, q.ID, q.Path, n, q.Type)

	//Gather enforcement status
	if q.Enforced {
		enforced = 1
	} else {
		enforced = 0
	}
	ch <- prometheus.MustNewConstMetric(c.quotaEnforced, prometheus.GaugeValue, enforced, q.ID, q.Path, n, q.Type)

	//Gather include_snapshots
	if q.IncludeSnapshots {
		includeSnapshot = 1
	} else {
		includeSnapshot = 0
	}
	ch <- prometheus.MustNewConstMetric(c.quotaIncludeSnapshots, prometheus.GaugeValue, includeSnapshot, q.ID, q.Path, n, q.Type)

	return nil
}

func (c *quotaCollector) updateUsage(ch chan<- prometheus.Metric, q isiclient.IsiQuota, n string) error {
	// Update logical
	ch <- prometheus.MustNewConstMetric(c.quotaUsageLogical, prometheus.GaugeValue, q.Usage.Logical, q.ID, q.Path, n, q.Type)
	ch <- prometheus.MustNewConstMetric(c.quotaUsageInodes, prometheus.GaugeValue, q.Usage.Inodes, q.ID, q.Path, n, q.Type)
	ch <- prometheus.MustNewConstMetric(c.quotaUsagePhysical, prometheus.GaugeValue, q.Usage.Physical, q.ID, q.Path, n, q.Type)
	return nil
}

func (c *quotaCollector) updateThresholds(ch chan<- prometheus.Metric, q isiclient.IsiQuota, n string) error {
	//gather advisory thresholds
	var (
		ae float64
		he float64
		se float64
	)
	ch <- prometheus.MustNewConstMetric(c.quotaThresholdAdvisory, prometheus.GaugeValue, q.Thresholds.Advisory, q.ID, q.Path, n, q.Type)
	if q.Thresholds.AdvisoryExceeded {
		ae = 1
	} else {
		ae = 0
	}
	ch <- prometheus.MustNewConstMetric(c.quotaThresholdAdvisoryExceeded, prometheus.GaugeValue, ae, q.ID, q.Path, n, q.Type)

	//gather hard thresholds
	ch <- prometheus.MustNewConstMetric(c.quotaThresholdHard, prometheus.GaugeValue, q.Thresholds.Hard, q.ID, q.Path, n, q.Type)
	if q.Thresholds.HardExceeded {
		he = 1
	} else {
		he = 0
	}
	ch <- prometheus.MustNewConstMetric(c.quotaThresholdHardExceeded, prometheus.GaugeValue, he, q.ID, q.Path, n, q.Type)

	//gather soft thresholds
	ch <- prometheus.MustNewConstMetric(c.quotaThresholdSoft, prometheus.GaugeValue, q.Thresholds.Soft, q.ID, q.Path, n, q.Type)
	if q.Thresholds.SoftExceeded {
		se = 1
	} else {
		se = 0
	}
	ch <- prometheus.MustNewConstMetric(c.quotaThresholdSoftExceeded, prometheus.GaugeValue, se, q.ID, q.Path, n, q.Type)
	ch <- prometheus.MustNewConstMetric(c.quotaThresholdSoftGrace, prometheus.GaugeValue, q.Thresholds.SoftGrace, q.ID, q.Path, n, q.Type)
	return nil
}
