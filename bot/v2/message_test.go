package v2

import (
	"encoding/json"
	"fmt"
	v2 "github.com/f97one/LineWorksBotMessenger/utils/v2"
	"github.com/jarcoal/httpmock"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func TestNewJaTextPayload(t *testing.T) {
	payload := NewJaTextPayload("メッセージ")
	b, err := json.Marshal(payload)
	if err != nil {
		t.Errorf("%+v", err)
	}
	expect := `{"content":{"type":"text","text":"メッセージ","i18nTexts":[{"language":"ja_JP","text":"メッセージ"}]}}`
	if string(b) != expect {
		t.Errorf("expect : %s, but got : %s", expect, string(b))
	}
}

func TestNewLocalizedTextPayload(t *testing.T) {
	i18NTexts := []I18NText{
		NewI18nText(SupportedLocaleJapanese, "こんにちは"),
		NewI18nText(SupportedLocaleEnglish, "hello"),
		NewI18nText(SupportedLocaleChineseSimplified, "你好"),
		NewI18nText(SupportedLocaleChineseTraditional, "你好"),
		NewI18nText(SupportedLocaleKorean, "안녕하세요"),
	}

	payload := NewLocalizedTextPayload(i18NTexts, SupportedLocaleEnglish)
	b, err := json.Marshal(payload)
	if err != nil {
		t.Errorf("%+v", err)
	}
	expect := `{"content":{"type":"text","text":"hello","i18nTexts":[{"language":"ja_JP","text":"こんにちは"},{"language":"en_US","text":"hello"},{"language":"zh_CN","text":"你好"},{"language":"zh_TW","text":"你好"},{"language":"ko_KR","text":"안녕하세요"}]}}`
	if string(b) != expect {
		t.Errorf("expect : %s, but got : %s", expect, string(b))
	}
}

func TestTextPayload_SendToUser_Valid(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// JWT
	mockJwt := "DUMMY_JWT"

	// リクエスト
	payload := NewJaTextPayload("メッセージ")
	b, err := json.Marshal(payload)
	if err != nil {
		t.Errorf("%+v", err)
	}
	expectedReq := string(b)

	// レスポンスは、単に 201を返すだけなので何も設定しない

	// endpoint
	botId := "12345678"
	mockUserId := "DUMMY_USER_ID@example"
	endpoint := fmt.Sprintf("https://www.worksapis.com/v1.0/bots/%s/users/%s/messages", url.PathEscape(botId), url.PathEscape(mockUserId))

	actualReq := ""
	httpmock.RegisterResponder(
		http.MethodPost,
		endpoint,
		func(req *http.Request) (*http.Response, error) {
			b, err := ioutil.ReadAll(req.Body)
			if err != nil {
				t.Fatal(err)
			}
			actualReq = string(b)

			return httpmock.NewStringResponse(http.StatusCreated, ""), nil
		},
	)

	err = payload.SendToUser(mockUserId, mockJwt, conf())
	if err != nil {
		t.Errorf("%+v", err)
	}

	if actualReq != expectedReq {
		t.Errorf("expect : %s, but got : %s", expectedReq, actualReq)
	}
}

func TestTextPayload_SendToRoom_Valid(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// JWT
	mockJwt := "DUMMY_JWT"

	// リクエスト
	payload := NewJaTextPayload("メッセージ")
	b, err := json.Marshal(payload)
	if err != nil {
		t.Errorf("%+v", err)
	}
	expectedReq := string(b)

	// レスポンスは、単に 201を返すだけなので何も設定しない

	// endpoint
	botId := "12345678"
	mockRoomId := "98765432"
	endpoint := fmt.Sprintf("https://www.worksapis.com/v1.0/bots/%s/channels/%s/messages", url.PathEscape(botId), url.PathEscape(mockRoomId))

	actualReq := ""
	httpmock.RegisterResponder(
		http.MethodPost,
		endpoint,
		func(req *http.Request) (*http.Response, error) {
			b, err := ioutil.ReadAll(req.Body)
			if err != nil {
				t.Fatal(err)
			}
			actualReq = string(b)

			return httpmock.NewStringResponse(http.StatusCreated, ""), nil
		},
	)

	err = payload.SendToRoom(mockRoomId, mockJwt, conf())
	if err != nil {
		t.Errorf("%+v", err)
	}

	if actualReq != expectedReq {
		t.Errorf("expect : %s, but got : %s", expectedReq, actualReq)
	}
}

func conf() v2.ConfigV2 {
	return v2.ConfigV2{
		ClientId:       "DUMMY_CLIENT_ID",
		ClientSecret:   "dummmy_client_secret",
		ServiceAccount: "xxxxx.serviceaccount@example",
		BotId:          12345678,
	}
}
