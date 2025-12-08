package headers

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

type Headers map[string]string

// EXMAPLES
// "Host: localhost:42069\r\n\r\n"
// Host: localhost:42069\r\nToken: 50694839\r\n\r\n
// WHERE \r\n\r\n is the end of headers.

const crlf = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))

	if idx == -1 {
		return 0, false, nil
	}

	//End of headers, return proper data
	if idx == 0 || idx == 1 {
		return len(data), true, nil
	}

	vdata := data[:idx]

	idxColon := bytes.Index(vdata, []byte(":"))

	fieldName := strings.ToLower(string(vdata[:idxColon]))

	fieldValue := string(vdata[1+idxColon:])
	fmt.Println(fieldValue, "|")
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
