package handlers

import (
	"awsx-api/log"
	"encoding/json"
	"fmt"
	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-lambda/controllers"
	"net/http"
)

func GetLambdas(w http.ResponseWriter, r *http.Request) {
	log.Info("Starting /awsx/lambda api")
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
		result := controllers.AllLambdaListController(*clientAuth)
		json.NewEncoder(w).Encode(result)

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
		result := controllers.AllLambdaListController(*clientAuth)
		json.NewEncoder(w).Encode(result)
	}

	log.Info("/awsx/lambda completed")
}

func GetNumberOfLambdas(w http.ResponseWriter, r *http.Request) {
	//log.Info("Starting GetAppConfig api")
	//
	//raw, done := authentication.Auth(w, r)
	//
	//if done {
	//	return
	//}
	//
	//w.Header().Set("Content-Type", "application/json")
	//lambdas := commands.GetLambdaList(raw.AccessDetails.Region, raw.AccessDetails.CrossAccountRoleArn, raw.AccessDetails.AccessKey, raw.AccessDetails.SecretKey, raw.AccessDetails.ExternalId)
	//
	//json.NewEncoder(w).Encode(len(lambdas.Functions))
	//
	//log.Info("GetAppconfig completed")
}
