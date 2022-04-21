package v2

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestGrantType_String(t *testing.T) {
	tests := []struct {
		name string
		gt   GrantType
		want string
	}{
		{
			name: "initial",
			gt:   GrantTypeInitial,
			want: "urn:ietf:params:oauth:grant-type:jwt-bearer",
		},
		{
			name: "refresh",
			gt:   GrantTypeRefresh,
			want: "refresh_token",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.gt.String(); got != tt.want {
				t.Errorf("GrantType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenRequest_ToForm_Initial(t *testing.T) {
	initReq := &TokenRequest{
		Assertion:    "dummy_assertion",
		RefreshToken: "",
		GrantType:    GrantTypeInitial.String(),
		ClientId:     "dummy_client_id",
		ClientSecret: "dummy_client_secret",
		Scope:        "Bot,Bot.read",
	}

	initValues := initReq.ToForm()

	if initValues.Get("assertion") != "dummy_assertion" {
		t.Errorf("TokenRequest.ToForm() = %v, want %v", initValues.Get("assertion"), "dummy_assertion")
	}
	if initValues.Get("grant_type") != "urn:ietf:params:oauth:grant-type:jwt-bearer" {
		t.Errorf("TokenRequest.ToForm() = %v, want %v", initValues.Get("grant_type"), "urn:ietf:params:oauth:grant-type:jwt-bearer")
	}
	if initValues.Get("client_id") != "dummy_client_id" {
		t.Errorf("TokenRequest.ToForm() = %v, want %v", initValues.Get("client_id"), "dummy_client_id")
	}
	if initValues.Get("client_secret") != "dummy_client_secret" {
		t.Errorf("TokenRequest.ToForm() = %v, want %v", initValues.Get("client_secret"), "dummy_client_secret")
	}
	if initValues.Get("scope") != "Bot,Bot.read" {
		t.Errorf("TokenRequest.ToForm() = %v, want %v", initValues.Get("scope"), "Bot,Bot.read")
	}
	if initValues.Get("refresh_token") != "" {
		t.Errorf("TokenRequest.ToForm() = %v, want %v", initValues.Get("refresh_token"), "")
	}
}

func TestTokenRequest_ToForm_Refresh(t *testing.T) {
	refreshReq := &TokenRequest{
		Assertion:    "",
		RefreshToken: "dummy_refresh_token",
		GrantType:    GrantTypeRefresh.String(),
		ClientId:     "dummy_client_id",
		ClientSecret: "dummy_client_secret",
		Scope:        "",
	}

	refreshValues := refreshReq.ToForm()

	if refreshValues.Get("refresh_token") != "dummy_refresh_token" {
		t.Errorf("TokenRequest.ToForm() = %v, want %v", refreshValues.Get("refresh_token"), "dummy_refresh_token")
	}
	if refreshValues.Get("grant_type") != "refresh_token" {
		t.Errorf("TokenRequest.ToForm() = %v, want %v", refreshValues.Get("grant_type"), "refresh_token")
	}
	if refreshValues.Get("client_id") != "dummy_client_id" {
		t.Errorf("TokenRequest.ToForm() = %v, want %v", refreshValues.Get("client_id"), "dummy_client_id")
	}
	if refreshValues.Get("client_secret") != "dummy_client_secret" {
		t.Errorf("TokenRequest.ToForm() = %v, want %v", refreshValues.Get("client_secret"), "dummy_client_secret")
	}
	if refreshValues.Get("scope") != "" {
		t.Errorf("TokenRequest.ToForm() = %v, want %v", refreshValues.Get("scope"), "")
	}
	if refreshValues.Get("assertion") != "" {
		t.Errorf("TokenRequest.ToForm() = %v, want %v", refreshValues.Get("assertion"), "")
	}
}

func TestTokenRequest_GetAccessToken_Valid(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// mock response
	mockResp := &TokenResponse{
		AccessToken:  "dummy_access_token",
		RefreshToken: "dummy_refresh_token",
		TokenType:    "Bearer",
		ExpiresIn:    86400,
		Scope:        "Bot",
	}
	// リクエスト検証変数
	var reqForm url.Values

	httpmock.RegisterResponder(http.MethodPost, "https://auth.worksmobile.com/oauth2/v2.0/token",
		func(req *http.Request) (*http.Response, error) {
			if err := req.ParseForm(); err != nil {
				t.Fatal(err)
			}
			reqForm = req.Form

			return httpmock.NewJsonResponse(200, mockResp)
		})

	req := &TokenRequest{
		Assertion:    "dummy_assertion",
		RefreshToken: "",
		GrantType:    GrantTypeInitial.String(),
		ClientId:     "dummy_client_id",
		ClientSecret: "dummy_client_secret",
		Scope:        "Bot,Bot.read",
	}

	response, err := req.GetAccessToken()

	// リクエストパラメータが正しいことを確認
	if reqForm.Get("assertion") != "dummy_assertion" {
		t.Errorf("TokenRequest.GetAccessToken() = %v, want %v", reqForm.Get("assertion"), "dummy_assertion")
	}
	if reqForm.Get("grant_type") != GrantTypeInitial.String() {
		t.Errorf("TokenRequest.GetAccessToken() = %v, want %v", reqForm.Get("grant_type"), "assertion")
	}
	if reqForm.Get("client_id") != "dummy_client_id" {
		t.Errorf("TokenRequest.GetAccessToken() = %v, want %v", reqForm.Get("client_id"), "dummy_client_id")
	}
	if reqForm.Get("client_secret") != "dummy_client_secret" {
		t.Errorf("TokenRequest.GetAccessToken() = %v, want %v", reqForm.Get("client_secret"), "dummy_client_secret")
	}
	if reqForm.Get("scope") != "Bot,Bot.read" {
		t.Errorf("TokenRequest.GetAccessToken() = %v, want %v", reqForm.Get("scope"), "Bot,Bot.read")
	}
	if reqForm.Get("refresh_token") != "" {
		t.Errorf("TokenRequest.GetAccessToken() = %v, want %v", reqForm.Get("refresh_token"), "")
	}

	// レスポンスが正しいことを確認
	if err != nil {
		t.Errorf("TokenRequest.GetAccessToken() error = %v", err)
	}
	if response.AccessToken != "dummy_access_token" {
		t.Errorf("TokenRequest.GetAccessToken() = %v, want %v", response.AccessToken, "dummy_access_token")
	}
	if response.RefreshToken != "dummy_refresh_token" {
		t.Errorf("TokenRequest.GetAccessToken() = %v, want %v", response.RefreshToken, "dummy_refresh_token")
	}
	if response.TokenType != "Bearer" {
		t.Errorf("TokenRequest.GetAccessToken() = %v, want %v", response.TokenType, "Bearer")
	}
	if response.ExpiresIn != 86400 {
		t.Errorf("TokenRequest.GetAccessToken() = %v, want %v", response.ExpiresIn, 86400)
	}
	if response.Scope != "Bot" {
		t.Errorf("TokenRequest.GetAccessToken() = %v, want %v", response.Scope, "Bot,Bot.read")
	}
}
