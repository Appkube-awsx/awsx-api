package handlers

import (
	"awsx-api/config"
	"awsx-api/log"
	"awsx-api/models"
	"awsx-api/util"
	"encoding/json"
	"fmt"
	"github.com/Appkube-awsx/awsx-cloudelements/cmd"
	"io/ioutil"
	"net/http"
)

func GetAppconfigByAccessId(w http.ResponseWriter, r *http.Request) {
	log.Info("Starting GetAppconfigByAccessId api")

	body, err := ioutil.ReadAll(r.Body)
	var raw models.AwsCredential
	err = json.Unmarshal(body, &raw)
	if err != nil {
		util.CommonError(err)
		return
	}
	if raw.Region == "" && raw.AccessKey == "" && raw.SecretKey == "" && raw.CrossAccountRoleArn == "" && raw.ExternalId == "" {
		fmt.Println("AWS credentials like account accesskey, secretkey, region and crossAccountRoleArn not provided")
		return
	}

	json.NewEncoder(w).Encode(cmd.GetConfig(raw.Region, raw.CrossAccountRoleArn, raw.AccessKey, raw.SecretKey, raw.ExternalId))

	log.Info("GetAppconfigByAccessId completed")

}

func GetAppconfig(w http.ResponseWriter, r *http.Request) {
	log.Info("Starting GetAppconfig api")

	accountId := r.URL.Query().Get("accountId")
	if accountId == "" {
		log.Error("AccountId not provided")
		http.Error(w, fmt.Sprintf("AccountId not provided"), http.StatusBadRequest)
		return
	}

	vaultUrl := r.URL.Query().Get("vaultUrl")
	if vaultUrl == "" {
		log.Infof("Calling default vault API: ", vaultUrl)
		vaultUrl = config.Get().Vault.Url
	}
	vaultUrl = vaultUrl + "?accountNo=" + accountId
	respBody, statusCode, err := util.HandleHttpRequest("GET", vaultUrl, "", nil)
	if err != nil {
		util.Error("Http request failed: ", err)
		http.Error(w, fmt.Sprintf("%s", err), statusCode)
		return
	}
	var raw models.AccessCredential
	err = json.Unmarshal(respBody, &raw)
	if err != nil {
		util.CommonError(err)
		http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}
	if raw.AccessDetails.Region == "" && raw.AccessDetails.AccessKey == "" && raw.AccessDetails.SecretKey == "" && raw.AccessDetails.CrossAccountRoleArn == "" && raw.AccessDetails.ExternalId == "" {
		fmt.Println("AWS credentials like account accesskey, secretkey, region and crossAccountRoleArn not provided")
		http.Error(w, fmt.Sprintf("AWS credentials like account accesskey, secretkey, region and crossAccountRoleArn not provided"), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(cmd.GetConfig(raw.AccessDetails.Region, raw.AccessDetails.CrossAccountRoleArn, raw.AccessDetails.AccessKey, raw.AccessDetails.SecretKey, raw.AccessDetails.ExternalId))

	log.Info("GetAppconfig completed")

}
