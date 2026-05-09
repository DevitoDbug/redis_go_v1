package resp

import (
	"fmt"
	"io"
)

type Writer struct {
	writer io.Writer
}

func NewWriter(writer io.Writer) *Writer {
	return &Writer{
		writer: writer,
	}
}

func (w *Writer) Write(v Value) error {
	response := v.Marshal()

	_, err := w.writer.Write(response)
	if err != nil {
		fmt.Printf("This is where the program crashes: \n%v\n", err)
		return err
	}

	return nil
}
