package v2

import "testing"

func TestLoad(t *testing.T) {
	conf, err := Load("../../testdata/v2settings-sample-test.json")
	if err != nil {
		t.Error(err)
	}
	if conf.ClientId != "123456789" {
		t.Error("ClientId is not correct")
	}
	if conf.ClientSecret != "abcdefghijklmnopqrstuvwxyz" {
		t.Error("ClientSecret is not correct")
	}
	if conf.ServiceAccount != "xxxxx.serviceaccount@example-company" {
		t.Error("ServiceAccount is not correct")
	}
	if conf.BotId != 9999999 {
		t.Error("BotId is not correct")
	}
}
