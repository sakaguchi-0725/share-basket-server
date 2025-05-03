package config

import "share-basket-server/core/util"

type AWS struct {
	Region              string
	CognitoClientID     string
	CognitoClientSecret string
	CognitoUserPoolID   string
}

func newAWSConfig() AWS {
	return AWS{
		Region:              util.GetEnv("AWS_REGION", "ap-northeast-1"),
		CognitoClientID:     util.GetEnv("COGNITO_CLIENT_ID", ""),
		CognitoClientSecret: util.GetEnv("COGNITO_CLIENT_SECRET", ""),
		CognitoUserPoolID:   util.GetEnv("COGNITO_USER_POOL_ID", ""),
	}
}
