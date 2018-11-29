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

// IsiConfig is used to unmarshal the config api response
type IsiConfig struct {
	Description string `json:"description"`
	Devices     []struct {
		Devid int    `json:"devid"`
		GUID  string `json:"guid"`
		IsUp  bool   `json:"is_up"`
		Lnn   int    `json:"lnn"`
	} `json:"devices"`
	Encoding     string `json:"encoding"`
	GUID         string `json:"guid"`
	HasQuorum    bool   `json:"has_quorum"`
	IsCompliance bool   `json:"is_compliance"`
	IsVirtual    bool   `json:"is_virtual"`
	IsVonefs     bool   `json:"is_vonefs"`
	JoinMode     string `json:"join_mode"`
	LocalDevid   int    `json:"local_devid"`
	LocalLnn     int    `json:"local_lnn"`
	LocalSerial  string `json:"local_serial"`
	Name         string `json:"name"`
	OnefsVersion struct {
		Build     string `json:"build"`
		Copyright string `json:"copyright"`
		Reldate   int    `json:"reldate"`
		Release   string `json:"release"`
		Revision  string `json:"revision"`
		Type      string `json:"type"`
		Version   string `json:"version"`
	} `json:"onefs_version"`
	Timezone struct {
		Abbreviation string `json:"abbreviation"`
		Custom       string `json:"custom"`
		Name         string `json:"name"`
		Path         string `json:"path"`
	} `json:"timezone"`
	UpgradeType interface{} `json:"upgrade_type"`
}

// IsiIdentity is used to unmarshal the identity api response
type IsiIdentity struct {
	Description string `json:"description"`
	Logon       struct {
		Motd       string `json:"motd"`
		MotdHeader string `json:"motd_header"`
	} `json:"logon"`
	Name string `json:"name"`
}

//IsiMultiVal is the struct used to unmarshal a multi value stat
type IsiMultiVal struct {
	Stats []struct {
		Devid     int                  `json:"devid"`
		Error     interface{}          `json:"error"`
		ErrorCode interface{}          `json:"error_code"`
		Key       string               `json:"key"`
		Time      int                  `json:"time"`
		ValueSet  []map[string]float64 `json:"value"`
	} `json:"stats"`
}

//IsiSingleVal is the struct used to unmarshal a single value stat
type IsiSingleVal struct {
	Stats []struct {
		Devid     int         `json:"devid"`
		Error     interface{} `json:"error"`
		ErrorCode interface{} `json:"error_code"`
		Key       string      `json:"key"`
		Time      int         `json:"time"`
		Value     float64     `json:"value"`
	} `json:"stats"`
}

type IsiQuotas struct {
	Quotas []IsiQuota `json:"quotas"`
	Resume string     `json:"resume"`
}

type IsiQuota struct {
	Container                 bool              `json:"container"`
	Enforced                  bool              `json:"enforced"`
	ID                        string            `json:"id"`
	IncludeSnapshots          bool              `json:"include_snapshots"`
	Linked                    bool              `json:"linked"`
	Notifications             string            `json:"notifications"`
	Path                      string            `json:"path"`
	Persona                   interface{}       `json:"persona"`
	Ready                     bool              `json:"ready"`
	Thresholds                IsiQuotaThreshold `json:"thresholds"`
	ThresholdsIncludeOverhead bool              `json:"thresholds_include_overhead"`
	Type                      string            `json:"type"`
	Usage                     IsiQuotaUsage     `json:"usage"`
}

type IsiQuotaUsage struct {
	Inodes   float64 `json:"inodes"`
	Logical  float64 `json:"logical"`
	Physical float64 `json:"physical"`
}

type IsiQuotaThreshold struct {
	Advisory             float64     `json:"advisory"`
	AdvisoryExceeded     bool        `json:"advisory_exceeded"`
	AdvisoryLastExceeded interface{} `json:"advisory_last_exceeded"`
	Hard                 float64     `json:"hard"`
	HardExceeded         bool        `json:"hard_exceeded"`
	HardLastExceeded     interface{} `json:"hard_last_exceeded"`
	Soft                 float64     `json:"soft"`
	SoftExceeded         bool        `json:"soft_exceeded"`
	SoftGrace            float64     `json:"soft_grace"`
	SoftLastExceeded     interface{} `json:"soft_last_exceeded"`
}

type IsiQuotaSummaryResp struct {
	Summary IsiQuotaSummary `json:"summary"`
}

type IsiQuotaSummary struct {
	Count                   float64 `json:"count"`
	DefaultGroupQuotasCount float64 `json:"default_group_quotas_count"`
	DefaultUserQuotasCount  float64 `json:"default_user_quotas_count"`
	DirectoryQuotasCount    float64 `json:"directory_quotas_count"`
	GroupQuotasCount        float64 `json:"group_quotas_count"`
	LinkedQuotasCount       float64 `json:"linked_quotas_count"`
	UserQuotasCount         float64 `json:"user_quotas_count"`
}

type IsiProtoStat struct {
	Stats []struct {
		Devid     int         `json:"devid"`
		Error     interface{} `json:"error"`
		ErrorCode interface{} `json:"error_code"`
		Key       string      `json:"key"`
		Time      int         `json:"time"`
		Value     interface{} `json:"value"`
	} `json:"stats"`
}

type IsiProtoStatOp struct {
	ClassName string  `json:"class_name"`
	InMax     float64 `json:"in_max"`
	InMin     float64 `json:"in_min"`
	InRate    float64 `json:"in_rate"`
	OpCount   float64 `json:"op_count"`
	OpID      float64 `json:"op_id"`
	OpName    string  `json:"op_name"`
	OpRate    float64 `json:"op_rate"`
	OutMax    float64 `json:"out_max"`
	OutMin    float64 `json:"out_min"`
	OutRate   float64 `json:"out_rate"`
	TimeAvg   float64 `json:"time_avg"`
	TimeMax   float64 `json:"time_max"`
	TimeMin   float64 `json:"time_min"`
}

type IsiProtoStatTotal struct {
	InMax   float64 `json:"in_max"`
	InMin   float64 `json:"in_min"`
	InRate  float64 `json:"in_rate"`
	OpCount float64 `json:"op_count"`
	OpID    float64 `json:"op_id"`
	OpRate  float64 `json:"op_rate"`
	OutMax  float64 `json:"out_max"`
	OutMin  float64 `json:"out_min"`
	OutRate float64 `json:"out_rate"`
	TimeAvg float64 `json:"time_avg"`
	TimeMax float64 `json:"time_max"`
	TimeMin float64 `json:"time_min"`
}

type IsiSyncPolicies struct {
	Policies []struct {
		AcceleratedFailback       bool          `json:"accelerated_failback"`
		Action                    string        `json:"action"`
		BurstMode                 bool          `json:"burst_mode"`
		Changelist                bool          `json:"changelist"`
		CheckIntegrity            bool          `json:"check_integrity"`
		CloudDeepCopy             string        `json:"cloud_deep_copy"`
		Conflicted                bool          `json:"conflicted"`
		Description               string        `json:"description"`
		DisableFileSplit          bool          `json:"disable_file_split"`
		DisableFofb               bool          `json:"disable_fofb"`
		DisableStf                bool          `json:"disable_stf"`
		Enabled                   bool          `json:"enabled"`
		ExpectedDataloss          bool          `json:"expected_dataloss"`
		FileMatchingPattern       struct{}      `json:"file_matching_pattern"`
		ForceInterface            bool          `json:"force_interface"`
		HasSyncState              bool          `json:"has_sync_state"`
		ID                        string        `json:"id"`
		JobDelay                  interface{}   `json:"job_delay"`
		LastJobState              string        `json:"last_job_state"`
		LastStarted               float64       `json:"last_started"`
		LastSuccess               float64       `json:"last_success"`
		LogLevel                  string        `json:"log_level"`
		LogRemovedFiles           bool          `json:"log_removed_files"`
		Name                      string        `json:"name"`
		NextRun                   float64       `json:"next_run"`
		PasswordSet               bool          `json:"password_set"`
		Priority                  float64       `json:"priority"`
		ReportMaxAge              float64       `json:"report_max_age"`
		ReportMaxCount            float64       `json:"report_max_count"`
		RestrictTargetNetwork     bool          `json:"restrict_target_network"`
		RpoAlert                  interface{}   `json:"rpo_alert"`
		Schedule                  string        `json:"schedule"`
		SkipWhenSourceUnmodified  bool          `json:"skip_when_source_unmodified"`
		SnapshotSyncExisting      bool          `json:"snapshot_sync_existing"`
		SnapshotSyncPattern       string        `json:"snapshot_sync_pattern"`
		SourceExcludeDirectories  []interface{} `json:"source_exclude_directories"`
		SourceIncludeDirectories  []interface{} `json:"source_include_directories"`
		SourceNetwork             interface{}   `json:"source_network"`
		SourceRootPath            string        `json:"source_root_path"`
		SourceSnapshotArchive     bool          `json:"source_snapshot_archive"`
		SourceSnapshotExpiration  float64       `json:"source_snapshot_expiration"`
		SourceSnapshotPattern     string        `json:"source_snapshot_pattern"`
		TargetCompareInitialSync  bool          `json:"target_compare_initial_sync"`
		TargetDetectModifications bool          `json:"target_detect_modifications"`
		TargetHost                string        `json:"target_host"`
		TargetPath                string        `json:"target_path"`
		TargetSnapshotAlias       string        `json:"target_snapshot_alias"`
		TargetSnapshotArchive     bool          `json:"target_snapshot_archive"`
		TargetSnapshotExpiration  float64       `json:"target_snapshot_expiration"`
		TargetSnapshotPattern     string        `json:"target_snapshot_pattern"`
		WorkersPerNode            float64       `json:"workers_per_node"`
	} `json:"policies"`
	Resume interface{} `json:"resume"`
	Total  int         `json:"total"`
}

type IsiSnapshotsSummary struct {
	Summary struct {
		ActiveCount   float64 `json:"active_count"`
		ActiveSize    float64 `json:"active_size"`
		AliasesCount  float64 `json:"aliases_count"`
		Count         float64 `json:"count"`
		DeletingCount float64 `json:"deleting_count"`
		DeletingSize  float64 `json:"deleting_size"`
		ShadowBytes   float64 `json:"shadow_bytes"`
		Size          float64 `json:"size"`
	} `json:"summary"`
}

type IsiSnapshots struct {
	Resume    string `json:"resume"`
	Snapshots []struct {
		Created       int64       `json:"created"`
		Expires       int64       `json:"expires"`
		HasLocks      bool        `json:"has_locks"`
		ID            float64     `json:"id"`
		Name          string      `json:"name"`
		Path          string      `json:"path"`
		PctFilesystem float64     `json:"pct_filesystem"`
		PctReserve    float64     `json:"pct_reserve"`
		Schedule      string      `json:"schedule"`
		ShadowBytes   float64     `json:"shadow_bytes"`
		Size          float64     `json:"size"`
		State         string      `json:"state"`
		TargetID      interface{} `json:"target_id"`
		TargetName    interface{} `json:"target_name"`
	} `json:"snapshots"`
	Total float64 `json:"total"`
}

type IsiNodesStatus struct {
	Errors []interface{} `json:"errors"`
	Nodes  []struct {
		Batterystatus struct {
			LastTestTime1 string `json:"last_test_time1"`
			LastTestTime2 string `json:"last_test_time2"`
			NextTestTime1 string `json:"next_test_time1"`
			NextTestTime2 string `json:"next_test_time2"`
			Present       bool   `json:"present"`
			Result1       string `json:"result1"`
			Result2       string `json:"result2"`
			Status1       string `json:"status1"`
			Status2       string `json:"status2"`
			Supported     bool   `json:"supported"`
		} `json:"batterystatus"`
		Capacity []struct {
			Bytes float64 `json:"bytes"`
			Count float64 `json:"count"`
			Type  string  `json:"type"`
		} `json:"capacity"`
		CPU struct {
			Model      string `json:"model"`
			Overtemp   string `json:"overtemp"`
			Proc       string `json:"proc"`
			SpeedLimit string `json:"speed_limit"`
		} `json:"cpu"`
		ID    float64 `json:"id"`
		Lnn   float64 `json:"lnn"`
		Nvram struct {
			Batteries []struct {
				Color   string  `json:"color"`
				ID      float64 `json:"id"`
				Status  string  `json:"status"`
				Voltage string  `json:"voltage"`
			} `json:"batteries"`
			BatteryCount       float64 `json:"battery_count"`
			ChargeStatus       string  `json:"charge_status"`
			ChargeStatusNumber float64 `json:"charge_status_number"`
			Device             string  `json:"device"`
			Present            bool    `json:"present"`
			PresentFlash       bool    `json:"present_flash"`
			PresentSize        float64 `json:"present_size"`
			PresentType        string  `json:"present_type"`
			ShipMode           float64 `json:"ship_mode"`
			Supported          bool    `json:"supported"`
			SupportedFlash     bool    `json:"supported_flash"`
			SupportedSize      float64 `json:"supported_size"`
			SupportedType      string  `json:"supported_type"`
		} `json:"nvram"`
		Powersupplies struct {
			Count    float64 `json:"count"`
			Failures float64 `json:"failures"`
			HasCff   bool    `json:"has_cff"`
			Status   string  `json:"status"`
			Supplies []struct {
				Chassis  float64 `json:"chassis"`
				Firmware string  `json:"firmware"`
				Good     string  `json:"good"`
				ID       float64 `json:"id"`
				Name     string  `json:"name"`
				Status   string  `json:"status"`
				Type     string  `json:"type"`
			} `json:"supplies"`
			SupportsCff bool `json:"supports_cff"`
		} `json:"powersupplies"`
		Release string  `json:"release"`
		Uptime  float64 `json:"uptime"`
		Version string  `json:"version"`
	} `json:"nodes"`
	Total float64 `json:"total"`
}

type IsiNodesHardware struct {
	Errors []interface{} `json:"errors"`
	Nodes  []struct {
		Chassis         string   `json:"chassis"`
		ChassisCode     string   `json:"chassis_code"`
		ChassisCount    string   `json:"chassis_count"`
		Class           string   `json:"class"`
		ConfigurationID string   `json:"configuration_id"`
		CPU             string   `json:"cpu"`
		DiskController  string   `json:"disk_controller"`
		DiskExpander    string   `json:"disk_expander"`
		FamilyCode      string   `json:"family_code"`
		FlashDrive      string   `json:"flash_drive"`
		GenerationCode  string   `json:"generation_code"`
		Hwgen           string   `json:"hwgen"`
		ID              float64  `json:"id"`
		ImbVersion      string   `json:"imb_version"`
		Infiniband      string   `json:"infiniband"`
		LcdVersion      string   `json:"lcd_version"`
		Lnn             float64  `json:"lnn"`
		Motherboard     string   `json:"motherboard"`
		NetInterfaces   string   `json:"net_interfaces"`
		Nvram           string   `json:"nvram"`
		Powersupplies   []string `json:"powersupplies"`
		Processor       string   `json:"processor"`
		Product         string   `json:"product"`
		RAM             float64  `json:"ram"`
		SerialNumber    string   `json:"serial_number"`
		Series          string   `json:"series"`
		StorageClass    string   `json:"storage_class"`
	} `json:"nodes"`
	Total float64 `json:"total"`
}

type IsiNodesPartitions struct {
	Errors []interface{} `json:"errors"`
	Nodes  []struct {
		Count      float64 `json:"count"`
		ID         float64 `json:"id"`
		Lnn        float64 `json:"lnn"`
		Partitions []struct {
			BlockSize        float64 `json:"block_size"`
			Capacity         float64 `json:"capacity"`
			ComponentDevices string  `json:"component_devices"`
			MountPoint       string  `json:"mount_point"`
			PercentUsed      string  `json:"percent_used"`
			Statfs           struct {
				FBavail      float64 `json:"f_bavail"`
				FBfree       float64 `json:"f_bfree"`
				FBlocks      float64 `json:"f_blocks"`
				FBsize       float64 `json:"f_bsize"`
				FFfree       float64 `json:"f_ffree"`
				FFiles       float64 `json:"f_files"`
				FFlags       float64 `json:"f_flags"`
				FFstypename  string  `json:"f_fstypename"`
				FIosize      float64 `json:"f_iosize"`
				FMntfromname string  `json:"f_mntfromname"`
				FMntonname   string  `json:"f_mntonname"`
				FNamemax     float64 `json:"f_namemax"`
				FOwner       float64 `json:"f_owner"`
				FType        float64 `json:"f_type"`
				FVersion     float64 `json:"f_version"`
			} `json:"statfs"`
			Used float64 `json:"used"`
		} `json:"partitions"`
	} `json:"nodes"`
	Total float64 `json:"total"`
}

type IsiNodesState struct {
	Errors []interface{} `json:"errors"`
	Nodes  []struct {
		ID       float64 `json:"id"`
		Lnn      float64 `json:"lnn"`
		Readonly struct {
			Allowed bool    `json:"allowed"`
			Enabled bool    `json:"enabled"`
			Mode    bool    `json:"mode"`
			Status  string  `json:"status"`
			Valid   bool    `json:"valid"`
			Value   float64 `json:"value"`
		} `json:"readonly"`
		Servicelight struct {
			Enabled   bool `json:"enabled"`
			Present   bool `json:"present"`
			Supported bool `json:"supported"`
			Valid     bool `json:"valid"`
		} `json:"servicelight"`
		Smartfail struct {
			Dead             bool `json:"dead"`
			Down             bool `json:"down"`
			InCluster        bool `json:"in_cluster"`
			Readonly         bool `json:"readonly"`
			ShutdownReadonly bool `json:"shutdown_readonly"`
			Smartfailed      bool `json:"smartfailed"`
		} `json:"smartfail"`
	} `json:"nodes"`
	Total float64 `json:"total"`
}

type IsiStatfs struct {
	FBavail      float64 `json:"f_bavail"`
	FBfree       float64 `json:"f_bfree"`
	FBlocks      float64 `json:"f_blocks"`
	FBsize       float64 `json:"f_bsize"`
	FFfree       float64 `json:"f_ffree"`
	FFiles       float64 `json:"f_files"`
	FFlags       float64 `json:"f_flags"`
	FFstypename  string  `json:"f_fstypename"`
	FIosize      float64 `json:"f_iosize"`
	FMntfromname string  `json:"f_mntfromname"`
	FMntonname   string  `json:"f_mntonname"`
	FNamemax     float64 `json:"f_namemax"`
	FOwner       float64 `json:"f_owner"`
	FType        float64 `json:"f_type"`
	FVersion     float64 `json:"f_version"`
}

type IsiExportSummary struct {
	Summary struct {
		Count float64 `json:"count"`
	} `json:"summary"`
}

type IsiSharesSummary struct {
	Summary struct {
		Count float64 `json:"count"`
	} `json:"summary"`
}

type IsiStoragePools struct {
	Storagepools []struct {
		Children    []string      `json:"children,omitempty"`
		HealthFlags []interface{} `json:"health_flags"`
		ID          float64       `json:"id"`
		Lnns        []float64     `json:"lnns"`
		Name        string        `json:"name"`
		Type        string        `json:"type"`
		Usage       struct {
			AvailBytes           string `json:"avail_bytes"`
			AvailSsdBytes        string `json:"avail_ssd_bytes"`
			Balanced             bool   `json:"balanced"`
			FreeBytes            string `json:"free_bytes"`
			FreeSsdBytes         string `json:"free_ssd_bytes"`
			TotalBytes           string `json:"total_bytes"`
			TotalSsdBytes        string `json:"total_ssd_bytes"`
			VirtualHotSpareBytes string `json:"virtual_hot_spare_bytes"`
		} `json:"usage"`
		CanDisableL3     bool   `json:"can_disable_l3,omitempty"`
		CanEnableL3      bool   `json:"can_enable_l3,omitempty"`
		L3               bool   `json:"l3,omitempty"`
		L3Status         string `json:"l3_status,omitempty"`
		Manual           bool   `json:"manual,omitempty"`
		ProtectionPolicy string `json:"protection_policy,omitempty"`
	} `json:"storagepools"`
	Total float64 `json:"total"`
}
