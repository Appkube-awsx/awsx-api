package models

type AwsCredential struct {
	Region              string `json:"region,omitempty"`
	AccessKey           string `json:"accessKey,omitempty"`
	SecretKey           string `json:"secretKey,omitempty"`
	CrossAccountRoleArn string `json:"crossAccountRoleArn,omitempty"`
	ExternalId          string `json:"externalId,omitempty"`
}

type AccessCredential struct {
	Id            int64         `json:"id,omitempty"`
	CloudType     string        `json:"cloudType,omitempty"`
	AccountId     string        `json:"accountId,omitempty"`
	AccessDetails AwsCredential `json:"accessDetails,omitempty"`
}
