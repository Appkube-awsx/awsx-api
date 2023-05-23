package handlers

import (
	"awsx-api/log"
	"encoding/json"
	"fmt"
	"github.com/Appkube-awsx/awsx-cloudelements/cmd"
	"net/http"
)

func GetAppconfig(w http.ResponseWriter, r *http.Request) {
	log.Info("Starting GetAppconfig api")

	vaultUrl := r.URL.Query().Get("vaultUrl")
	if vaultUrl == "" {
		log.Error("VaultUrl not provided")
		http.Error(w, fmt.Sprintf("VaultUrl not provided"), http.StatusBadRequest)
		return
	}

	accountId := r.URL.Query().Get("accountId")
	if accountId == "" {
		log.Error("AccountId not provided")
		http.Error(w, fmt.Sprintf("AccountId not provided"), http.StatusBadRequest)
		return
	}

	region := r.URL.Query().Get("zone")
	if region == "" {
		log.Error("Region/Zone not provided")
		http.Error(w, fmt.Sprintf("Region/Zone not provided"), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	result, err := cmd.GetCloudConfigSummary(region, vaultUrl, accountId)
	if err != nil {
		log.Error("Exception in getting cloud config summary: %v", err)
		http.Error(w, fmt.Sprintf("Exception in getting cloud config summary"), http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(result)

	//json.NewEncoder(w).Encode(cmd.GetConfig(raw.AccessDetails.Region, raw.AccessDetails.CrossAccountRoleArn, raw.AccessDetails.AccessKey, raw.AccessDetails.SecretKey, raw.AccessDetails.ExternalId))

	log.Info("GetAppconfig completed")

}
