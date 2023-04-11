package main

import (
	"fmt"
	"github.com/linqcod/cbr-currency-app/internal/currency"
	"github.com/linqcod/cbr-currency-app/internal/currency/model"
	"github.com/linqcod/cbr-currency-app/pkg/parser"
	"github.com/linqcod/cbr-currency-app/pkg/xmldecoder"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
)

func main() {
	// Парсим названия курсов валют и их идентификаторы
	p := parser.NewParser()
	decoder := xmldecoder.NewXMLDecoder()

	respBody, err := p.Parse("https://www.cbr.ru/scripts/XML_val.asp?d=0")
	if err != nil {
		log.Fatal(err)
	}
	defer respBody.Close()

	var currencyCodes model.Valuta
	if err = decoder.Decode(respBody, &currencyCodes); err != nil {
		log.Fatal(err)
	}

	// Создаем мапу идентификаторов и названий курсов валют для быстрого обращения по ID валюты
	currencyCodesMap := make(map[string]string, len(currencyCodes.Items))
	for _, item := range currencyCodes.Items {
		currencyCodesMap[item.ID] = item.EngName
	}

	// Создаем задачи на парсинг каждого идентификатора валюты период в 90 дней
	currentTime := time.Now()
	endDate := currentTime.Format("02/01/2006")
	startDate := currentTime.AddDate(0, 0, -90).Format("02/01/2006")

	parserTasks := make([]*currency.ParserTask, len(currencyCodes.Items))
	for i := 0; i < len(currencyCodes.Items); i++ {
		task := currency.NewParserTask(currencyCodes.Items[i].ID, startDate, endDate, p, decoder)

		parserTasks[i] = task
	}

	// Выполняем задачи на парсинг и собираем результат
	var maxCursVal float32
	maxCursID := ""
	maxCursDate := ""
	minCursVal := float32(math.MaxFloat32)
	minCursID := ""
	minCursDate := ""
	for _, task := range parserTasks {
		valCurs, err := task.Process()
		if err != nil {
			log.Printf("error while processing task: %v", err)
			continue
		}
		if len(valCurs.Records) != 0 {
			var averageCursVal float32
			for _, v := range valCurs.Records {
				valueStr := strings.Replace(v.Value, ",", ".", -1)
				value, err := strconv.ParseFloat(valueStr, 32)
				if err != nil {
					log.Fatal(err)
				}
				averageCursVal += float32(value)
				if maxCursVal < float32(value) {
					maxCursVal = float32(value)
					maxCursID = task.CurrencyId
					maxCursDate = v.Date
				}
				if minCursVal > float32(value) {
					minCursVal = float32(value)
					minCursID = task.CurrencyId
					minCursDate = v.Date
				}
			}
			averageCursVal /= float32(len(valCurs.Records))

			fmt.Printf("Среднее значение курса рубля по валюте %s = %.2f\n", currencyCodesMap[task.CurrencyId], averageCursVal)
		}
	}

	fmt.Printf("\nЗначение минимального курса валюты: %s = %.2f (%s)\n", currencyCodesMap[minCursID], minCursVal, minCursDate)
	fmt.Printf("Значение максимального курса валюты: %s = %.2f (%s)\n", currencyCodesMap[maxCursID], maxCursVal, maxCursDate)
}
