package client

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
)

type Client struct {
	Config aws.Config
}

// Client instantiates a new client
func New(profile string, region string) (Client, error) {
	var cfg aws.Config
	var err error
	var client Client

	if profile == "" {
		cfg, err = external.LoadDefaultAWSConfig()
	} else {
		cfg, err = external.LoadDefaultAWSConfig(external.WithSharedConfigProfile(profile))
	}

	if err != nil {
		return client, err
	}

	client.Config = cfg

	return Client{Config: cfg}, nil
}
