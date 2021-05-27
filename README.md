# LineWorksBotMessenger
CLI app that make Bot speak mseeages to LINE WORKS / LINE WORKSにBot投稿させるCLIアプリケーション

## Requirements

- Go 1.15 or later

## How to build

```bash
$ git clone git@github.com:f97one/LineWorksBotMessenger.git
$ cd LineWorksBotMessenger
$ go build
```

## App usage

cli format is below:

```
LineWorksBotMessenger [options] messages
  -c configFilePath
        configuration file path
  -d userId
        Destination username to speak
  -k authorizationKeyPath
        Authorization Key file path
  messages
        messages to make LINE WORKS Bot speak
```

configuration file format is below:

```json
{
  "api_id": "dev_console_api_id",
  "consumer_key": "non_redirect_server_api_consumer_key",
  "server_id": "server_list_bot_id",
  "bot_no": 9999999
}
```

- api_id : API ID supplied from LINE WORKS Developer Console
- consumer_key : non-redirect server API Consumer Key supplied from LINE WORKS Developer Console
- server_id : Server List Bot ID supplied from LINE WORKS Developer Console
- bot_no : Bot Number supplied from LINE WORKS administration console

## Licensing

This app is licensed under the GPL v3.