package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/f97one/LineWorksBotMessenger/utils"
	"io/ioutil"
	"net/http"
	"time"
)

type TalkPayload struct {
	AccountId string      `json:"accountId"`
	Content   TalkContent `json:"content"`
}

type TalkContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func NewTalkPayload(accountId string, msg string) TalkPayload {
	return TalkPayload{
		AccountId: accountId,
		Content: TalkContent{
			Type: "text",
			Text: msg,
		},
	}
}

func sendToUser(accessToken string, conf utils.Config, accountId string, msg string) error {
	textEndpoint := fmt.Sprintf("https://apis.worksmobile.com/r/%s/message/v1/bot/%d/message/push", conf.ApiId, conf.BotNo)

	talkPayload := NewTalkPayload(accountId, msg)
	body, err := json.Marshal(talkPayload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, textEndpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Add("content-type", "application/json; charset=UTF-8")
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("consumerKey", conf.ConsumerKey)

	client := &http.Client{}
	client.Timeout = time.Second * 15
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	err = parseStatusError(resp)
	if err != nil {
		return err
	}

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return nil
}
