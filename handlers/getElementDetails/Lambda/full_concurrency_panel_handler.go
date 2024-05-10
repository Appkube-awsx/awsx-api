package Lambda

import (
	"awsx-api/cache"
	"awsx-api/log"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Appkube-awsx/awsx-getelementdetails/handler/Lambda"

	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/spf13/cobra"
)

func GetFullConcurrencyPanel(w http.ResponseWriter, r *http.Request) {
	fmt.Print("HERE HERE")
	w.Header().Set("Content-Type", "application/json")

	region := r.URL.Query().Get("zone")
	elementId := r.URL.Query().Get("elementId")
	elementApiUrl := r.URL.Query().Get("cmdbApiUrl")
	elementType := r.URL.Query().Get("elementType")

	crossAccountRoleArn := r.URL.Query().Get("crossAccountRoleArn")
	externalId := r.URL.Query().Get("externalId")
	responseType := r.URL.Query().Get("responseType")
	instanceId := r.URL.Query().Get("instanceId")
	startTime := r.URL.Query().Get("startTime")
	endTime := r.URL.Query().Get("endTime")

	commandParam := model.CommandParam{}

	if elementId != "" {
		commandParam.CloudElementId = elementId
		commandParam.CloudElementApiUrl = elementApiUrl
		commandParam.Region = region
	} else {
		commandParam.CrossAccountRoleArn = crossAccountRoleArn
		commandParam.ExternalId = externalId
		commandParam.Region = region
	}

	clientAuth, awsClient, err := cache.GetAwsCredsAndClient(commandParam, awsclient.LAMBDA_CLIENT)
	if err != nil {
		http.Error(w, fmt.Sprintf("Authentication failed: %s", err), http.StatusInternalServerError)
		return
	}

	lambdaClient := awsClient.(*lambda.Lambda)
	if err != nil {
		http.Error(w, fmt.Sprintf("lambdaClient client creation/store in cache failed: %s", err), http.StatusInternalServerError)
		return
	}

	if clientAuth != nil {
		cmd := &cobra.Command{}
		cmd.PersistentFlags().StringVar(&elementId, "elementId", r.URL.Query().Get("elementId"), "Description of the elementId flag")
		cmd.PersistentFlags().StringVar(&instanceId, "instanceId", r.URL.Query().Get("instanceId"), "Description of the instanceId flag")
		cmd.PersistentFlags().StringVar(&elementType, "elementType", r.URL.Query().Get("elementType"), "Description of the elementType flag")
		cmd.PersistentFlags().StringVar(&startTime, "startTime", r.URL.Query().Get("startTime"), "Description of the startTime flag")
		cmd.PersistentFlags().StringVar(&endTime, "endTime", r.URL.Query().Get("endTime"), "Description of the endTime flag")
		cmd.PersistentFlags().StringVar(&responseType, "responseType", r.URL.Query().Get("responseType"), "responseType flag - json/frame")

		_, cloudwatchMetricData, err := Lambda.GetLambdaFullConcurrencyData(cmd, clientAuth, lambdaClient)
		if err != nil {
			log.Infof("error found in GetColdStartDurationPanel: ", err)
			var awsErr awserr.Error
			if errors.As(err, &awsErr) {
				if awsErr.Code() == "ExpiredToken" {
					log.Infof("aws session expired. resetting connection cache")
					clientAuth, awsClient, err = cache.SetAwsCredsAndClientInCache(commandParam, awsclient.LAMBDA_CLIENT)
				}
			}
		}
		log.Infof("response type :" + responseType)

		if responseType == "frame" {
			err = json.NewEncoder(w).Encode(cloudwatchMetricData)
			if err != nil {
				http.Error(w, fmt.Sprintf("Exception: %s", err), http.StatusInternalServerError)
				return
			}
		} else {
			jsonBytes, err := json.Marshal(cloudwatchMetricData)
			if err != nil {
				http.Error(w, fmt.Sprintf("Exception: %s", err), http.StatusInternalServerError)
				return
			}

			// Set Content-Type header and write the JSON response
			w.Header().Set("Content-Type", "application/json")
			_, err = w.Write(jsonBytes)
			if err != nil {
				http.Error(w, fmt.Sprintf("Exception: %s", err), http.StatusInternalServerError)
				return
			}
		}
	}
}
