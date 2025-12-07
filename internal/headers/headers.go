package headers

import (
	"bytes"
	"errors"
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
	if idx == 0 {
		return len(data), true, nil
	}

	idxColon := bytes.Index(data, []byte(":"))

	fieldName := strings.ToLower(string(data[:idxColon]))

	fieldValue := string(data[1+idxColon:])

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

	//WARN: Boot.Dev has this set not to include the last two bytes or the last crlf. I included it as it's techincally bytes consumed.
	return len(data), false, nil
}

func isValidName(s string) bool {
	for _, r := range s {
		if strings.ContainsAny(string(r), "abcdefghijklmnopqrstuvwxyz0123456789!#$%&*-+.^_`|~") != true {
			return false
		}

	}

	return true
}
