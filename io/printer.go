package io

import (
	"encoding/json"
	"io"
)

type ResponsePrinter interface {
	Print(question, answer string)
}

type NoOpResponsePrinter struct{}

func (NoOpResponsePrinter) Print(question, answer string) {}

type PlainResponsePrinter struct {
	Target io.Writer
}

func (p PlainResponsePrinter) Print(question, answer string) {
	p.Target.Write([]byte("SOQ\n"))
	p.Target.Write([]byte(question + "\n"))
	p.Target.Write([]byte("EOQ\nSOA\n"))
	p.Target.Write([]byte(answer + "\n"))
	p.Target.Write([]byte("EOA\n"))
}

type JsonResponsePrinter struct {
	Target io.Writer
}

func (j JsonResponsePrinter) Print(question, answer string) {
	json.NewEncoder(j.Target).Encode(struct {
		Question string `json:"q"`
		Answer   string `json:"a"`
	}{
		Question: question,
		Answer:   answer,
	})
}

type MultiResponsePrinter struct {
	Printers []ResponsePrinter
}

func (m MultiResponsePrinter) Print(question, answer string) {
	for _, d := range m.Printers {
		d.Print(question, answer)
	}
}
