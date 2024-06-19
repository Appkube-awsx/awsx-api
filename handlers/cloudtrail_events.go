package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudtrail"
	"net/http"
)

func GetAwsEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	instanceId := r.URL.Query().Get("instanceId")
	landingZoneId := r.URL.Query().Get("landingZoneId")
	commandParam := model.CommandParam{
		LandingZoneId: landingZoneId,
	}
	_, clientAuth, err := authenticate.DoAuthenticate(commandParam)
	if err != nil {
		http.Error(w, fmt.Sprintf("error in getting aws credentials: %s", err), http.StatusInternalServerError)
		return
	}
	client := awsclient.GetClient(*clientAuth, awsclient.CLOUDTRAIL_CLIENT).(*cloudtrail.CloudTrail)
	input := &cloudtrail.LookupEventsInput{
		LookupAttributes: []*cloudtrail.LookupAttribute{
			{
				AttributeKey:   aws.String("ResourceName"),
				AttributeValue: aws.String(instanceId),
			},
		},
	}
	// Call LookupEvents with the input.
	result, err := client.LookupEvents(input)
	if err != nil {
		fmt.Println("Error looking up events:", err)
		return
	}
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, fmt.Sprintf("errror in json encoding %s ", err), http.StatusInternalServerError)
		return
	}
}
