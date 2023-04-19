package handlers

import (
	"awsx-api/authentication"
	"awsx-api/log"
	"encoding/json"
	"fmt"
	"github.com/Appkube-awsx/awsx-lambda/actuator"
	"net/http"
)

func GetLambdas(w http.ResponseWriter, r *http.Request) {
	log.Info("Starting GetAppConfig api")

	raw, done := authentication.Auth(w, r)
	fmt.Println("this")
	if done {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(actuator.LambdaListActuator("", true, "", "", raw.AccessDetails.Region, raw.AccessDetails.AccessKey, raw.AccessDetails.SecretKey, raw.AccessDetails.CrossAccountRoleArn, raw.AccessDetails.ExternalId))

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
