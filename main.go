package main

import (
	"flag"
	"fmt"
	authV2 "github.com/f97one/LineWorksBotMessenger/auth/v2"
	botV2 "github.com/f97one/LineWorksBotMessenger/bot/v2"
	utilsV2 "github.com/f97one/LineWorksBotMessenger/utils/v2"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"syscall"
)

func main() {
	// 引数
	//   -c 設定ファイルのパス
	//   -k 認証キーファイルのパス
	//   -d 宛先のユーザー名
	//   -r 宛先のトークルームID
	var confFilePath string
	var authKeyPath string
	var destUsername string
	var destRoomId string
	var messages string
	flag.StringVar(&confFilePath, "c", "", "configuration file path")
	flag.StringVar(&authKeyPath, "k", "", "Authorization Key file path")
	flag.StringVar(&destUsername, "d", "", "Destination username to speak")
	flag.StringVar(&destRoomId, "r", "", "Destination talk room id to speak")

	flag.Parse()

	// しゃべらせるメッセージは名前なし引数にする
	if args := flag.Args(); len(args) > 0 {
		messages = args[0]
	}

	if len(messages) == 0 {
		if terminal.IsTerminal(syscall.Stdin) {
			exitOnEmpty("Messages to speak must not be empty")
		}
		messages = readFromStdin()
	} else if messages == "-" {
		messages = readFromStdin()
	}

	if len(confFilePath) == 0 {
		exitOnEmpty("configuration file path must not be empty")
	}

	if len(authKeyPath) == 0 {
		exitOnEmpty("Authorization Key file path must not be empty")
	}

	if len(destUsername) == 0 && len(destRoomId) == 0 {
		exitOnEmpty("Destination username or talk room id to speak must not be empty")
	}

	if len(destUsername) > 0 && len(destRoomId) > 0 {
		exitOnEmpty("Destination username and talk room id to speak must not be both set")
	}

	if len(messages) == 0 {
		exitOnEmpty("Messages to speak must not be empty")
	}

	//log.Printf("得られた値 : confFilePath = %s\n", confFilePath)
	//log.Printf("得られた値 : authKeyPath = %s\n", authKeyPath)
	//log.Printf("得られた値 : destUsername = %s\n", destUsername)
	//log.Printf("得られた値 : msg = %s\n", messages)

	conf, err := utilsV2.Load(filepath.Clean(confFilePath))
	if err != nil {
		log.Fatalln(err)
	}

	assertion, err := authV2.GenerateAuthToken(conf.ClientId, conf.ServiceAccount, authKeyPath)
	if err != nil {
		log.Fatalln(err)
	}

	tokenRequest := authV2.TokenRequest{
		Assertion:    assertion,
		GrantType:    authV2.GrantTypeInitial.String(),
		ClientId:     conf.ClientId,
		ClientSecret: conf.ClientSecret,
		Scope:        "Bot,Bot.read",
	}

	tokenResp, err := tokenRequest.GetAccessToken()
	if err != nil {
		log.Fatalln(err)
	}

	text := botV2.NewI18nText(botV2.SupportedLocaleEnglish, messages)
	i18NTexts := []botV2.I18NText{text}
	payload := botV2.NewLocalizedTextPayload(i18NTexts, botV2.SupportedLocaleEnglish)

	if len(destUsername) == 0 {
		err = payload.SendToRoom(destRoomId, tokenResp.AccessToken, conf)
	} else {
		err = payload.SendToUser(destUsername, tokenResp.AccessToken, conf)
	}
	if err != nil {
		log.Fatalln(err)
	}

	os.Exit(0)
}

func readFromStdin() string {
	msgBytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		exitOnEmpty("Could not read messages to speak")
	}
	return string(msgBytes)
}

func exitOnEmpty(msg string) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", msg)
	fmt.Fprintf(os.Stderr, "Usage: %s [options] messages\n", filepath.Base(os.Args[0]))
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "  messages")
	fmt.Fprintln(os.Stderr, "        messages to make LINE WORKS Bot speak")
	os.Exit(2)
}
