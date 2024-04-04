package handlers

import (
	"awsx-api/handlers/getElementDetails/ApiGateway"
	"awsx-api/handlers/getElementDetails/EC2"
	"awsx-api/handlers/getElementDetails/ECS"
	"awsx-api/handlers/getElementDetails/EKS"
	"awsx-api/handlers/getElementDetails/Lambda"
	"awsx-api/handlers/getElementDetails/RDS"

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
	if elementType == "ECS" || elementType == "AWS/ECS" && query == "top_events_panel" {
		ECS.GetTopEventsPanel(w, r)
	}
	if elementType == "ECS" || elementType == "AWS/ECS" && query == "registration_events_panel" {
		ECS.GetRegistrationEventsPanel(w, r)
	}
	if elementType == "ECS" || elementType == "AWS/ECS" && query == "deregistration_events_panel" {
		ECS.GetDeRegistrationEventsPanel(w, r)
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
	if elementType == "EC2" && query == "instance_hours_stopped_panel" {
		EC2.InstanceHourStoppedPanel(w, r)
	}
	if elementType == "EC2" && query == "instance_running_hour_panel" {
		EC2.InstanceRunningHourPanelHandler(w, r)
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
	if elementType == "EC2" && query == "storage_utilization_panel" {
		EC2.GetStorageUtilizationPanel(w, r)
	}
	if elementType == "EC2" && query == "disk_io_panel" {
		EC2.GetDiskIOPanel(w, r)
	}
	if elementType == "EC2" && query == "cpu_utilization_graph_panel" {
		EC2.GetCPUUtilizationPanel(w, r)
	}
	if elementType == "EC2" && query == "memory_utilization_graph_panel" {
		EC2.GetMemoryUtilizationPaneel(w, r)
	}
	if elementType == "EC2" && query == "network_traffic_panel" {
		EC2.GetNetworkTrafficPanel(w, r)
	}
	if elementType == "EKS" && query == "resource_utilization_patterns_panel" {
		EKS.GetResourceUtilizationPanel(w, r)
	}
	if elementType == "EKS" && query == "node_stability_index_panel" {
		EKS.GetNodeStabilityIndexPanel(w, r)
	}
	if elementType == "EKS" && query == "disk_utilization_panel" {
		EKS.GetEKSDiskUtilizationPanel(w, r)
	}
	if elementType == "EKS" && query == "disk_io_performance_panel" {
		EKS.GetEKSDiskIoPerformancePanel(w, r)
	}
	if elementType == "EKS" && query == "node_condition_panel" {
		EKS.GetEKSNodeConditionPanel(w, r)
	}
	if elementType == "EKS" && query == "storage_utilization_panel" {
		EKS.GetStorageUtilizationPanell(w, r)
	}
	if elementType == "EKS" && query == "node_failure_panel" {
		EKS.GetNodeFailurePanel(w, r)
	}
	if elementType == "EKS" && query == "incident_response_time_panel" {
		EKS.GetIncidentResponseTimePanel(w, r)
	}
	if elementType == "Lambda" && query == "used_and_unused_memory_data_panel" {
		Lambda.GetUsedAndUnusedMemoryDataPanel(w, r)
	}
	if elementType == "Lambda" && query == "max_memory_used_panel" {
		Lambda.GetMaxMemoryUsedPanel(w, r)
	}
	if elementType == "Lambda" && query == "execution_time_panel" {
		Lambda.GetExecutionTimePanel(w, r)
	}
	if elementType == "Lambda" && query == "max_memory_used_graph_panel" {
		Lambda.GetMaxMemoryUsedPanell(w, r)
	}
	if elementType == "Lambda" && query == "cold_start_duration_panel" {
		Lambda.GetColdStartDurationPanel(w, r)
	}
	if elementType == "Lambda" && query == "concurrency_panel" {
		Lambda.GetConcurrencyPanel(w, r)
	}
	if elementType == "Lambda" && query == "functions_by_region_panel" {
		Lambda.GetFunctionByRegionPanel(w, r)
	}
	if elementType == "Lambda" && query == "throttles_panel" {
		Lambda.GetThrottlesPanel(w, r)
	}
	if elementType == "Lambda" && query == "number_of_calls_panel" {
		Lambda.GetNumberOfCallsPanel(w, r)
	}
	if elementType == "Lambda" && query == "error_messages_count_panel" {
		Lambda.GetErrorMsgCountPanel(w, r)
	}
	if elementType == "Lambda" && query == "throttling_trends_panel" {
		Lambda.GetThrottlingTrendsPanel(w, r)
	}
	if elementType == "Lambda" && query == "invocation_trend_panel" {
		Lambda.GetInvocationTrendPanel(w, r)
	}
	if elementType == "Lambda" && query == "error_and_warning_events_panel" {
		Lambda.GetErrorAndWarningEventsPanel(w, r)
	}
	if elementType == "Lambda" && query == "success_and_failed_function_panel" {
		Lambda.GetSuccessAndFailedFunctionPanel(w, r)
	}
	if elementType == "RDS" && query == "cpu_utilization_panel" {
		RDS.GetCpuUtilizationPanel(w, r)
	}
	if elementType == "RDS" && query == "network_utilization_panel" {
		RDS.GetNetworkUtilizationPanel(w, r)
	}
	if elementType == "RDS" && query == "cpu_utilization_graph_panel" {
		RDS.GetCPUUtilizationPanel(w, r)
	}
	if elementType == "RDS" && query == "alert_and_notification_panel" {
		RDS.GetAlertsAndNotificationsPanel(w, r)
	}
	if elementType == "RDS" && query == "instance_health_check_panel" {
		RDS.GetInstanceHealthCheck(w, r)
	}
	if elementType == "RDS" && query == "freeable_memory_panel" {
		RDS.GetFreeableMemoryPanel(w, r)
	}
	if elementType == "RDS" && query == "cpu_credit_balance_panel" {
		RDS.GetCpuCreditBalancePanel(w, r)
	}
	if elementType == "RDS" && query == "cpu_credit_usage_panel" {
		RDS.GetCpuCreditUsagePanel(w, r)
	}
	if elementType == "RDS" && query == "cpu_surplus_credit_balance_panel" {
		RDS.GetCPUSurplusCreditBalancePanel(w, r)
	}
	if elementType == "RDS" && query == "cpu_surplus_credits_charged_panel" {
		RDS.GetCPUSurplusCreditChargedPanel(w, r)
	}
	if elementType == "RDS" && query == "database_connections_panel" {
		RDS.GetDatabaseConnectionPanel(w, r)
	}
	if elementType == "RDS" && query == "database_workload_overview_panel" {
		RDS.GetDatabaseWorkloadOverviewPanel(w, r)
	}
	if elementType == "RDS" && query == "db_load_cpu_panel" {
		RDS.GetDBLoadCPULoadPanel(w, r)
	}
	if elementType == "RDS" && query == "db_load_non_cpu_panel" {
		RDS.GetDBLoadNonCPUPanel(w, r)
	}
	if elementType == "RDS" && query == "disk_queue_depth_panel" {
		RDS.GetDiskQueueDepthPanel(w, r)
	}
	if elementType == "RDS" && query == "free_storage_space_panel" {
		RDS.GetFreeStorageSpacePanel(w, r)
	}
	if elementType == "RDS" && query == "index_size_panel" {
		RDS.GetIndexSizePanel(w, r)
	}
	if elementType == "RDS" && query == "iops_panel" {
		RDS.GetIOPPanel(w, r)
	}
	if elementType == "RDS" && query == "network_receive_throughput_panel" {
		RDS.GetNetworkReceiveThroughputPanel(w, r)
	}
	if elementType == "RDS" && query == "network_traffic_panel" {
		RDS.GetNetworkTrafficPanel(w, r)
	}
	if elementType == "RDS" && query == "network_transmit_throughput_panel" {
		RDS.GetNetworkTransmitThroughputPanel(w, r)
	}
	if elementType == "RDS" && query == "replication_slot_disk_usage" {
		RDS.GetReplicationSlotDiskUsagePanel(w, r)
	}
	if elementType == "RDS" && query == "read_iops_panel" {
		RDS.GetReadIOPSPanel(w, r)
	}
	if elementType == "RDS" && query == "storage_utilization_panel" {
		RDS.GetStorageUtilizationPanel(w, r)
	}
	if elementType == "RDS" && query == "latency_analysis_panel" {
		RDS.GetLatencyAnalysisPanel(w, r)
	}
	if elementType == "RDS" && query == "write_iops_panel" {
		RDS.GetWriteIOPSPanel(w, r)
	}
	if elementType == "RDS" && query == "transaction_logs_generation_panel" {
		RDS.GetTransactionLogsGenerationPanel(w, r)
	}
	if elementType == "RDS" && query == "transaction_logs_disk_usage_panel" {
		RDS.GetTransactionLogsDiskPanel(w, r)
	}
	if elementType == "RDS" && query == "maintenance_schedule_overview_panel" {
		RDS.ScheduleOverviewPanel(w, r)
	}
	if elementType == "ApiGateway" && query == "uptime_percentage_panel" {
		ApiGateway.GetUptimePercentagePanel(w, r)
	}
	if elementType == "ApiGateway" && query == "uptime_of_deployment_stages" {
		ApiGateway.GetUptimeOfDeploymentPanel(w, r)
	}
	if elementType == "ApiGateway" && query == "4xx_errors_panel" {
		ApiGateway.Get4XXErrorsPanel(w, r)
	}
	if elementType == "ApiGateway" && query == "5xx_errors_panel" {
		ApiGateway.GetApi5xxErrorsPanel(w, r)
	}
	if elementType == "ApiGateway" && query == "total_api_calls_panel" {
		ApiGateway.GetTotalApiCallsPanel(w, r)
	}
	if elementType == "ApiGateway" && query == "latency_panel" {
		ApiGateway.GetLatencyPanel(w, r)
	}
	if elementType == "ApiGateway" && query == "integration_latency_panel" {
		ApiGateway.GetIntegrationLatencyPanel(w, r)
	}
	if elementType == "ApiGateway" && query == "cache_hit_count_panel" {
		ApiGateway.GetCacheHitsPanel(w, r)
	}
	if elementType == "ApiGateway" && query == "cache_miss_count_panel" {
		ApiGateway.GetCacheMissPanel(w, r)
	}
	// if elementType == "ApiGateway" && query == "downtime_incident_panel" {
	// 	ApiGateway.GetDowntimeIncidentPanel(w, r)
	// }
}
