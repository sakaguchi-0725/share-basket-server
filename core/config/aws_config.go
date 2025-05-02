package config

type AWS struct {
	Region              string
	CognitoClientID     string
	CognitoClientSecret string
	CognitoUserPoolID   string
}

func newAWSConfig() AWS {
	return AWS{
		Region:              getEnv("AWS_REGION", "ap-northeast-1"),
		CognitoClientID:     getEnv("COGNITO_CLIENT_ID", ""),
		CognitoClientSecret: getEnv("COGNITO_CLIENT_SECRET", ""),
		CognitoUserPoolID:   getEnv("COGNITO_USER_POOL_ID", ""),
	}
}
