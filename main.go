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
	var confFilePath string
	var authKeyPath string
	var destUsername string
	flag.StringVar(&confFilePath, "c", "", "configuration file path")
	flag.StringVar(&authKeyPath, "k", "", "Authorization Key file path")
	flag.StringVar(&destUsername, "d", "", "Destination username to speak")

	// しゃべらせるメッセージは名前なし引数にする
	if len(flag.Args()) == 0 {
		exitOnEmpty("Messages to speak must not be empty")
	}
	msg := flag.Args()[0]

	flag.Parse()

	if len(confFilePath) == 0 {
		exitOnEmpty("configuration file path must not be empty")
	}

	if len(authKeyPath) == 0 {
		exitOnEmpty("Authorization Key file path must not be empty")
	}

	if len(destUsername) == 0 {
		exitOnEmpty("Destination username to speak must not be empty")
	}

	if len(msg) == 0 {
		exitOnEmpty("Messages to speak must not be empty")
	}

	log.Printf("得られた値 : confFilePath = %s\n", confFilePath)
	log.Printf("得られた値 : authKeyPath = %s\n", authKeyPath)
	log.Printf("得られた値 : destUsername = %s\n", destUsername)
	log.Printf("得られた値 : msg = %s\n", msg)

	conf, err := utils.Load(filepath.Clean(confFilePath))
	if err != nil {
		log.Fatalln(err)
	}

	authToken, err := createAuthToken(conf, authKeyPath)
	if err != nil {
		log.Fatalln(err)
	}
	accessToken, err := getAccessToken(conf, authToken)
	if err != nil {
		log.Fatalln(err)
	}

	err = sendToUser(accessToken, conf, destUsername, msg)
	if err != nil {
		log.Fatalln(err)
	}

	os.Exit(0)
}

func exitOnEmpty(msg string) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", msg)
	fmt.Fprintf(os.Stderr, "Usage: %s [options] messages\n", filepath.Base(os.Args[0]))
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "  messages")
	fmt.Fprintln(os.Stderr, "        messages to make LINE WORKS Bot speak")
	os.Exit(2)
}
