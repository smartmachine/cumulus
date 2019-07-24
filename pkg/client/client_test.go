package client

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"os"
	"runtime"
	"testing"
)

func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

func TestClient(t *testing.T) {
	client,err := New("", "")
	if err != nil {
		t.Errorf("unable to obtain client: %+v", err)
	}
	svc := sts.New(client.Config)
	ident := svc.GetCallerIdentityRequest(&sts.GetCallerIdentityInput{})
	resp, err := ident.Send(context.Background())
	if err != nil {
		t.Errorf("unable to get caller identity: %+v", err)
	}
	t.Log(resp)
}

func TestClientWithProfile(t *testing.T)  {
	client,err := New("admin", "")
	if err != nil {
		t.Errorf("unable to obtain client: %+v", err)
	}
	svc := sts.New(client.Config)
	ident := svc.GetCallerIdentityRequest(&sts.GetCallerIdentityInput{})
	resp, err := ident.Send(context.Background())
	if err != nil {
		t.Errorf("unable to get caller identity: %+v", err)
	}
	t.Log(resp)
}
