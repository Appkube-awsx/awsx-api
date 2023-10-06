package handlers

import (
	"awsx-api/log"
	"encoding/json"
	"fmt"
	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/client"
	"github.com/Appkube-awsx/awsx-lambda/controllers"
	"github.com/aws/aws-sdk-go/service/lambda"
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
		response := processResponse(w, err, authFlag, clientAuth)
		json.NewEncoder(w).Encode(response)
	} else {
		accessKey := r.URL.Query().Get("accessKey")
		secretKey := r.URL.Query().Get("secretKey")
		crossAccountRoleArn := r.URL.Query().Get("crossAccountRoleArn")
		externalId := r.URL.Query().Get("externalId")
		authFlag, clientAuth, err := authenticate.AuthenticateData("", "", "", region, accessKey, secretKey, crossAccountRoleArn, externalId)
		response := processResponse(w, err, authFlag, clientAuth)
		json.NewEncoder(w).Encode(response)
	}

	log.Info("/awsx/lambda completed")
}

func processResponse(w http.ResponseWriter, err error, authFlag bool, clientAuth *client.Auth) []*lambda.FunctionConfiguration {
	if err != nil || !authFlag {
		log.Error(err.Error())
		http.Error(w, fmt.Sprintf("Exception: "+err.Error()), http.StatusInternalServerError)
		return nil
	}
	result := controllers.AllLambdaListController(*clientAuth)
	return result
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

func GetLambdaWithTags(w http.ResponseWriter, r *http.Request) {
	log.Info("Starting /awsx/lambda/tag/function api")
	w.Header().Set("Content-Type", "application/json")
	functionName := r.URL.Query().Get("functionName")
	if functionName == "" {
		log.Error("Lambda function name not provided")
		http.Error(w, fmt.Sprintf("Lambda function name not provided"), http.StatusBadRequest)
		return
	}
	region := r.URL.Query().Get("zone")
	vaultUrl := r.URL.Query().Get("vaultUrl")
	if vaultUrl != "" {
		accountId := r.URL.Query().Get("accountId")
		vaultToken := r.URL.Query().Get("vaultToken")
		authFlag, clientAuth, err := authenticate.AuthenticateData(vaultUrl, vaultToken, accountId, region, "", "", "", "")
		lambdaObj, isError := processResponseWithTags(w, err, authFlag, clientAuth, functionName)
		if isError {
			return
		}
		json.NewEncoder(w).Encode(lambdaObj)
	} else {
		accessKey := r.URL.Query().Get("accessKey")
		secretKey := r.URL.Query().Get("secretKey")
		crossAccountRoleArn := r.URL.Query().Get("crossAccountRoleArn")
		externalId := r.URL.Query().Get("externalId")
		authFlag, clientAuth, err := authenticate.AuthenticateData("", "", "", region, accessKey, secretKey, crossAccountRoleArn, externalId)
		lambdaObj, isError := processResponseWithTags(w, err, authFlag, clientAuth, functionName)
		if isError {
			return
		}
		json.NewEncoder(w).Encode(lambdaObj)
	}

	log.Info("/awsx/lambda/tag/function completed")
}

func processResponseWithTags(w http.ResponseWriter, err error, authFlag bool, clientAuth *client.Auth, functionName string) (*controllers.LambdaObj, bool) {
	if err != nil || !authFlag {
		log.Error(err.Error())
		http.Error(w, fmt.Sprintf("Exception: "+err.Error()), http.StatusInternalServerError)
		return nil, true
	}

	result, respErr := controllers.LambdaFunctionWithTagsController(functionName, *clientAuth)
	if respErr != nil {
		log.Error(respErr.Error())
		http.Error(w, fmt.Sprintf("Exception: "+respErr.Error()), http.StatusInternalServerError)
		return nil, true
	}
	var lambdaObj *controllers.LambdaObj
	unMarshalErr := json.Unmarshal([]byte(result), &lambdaObj)
	if unMarshalErr != nil {
		log.Error(unMarshalErr.Error())
		http.Error(w, fmt.Sprintf("Exception: "+unMarshalErr.Error()), http.StatusInternalServerError)
		return nil, true
	}
	return lambdaObj, false
}
