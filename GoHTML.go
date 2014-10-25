package main

import (
	"bytes"
	"fmt"
	//"log"
	"regexp"
	"strings"
)

func main() {
	t := `package tpl

import (
	"fmt"
	"bytes"
)

func Index() string {
	_buffer := new(bytes.Buffer)

	<a>Text</a>

	for i:=0; i<10; i++ {
		<a>{{string(i)}}</a>
	}

	var t string = "2"
	fmt.Println(t)
	{{t}}

	print:=func() {
		<a>text</a>
	}
	print()

	test1(_buffer)
	test2(_buffer, t)

	return _buffer.String()
}

func test1(_buffer bytes.Buffer) {
	<a>test1</a>
}

func test2(_buffer bytes.Buffer, t string) {
	<a>test{{t}}</a>
}`
	fmt.Println(generate(t))
	/*
		package tpl

		import (
			"fmt"
			"bytes"
		)

		func Index() string {
			_buffer := new(bytes.Buffer)

			_buffer.WriteString("<a>Text</a>")

			for i:=0; i<10; i++ {
				_buffer.WriteString("<a>")
				_buffer.WriteString(string(i))
				_buffer.WriteString("</a>")
			}

			var t string = "2"
			fmt.Println(t)
			_buffer.WriteString(t)

			print:=func() {
				_buffer.WriteString("<a>text</a>")
			}
			print()

			test1(_buffer)
			test2(_buffer,t)

			return _buffer.String()
		}

		func test1(_buffer bytes.Buffer) {
			_buffer.WriteString("<a>test1</a>")
		}
	
		func test2(_buffer bytes.Buffer, t string) {
			_buffer.WriteString("<a>test")
			_buffer.WriteString(t)
			_buffer.WriteString("</a>")
		}
	*/
}

func generate(in string) (string, error) {
	var _buffer bytes.Buffer

	isHTML := regexp.MustCompile(`^\s*<.*>\s*$`)
	space := regexp.MustCompile(`^\s*$`)
	symbolLeft := "{{"
	symbolRight := "}}"
	symbolLeftLen := len(symbolLeft)
	symbolRightLen := len(symbolRight)
	// 分隔符
	symbol := regexp.MustCompile(symbolLeft + ".*" + symbolRight)

	a := regexp.MustCompile(`^\s*`)
	z := regexp.MustCompile(`\s*$`)

	var htmlBUF string

	r := new(readLine)
	for {
		buf, ok := r.read(in)

		switch {
		case isHTML.MatchString(buf):
			// 去首尾空
			buf = a.ReplaceAllLiteralString(buf, "")
			buf = z.ReplaceAllLiteralString(buf, "")
			// 转义双引号
			buf = strings.Replace(buf, `"`, `\"`, -1)
			// 替换分隔符
			if symbol.MatchString(buf) {
				// 找到本行全部分隔符
				symbolBUF := symbol.FindAllString(buf, -1)
				for _, v := range symbolBUF {
					// 去掉两边分隔符
					vBUF := v[symbolLeftLen : len(v)-symbolRightLen]
					// 组合插入
					vBUF = "\")\n_buffer.WriteString(" + vBUF + ")\n_buffer.WriteString(\""
					// 替换
					buf = strings.Replace(buf, v, vBUF, -1)
					fmt.Println(v)
				}
			}

			// 将本行添加到缓存中，输出为一行
			// 就是将两行
			// _buffer.WriteString('<a>xxx</a>')
			// _buffer.WriteString('<a>xxx</a>')
			// 简化为一行
			// _buffer.WriteString('<a>xxx</a><a>xxx</a>')
			htmlBUF += buf
		case space.MatchString(buf):
			// 空行，跳过
		default:
			// 检查是否有html需要输出
			if htmlBUF != "" {
				// 输出html
				_buffer.WriteString("_buffer.WriteString(\"" + htmlBUF + "\")\n")
				// 清空缓存
				htmlBUF = ""
			}
			_buffer.WriteString(buf + "\n")
		}

		if ok {
			break
		}
	}
	return _buffer.String(), nil
}

type readLine struct {
	start int
	end   int
}

// 按行输出字符串
func (r *readLine) read(str string) (string, bool) {
	// 找到换行符位置
	n := strings.Index(str[r.start:], "\n")
	if n == -1 {
		// 无法找到换行符，已经到达最后一行
		// 返回ok
		line := string(str[r.start:])
		r.start = 0
		r.end = 0
		return line, true
	}
	r.end = r.start + n
	line := string(str[r.start:r.end])
	r.start = r.end + 1 // 跳过换行符,下次循环时从换行符后开始

	return line, false
}
