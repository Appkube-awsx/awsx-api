package handlers

import (
	"awsx-api/log"
	"awsx-api/models"
	"awsx-api/util"
	"encoding/json"
	"fmt"
	"github.com/Appkube-awsx/appconfig/cmd"
	"io/ioutil"
	"net/http"
)

func GetAppconfigHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Starting GetAppconfigHandler...")

	body, err := ioutil.ReadAll(r.Body)
	var raw models.AwsCredential
	err = json.Unmarshal(body, &raw)
	if err != nil {
		util.CommonError(err)
		return
	}
	if raw.Region == "" && raw.AccessKey == "" && raw.SecretKey == "" && raw.CrossAccountRoleArn == "" {
		fmt.Println("AWS credentials like account accesskey, secretkey, region and crossAccountRoleArn not provided")
		return
	}

	json.NewEncoder(w).Encode(cmd.GetConfig(raw.Region, raw.CrossAccountRoleArn, raw.AccessKey, raw.SecretKey))

	log.Info("GetAppconfigHandler completed")

}
