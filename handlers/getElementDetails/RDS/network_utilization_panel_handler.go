package RDS

import (
	"awsx-api/log"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/aws/aws-sdk-go/service/cloudwatch"

	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/handler/RDS"
	"github.com/spf13/cobra"
)

var (
	netauthCache       sync.Map
	netclientCache     sync.Map
	netauthCacheLock   sync.RWMutex
	netclientCacheLock sync.RWMutex
)

func GetNetworkUtilizationPanel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	region := r.URL.Query().Get("zone")
	elementId := r.URL.Query().Get("elementId")
	elementApiUrl := r.URL.Query().Get("cmdbApiUrl")
	elementType := r.URL.Query().Get("elementType")

	crossAccountRoleArn := r.URL.Query().Get("crossAccountRoleArn")
	externalId := r.URL.Query().Get("externalId")
	responseType := r.URL.Query().Get("responseType")
	filter := r.URL.Query().Get("filter")
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
	clientAuth, err := netauthenticateAndCache(commandParam)
	if err != nil {
		http.Error(w, fmt.Sprintf("Authentication failed: %s", err), http.StatusInternalServerError)
		return
	}
	cloudwatchClient, err := netcloudwatchClientCache(*clientAuth)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cloudwatch client creation/store in cache failed: %s", err), http.StatusInternalServerError)
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
		jsonString, cloudwatchMetricData, err := RDS.GetRDSNetworkUtilizationPanel(cmd, clientAuth, cloudwatchClient)
		if err != nil {
			http.Error(w, fmt.Sprintf("Exception: %s", err), http.StatusInternalServerError)
			return
		}
		log.Infof("response type :" + responseType)

		if responseType == "frame" {
			log.Infof("creating response frame")
			log.Infof("response type :" + responseType)
			if filter == "InboundTraffic" {
				err = json.NewEncoder(w).Encode(cloudwatchMetricData["InboundTraffic"])
				if err != nil {
					http.Error(w, fmt.Sprintf("Exception: %s ", err), http.StatusInternalServerError)
					return
				}
			} else if filter == "OutboundTraffic" {
				err = json.NewEncoder(w).Encode(cloudwatchMetricData["OutboundTraffic"])
				if err != nil {
					http.Error(w, fmt.Sprintf("Exception: %s ", err), http.StatusInternalServerError)
					return
				}
			} else if filter == "DataTransferred" {
				// Calculate Data Transferred (sum of inbound and outbound)
				if cloudwatchMetricData["InboundTraffic"] != nil && cloudwatchMetricData["OutboundTraffic"] != nil {
					inbound := extractMetricData(cloudwatchMetricData["InboundTraffic"])
					outbound := extractMetricData(cloudwatchMetricData["OutboundTraffic"])

					dataTransferred := make(map[string]float64)
					for timestamp, value := range inbound {
						dataTransferred[timestamp] = value + outbound[timestamp]
					}
					err = json.NewEncoder(w).Encode(dataTransferred)
					if err != nil {
						http.Error(w, fmt.Sprintf("Exception: %s ", err), http.StatusInternalServerError)
						return
					}
				} else {
					// Handle case where one or both metrics are missing
					http.Error(w, "Inbound or Outbound traffic metrics are not available", http.StatusInternalServerError)
					return
				}
			} else {
				err = json.NewEncoder(w).Encode(cloudwatchMetricData)
				if err != nil {
					http.Error(w, fmt.Sprintf("Exception: %s ", err), http.StatusInternalServerError)
					return
				}
			}
		} else {
			log.Infof("creating response json")
			type UsageData struct {
				InboundTraffic  float64 `json:"Network RX"`
				OutboundTraffic float64 `json:"Network TX"`
				DataTransferred float64 `json:"DataTransferred"`
			}
			var data UsageData
			err := json.Unmarshal([]byte(jsonString), &data)
			if err != nil {
				http.Error(w, fmt.Sprintf("Exception: %s", err), http.StatusInternalServerError)
				return
			}

			// Marshal the struct back to JSON
			jsonBytes, err := json.Marshal(data)
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

// Function to extract metric data from GetMetricDataOutput
func extractMetricData(metricData *cloudwatch.GetMetricDataOutput) map[string]float64 {
	data := make(map[string]float64)
	for i := range metricData.MetricDataResults {
		for j := range metricData.MetricDataResults[i].Timestamps {
			timestamp := metricData.MetricDataResults[i].Timestamps[j].String()
			value := *metricData.MetricDataResults[i].Values[j]
			data[timestamp] = value
		}
	}
	return data
}

func netauthenticateAndCache(commandParam model.CommandParam) (*model.Auth, error) {
	cacheKey := commandParam.CloudElementId

	netauthCacheLock.Lock()
	if auth, ok := netauthCache.Load(cacheKey); ok {
		log.Infof("client credentials found in cache")
		netauthCacheLock.Unlock()
		return auth.(*model.Auth), nil
	}

	// If not in cache, perform authentication
	log.Infof("getting client credentials from vault/db")
	_, clientAuth, err := authenticate.DoAuthenticate(commandParam)
	if err != nil {
		return nil, err
	}

	netauthCache.Store(cacheKey, clientAuth)
	netauthCacheLock.Unlock()

	return clientAuth, nil
}

func netcloudwatchClientCache(clientAuth model.Auth) (*cloudwatch.CloudWatch, error) {
	cacheKey := clientAuth.CrossAccountRoleArn

	netclientCacheLock.Lock()
	if client, ok := netclientCache.Load(cacheKey); ok {
		log.Infof("cloudwatch client found in cache for given cross account role: %s", cacheKey)
		netclientCacheLock.Unlock()
		return client.(*cloudwatch.CloudWatch), nil
	}

	// If not in cache, create new cloud watch client
	log.Infof("creating new cloudwatch client for given cross account role: %s", cacheKey)
	cloudWatchClient := awsclient.GetClient(clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)

	netclientCache.Store(cacheKey, cloudWatchClient)
	netclientCacheLock.Unlock()

	return cloudWatchClient, nil
}
