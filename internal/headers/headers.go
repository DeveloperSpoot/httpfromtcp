package headers

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strings"
)

type Headers map[string]string

// EXMAPLES
// "Host: localhost:42069\r\n\r\n"
// Host: localhost:42069\r\nToken: 50694839\r\n\r\n
// WHERE \r\n\r\n is the end of headers.

const crlf = "\r\n"
const rnrn = "\r\n\r\n"

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) SetHeader(name string, value string) error {
	name = strings.ToLower(name)

	if isValidName(name) == false {
		return errors.New("Filed-Name contains invalid characters.")
	}

	name = strings.TrimSpace(name)

	h[name] = value

	return nil
}

func (head Headers) GetDefualtHeaders(contentLen int) {

	head.SetHeader("Content-Length", fmt.Sprintf("%v", contentLen))
	head.SetHeader("Connection", "close")
	head.SetHeader("Content-Type", "text/plain")
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	iDX := bytes.Index(data, []byte(rnrn))

	log.Println("HEADERS PARSE")
	log.Println("IDEX: ", iDX, "ID:", idx)

	if idx == -1 {
		log.Println("HP HP HP Returning")
		return 0, false, nil
	}

	//End of headers, return proper data
	if idx == 0 || idx == 1 {
		log.Println("HP HP HP END END END END")
		return len(data[:idx]) + 2, true, nil
	}

	vdata := data[:idx]

	idxColon := bytes.Index(vdata, []byte(":"))

	fieldName := strings.ToLower(string(vdata[:idxColon]))

	fieldValue := string(vdata[1+idxColon:])

	idxSpace := strings.LastIndex(fieldName, " ")

	if idxSpace == len(fieldName)-1 {
		// NOTE:RFC requires this to return a status code 400
		return 0, false, errors.New("Invalid spacing header; specifcally spacing in the field name")
	}

	fieldName = strings.TrimSpace(fieldName)

	if len(fieldName) == 0 {
		return 0, false, errors.New("Invalid Header Length.")
	}

	if isValidName(fieldName) == false {
		return 0, false, errors.New("Invalid Header Field Name; Field Name Contains Invalid Character.")
	}

	fieldValue = strings.TrimSpace(fieldValue)

	if h[fieldName] != "" {
		h[fieldName] = h[fieldName] + ", " + fieldValue
	} else {
		h[fieldName] = fieldValue
	}

	log.Println("HP HP RETURN RETURN RETURN")
	log.Println("FN: ", fieldName, "\nFV: ", fieldValue)
	if iDX != -1 {
		log.Println(string(data[:iDX]), data[:iDX])
	}

	return len(vdata) + 1, false, nil
}

func isValidName(s string) bool {
	for _, r := range s {
		if strings.ContainsAny(string(r), "abcdefghijklmnopqrstuvwxyz0123456789!#$%&*-+.^_`|~") != true {
			return false
		}

	}

	return true
}
