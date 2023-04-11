package main

import (
	"fmt"
	"github.com/linqcod/cbr-currency-app/cmd/internal/currency/model"
	"github.com/linqcod/cbr-currency-app/cmd/internal/currency/workerpool"
	"github.com/linqcod/cbr-currency-app/pkg/parser"
	"github.com/linqcod/cbr-currency-app/pkg/xmldecoder"
	"log"
	"math"
	"strconv"
	"strings"
)

func main() {
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

	currencyCodesMap := make(map[string]string, len(currencyCodes.Items))
	for _, item := range currencyCodes.Items {
		currencyCodesMap[item.ID] = item.EngName
	}

	parserTasks := make([]*workerpool.Task, len(currencyCodes.Items))
	for i := 0; i < len(currencyCodes.Items); i++ {
		task := workerpool.NewTask(currencyCodes.Items[i].ID, p, decoder)

		parserTasks[i] = task
	}

	resChan := make(chan *model.ValResults, len(currencyCodes.Items))

	pool := workerpool.NewPool(parserTasks, resChan, 5)
	pool.Run()

	var maxCursVal float32
	maxCursID := ""
	maxCursDate := ""
	minCursVal := float32(math.MaxFloat32)
	minCursID := ""
	minCursDate := ""
	for res := range resChan {
		var averageCursVal float32
		for _, v := range res.ValCurs.Records {
			valueStr := strings.Replace(v.Value, ",", ".", -1)
			value, err := strconv.ParseFloat(valueStr, 32)
			if err != nil {
				log.Fatal(err)
			}
			averageCursVal += float32(value)
			if maxCursVal < float32(value) {
				maxCursVal = float32(value)
				maxCursID = res.CurrencyId
				maxCursDate = v.Date
			}
			if minCursVal > float32(value) {
				minCursVal = float32(value)
				minCursID = res.CurrencyId
				minCursDate = v.Date
			}
		}
		averageCursVal /= float32(len(res.ValCurs.Records))
		fmt.Printf("Среднее значение курса рубля по валюте %s = %f\n", currencyCodesMap[res.CurrencyId], averageCursVal)
	}

	fmt.Printf("\nЗначение минимального курса валюты: %s = %f (%s)\n", currencyCodesMap[minCursID], minCursVal, minCursDate)
	fmt.Printf("Значение максимального курса валюты: %s = %f (%s)\n", currencyCodesMap[maxCursID], maxCursVal, maxCursDate)

	close(resChan)
}
