package main

import (
	"fmt"
	"log"
	"main.go/db"
	"main.go/models"
	"main.go/parser"
	"main.go/telegram"
	"os"
	"time"

	"github.com/pkg/errors"
)

func main() {
	fmt.Println("main started")

	if err := search(true); err != nil {
		log.Println("error during initial search:", err)
	}

	ticker := time.NewTicker(60 * time.Second)
	quit := make(chan struct{})

	for {
		select {
		case <-ticker.C:
			if err := search(false); err != nil {
				log.Println("error during search:", err)
			}
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func search(isFirst bool) error {
	dbConnString := os.Getenv("DB_CONNECTION_STRING")
	println(dbConnString)
	if dbConnString == "" {
		return errors.New("DB_CONNECTION_STRING environment variable is not set")
	}
	database, err := db.NewDatabase(dbConnString)
	if err != nil {
		return errors.Wrap(err, "ошибка получения соединения с базой данных")
	}
	defer func() {
		if cerr := database.Close(); cerr != nil {
			log.Println("error closing database connection:", cerr)
		}
		fmt.Println("Подключение к базе данных закрыто\n---------------------------------")
	}()

	dt := time.Now()
	log.Println("Поиск начат", dt.Format("01-02-2006 15:04:05"))

	flats, err := parseFlats()
	if err != nil {
		return errors.Wrap(err, "ошибка при разборе объявлений")
	}

	for i, flat := range flats {
		isFlatNew, err := database.IsFlatNew(flat)
		if err != nil {
			log.Println("ошибка при проверке нового объявления:", err)
			continue
		}
		if isFlatNew {
			if !isFirst {
				log.Println("[", i, "]", flat.FoundOnSite, flat.Href)
				message := fmt.Sprintf("Новое объявление на сайте %s:\n%s", flat.FoundOnSite, flat.Href)
				if err := telegram.SendMessage(message); err != nil {
					log.Println("ошибка отправки сообщения в Telegram:", err)
					continue
				}
			}
		}
	}

	if isFirst {
		message := "Бот запущен! Поиск начат.\n\nПожалуйста, обратите внимание, что отображаются только новые объявления, появившиеся после запуска бота.\n\nУдачи в поиске!"
		if err := telegram.SendMessage(message); err != nil {
			return errors.Wrap(err, "ошибка отправки сообщения в Telegram")
		}
	}

	log.Println("Поиск завершен")
	return nil
}

func parseFlats() ([]models.Flat, error) {
	olxFlats, err := parser.ParseOlx("https://www.olx.pl/nieruchomosci/mieszkania/wynajem/warszawa/?search%5Bfilter_enum_furniture%5D%5B0%5D=yes&search%5Bfilter_enum_rooms%5D%5B0%5D=three&search%5Bfilter_enum_rooms%5D%5B1%5D=four&search%5Bfilter_float_m%3Afrom%5D=60&search%5Bfilter_float_price%3Ato%5D=4500&search%5Border%5D=created_at%3Adesc&search%5Bphotos%5D=1&view=list")
	if err != nil {
		return nil, errors.Wrap(err, "ошибка при парсинге объявлений OLX")
	}

	return olxFlats, nil
}
