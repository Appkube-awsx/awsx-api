package handlers

import (
	"awsx-api/log"
	"encoding/json"
	"fmt"
	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-s3/command/bucketcmd"
	"github.com/Appkube-awsx/awsx-s3/controller"
	"net/http"
)

func GetS3(w http.ResponseWriter, r *http.Request) {
	log.Info("Starting /awsx/s3 api")
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
		result, respErr := controller.S3BucketListController(*clientAuth)
		if respErr != nil {
			log.Error(respErr.Error())
			http.Error(w, fmt.Sprintf("Exception: "+respErr.Error()), http.StatusInternalServerError)
			return
		}
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
		result, respErr := controller.S3BucketListController(*clientAuth)
		if respErr != nil {
			log.Error(respErr.Error())
			http.Error(w, fmt.Sprintf("Exception: "+respErr.Error()), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(result)
	}
	log.Info("/awsx/s3 completed")
}

func GetS3WithTags(w http.ResponseWriter, r *http.Request) {
	log.Info("Starting /awsx/s3/bucket-with-tag api")
	w.Header().Set("Content-Type", "application/json")
	bucketName := r.URL.Query().Get("bucketName")
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
		result, respErr := controller.S3BucketWithTagsController(bucketName, *clientAuth)
		if respErr != nil {
			log.Error(respErr.Error())
			http.Error(w, fmt.Sprintf("Exception: "+respErr.Error()), http.StatusExpectationFailed)
			return
		}
		var bucketObj *bucketcmd.S3Bucket
		unMarshalErr := json.Unmarshal([]byte(result), &bucketObj)
		if unMarshalErr != nil {
			log.Error(unMarshalErr.Error())
			http.Error(w, fmt.Sprintf("Exception: "+unMarshalErr.Error()), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(bucketObj)
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
		result, respErr := controller.S3BucketWithTagsController(bucketName, *clientAuth)
		if respErr != nil {
			log.Error(respErr.Error())
			http.Error(w, fmt.Sprintf("Exception: "+respErr.Error()), http.StatusInternalServerError)
			return
		}

		var bucketObj *bucketcmd.S3Bucket
		unMarshalErr := json.Unmarshal([]byte(result), &bucketObj)
		if unMarshalErr != nil {
			log.Error(unMarshalErr.Error())
			http.Error(w, fmt.Sprintf("Exception: "+unMarshalErr.Error()), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(bucketObj)
	}
	log.Info("/awsx/s3/bucket-with-tag completed")
}
