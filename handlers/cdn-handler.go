package handlers

import (
	"awsx-api/log"
	"encoding/json"
	"fmt"
	"github.com/Appkube-awsx/awsx-cdn/command/cloudfrontcmd"
	"github.com/Appkube-awsx/awsx-cdn/controller"
	"github.com/Appkube-awsx/awsx-common/authenticate"
	"net/http"
)

func GetCdn(w http.ResponseWriter, r *http.Request) {
	log.Info("Starting /awsx/cdn api")
	w.Header().Set("Content-Type", "application/json")

	region := r.URL.Query().Get("zone")
	vaultUrl := r.URL.Query().Get("vaultUrl")

	if vaultUrl != "" {
		accountId := r.URL.Query().Get("accountId")
		vaultToken := r.URL.Query().Get("vaultToken")
		authFlag, clientAuth, err := authenticate.AuthenticateData(vaultUrl, vaultToken, accountId, region, "", "", "", "")
		if err != nil || !authFlag {
			log.Error(err.Error())
			http.Error(w, fmt.Sprintf("Exception: "+err.Error()), http.StatusInternalServerError)
			return
		}
		result, respErr := controller.GetCdnDistributionConfigWithTags(clientAuth)
		if respErr != nil {
			log.Error(respErr.Error())
			http.Error(w, fmt.Sprintf("Exception: "+respErr.Error()), http.StatusInternalServerError)
			return
		}
		var cdn []cloudfrontcmd.Cdn
		unMarshalErr := json.Unmarshal([]byte(result), &cdn)
		if unMarshalErr != nil {
			log.Error(unMarshalErr.Error())
			http.Error(w, fmt.Sprintf("Exception: "+unMarshalErr.Error()), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(cdn)

	} else {
		accessKey := r.URL.Query().Get("accessKey")
		secretKey := r.URL.Query().Get("secretKey")
		crossAccountRoleArn := r.URL.Query().Get("crossAccountRoleArn")
		externalId := r.URL.Query().Get("externalId")
		authFlag, clientAuth, err := authenticate.AuthenticateData("", "", "", region, accessKey, secretKey, crossAccountRoleArn, externalId)
		if err != nil || !authFlag {
			log.Error(err.Error())
			http.Error(w, fmt.Sprintf("Exception: "+err.Error()), http.StatusInternalServerError)
			return
		}
		result, respErr := controller.GetCdnDistributionConfigWithTags(clientAuth)
		if respErr != nil {
			log.Error(respErr.Error())
			http.Error(w, fmt.Sprintf("Exception: "+respErr.Error()), http.StatusInternalServerError)
			return
		}
		var cdn []cloudfrontcmd.Cdn
		//all := []Cdn{}
		unMarshalErr := json.Unmarshal([]byte(result), &cdn)
		if unMarshalErr != nil {
			log.Error(unMarshalErr.Error())
			http.Error(w, fmt.Sprintf("Exception: "+unMarshalErr.Error()), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(cdn)
	}

	log.Info("/awsx/cdn completed")

}
