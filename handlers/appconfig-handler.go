package handlers

import (
	"awsx-api/log"
	"encoding/json"
	"fmt"
	"github.com/Appkube-awsx/awsx-cloudelements/controller"
	"net/http"
)

func GetAppconfig(w http.ResponseWriter, r *http.Request) {
	log.Info("Starting GetAppconfig api")

	region := r.URL.Query().Get("zone")
	if region == "" {
		log.Error("Zone(Region) not provided")
		http.Error(w, fmt.Sprintf("Zone(Region) not provided"), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	vaultUrl := r.URL.Query().Get("vaultUrl")

	if vaultUrl != "" {
		accountId := r.URL.Query().Get("accountId")
		if accountId == "" {
			log.Error("AccountId not provided")
			http.Error(w, fmt.Sprintf("AccountId not provided"), http.StatusBadRequest)
			return
		}

		vaultToken := r.URL.Query().Get("vaultToken")
		if vaultToken == "" {
			log.Error("Vault token not provided")
			http.Error(w, fmt.Sprintf("Vault token not provided"), http.StatusBadRequest)
			return
		}

		result, err := controller.GetDiscoveredResourceByAccountNo(vaultUrl, vaultToken, accountId, region)
		if err != nil {
			log.Error("Exception in getting cloud config summary: %v", err)
			http.Error(w, fmt.Sprintf("Exception in getting cloud config summary"), http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(result)
	} else {
		accessKey := r.URL.Query().Get("accessKey")
		if accessKey == "" {
			log.Error("AccessKey not provided")
			http.Error(w, fmt.Sprintf("AccessKey not provided"), http.StatusBadRequest)
			return
		}
		secretKey := r.URL.Query().Get("secretKey")
		if secretKey == "" {
			log.Error("SecretKey not provided")
			http.Error(w, fmt.Sprintf("SecretKey not provided"), http.StatusBadRequest)
			return
		}
		crossAccountRoleArn := r.URL.Query().Get("crossAccountRoleArn")
		if crossAccountRoleArn == "" {
			log.Error("CrossAccountRoleArn not provided")
			http.Error(w, fmt.Sprintf("CrossAccountRoleArn not provided"), http.StatusBadRequest)
			return
		}
		externalId := r.URL.Query().Get("externalId")
		if externalId == "" {
			log.Error("ExternalId not provided")
			http.Error(w, fmt.Sprintf("ExternalId not provided"), http.StatusBadRequest)
			return
		}
		result, err := controller.GetDiscoveredResourceByUserCreds(region, accessKey, secretKey, crossAccountRoleArn, externalId)
		if err != nil {
			log.Error("Exception in getting cloud config summary: %v", err)
			http.Error(w, fmt.Sprintf("Exception in getting cloud config summary"), http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(result)
	}

	log.Info("GetAppconfig completed")

}
