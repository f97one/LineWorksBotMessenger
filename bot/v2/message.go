package v2

import (
	"bytes"
	"encoding/json"
	"fmt"
	v2 "github.com/f97one/LineWorksBotMessenger/utils/v2"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type LocaleCode int

const (
	// SupportedLocaleJapanese は言語コード "ja_JP" を表す。
	SupportedLocaleJapanese LocaleCode = iota
	// SupportedLocaleEnglish は言語コード "en_US" を表す。
	SupportedLocaleEnglish
	// SupportedLocaleChineseSimplified は言語コード "zh_CN" を表す。
	SupportedLocaleChineseSimplified
	// SupportedLocaleChineseTraditional は言語コード "zh_TW" を表す。
	SupportedLocaleChineseTraditional
	// SupportedLocaleKorean は言語コード "ko_KR" を表す。
	SupportedLocaleKorean
)

// String は LocaleCode l に紐づく言語コードを string として返す。
func (l LocaleCode) String() string {
	switch l {
	case SupportedLocaleJapanese:
		return "ja_JP"
	case SupportedLocaleEnglish:
		return "en_US"
	case SupportedLocaleChineseSimplified:
		return "zh_CN"
	case SupportedLocaleChineseTraditional:
		return "zh_TW"
	case SupportedLocaleKorean:
		return "ko_KR"
	default:
		return ""
	}
}

type TextPayload struct {
	Content Content `json:"content"`
}

type Content struct {
	TalkType  string     `json:"type"`
	Text      string     `json:"text"`
	I18NTexts []I18NText `json:"i18nTexts"`
}

type I18NText struct {
	Language string `json:"language"`
	Text     string `json:"text"`
}

func NewJaTextPayload(text string) TextPayload {
	ja := &I18NText{
		Language: SupportedLocaleJapanese.String(),
		Text:     text,
	}

	i18NTexts := []I18NText{*ja}

	return TextPayload{
		Content: Content{
			TalkType:  "text",
			Text:      text,
			I18NTexts: i18NTexts,
		},
	}
}

func NewI18nText(locale LocaleCode, text string) I18NText {
	return I18NText{
		Language: locale.String(),
		Text:     text,
	}
}

func NewLocalizedTextPayload(i18NTexts []I18NText, primaryLocale LocaleCode) TextPayload {
	// デフォルトを空にしておく
	primaryText := ""
	for _, text := range i18NTexts {
		if text.Language == primaryLocale.String() {
			primaryText = text.Text
		}
	}

	return TextPayload{
		Content: Content{
			TalkType:  "text",
			Text:      primaryText,
			I18NTexts: i18NTexts,
		},
	}
}

func (p *TextPayload) SendToUser(userID string, authToken string, config v2.ConfigV2) error {
	endpoint := createBotEndpointUrl(config.BotId, true, userID)

	return sendMessage(p, endpoint, authToken)
}

func sendMessage(p *TextPayload, endpoint string, authToken string) error {
	b, err := json.Marshal(p)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	client.Timeout = time.Second * 10
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (p *TextPayload) SendToRoom(roomID string, authToken string, config v2.ConfigV2) error {
	endpoint := createBotEndpointUrl(config.BotId, false, roomID)
	return sendMessage(p, endpoint, authToken)
}

func createBotEndpointUrl(botId int, toUser bool, dest string) string {
	target := "channels"
	if toUser {
		target = "users"
	}
	return fmt.Sprintf("https://www.worksapis.com/v1.0/bots/%s/%s/%s/messages", url.PathEscape(strconv.Itoa(botId)), target, url.PathEscape(dest))
}
