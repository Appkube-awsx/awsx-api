package handlers

import (
	"awsx-api/log"
	"encoding/json"
	"fmt"
	"github.com/Appkube-awsx/awsx-lambda/authenticater"
	"github.com/Appkube-awsx/awsx-lambda/controllers"
	"net/http"
)

func GetLambdas(w http.ResponseWriter, r *http.Request) {
	log.Info("Starting aws lambda api")

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
		authFlag, clientAuth := authenticater.ApiAuth(vaultUrl, vaultToken, accountId, region, "", "", "", "")

		if !authFlag {
			http.Error(w, fmt.Sprintf("Exception in getting aws lambda by vault url"), http.StatusInternalServerError)
		}
		result := controllers.AllLambdaListController(*clientAuth)
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
		authFlag, clientAuth := authenticater.ApiAuth("", "", "", region, accessKey, secretKey, crossAccountRoleArn, externalId)

		if !authFlag {
			http.Error(w, fmt.Sprintf("Exception in getting aws lambda with access/secret key"), http.StatusInternalServerError)
		}
		result := controllers.AllLambdaListController(*clientAuth)
		json.NewEncoder(w).Encode(result)
	}

	log.Info("Get lambdas completed")
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
