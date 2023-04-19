package handlers

import (
	"awsx-api/authentication"
	"awsx-api/config"
	"awsx-api/log"
	"awsx-api/models"
	"awsx-api/util"
	"encoding/json"
	"fmt"
	"github.com/Appkube-awsx/awsx-cloudelements/cmd"
	"io/ioutil"
	"net/http"
	"strings"
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

	raw, done := authentication.Auth(w, r)

	if done {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cmd.GetConfig(raw.AccessDetails.Region, raw.AccessDetails.CrossAccountRoleArn, raw.AccessDetails.AccessKey, raw.AccessDetails.SecretKey, raw.AccessDetails.ExternalId))

	log.Info("GetAppconfig completed")

}

func CreateCloudElements(w http.ResponseWriter, r *http.Request) {
	log.Info("Starting CloudElements api")

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

	var cloudElement models.CloudElement

	log.Infof("Getting app config details")
	appConfigJson := cmd.GetConfig(raw.AccessDetails.Region, raw.AccessDetails.CrossAccountRoleArn, raw.AccessDetails.AccessKey, raw.AccessDetails.SecretKey, raw.AccessDetails.ExternalId)
	appConfigByteArray, err := json.Marshal(appConfigJson)
	if err != nil {
		util.Error("Appconfig details could not be converted to byte array: ", err)
		http.Error(w, fmt.Sprintf("%s", err), statusCode)
		return
	}

	err = json.Unmarshal(appConfigByteArray, &cloudElement.ViewJson)
	if err != nil {
		util.Error("Byte array could not be marshalled to view json: ", err)
		http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	elementKey := r.URL.Query().Get("elementKey")
	if elementKey == "" {
		log.Infof("Getting default element key: ", vaultUrl)
		elementKey = config.Get().CloudElement.ElementKey
	}
	cloudElement.Name = elementKey
	cloudElement.AccountId = accountId

	cloudElementByteArray, err := json.Marshal(cloudElement)
	if err != nil {
		util.Error("CloudElement could not be converted to json: ", err)
		http.Error(w, fmt.Sprintf("%s", err), statusCode)
		return
	}
	payload := strings.NewReader(string(cloudElementByteArray))

	cloudElementUrl := r.URL.Query().Get("cloudElementUrl")
	if cloudElementUrl == "" {
		log.Infof("Calling default cloud-element API: ", vaultUrl)
		cloudElementUrl = config.Get().CloudElement.Url
	}
	cloudElementResp, statusCode, err := util.HandleHttpRequest("POST", cloudElementUrl, "", payload)
	if err != nil {
		util.Error("Http request failed: ", err)
		http.Error(w, fmt.Sprintf("%s", err), statusCode)
		return
	}
	var rawResp models.CloudElement
	err = json.Unmarshal(cloudElementResp, &rawResp)
	json.NewEncoder(w).Encode(rawResp)

	log.Info("CloudElements completed")

}
