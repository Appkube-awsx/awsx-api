package handlers

import (
	"awsx-api/handlers/getElementDetails/EC2"
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
		EKS.GetEKSCPURequestPanel(w, r)
	}
}
