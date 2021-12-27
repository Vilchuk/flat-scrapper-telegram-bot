package telegram

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

const (
	chatID       = 505172736
	chatID2      = 307669920
	BOT_TOKEN    = "5059878186:AAHZ0GgLS6qDVajWQ7lO-Et03t7duKWJ8bE"
	TELEGRAM_URL = "https://api.telegram.org/bot"
)

type BotSendMessageID = struct {
	Result struct {
		Message_id string
	}
}

func SendMessage(text string) {
	sendTextToTelegramChat(chatID, text)
	sendTextToTelegramChat(chatID2, text)
}

// sendTextToTelegramChat sends a text message to the Telegram chat identified by its chat Id
func sendTextToTelegramChat(chatId int, text string) (string, error) {

	log.Printf("Sending %s to chat_id: %d", text, chatId)
	var telegramApi string = TELEGRAM_URL + BOT_TOKEN + "/sendMessage"
	response, err := http.PostForm(
		telegramApi,
		url.Values{
			"chat_id": {strconv.Itoa(chatId)},
			"text":    {text},
		})

	if err != nil {
		log.Printf("error when posting text to the chat: %s", err.Error())
		return "", err
	}
	defer response.Body.Close()

	var bodyBytes, errRead = ioutil.ReadAll(response.Body)
	if errRead != nil {
		log.Printf("error in parsing telegram answer %s", errRead.Error())
		return "", err
	}
	bodyString := string(bodyBytes)
	//log.Printf("Body of Telegram Response: %s", bodyString)

	return bodyString, nil
}
