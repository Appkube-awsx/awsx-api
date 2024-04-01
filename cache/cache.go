package cache

import (
	"awsx-api/log"
	"fmt"
	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/cmdb"
	"github.com/Appkube-awsx/awsx-common/model"
	"strconv"
	"sync"
)

var (
	credentialCache sync.Map
	awsClientCache  sync.Map
	cacheLock       sync.RWMutex
)

func GetLandingZone(commandParam model.CommandParam) (*model.Landingzone, error) {
	log.Infof("getting cloud-element data to do aws connection caching. cloudElementId: " + commandParam.CloudElementId)
	cloudElementResp, err := cmdb.GetCloudElement(commandParam)
	if err != nil {
		return nil, fmt.Errorf("cmdb api failed to get cloud-element response in local caching", err)
	}
	log.Infof("getting landing-zone data to do aws connection caching. landingZoneId: " + strconv.FormatInt(cloudElementResp.LandingzoneId, 10))
	landingZoneResp, err := cmdb.GetLandingZone(commandParam, int(cloudElementResp.LandingzoneId))
	if err != nil {
		return nil, fmt.Errorf("cmdb api failed to get landing-zone response in local caching", err)
	}
	return landingZoneResp, nil
}

func GetAwsCredsAndClient(commandParam model.CommandParam, clientType string) (*model.Auth, interface{}, error) {
	landingZoneResp, err := GetLandingZone(commandParam)
	if err != nil {
		return nil, nil, err
	}

	var awsCredsAuth *model.Auth
	var awsClient interface{}

	cacheLock.Lock()
	if credAuth, ok := credentialCache.Load(landingZoneResp.RoleArn); ok {
		log.Infof("client credentials found in cache")
		awsCredsAuth = credAuth.(*model.Auth)
	} else {
		log.Infof("storing new aws credential reference in cache")
		_, awsCredsAuth, err = authenticate.DoAuthenticate(commandParam)
		if err != nil {
			cacheLock.Unlock()
			return nil, nil, err
		}
		credentialCache.Store(landingZoneResp.RoleArn, awsCredsAuth)
	}

	if awsClientAuth, ok := awsClientCache.Load(landingZoneResp.RoleArn + "$$" + clientType); ok {
		log.Infof("aws client found in cache")
		awsClient = awsClientAuth
	} else {
		log.Infof("storing new client connection reference in cache")
		awsClient = awsclient.GetClient(*awsCredsAuth, clientType)
		awsClientCache.Store(landingZoneResp.RoleArn+"$$"+clientType, awsClient)
	}
	cacheLock.Unlock()

	return awsCredsAuth, awsClient, nil
}

func SetAwsCredsAndClientInCache(commandParam model.CommandParam, clientType string) (*model.Auth, interface{}, error) {
	log.Infof("storing aws credentials and client of a landing-zone in cache")
	cacheLock.Lock()
	landingZoneResp, err := GetLandingZone(commandParam)
	if err != nil {
		cacheLock.Unlock()
		return nil, nil, err
	}
	_, awsCredsAuth, err := authenticate.DoAuthenticate(commandParam)
	if err != nil {
		cacheLock.Unlock()
		return nil, nil, err
	}
	credentialCache.Store(landingZoneResp.RoleArn, awsCredsAuth)

	awsClient := awsclient.GetClient(*awsCredsAuth, clientType)
	awsClientCache.Store(landingZoneResp.RoleArn+"$$"+clientType, awsClient)
	cacheLock.Unlock()
	return awsCredsAuth, awsClient, nil
}
