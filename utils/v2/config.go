package v2

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

// ConfigV2 :API v2向けの設定
type ConfigV2 struct {
	// client id (required)
	ClientId string `json:"client_id"`
	// client secret (required)
	ClientSecret string `json:"client_secret"`
	// API service account (required)
	ServiceAccount string `json:"service_account"`
	// bot id to make message post (required)
	BotId int `json:"bot_id"`
}

var config *ConfigV2

// Load : ファイルシステム上の path から設定ファイルを読み込む
func Load(path string) (ConfigV2, error) {
	absPath, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		return ConfigV2{}, err
	}
	f, err := os.Open(absPath)
	if err != nil {
		return ConfigV2{}, err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return ConfigV2{}, err
	}

	var ret ConfigV2
	err = json.Unmarshal(b, &ret)
	if err != nil {
		return ConfigV2{}, err
	}

	config = &ret
	return ret, nil
}

// GetConfig : 読み込み済み設定を返す
func GetConfig() ConfigV2 {
	return *config
}
