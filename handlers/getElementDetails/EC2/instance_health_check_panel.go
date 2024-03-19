package EC2

import (
	"awsx-api/log"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type InstanceHealthCheckData struct {
	InstanceID        string `json:"instance_id"`
	InstanceType      string `json:"instance_type"`
	AvailabilityZone  string `json:"availability_zone"`
	InstanceStatus    string `json:"instance_status"`
	SystemCheck       string `json:"system_check"`
	InstanceCheck     string `json:"instance_check"`
	Alarm             bool   `json:"alarm"`
	SystemCheckTime   string `json:"system_check_time"`
	InstanceCheckTime string `json:"instance_check_time"`
}

var (
	authCacheInstance       sync.Map
	clientCacheInstance     sync.Map
	authCacheLockInstance   sync.RWMutex
	clientCacheLockInstance sync.RWMutex
)

func GetInstanceHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	region := r.URL.Query().Get("zone")
	instanceID := r.URL.Query().Get("instanceId")
	//startTime := r.URL.Query().Get("startTime")
	//endTime := r.URL.Query().Get("endTime")
	responseType := r.URL.Query().Get("responseType")

	commandParam := model.CommandParam{
		Region: region,
	}

	clientAuth, err := authenticateAndCacheInstance(commandParam)
	if err != nil {
		http.Error(w, fmt.Sprintf("Authentication failed: %s", err), http.StatusInternalServerError)
		return
	}

	if clientAuth != nil {
		ec2Client := awsclient.GetClient(*clientAuth, awsclient.EC2_CLIENT).(*ec2.EC2)
		cloudWatchClient, err := cloudwatchClientCacheInstance(*clientAuth)
		if err != nil {
			http.Error(w, fmt.Sprintf("Cloudwatch client creation/store in cache failed: %s", err), http.StatusInternalServerError)
			return
		}

		resp, err := ec2Client.DescribeInstances(&ec2.DescribeInstancesInput{
			InstanceIds: []*string{&instanceID},
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("Error describing instances: %s", err), http.StatusInternalServerError)
			return
		}

		var instanceHealthCheckList []InstanceHealthCheckData
		for _, reservation := range resp.Reservations {
			for _, instance := range reservation.Instances {
				instanceID := *instance.InstanceId
				instanceType := *instance.InstanceType
				availabilityZone := *instance.Placement.AvailabilityZone
				instanceStatus := *instance.State.Name
				systemCheck, instanceCheck, alarm, systemCheckTime, instanceCheckTime, err := getInstanceHealthCheckStatus(cloudWatchClient, instanceID)
				if err != nil {
					http.Error(w, fmt.Sprintf("Error getting instance health check status: %s", err), http.StatusInternalServerError)
					return
				}
				instanceHealthCheck := InstanceHealthCheckData{
					InstanceID:        instanceID,
					InstanceType:      instanceType,
					AvailabilityZone:  availabilityZone,
					InstanceStatus:    instanceStatus,
					SystemCheck:       systemCheck,
					InstanceCheck:     instanceCheck,
					Alarm:             alarm,
					SystemCheckTime:   systemCheckTime,
					InstanceCheckTime: instanceCheckTime,
				}
				instanceHealthCheckList = append(instanceHealthCheckList, instanceHealthCheck)
			}
		}

		log.Infof("response type: %s", responseType)
		if responseType == "json" {
			err = json.NewEncoder(w).Encode(instanceHealthCheckList)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error encoding JSON response: %s", err), http.StatusInternalServerError)
				return
			}
		} else {
			// Provide valid JSON data to unmarshal
			jsonData := `{"InstanceHealthCheckData": []}` // Example empty JSON object
			var data InstanceHealthCheckData
			err := json.Unmarshal([]byte(jsonData), &data)
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
	}
}

func getInstanceHealthCheckStatus(cloudWatchClient *cloudwatch.CloudWatch, instanceID string) (string, string, bool, string, string, error) {
	// Placeholder logic to retrieve system check, instance check, alarm status, system check time, and instance check time
	return "SystemCheck", "InstanceCheck", false, time.Now().Format(time.RFC3339), time.Now().Format(time.RFC3339), nil
}

func authenticateAndCacheInstance(commandParam model.CommandParam) (*model.Auth, error) {
	cacheKey := commandParam.Region

	authCacheLockInstance.Lock()
	if auth, ok := authCacheInstance.Load(cacheKey); ok {
		log.Infof("client credentials found in cache")
		authCacheLockInstance.Unlock()
		return auth.(*model.Auth), nil
	}

	// If not in cache, perform authentication
	log.Infof("getting client credentials from vault/db")
	_, clientAuth, err := authenticate.DoAuthenticate(commandParam)
	if err != nil {
		return nil, err
	}

	authCacheInstance.Store(cacheKey, clientAuth)
	authCacheLockInstance.Unlock()

	return clientAuth, nil
}

func cloudwatchClientCacheInstance(clientAuth model.Auth) (*cloudwatch.CloudWatch, error) {
	cacheKey := clientAuth.CrossAccountRoleArn

	clientCacheLockInstance.Lock()
	if client, ok := clientCacheInstance.Load(cacheKey); ok {
		log.Infof("cloudwatch client found in cache for given cross acount role: %s", cacheKey)
		clientCacheLockInstance.Unlock()
		return client.(*cloudwatch.CloudWatch), nil
	}

	// If not in cache, create new cloud watch client
	log.Infof("creating new cloudwatch client for given cross acount role: %s", cacheKey)
	cloudWatchClient := awsclient.GetClient(clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)

	clientCacheInstance.Store(cacheKey, cloudWatchClient)
	clientCacheLockInstance.Unlock()

	return cloudWatchClient, nil
}
