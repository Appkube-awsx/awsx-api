package handlers

import (
	"awsx-api/handlers/getElementDetails/EC2"
	"awsx-api/handlers/getElementDetails/ECS"
	"awsx-api/handlers/getElementDetails/EKS"
	"awsx-api/log"
	"net/http"
)

func ExecuteQuery(w http.ResponseWriter, r *http.Request) {
	log.Info("Starting /awsx-api/execute-query api")
	query := r.URL.Query().Get("query")
	elementType := r.URL.Query().Get("elementType")
	if elementType == "EC2" && query == "cpu_utilization_panel" {
		EC2.GetCpuUtilizationPanel(w, r)
	}
	if elementType == "EKS" && query == "cpu_utilization_panel" {
		EKS.GetEKScpuUtilizationPanel(w, r)
	}
	if elementType == "EKS" && query == "memory_utilization_panel" {
		EKS.GetEKSMemoryUtilizationPanel(w, r)
	}
	if elementType == "EKS" && query == "network_utilization_panel" {
		EKS.GetEKSNetworkUtilizationPanel(w, r)
	}
	if elementType == "EKS" && query == "cpu_requests_panel" {
		EKS.GetCpuRequestsPanel(w, r)
	}
	if elementType == "ECS" && query == "cpu_utilization_panel" {
		ECS.GetECScpuUtilizationPanel(w, r)
	}
	if elementType == "ECS" && query == "memory_utilization_panel" {
		ECS.GetECSMemoryUtilizationPanel(w, r)
	}
	if elementType == "EC2" && query == "memory_utilization_panel" {
		EC2.GetMemoryUtilizationPanel(w, r)
	}
	if elementType == "EC2" && query == "network_utilization_panel" {
		EC2.GetNetworkUtilizationPanel(w, r)
	}
	if elementType == "EC2" && query == "cpu_usage_user_panel" {
		EC2.GetCPUUsageUserPanel(w, r)
	}
	if elementType == "EC2" && query == "cpu_usage_sys_panel" {
		EC2.GetCPUUsageSysPanel(w, r)
	}
	if elementType == "EC2" && query == "cpu_usage_nice_panel" {
		EC2.GetCPUUsageNicePanel(w, r)
	}
	if elementType == "EC2" && query == "cpu_usage_idle_panel" {
		EC2.GetCPUUsageIdlePanel(w, r)
	}
	if elementType == "EC2" && query == "mem_usage_free_panel" {
		EC2.GetMemUsageFreePanel(w, r)
	}
	if elementType == "EC2" && query == "mem_cached_panel" {
		EC2.GetMemCachePanel(w, r)
	}
	if elementType == "EC2" && query == "mem_usage_total_panel" {
		EC2.GetMemUsageTotal(w, r)
	}

	if elementType == "EC2" && query == "mem_usage_used_panel" {
		EC2.GetMemUsageUsed(w, r)
	}

	if elementType == "EC2" && query == "disk_writes_panel" {
		EC2.GetDiskWritePanel(w, r)
	}

	if elementType == "EC2" && query == "disk_reads_panel" {
		EC2.GetDiskReadPanel(w, r)
	}

	if elementType == "EC2" && query == "disk_available_panel" {
		EC2.GetDiskAvailablePanel(w, r)
	}

	if elementType == "EC2" && query == "disk_used_panel" {
		EC2.GetDiskUsedPanel(w, r)
	}

	if elementType == "EC2" && query == "net_inpackets_panel" {
		EC2.GetNetworkInPacketsPanel(w, r)
	}

	if elementType == "EC2" && query == "net_inbytes_panel" {
		EC2.GetNetworkInBytesPanel(w, r)
	}

	if elementType == "EC2" && query == "net_outbytes_panel" {
		EC2.GetNetworkOutBytesPanel(w, r)
	}

	if elementType == "EC2" && query == "net_outpackets_panel" {
		EC2.GetNetworkOutPacketsPanel(w, r)
	}

}
