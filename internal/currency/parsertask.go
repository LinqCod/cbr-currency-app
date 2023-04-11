package currency

import (
	"fmt"
	"github.com/linqcod/cbr-currency-app/internal/currency/model"
	"github.com/linqcod/cbr-currency-app/pkg/parser"
	"github.com/linqcod/cbr-currency-app/pkg/xmldecoder"
)

type ParserTask struct {
	CurrencyId string
	startDate  string
	endDate    string
	parser     *parser.Parser
	xmlDecoder *xmldecoder.XMLDecoder
}

func NewParserTask(currencyId string, startDate string, endDate string, parser *parser.Parser, xmlDecoder *xmldecoder.XMLDecoder) *ParserTask {
	return &ParserTask{
		CurrencyId: currencyId,
		startDate:  startDate,
		endDate:    endDate,
		parser:     parser,
		xmlDecoder: xmlDecoder,
	}
}

func (t *ParserTask) Process() (result model.ValCurs, err error) {
	respBody, err := t.parser.Parse(fmt.Sprintf(
		"https://www.cbr.ru/scripts/XML_dynamic.asp?date_req1=%s&date_req2=%s&VAL_NM_RQ=%s",
		t.startDate, t.endDate, t.CurrencyId))
	if err != nil {
		return model.ValCurs{}, err
	}
	defer respBody.Close()

	if err = t.xmlDecoder.Decode(respBody, &result); err != nil {
		return model.ValCurs{}, err
	}

	return result, nil
}
