package EC2

import (
	"awsx-api/cache"
	"awsx-api/log"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/handler/EC2"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
)

type Ec2CpuUtilizationResult struct {
	InstanceType string
	Items        map[time.Time]float64
}

func GetCpuUtilizationPerInstancePanel(w http.ResponseWriter, r *http.Request) {
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

	clientAuth, awsClient, err := cache.GetAwsCredsAndClient(commandParam, awsclient.CLOUDWATCH)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cloudwatch client creation/store in cache failed: %s", err), http.StatusInternalServerError)
		return
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("EC2 client creation/store in cache failed: %s", err), http.StatusInternalServerError)
		return
	}
	ec2Svc := awsclient.GetClient(*clientAuth, awsclient.EC2_CLIENT).(*ec2.EC2)
	cloudwatchClient := awsClient.(*cloudwatch.CloudWatch)
	cmd := &cobra.Command{}
	cmd.PersistentFlags().StringVar(&elementId, "elementId", r.URL.Query().Get("elementId"), "Description of the elementId flag")
	cmd.PersistentFlags().StringVar(&instanceId, "instanceId", r.URL.Query().Get("instanceId"), "Description of the instanceId flag")
	cmd.PersistentFlags().StringVar(&elementType, "elementType", r.URL.Query().Get("elementType"), "Description of the elementType flag")
	cmd.PersistentFlags().StringVar(&startTime, "startTime", r.URL.Query().Get("startTime"), "Description of the startTime flag")
	cmd.PersistentFlags().StringVar(&endTime, "endTime", r.URL.Query().Get("endTime"), "Description of the endTime flag")
	cmd.PersistentFlags().StringVar(&responseType, "responseType", r.URL.Query().Get("responseType"), "responseType flag - json/frame")

	jsonString, _, err := EC2.CpuUtilizationPerInstanceType(cmd, clientAuth, ec2Svc, cloudwatchClient)
	if err != nil {
		log.Infof("Error found in CpuUtilizationPerInstanceType: %v", err)
		var awsErr awserr.Error
		if errors.As(err, &awsErr) && awsErr.Code() == "ExpiredToken" {
			log.Infof("AWS session expired. Resetting connection cache")
			clientAuth, awsClient, err = cache.SetAwsCredsAndClientInCache(commandParam, awsclient.CLOUDWATCH)
			if err != nil {
				http.Error(w, fmt.Sprintf("CloudWatch client re-creation/store in cache failed: %s", err), http.StatusInternalServerError)
				return
			}
			cloudwatchClient = awsClient.(*cloudwatch.CloudWatch)
			jsonString, _, err = EC2.CpuUtilizationPerInstanceType(cmd, clientAuth, ec2Svc, cloudwatchClient)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error getting CPU utilization data: %s", err), http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, fmt.Sprintf("Error getting CPU utilization data: %s", err), http.StatusInternalServerError)
			return
		}
	}

	log.Infof("Response type: %s", responseType)

	if responseType == "frame" {
		log.Infof("Creating response frame")
		// Implement frame response if needed
		return
	}

	log.Infof("Creating response JSON")
	var data []Ec2CpuUtilizationResult
	err = json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Exception: %s", err), http.StatusInternalServerError)
		return
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Exception: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonBytes)
	if err != nil {
		http.Error(w, fmt.Sprintf("Exception: %s", err), http.StatusInternalServerError)
		return
	}
}
