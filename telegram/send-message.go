package telegram

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

const (
	ARTSEM_CHAT_ID  = 505172736
	VALERYA_CHAT_ID = 307669920
	BOT_TOKEN       = "5059878186:AAHZ0GgLS6qDVajWQ7lO-Et03t7duKWJ8bE"
	TELEGRAM_URL    = "https://api.telegram.org/bot"
)

type BotSendMessageID struct {
	Result struct {
		Message_id string
	}
}

func SendMessage(text string) error {
	err := sendTextToTelegramChat(ARTSEM_CHAT_ID, text)
	if err != nil {
		return errors.Wrap(err, "error sending message to Telegram")
	}
	//err = sendTextToTelegramChat(VALERYA_CHAT_ID, text)
	//if err != nil {
	//	return err
	//}
	return nil
}

func sendTextToTelegramChat(chatID int, text string) error {
	telegramAPI := TELEGRAM_URL + BOT_TOKEN + "/sendMessage"
	response, err := http.PostForm(
		telegramAPI,
		url.Values{
			"chat_id": {strconv.Itoa(chatID)},
			"text":    {text},
		})

	if err != nil {
		log.Printf("error when posting text to the chat: %s", err.Error())
		return errors.Wrap(err, "error posting text to the chat")
	}
	defer func() {
		err := response.Body.Close()
		if err != nil {
			log.Printf("error when closing response body: %s", err.Error())
		}
	}()

	_, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("error in parsing telegram answer: %s", err.Error())
		return errors.Wrap(err, "error parsing telegram answer")
	}
	return nil
}
