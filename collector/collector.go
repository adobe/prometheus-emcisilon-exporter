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
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Defined top level common namespace that all metrics use.
const (
	defaultEnabled  = true
	defaultDisabled = false
	namespace       = "isilon"
)

var (
	factories      = make(map[string]func() (Collector, error))
	collectorState = make(map[string]*bool)
)

func registerCollector(collector string, isDefaultEnabled bool, factory func() (Collector, error)) {
	var helpDefaultState string
	if isDefaultEnabled {
		helpDefaultState = "enabled"
	} else {
		helpDefaultState = "disabled"
	}

	flagName := fmt.Sprintf("collector.%s", collector)
	flagHelp := fmt.Sprintf("Enable the %s collector (default: %s).", collector, helpDefaultState)
	defaultValue := fmt.Sprintf("%v", isDefaultEnabled)

	flag := kingpin.Flag(flagName, flagHelp).Default(defaultValue).Bool()
	collectorState[collector] = flag

	factories[collector] = factory
}

//IsilonCollector implements the prometheus.Collector interface.
type isilonCollector struct {
	Collectors map[string]Collector
}

// NewIsilonCollector creates a new IsilonCollector
func NewIsilonCollector(fqdn string, port string, uname string, pwdenv string, site string, auth bool, qOnly bool, filters ...string) (*isilonCollector, error) {
	if auth {
		// Take the struct that was generated in main and use it as the configuration for connecting to the clusters.
		IsiCluster.FQDN = fqdn
		IsiCluster.Port = port
		IsiCluster.Username = uname
		IsiCluster.PasswordEnv = pwdenv
		IsiCluster.Site = site
		IsiCluster.QuotaOnly = qOnly

		// Get the the goisilon connector and put it into the shared IsiClusterConfig struct.
		log.Debugf("Creating connection to the cluster endpoint %s", IsiCluster.FQDN)
		err := GetClusterConnector()
		if err != nil {
			return nil, fmt.Errorf("Unable to connect to the isilon cluster %s: %s", IsiCluster.FQDN, err)
		}

		log.Debug("Getting isi config cluster name from identity endpoint.")
		//Get the clusster name from the isilon client.
		err = SetClusterConfigName()
		if err != nil {
			return nil, fmt.Errorf("Unable to get the cluster config name from the identity endpoint: %s", err)
		}

		if IsiCluster.QuotaOnly {
			log.Debug("Setting up collector to only collect quota info.")
			err := GetNumQuotas()
			if err != nil {
				return nil, fmt.Errorf("Unable to get count of quotas from the system. %s", err)
			}

			flag := kingpin.Flag("collector.quota.retry", "Number of time to attempt collection of quota metrics (default: 3).").Default("3").Int64()
			IsiCluster.Quotas.Retry = *flag
		}

		// Create descriptors for collector leve metrics.
		scrapeDurationDesc = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "scrape", "collector_duration_seconds"),
			"isilon_exporter: Duration of a collector scrape,",
			[]string{"collector"}, ConstLabels,
		)
		scrapeSuccessDesc = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "scrape", "collector_success"),
			"isilon_exporter: Whether a collector succeeded.",
			[]string{"collector"}, ConstLabels,
		)
		exporterDurationDesc = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "exporter", "duration_seconds"),
			"Duration in second of the entire exporter run.",
			nil, ConstLabels,
		)
		statsEngineCallFailure = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "stats_engine", "call_success"),
			"0 = Successful, 1 = Failure.  Represent the successful call or failure to the stats engine.",
			[]string{"stat_key"}, ConstLabels,
		)
		statsEngineCallDuration = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "stats_engine", "call_duration_seconds"),
			"Duration in seconds a call to the stats engine takes.",
			[]string{"stat_key"}, ConstLabels,
		)
	}

	//If qOnly then set all collectors to disabled except for quotas
	if qOnly {
		var disabled = false
		var enabled = true
		for key := range collectorState {
			if key == "quota" {
				collectorState[key] = &enabled
			} else {
				collectorState[key] = &disabled
			}
		}
	}

	f := make(map[string]bool)
	for _, filter := range filters {
		enabled, exist := collectorState[filter]
		if !exist {
			return nil, fmt.Errorf("missing collector: %s", filter)
		}
		if !*enabled {
			return nil, fmt.Errorf("disabled collector: %s", filter)
		}
		f[filter] = true
	}
	collectors := make(map[string]Collector)
	for key, enabled := range collectorState {
		if *enabled {
			collector, err := factories[key]()
			if err != nil {
				return nil, err
			}
			if len(f) == 0 || f[key] {
				collectors[key] = collector
			}
		}
	}
	return &isilonCollector{Collectors: collectors}, nil
}

// Descibe implements the prometheus.Collector interface.
func (n isilonCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- scrapeDurationDesc
	ch <- scrapeSuccessDesc
	ch <- exporterDurationDesc
	ch <- statsEngineCallDuration
	ch <- statsEngineCallFailure
}

// Collect implements the prometheus.Collector interface.
func (n isilonCollector) Collect(ch chan<- prometheus.Metric) {
	begin := time.Now()
	wg := sync.WaitGroup{}
	wg.Add(len(n.Collectors))
	for name, c := range n.Collectors {
		go func(name string, c Collector) {
			execute(name, c, ch)
			wg.Done()
		}(name, c)
	}
	wg.Wait()
	duration := time.Since(begin)
	log.Debugf("Exporter finished after %fs", duration.Seconds())
	ch <- prometheus.MustNewConstMetric(exporterDurationDesc, prometheus.GaugeValue, duration.Seconds())
}

func execute(name string, c Collector, ch chan<- prometheus.Metric) {
	begin := time.Now()
	err := c.Update(ch)
	duration := time.Since(begin)
	var success float64

	if err != nil {
		log.Errorf("ERROR: %s collector failed after %fs: %s", name, duration.Seconds(), err)
		success = 0
	} else {
		log.Debugf("OK: %s collector succeeded after %fs.", name, duration.Seconds())
		success = 1
	}
	ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, duration.Seconds(), name)
	ch <- prometheus.MustNewConstMetric(scrapeSuccessDesc, prometheus.GaugeValue, success, name)
}

// Collector is the interface a collector has to implement.
type Collector interface {
	// Get new metrics and expose them via prometheus registry.
	Update(ch chan<- prometheus.Metric) error
}
