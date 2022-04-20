package v1

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	ApiId       string `json:"api_id"`
	ConsumerKey string `json:"consumer_key"`
	ServerId    string `json:"server_id"`
	BotNo       int    `json:"bot_no"`
}

var config *Config

func Load(path string) (Config, error) {
	absPath, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		return Config{}, err
	}
	f, err := os.Open(absPath)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return Config{}, err
	}

	var ret Config
	err = json.Unmarshal(b, &ret)
	if err != nil {
		return Config{}, err
	}

	//log.Println(ret)

	config = &ret
	return ret, nil
}

func GetConfig() Config {
	return *config
}
