package tpl

import (
	"bytes"
)

func Index(name string) string {
	_buffer := new(bytes.Buffer)

	if name == "" {
		_buffer.WriteString("<a>Hello World!</a>")
	} else {
		_buffer.WriteString("<a>Hello ")
		_buffer.WriteString(name)
		_buffer.WriteString("</a>")
	}

	return _buffer.String()
}
