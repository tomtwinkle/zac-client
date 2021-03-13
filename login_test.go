package zacclient

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestZACClient_Login(t *testing.T) {
	if ci := os.Getenv("CI"); ci != "" {
		t.Log("Running on CI")
		t.SkipNow()
	}
	if err := godotenv.Load(".env"); err != nil {
		t.Error(err)
		t.FailNow()
	}
	var config ZACConfig
	if err := envconfig.Process("", &config); err != nil {
		t.Error(err)
		t.FailNow()
	}

	c, err := NewClient(&config, WithDebug())
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	t.Run("login success", func(t *testing.T) {
		err := c.Login()
		if !assert.NoError(t, err) {
			t.FailNow()
		}
		logged := c.IsLoggedIn()
		assert.True(t, logged)
	})
}
