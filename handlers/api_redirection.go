package handlers

import (
	"awsx-api/handlers/getElementDetails/EC2"
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
}
