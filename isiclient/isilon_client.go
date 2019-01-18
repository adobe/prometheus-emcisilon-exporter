/*
Copyright 2018 Adobe
All Rights Reserved.

NOTICE: Adobe permits you to use, modify, and distribute this file in
accordance with the terms of the Adobe license agreement accompanying
it. If you have received this file from a source other than Adobe,
then your use, modification, or distribution of it requires the prior
written permission of Adobe.
*/
package isiclient

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/prometheus/common/log"
	"github.com/hpanike/goisilon"
	"github.com/hpanike/goisilon/api"
)

// NewIsilonClient creates and isilon client from goisilon.NewClientsWithArgs.
func NewIsilonClient(fqdn string, port string, username string, passwordEnv string) (*goisilon.Client, error) {
	// Setup the client from the cluster info and return the client.
	//Build endpoint from fqdn and port. Force HTTPS as we are using basic auth.
	endpoint := fmt.Sprintf("https://%s:%s", fqdn, port)

	//Get the password from the provided environment variable.
	password, ok := os.LookupEnv(passwordEnv)
	if !ok {
		mesg := fmt.Sprintf("Unabled to retrieve password from env variable: %s", passwordEnv)
		log.Warn(mesg)
		err := errors.New(mesg)
		return nil, err
	}

	// Instantiate new isilon connector
	c, err := goisilon.NewClientWithArgs(
		context.Background(),
		endpoint,
		true,
		username,
		"",
		password,
		"",
	)
	if err != nil {
		log.Warnf("Could not create connection to Isilon Cluster %s: %s", endpoint, err)
		return nil, err
	}
	return c, nil
}

//GetClusterName is used to get the cluster name from the api call to isi config
func GetClusterName(c *goisilon.Client) (string, error) {
	var (
		path string
		resp IsiIdentity
	)

	//Static set path to platform api
	path = "/platform/3/cluster/identity"

	err := c.API.Get(context.Background(), path, "", nil, nil, &resp)
	if err != nil {
		log.Warnf("Unable to get cluster identity from api: %s", err)
	}
	if resp.Name != "" {
		return resp.Name, nil
	}
	mesg := fmt.Sprintf("Could not retrieve name from identity response: %v", resp.Name)
	err = errors.New(mesg)
	return "", err
}

//QueryStatsEngineMultiVal is used to unmarshal a stat with muliple values.
func QueryStatsEngineMultiVal(c *goisilon.Client, key string) (IsiMultiVal, error) {
	var (
		path   string
		resp   IsiMultiVal
		params api.OrderedValues
	)

	//Static set path to stats engine
	path = "/platform/1/statistics/current"

	//Set params to be key passed, and query all device ids at the same time
	params = api.NewOrderedValues([][]string{
		{"keys", key},
		{"devid", "all"},
	})

	//Query API
	err := c.API.Get(context.Background(), path, "", params, nil, &resp)
	if err != nil {
		log.Warnf("Unable to retrieve stats for %s: %s", key, err)
		return resp, err
	}
	return resp, nil
}

//QueryStatsEngineSingleVal is used to unmarshal a stat with a single value.
func QueryStatsEngineSingleVal(c *goisilon.Client, key string) (IsiSingleVal, error) {
	var (
		path      string
		statsResp IsiSingleVal
		params    api.OrderedValues
	)

	//Static set path to stats engine
	path = "/platform/1/statistics/current"

	//Set params to be key passed, and query all device ids at the same time
	params = api.NewOrderedValues([][]string{
		{"keys", key},
		{"devid", "all"},
	})

	//Query API
	err := c.API.Get(context.Background(), path, "", params, nil, &statsResp)
	if err != nil {
		log.Warnf("Unable to retrieve stats for %s: %s", key, err)
		return statsResp, err
	}
	return statsResp, nil
}

//GetOneFsVersion will grab the config from api and unmarshal the struct
func GetOneFsVersion(c *goisilon.Client) (string, error) {
	var (
		path = "/platform/3/cluster/config"
		resp IsiConfig
	)

	err := c.API.Get(context.Background(), path, "", nil, nil, &resp)
	if err != nil {
		log.Warnf("Unable to get cluster config from api: %s", err)
	}
	return resp.OnefsVersion.Release, nil
}

//GetQuotas returns all quotas from the api
func GetQuotas(c *goisilon.Client, exceeded bool) (IsiQuotas, error) {
	var (
		path   = "/platform/1/quota/quotas"
		resp   IsiQuotas
		params api.OrderedValues
		token  = ""
	)
	params = api.NewOrderedValues([][]string{
		{"resolve_names", "true"},
		{"resume", token},
	})

	if exceeded {
		params.StringSet("exceeded", "true")
	}

	fmt.Printf("%v", params)

	err := c.API.Get(context.Background(), path, "", params, nil, &resp)
	if err != nil {
		log.Warnf("Unable to get quotas from api: %s", err)
	}
	return resp, nil
}

func GetAllQuotas(c *goisilon.Client, exceeded bool) (IsiQuotas, error) {
	var (
		path   = "/platform/1/quota/quotas"
		resp   IsiQuotas
		params api.OrderedValues
	)
	params = api.NewOrderedValues([][]string{
		{"resolve_names", "true"},
	})
	if exceeded {
		params.StringSet("exceeded", "true")
	}

	err := c.API.Get(context.Background(), path, "", params, nil, &resp)
	if err != nil {
		log.Warnf("Unable to get quotas from api: %s", err)
	}
	return resp, nil
}

func GetQuotasOfType(c *goisilon.Client, exceeded bool, qtype string) (IsiQuotas, error) {
	var (
		path   = "/platform/1/quota/quotas"
		resp   IsiQuotas
		params api.OrderedValues
	)
	params = api.NewOrderedValues([][]string{
		{"resolve_names", "true"},
		{"type", qtype},
	})

	if exceeded {
		params.StringSet("exceeded", "true")
	}

	err := c.API.Get(context.Background(), path, "", params, nil, &resp)
	if err != nil {
		log.Warnf("Unable to get quotas from api: %s", err)
	}
	return resp, nil
}

//GetQuotaSummary will return a IsiQuotaSummary struct with information from /platform/1/quota/quotas-summary
func GetQuotaSummary(c *goisilon.Client) (IsiQuotaSummary, error) {
	var (
		path = "/platform/1/quota/quotas-summary"
		resp IsiQuotaSummaryResp
	)
	err := c.API.Get(context.Background(), path, "", nil, nil, &resp)
	if err != nil {
		log.Warnf("Unable to get summary about quotas from api: %s", err)
	}
	return resp.Summary, nil
}

func GetQuotasWithResume(c *goisilon.Client, rtoken string) (IsiQuotas, error) {
	var (
		path   = "/platform/1/quota/quotas"
		resp   IsiQuotas
		params api.OrderedValues
	)
	params = api.NewOrderedValues([][]string{
		{"resume", rtoken},
	})

	err := c.API.Get(context.Background(), path, "", params, nil, &resp)
	if err != nil {
		log.Warnf("Unable to get quotas from api with resume token %s : %s ", rtoken, err)
	}
	return resp, nil
}

//GetProtoStat for protocol level information
func GetProtoStat(c *goisilon.Client, key string) (IsiProtoStat, error) {
	var (
		resp   IsiProtoStat
		params api.OrderedValues
	)

	const path = "/platform/1/statistics/current"

	//Set params to be key passed, and query all device ids at the same time
	params = api.NewOrderedValues([][]string{
		{"keys", key},
		{"devid", "all"},
	})

	//Query API
	err := c.API.Get(context.Background(), path, "", params, nil, &resp)
	if err != nil {
		log.Warnf("Unable to retrieve stats for %s: %s", key, err)
		return resp, err
	}
	return resp, nil
}

//GetSyncPolicies retrieve all sync iq policies
func GetSyncPolicies(c *goisilon.Client) (IsiSyncPolicies, error) {
	const path = "/platform/3/sync/policies"
	var resp IsiSyncPolicies

	err := c.API.Get(context.Background(), path, "", nil, nil, &resp)
	if err != nil {
		log.Warnf("Unable to retrieve sync")
		return resp, err
	}
	return resp, nil
}

//GetSnapshotsSummary retrieves summary statistics for snapshots
func GetSnapshotsSummary(c *goisilon.Client) (IsiSnapshotsSummary, error) {
	const path = "/platform/1/snapshot/snapshots-summary"
	var resp IsiSnapshotsSummary

	err := c.API.Get(context.Background(), path, "", nil, nil, &resp)
	if err != nil {
		log.Warnf("Unable to retrieve snapshot summary.")
		return resp, err
	}
	return resp, nil
}

//GetSnapshots returns a struct with all snapshots.
func GetSnapshots(c *goisilon.Client) (IsiSnapshots, error) {
	const path = "/platform/1/snapshot/snapshots"
	var resp IsiSnapshots

	err := c.API.Get(context.Background(), path, "", nil, nil, &resp)
	if err != nil {
		log.Warnf("Unable to retrieve snapshots. Err: %s", err)
		return resp, err
	}
	return resp, nil
}

func GetNodesPartitions(c *goisilon.Client) (IsiNodesPartitions, error) {
	const path = "/platform/3/cluster/nodes/all/partitions"
	var resp IsiNodesPartitions

	err := c.API.Get(context.Background(), path, "", nil, nil, &resp)
	if err != nil {
		log.Warn("Unabled to retrieve node partitions.")
		return resp, err
	}
	return resp, nil
}

func GetNodesStatus(c *goisilon.Client) (IsiNodesStatus, error) {
	const path = "/platform/3/cluster/nodes/all/status"
	var resp IsiNodesStatus

	err := c.API.Get(context.Background(), path, "", nil, nil, &resp)
	if err != nil {
		log.Warn("Unable to get nodes status")
		return resp, err
	}
	return resp, nil
}

func GetNodesHardware(c *goisilon.Client) (IsiNodesHardware, error) {
	const path = "/platform/3/cluster/nodes/all/hardware"
	var resp IsiNodesHardware

	err := c.API.Get(context.Background(), path, "", nil, nil, &resp)
	if err != nil {
		log.Warn("Unable to get nodes hardware information.")
		return resp, err
	}
	return resp, nil
}

func GetNodesState(c *goisilon.Client) (IsiNodesState, error) {
	const path = "/platform/3/cluster/nodes/all/state"
	var resp IsiNodesState

	err := c.API.Get(context.Background(), path, "", nil, nil, &resp)
	if err != nil {
		log.Warn("Unable to get nodes state.")
		return resp, err
	}
	return resp, nil
}

func GetStatfs(c *goisilon.Client) (IsiStatfs, error) {
	const path = "/platform/1/cluster/statfs"
	var resp IsiStatfs

	err := c.API.Get(context.Background(), path, "", nil, nil, &resp)
	if err != nil {
		log.Warn("Unabled to get statfs.")
		return resp, err
	}
	return resp, nil
}

func GetExportSummary(c *goisilon.Client) (IsiExportSummary, error) {
	const path = "/platform/2/protocols/nfs/exports-summary"
	var resp IsiExportSummary

	err := c.API.Get(context.Background(), path, "", nil, nil, &resp)
	if err != nil {
		log.Warn("Unabled to get nfs exports summary.")
		return resp, err
	}
	return resp, nil
}

func GetSharesSummary(c *goisilon.Client) (IsiSharesSummary, error) {
	const path = "/platform/3/protocols/smb/shares-summary"
	var resp IsiSharesSummary

	err := c.API.Get(context.Background(), path, "", nil, nil, &resp)
	if err != nil {
		log.Warn("Unabled to get smb shares summary.")
		return resp, err
	}
	return resp, nil
}

func GetStoragePools(c *goisilon.Client) (IsiStoragePools, error) {
	const path = "/platform/3/storagepool/storagepools"
	var resp IsiStoragePools

	err := c.API.Get(context.Background(), path, "", nil, nil, &resp)
	if err != nil {
		log.Warn("Unable to get storagepools info.")
		return resp, err
	}
	return resp, nil
}

func GetDriveInfo(c *goisilon.Client) (IsiNodesDrives, error) {
	const path = "/platform/3/cluster/nodes/all/drives"
	var resp IsiNodesDrives
	err := c.API.Get(context.Background(), path, "", nil, nil, &resp)
	if err != nil {
		log.Warn("Unable to get drive info.")
		return resp, err
	}
	return resp, nil
}
