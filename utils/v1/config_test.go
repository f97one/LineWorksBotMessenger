package v1

import "testing"

func TestLoad(t *testing.T) {
	conf, err := Load("../../testdata/settings-sample-test.json")
	if err != nil {
		t.Error(err)
	}
	if conf.ApiId != "dev_console_api_id" {
		t.Error("ServerId is not correct")
	}
	if conf.ConsumerKey != "non_redirect_server_api_consumer_key" {
		t.Error("ConsumerKey is not correct")
	}
	if conf.ServerId != "server_list_bot_id" {
		t.Error("ServerId is not correct")
	}
	if conf.BotNo != 9999999 {
		t.Error("BotNo is not correct")
	}

}
