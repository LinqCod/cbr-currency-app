package xmldecoder

import (
	"encoding/xml"
	"fmt"
	"io"

	"golang.org/x/text/encoding/charmap"
)

type XMLDecoder struct{}

func NewXMLDecoder() *XMLDecoder {
	return &XMLDecoder{}
}

func (d *XMLDecoder) Decode(data io.ReadCloser, v interface{}) error {
	dec := xml.NewDecoder(data)
	dec.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch charset {
		case "windows-1251":
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		default:
			return nil, fmt.Errorf("unknown charset: %s", charset)
		}
	}

	if err := dec.Decode(&v); err != nil {
		return err
	}

	return nil
}
