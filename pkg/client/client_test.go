package client

import (
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
	want := "SharedConfigCredentials: " + UserHomeDir() + "/.aws/credentials"
	got, err := Client()

	if  err != nil {
		t.Errorf("Error getting client: %+v", err)
	}

	if got != want {
		t.Errorf("Client() = %q, want %q", got, want)
	}
}
