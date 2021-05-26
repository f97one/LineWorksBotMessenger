package main

import (
	"flag"
	"fmt"
	"github.com/f97one/LineWorksBotMessenger/utils"
	"log"
	"os"
	"path/filepath"
)

func main() {
	// 引数
	//   -c 設定ファイルのパス
	//   -k 認証キーファイルのパス
	//   -d 宛先のユーザー名
	//   -m しゃべらせるメッセージ
	confFilePath := flag.String("c", "", "configuration file path")
	authKeyPath := flag.String("k", "", "Authorization Key file path")
	destUsername := flag.String("d", "", "Destination username to speak")
	msg := flag.String("m", "", "Messages to speak")

	flag.Parse()

	if confFilePath == nil || len(*confFilePath) == 0 {
		exitOnEmpty("configuration file path must not be empty")
	}

	if authKeyPath == nil || len(*authKeyPath) == 0 {
		exitOnEmpty("Authorization Key file path must not be empty")
	}

	if destUsername == nil || len(*destUsername) == 0 {
		exitOnEmpty("Destination username to speak must not be empty")
	}

	if msg == nil || len(*msg) == 0 {
		exitOnEmpty("Messages to speak must not be empty")
	}

	//log.Printf("得られた値 : confFilePath = %s\n", *confFilePath)
	//log.Printf("得られた値 : authKeyPath = %s\n", *authKeyPath)
	//log.Printf("得られた値 : destUsername = %s\n", *destUsername)
	//log.Printf("得られた値 : msg = %s\n", *msg)

	conf, err := utils.Load(filepath.Clean(*confFilePath))
	if err != nil {
		log.Fatalln(err)
	}

	authToken, err := createAuthToken(conf, *authKeyPath)
	if err != nil {
		log.Fatalln(err)
	}
	accessToken, err := getAccessToken(conf, authToken)
	if err != nil {
		log.Fatalln(err)
	}

	err = sendToUser(accessToken, conf, *destUsername, *msg)
	if err != nil {
		log.Fatalln(err)
	}

	os.Exit(0)
}

func exitOnEmpty(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	flag.Usage()
	os.Exit(2)
}