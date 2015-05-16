package tpl

import (
	"bytes"
)

func head(_buffer *bytes.Buffer) {
	_buffer.WriteString("<head><title>Hello from GoHTML</title></head>")
}

func Index(name string) []byte {
	_buffer := new(bytes.Buffer)

	_buffer.WriteString("<html>")
	head(_buffer)
	if name == "" {
		_buffer.WriteString("<a>Hello World!</a>")
	} else {
		_buffer.WriteString("<a>Hello ")
		_buffer.WriteString(name)
		_buffer.WriteString("</a>")
	}
	_buffer.WriteString("<html>")

	return _buffer.Bytes()
}
