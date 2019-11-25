/*
Copyright 2018 Adobe
All Rights Reserved.

NOTICE: Adobe permits you to use, modify, and distribute this file in
accordance with the terms of the Adobe license agreement accompanying
it. If you have received this file from a source other than Adobe,
then your use, modification, or distribution of it requires the prior
written permission of Adobe.
*/
package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"sort"
	"time"

	"github.com/adobe/prometheus-emcisilon-exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	fqdn   *string
	port   *string
	uname  *string
	pwdenv *string
	site   *string
	qOnly  *bool
)

// Registers the isilon_exporter as a prometheus collector
func init() {
	version.Version = "2.0.0"
	version.BuildDate = fmt.Sprintf("%v", time.Now())
	version.BuildUser = "panike"
	prometheus.MustRegister(version.NewCollector("prometheus_emcisilon_exporter"))

}

// Handler takes care of the local http traffic.
func handler(w http.ResponseWriter, r *http.Request) {
	filters := r.URL.Query()["collect[]"]
	log.Debugln("collect query:", filters)

	//Creates a new isilon collector with filters applied. (Kingpin flags)
	nc, err := collector.NewIsilonCollector(*fqdn, *port, *uname, *pwdenv, *site, true, *qOnly, filters...)
	if err != nil {
		log.Warnf("Could not create exporter: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Could not create exporter: %s", err)))
		return
	}

	registry := prometheus.NewRegistry()
	err = registry.Register(nc)
	if err != nil {
		log.Errorf("Could not register collector: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Could not register collector: %s", err)))
		return
	}

	gatherers := prometheus.Gatherers{
		prometheus.DefaultGatherer,
		registry,
	}
	// Delegate http serving to Prometheus client Librar, which will call collector.Collect.
	h := promhttp.InstrumentMetricHandler(
		registry,
		promhttp.HandlerFor(gatherers,
			promhttp.HandlerOpts{
				ErrorLog:      log.NewErrorLogger(),
				ErrorHandling: promhttp.ContinueOnError,
			}),
	)
	h.ServeHTTP(w, r)
}

func main() {
	var (
		//HTTP Variables
		listenAddress = kingpin.Flag("web.listen-address", "Address on which to expose metrics and web interface.").Default(":9300").String()
		metricsPath   = kingpin.Flag("web.telemtry-path", "Path under which to expose metrics.").Default("/metrics").String()

		//Isilon Specific Variables
		cFQDN     = kingpin.Flag("isilon.cluster.fqdn", "FQDN for the isilon cluster to be scraped.").Default("localhost").String()
		cPort     = kingpin.Flag("isilon.cluster.port", "Port to connect to the isilon cluster.").Default("8080").String()
		cUname    = kingpin.Flag("isilon.cluster.username", "Username for access the isilon API.").Default("").String()
		cPwdenv   = kingpin.Flag("isilon.cluster.password.env", "Environment variable that contains the password for the Isilon cluster user.").Default("ISILON_CLUSTER_PASSWORD").String()
		cSite     = kingpin.Flag("isilon.cluster.site", "Data Center site the cluster is located in.").Default("").String()
		quotaOnly = kingpin.Flag("quota-only", "Set exporter to only collect quota information.").Default("false").Bool()
	)

	log.AddFlags(kingpin.CommandLine)
	kingpin.Version(version.Print("prometheus-emcisilon-exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	//Create an IsilonCluster struct and pass it the infor from the kingpin flags.
	if *cUname == "" {
		log.Fatalf("No cluster username specified.")
	}

	fqdn = cFQDN
	port = cPort
	uname = cUname
	pwdenv = cPwdenv
	site = cSite
	qOnly = quotaOnly
	log.Infoln("Started prometheus-emcisilon-exporter", version.Info())

	log.Infof("Pointed to cluster %s", *fqdn)
	log.Infoln("Build context", version.BuildContext())

	// This instance is only used to check collector creation and logging.
	nc, err := collector.NewIsilonCollector(*fqdn, *port, *uname, *pwdenv, *site, false, *qOnly)
	if err != nil {
		log.Fatalf("Could not create collector: %s", err)
	}
	log.Infof("Enable collectors:")
	collectors := []string{}
	for n := range nc.Collectors {
		collectors = append(collectors, n)
	}
	sort.Strings(collectors)
	for _, n := range collectors {
		log.Infof(" - %s", n)
	}

	http.HandleFunc(*metricsPath, handler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Isilon Exporter</title></head>
			<body>
			<h1>Isilon Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	log.Infoln("Listening on", *listenAddress)
	err = http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
}
