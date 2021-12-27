package main

import (
	"fmt"
	"main.go/actions"
	"main.go/db"
	"main.go/telegram"
	"time"
)

func main() {
	fmt.Println("main started")
	Search(true)
	ticker := time.NewTicker(60 * time.Second)
	quit := make(chan struct{})

	for {
		select {
		case <-ticker.C:
			Search(false)
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func Search(isFirst bool) {
	dt := time.Now()
	fmt.Println("Search started " + dt.Format("01-02-2006 15:04:05"))

	otodomFlats := actions.OtodomGetFlats("https://www.otodom.pl/pl/oferty/wynajem/mieszkanie/warszawa?priceMax=4000&areaMin=60&roomsNumber=%5BTHREE%2CFOUR%5D&hasPhotos=true")
	olxFlats := actions.OlxGetFlats("https://www.olx.pl/nieruchomosci/mieszkania/wynajem/warszawa/?search%5Bfilter_float_price%3Ato%5D=4500&search%5Bfilter_enum_furniture%5D%5B0%5D=yes&search%5Bfilter_float_m%3Afrom%5D=60&search%5Bfilter_enum_rooms%5D%5B0%5D=three&search%5Bfilter_enum_rooms%5D%5B1%5D=four&search%5Bphotos%5D=1")

	for _, flat := range otodomFlats {
		isFlatNew := db.IsFlatNew("otodom", flat.Href, flat.Hash)
		if isFlatNew {
			if !isFirst {
				telegram.SendMessage(fmt.Sprintf("%s\n%s", "Обьявление найдено на otodom.pl", flat.Href))
				fmt.Println(flat.Href)
			}
		}
	}

	for _, flat := range olxFlats {
		isFlatNew := db.IsFlatNew("olx", flat.Href, flat.Hash)
		if isFlatNew {
			if !isFirst {
				telegram.SendMessage(fmt.Sprintf("%s\n%s", "Обьявление найдено на olx.pl", flat.Href))
				fmt.Println(flat.Href)
			}
		}
	}

	if isFirst {
		//telegram.SendMessage("Бот запущен! Поиск начат. \n\nХочу заметить, что обьявления, которые уже существуют, не появятся, будут появляться обьявления, появившиеся позже текущего момента.\n\nУдачи, у тебя всё обязательно получится!!! :)")
	}

	fmt.Println("Search ended\n")
}
