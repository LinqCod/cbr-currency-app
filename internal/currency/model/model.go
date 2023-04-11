package model

import "encoding/xml"

type ValResults struct {
	ValCurs    ValCurs
	CurrencyId string
}

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Records []Record `xml:"Record"`
}

type Record struct {
	Date    string `xml:"Date,attr"`
	Nominal int    `xml:"Nominal"`
	Value   string `xml:"Value"`
}

type Valuta struct {
	XMLName xml.Name `xml:"Valuta"`
	Items   []Item   `xml:"Item"`
}

type Item struct {
	ID      string `xml:"ID,attr"`
	EngName string `xml:"EngName"`
}
