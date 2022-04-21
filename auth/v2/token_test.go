package v2

import (
	"testing"
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
