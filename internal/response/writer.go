package response

import (
	"errors"
	"io"

	"github.com/DeveloperSpoot/httpfromtcp/internal/headers"
)

type Writer struct {
	writerState
	output io.Writer
}

type writerState int

const (
	writerNotStarted writerState = iota
	writerStatusLineDone
	writerHeadersDone
	writerBodyStarted
)

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writerState: writerNotStarted,
		output:      w,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.writerState > writerNotStarted {
		return errors.New("Status Line already written.")
	}

	err := WriteStatusLine(w.output, statusCode)

	if err != nil {
		return err
	}
	w.writerState = writerStatusLineDone
	return nil
}

func (w *Writer) WriteHeaders(head headers.Headers) error {
	if w.writerState == writerNotStarted {
		return errors.New("Ensure to write the response in order! Start with the Status Line.")
	}
	if w.writerState == writerHeadersDone {
		return errors.New("Headers already written.")
	}

	err := WriteHeaders(w.output, head)
	if err != nil {
		return err
	}
	w.writerState = writerHeadersDone
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.writerState == writerNotStarted {
		return 0, errors.New("Ensure to write the response in order! Start with the Status Line.")
	}

	if w.writerState == writerStatusLineDone {
		return 0, errors.New("Ensure to write headers before the body!")
	}

	idx, err := w.output.Write(p)

	if err != nil {
		return idx, err
	}

	w.writerState = writerBodyStarted

	return idx, nil
}
