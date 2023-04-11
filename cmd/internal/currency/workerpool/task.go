package workerpool

import (
	"fmt"
	"github.com/linqcod/cbr-currency-app/cmd/internal/currency/model"
	"github.com/linqcod/cbr-currency-app/pkg/parser"
	"github.com/linqcod/cbr-currency-app/pkg/xmldecoder"
)

type Task struct {
	currencyId string
	parser     *parser.Parser
	xmlDecoder *xmldecoder.XMLDecoder
}

func NewTask(currencyId string, parser *parser.Parser, xmlDecoder *xmldecoder.XMLDecoder) *Task {
	return &Task{
		currencyId: currencyId,
		parser:     parser,
		xmlDecoder: xmlDecoder,
	}
}

func (t *Task) process(workerID int) (result model.ValCurs, err error) {
	fmt.Printf("Parser worker %d parsing currency with id: %s\n", workerID, t.currencyId)

	respBody, err := t.parser.Parse(fmt.Sprintf(
		"https://www.cbr.ru/scripts/XML_dynamic.asp?date_req1=14/12/2022&date_req2=14/03/2023&VAL_NM_RQ=%s",
		t.currencyId))
	if err != nil {
		return model.ValCurs{}, err
	}
	defer respBody.Close()

	if err = t.xmlDecoder.Decode(respBody, &result); err != nil {
		return model.ValCurs{}, err
	}

	return result, nil
}
