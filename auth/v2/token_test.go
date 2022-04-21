package v2

import (
	"encoding/json"
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

func Test_GeneratesInitialRequestAsJson(t *testing.T) {
	tokenReq := TokenRequest{
		Assertion:    "DUMMY_ASSERTION",
		RefreshToken: "",
		GrantType:    GrantTypeInitial.String(),
		ClientId:     "DUMMY_CLIENT_ID",
		ClientSecret: "DUMMY_CLIENT_SECRET",
		Scope:        "Bot",
	}
	b, err := json.Marshal(tokenReq)
	if err != nil {
		t.Errorf("json.Marshal(%v) = %v", tokenReq, err)
	}
	expected := `{"assertion":"DUMMY_ASSERTION","grant_type":"urn:ietf:params:oauth:grant-type:jwt-bearer","client_id":"DUMMY_CLIENT_ID","client_secret":"DUMMY_CLIENT_SECRET","scope":"Bot"}`
	if string(b) != expected {
		t.Errorf("json.Marshal(%v) = %v, want %v", tokenReq, string(b), expected)
	}
}

func Test_GeneratesRefreshRequestAsJson(t *testing.T) {
	tokenReq := TokenRequest{
		Assertion:    "",
		RefreshToken: "DUMMY_REFRESH_TOKEN",
		GrantType:    GrantTypeRefresh.String(),
		ClientId:     "DUMMY_CLIENT_ID",
		ClientSecret: "DUMMY_CLIENT_SECRET",
		Scope:        "",
	}
	b, err := json.Marshal(tokenReq)
	if err != nil {
		t.Errorf("json.Marshal(%v) = %v", tokenReq, err)
	}
	expected := `{"refresh_token":"DUMMY_REFRESH_TOKEN","grant_type":"refresh_token","client_id":"DUMMY_CLIENT_ID","client_secret":"DUMMY_CLIENT_SECRET"}`
	if string(b) != expected {
		t.Errorf("json.Marshal(%v) = %v, want %v", tokenReq, string(b), expected)
	}
}
