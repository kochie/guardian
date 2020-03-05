package config

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

var sess *session.Session

func GetDefaultSession() *session.Session {
	return sess
}

func init() {
	sess = session.Must(session.NewSession(&aws.Config{
		Region: aws.String(DefaultRegion),
	}))
}
