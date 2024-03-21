package handlers

import (
	"awsx-api/handlers/getElementDetails/EC2"
	"awsx-api/handlers/getElementDetails/ECS"
	"awsx-api/handlers/getElementDetails/EKS"
	"awsx-api/handlers/getLandingZoneDetails"
	"awsx-api/log"
	"net/http"
)

func ExecuteQuery(w http.ResponseWriter, r *http.Request) {
	log.Info("Starting /awsx-api/execute-query api")
	query := r.URL.Query().Get("query")
	elementType := r.URL.Query().Get("elementType")
	if elementType == "landingZone" {
		getLandingZoneDetails.ExecuteLandingzoneQueries(w, r)
	}

	if elementType == "EC2" && query == "cpu_utilization_panel" {
		EC2.GetCpuUtilizationPanel(w, r)
	}
	if elementType == "EKS" && query == "cpu_utilization_panel" {
		EKS.GetEKScpuUtilizationPanel(w, r)
	}
	if elementType == "EKS" && query == "cpu_node_utilization_panel" {
		EKS.GetEKSCPUUtilizationNodeGraphPanel(w, r)
	}
	if elementType == "EKS" && query == "cpu_graph_utilization_panel" {
		EKS.GetEKSCPUUtilizationPanel(w, r)
	}
	if elementType == "EKS" && query == "memory_utilization_panel" {
		EKS.GetEKSMemoryUtilizationPanel(w, r)
	}
	if elementType == "EKS" && query == "network_utilization_panel" {
		EKS.GetEKSNetworkUtilizationPanel(w, r)
	}
	if elementType == "EKS" && query == "allocatable_cpu_panel" {
		EKS.GetEKSAllocatableCPUPanel(w, r)
	}
	if elementType == "EKS" && query == "allocatable_memory_panel" {
		EKS.GetEKSAllocatableMemoryPanel(w, r)
	}
	if elementType == "EKS" && query == "cpu_limits_panel" {
		EKS.GetEKSCPULimitsPanel(w, r)
	}
	if elementType == "EKS" && query == "cpu_requests_panel" {
		EKS.GetEKSCPURequestsPanel(w, r)
	}
	if elementType == "EKS" && query == "memory_limits_panel" {
		EKS.GetEKSMemoryLimitsPanel(w, r)
	}
	if elementType == "EKS" && query == "memory_requests_panel" {
		EKS.GetEKSMemoryRequestPanel(w, r)
	}
	if elementType == "EKS" && query == "memory_usage_panel" {
		EKS.GetEKSMemoryUsagePanel(w, r)
	}
	if elementType == "EKS" && query == "memory_graph_utilization_panel" {
		EKS.GetEKSMemoryUtilizationGraphPanel(w, r)
	}
	if elementType == "EKS" && query == "network_availability_panel" {
		EKS.GetEKSNetworkAvailabilityPanel(w, r)
	}
	if elementType == "EKS" && query == "network_in_out_panel" {
		EKS.GetEKSNetworkInOutPanel(w, r)
	}
	if elementType == "EKS" && query == "network_throughput_panel" {
		EKS.GetEKSNeworkThroughputPanel(w, r)
	}
	if elementType == "EKS" && query == "network_throughput_single_panel" {
		EKS.GetNetworkThroughputSinglePanel(w, r)
	}
	if elementType == "EKS" && query == "node_capacity_panel" {
		EKS.GetEKSNodeCapacityPanel(w, r)
	}
	if elementType == "EKS" && query == "node_downtime_panel" {
		EKS.GetEKSDowntimePanel(w, r)
	}
	if elementType == "EKS" && query == "node_uptime_panel" {
		EKS.NodeUptimePanelHandler(w, r)
	}
	if elementType == "EKS" && query == "node_event_logs_panel" {
		EKS.GetEKSEventLogsPanel(w, r)
	}
	if elementType == "EKS" && query == "service_availability_panel" {
		EKS.GetEKSServiceAvailabilityPanel(w, r)
	}

	if elementType == "ECS" && query == "cpu_utilization_panel" {
		ECS.GetECScpuUtilizationPanel(w, r)
	}
	if elementType == "ECS" && query == "memory_utilization_panel" {
		ECS.GetECSMemoryUtilizationPanel(w, r)
	}
	if elementType == "ECS" && query == "cpu_reservation_panel" {
		ECS.GetCPUReservationData(w, r)
	}
	if elementType == "ECS" && query == "memory_reservation_panel" {
		ECS.GetMemoryReservationData(w, r)
	}
	if elementType == "ECS" || elementType == "AWS/ECS" && query == "net_rxinbytes_panel" {
		ECS.GetECSNetworkRxInBytesPanel(w, r)
	}
	if elementType == "ECS" || elementType == "AWS/ECS" && query == "net_txinbytes_panel" {
		ECS.GetECSNetworkTxInBytesPanel(w, r)
	}
	if elementType == "ECS" || elementType == "AWS/ECS" && query == "volume_read_bytes_panel" {
		ECS.GetECSReadBytesPanel(w, r)
	}
	if elementType == "ECS" || elementType == "AWS/ECS" && query == "volume_write_bytes_panel" {
		ECS.GetECSWriteBytesPanel(w, r)
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
	if elementType == "EC2" && query == "net_throughput_panel" {
		EC2.GetNetworkThroughputPanel(w, r)
	}
	if elementType == "EC2" && query == "latency_panel" {
		EC2.GetLatencyPanel(w, r)
	}
	if elementType == "EC2" && query == "custom_alert_panel" {
		EC2.GetCustomAlert(w, r)
	}
	if elementType == "EC2" && query == "alerts_and_notifications_panel" {
		EC2.GetAlertsAndNotificationsPanel(w, r)
	}
	if elementType == "EC2" && query == "instance_start_count_panel" {
		EC2.InstanceStartCountPanelHandler(w, r)
	}
	if elementType == "EC2" && query == "instance_stop_count_panel" {
		EC2.InstanceStopCountPanelHandler(w, r)
	}
	// if elementType == "EC2" && query == "instance_running_hour_panel" {
	// 	EC2.InstanceStartPanelHandler(w, r)
	// }
	if elementType == "EC2" && query == "network_inbound_panel" {
		EC2.GetNetworkInboundPanell(w, r)
	}
	if elementType == "EC2" && query == "network_outbound_panel" {
		EC2.GetNetworkOutboundPanell(w, r)
	}
	if elementType == "EC2" && query == "instance_status_panel" {
		EC2.GetInstanceStatus(w, r)
	}
	if elementType == "EC2" && query == "instance_health_check_panel" {
		EC2.GetInstanceHealthCheck(w, r)
	}
	if elementType == "EC2" && query == "error_rate_panel" {
		EC2.GetInstanceErrorRatePanel(w, r)
	}
	if elementType == "EC2" && query == "error_tracking_panel" {
		EC2.ErrorTrackingHandler(w, r)
	}
	if elementType == "EC2" && query == "hosted_services_overview_panel" {
		EC2.HostedServicesOverviewHandler(w, r)
	}
	if elementType == "EKS" && query == "node_stability_index_panel" {
		EKS.GetNodeStabilityIndexPanel(w, r)
	}
}
