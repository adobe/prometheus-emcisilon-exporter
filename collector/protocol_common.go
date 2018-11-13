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
	"sort"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	protocolState          map[string]*bool
	protosUpdated          bool
	nfsClientstatGathered  bool
	smbClientstatsGathered bool
	// This is the list of protocols to collect stats for. (Used by both node_protocols and cluster_protocols)
	protocols = map[string]bool{
		"cifs":      defaultDisabled,
		"ftp":       defaultEnabled,
		"hdfs":      defaultDisabled,
		"http":      defaultDisabled,
		"irp":       defaultDisabled,
		"jobd":      defaultEnabled,
		"lsass_in":  defaultEnabled,
		"lsass_out": defaultEnabled,
		"nfs":       defaultEnabled,
		"nfs3":      defaultEnabled,
		"nfs4":      defaultDisabled,
		"nlm":       defaultDisabled,
		"papi":      defaultEnabled,
		"siq":       defaultEnabled,
		"smb1":      defaultDisabled,
		"smb2":      defaultEnabled,
	}
)

// GetProtos creates a map of protocols to be collected.
func GetProtos() {
	protocolState = make(map[string]*bool)
	var protoKeys []string
	for k := range protocols {
		protoKeys = append(protoKeys, k)
	}
	sort.Strings(protoKeys)

	for _, v := range protoKeys {
		key := v
		state := protocols[key]

		var helpDefaultState string
		if state {
			helpDefaultState = "enabled"
		} else {
			helpDefaultState = "disabled"
		}

		flagName := fmt.Sprintf("collector.protocol_common.%s", key)
		flagHelp := fmt.Sprintf("Enable colllection for the %s protocol (default: %s).", key, helpDefaultState)
		defaultValue := fmt.Sprintf("%v", state)

		flag := kingpin.Flag(flagName, flagHelp).Default(defaultValue).Bool()
		protocolState[key] = flag
		protosUpdated = true
	}
}
