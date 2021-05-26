package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/f97one/LineWorksBotMessenger/utils"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Detail  string `json:"detail"`
	Code    string `json:"code"`
}

func createAuthToken(conf utils.Config, authKeyPath string) (string, error) {
	absAuthKeyPath, err := filepath.Abs(filepath.Clean(authKeyPath))
	if err != nil {
		return "", err
	}

	// 認証キー抽出
	signBytes, err := ioutil.ReadFile(absAuthKeyPath)
	if err != nil {
		return "", err
	}
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return "", err
	}

	// 発行日を現在日時、有効期限はその30分後に設定
	issuedAt := time.Now().Unix()
	expiresAt := time.Now().Add(30 * time.Minute).Unix()

	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = conf.ServerId
	claims["iat"] = issuedAt
	claims["exp"] = expiresAt

	tokenString, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}

	//log.Printf("得られた TokenString = %s\n", tokenString)
	return tokenString, nil
}

func getAccessToken(conf utils.Config, authToken string) (string, error) {
	tokenEndpoint := fmt.Sprintf("https://auth.worksmobile.com/b/%s/server/token", conf.ApiId)

	// body
	//   form#Encode() でURLエンコードがかかるのでそのままでよい
	form := url.Values{}
	form.Add("grant_type", "urn:ietf:params:oauth:grant-type:jwt-bearer")
	form.Add("assertion", authToken)
	body := strings.NewReader(form.Encode())

	req, err := http.NewRequest(http.MethodPost, tokenEndpoint, body)
	if err != nil {
		return "", err
	}
	// content-type指定
	req.Header.Set("content-type", "application/x-www-form-urlencoded; charset=UTF-8")

	client := &http.Client{}
	client.Timeout = time.Second * 30
	resp, err := client.Do(req)
	if err != nil {
		// todo エンドポイントのエラー処理をちゃんと実装する
		return "", err
	}
	defer resp.Body.Close()

	err = parseStatusError(resp)
	if err != nil {
		return "", err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var tokenResp TokenResponse
	err = json.Unmarshal(respBody, &tokenResp)
	if err != nil {
		return "", err
	}
	err = parseStatusError(resp)
	if err != nil {
		return "", err
	}

	return tokenResp.AccessToken, nil
}

func parseStatusError(resp *http.Response) error {
	if resp.StatusCode >= 400 {
		errBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		var errResp ErrorResponse
		err = json.Unmarshal(errBody, &errResp)
		if err != nil {
			return err
		}
		msg := fmt.Sprintf("Error : %s : %s\ndetail : %s\n", errResp.Code, errResp.Message, errResp.Detail)
		return errors.New(msg)
	}
	return nil
}
