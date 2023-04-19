package authentication

import (
	"awsx-api/config"
	"awsx-api/log"
	"awsx-api/models"
	"awsx-api/util"
	"encoding/json"
	"fmt"
	"net/http"
)

func Auth(w http.ResponseWriter, r *http.Request) (models.AccessCredential, bool) {
	accountId := r.URL.Query().Get("accountId")
	if accountId == "" {
		log.Error("AccountId not provided")
		http.Error(w, fmt.Sprintf("AccountId not provided"), http.StatusBadRequest)
		return models.AccessCredential{}, true
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
		return models.AccessCredential{}, true
	}
	var raw models.AccessCredential
	err = json.Unmarshal(respBody, &raw)
	if err != nil {
		util.CommonError(err)
		http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return models.AccessCredential{}, true
	}
	if raw.AccessDetails.Region == "" && raw.AccessDetails.AccessKey == "" && raw.AccessDetails.SecretKey == "" && raw.AccessDetails.CrossAccountRoleArn == "" && raw.AccessDetails.ExternalId == "" {
		fmt.Println("AWS credentials like account accesskey, secretkey, region and crossAccountRoleArn not provided")
		http.Error(w, fmt.Sprintf("AWS credentials like account accesskey, secretkey, region and crossAccountRoleArn not provided"), http.StatusBadRequest)
		return models.AccessCredential{}, true
	}

	fmt.Println(raw)
	return raw, false
}
