package client

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

// Client instantiates a new client
func Client() (string, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile: "admin",
	})

	if err != nil {
		return "", err
	}

	ident, err := sts.New(sess).GetCallerIdentity(&sts.GetCallerIdentityInput{})

	if err != nil {
		return "", err
	}

	return ident.String(), nil
}
