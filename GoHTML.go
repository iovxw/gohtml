package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	delimiterLeft  = flag.String("dl", "{{", "Left delimiter")
	delimiterRight = flag.String("dr", "}}", "Right delimiter")
	suffix         = flag.String("suffix", "gohtml", "Suffix of the GoHTML template files")
	buffer         = flag.String("buffer", "_buffer", "Buffer name")
)

func main() {
	flag.Parse()
	// 如果命令行参数为空，那么输出帮助
	if flag.Arg(0) == "" {
		fmt.Println(`
 ██████╗         ██╗  ██╗████████╗███╗   ███╗██╗
██╔════╝  █████╗ ██║  ██║╚══██╔══╝████╗ ████║██║
██║  ███╗██╔══██╗███████║   ██║   ██╔████╔██║██║
██║   ██║██║  ██║██╔══██║   ██║   ██║╚██╔╝██║██║
╚██████╔╝╚█████╔╝██║  ██║   ██║   ██║ ╚═╝ ██║███████╗
 ╚═════╝  ╚════╝ ╚═╝  ╚═╝   ╚═╝   ╚═╝     ╚═╝╚══════╝
=====================================================
GoHTML使用帮助:
===============

命令结构:
  命令 [参数] <模板文件夹路径>

参数：
  -dl <字符串>     | 默认：{{      | 左分隔符样式
  -dr <字符串>     | 默认：}}      | 右分隔符样式
  -suffix <字符串> | 默认：gohtml  | GoHTML模板文件后缀
  -buffer <字符串> | 默认：_buffer | 缓冲器变量名称

举例：
  $ gohtml -dl <{ -dr }> -suffix temp -buffer buf /home/bluek404/gocode/web/view
   | 则会将
   | /home/bluek404/gocode/web/view
   | 里面所有temp为后缀的文件转换为go文件后放到同一文件夹内
   | 将缓冲器变量名称设为“buf”
   | 并设置左分隔符为“<{” ，右分割符为“}>”

备注：
  方括号[]为选填项目，尖括号<>为必填项目`)
		return
	}
	// 获取文件夹位置参数
	folder := flag.Arg(0)
	// 获取文件夹信息
	info, err := os.Lstat(folder)
	if err != nil {
		log.Println(err)
		return
	}
	// 检查所输入参数是否为文件夹
	if !info.IsDir() {
		log.Println("所输入路径不是文件夹")
		return
	}
	// 处理文件夹内的模板文件
	err = filepath.Walk(folder, walk)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("-==处理完成==-")
}

// 文件处理
func walk(path string, info os.FileInfo, err error) error {
	if info == nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	goHTML := regexp.MustCompile("." + *suffix + "$")
	// 检查是否为gohtml文件
	if !goHTML.MatchString(info.Name()) {
		return nil
	}
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	gen := []byte(generate(string(buf)))
	// 格式化转换后的文件
	buf, err = format.Source(gen)

	// 将文件后缀变为.go
	outPath := path[:len(path)-len(*suffix)] + "go"

	if err != nil {
		// 输出未格式化的错误文件
		ioutil.WriteFile(outPath, gen, 0700)
		return fmt.Errorf("Format error: %v %v", outPath, err)
	}
	// 输出文件
	ioutil.WriteFile(outPath, buf, 0700)
	fmt.Println(path, "==>", outPath)
	return nil
}

// 转换模板文件
func generate(in string) string {
	var _buffer bytes.Buffer

	isHTML := regexp.MustCompile(`^\s*<.*>\s*$`)
	// 分隔符
	delimiter := regexp.MustCompile(*delimiterLeft + "[^" + *delimiterLeft + *delimiterRight + "]+" + *delimiterRight)

	a := regexp.MustCompile(`^\s*`)
	z := regexp.MustCompile(`\s*$`)

	var htmlBUF string

	r := new(readLine)
	for {
		// 按行读取处理
		buf, ok := r.read(in)

		// 检查本行是否为HTML
		if isHTML.MatchString(buf) {
			// 去首尾空
			buf = a.ReplaceAllLiteralString(buf, "")
			buf = z.ReplaceAllLiteralString(buf, "")
			// 转义双引号
			buf = strings.Replace(buf, `"`, `\"`, -1)

			// 检查是否有分隔符需要替换
			if delimiter.MatchString(buf) {
				// 替换分隔符
				all := delimiter.FindAllString(buf, -1)
				for _, v := range all {
					// 去掉两边分隔符
					vBuf := v[len(*delimiterLeft) : len(v)-len(*delimiterRight)]
					// 组合插入
					vBuf = "\")\n" + *buffer + ".WriteString(" + vBuf + ")\n" + *buffer + ".WriteString(\""
					// 替换
					buf = strings.Replace(buf, v, vBuf, -1)
				}
			}

			// 将本行添加到缓存中，输出为一行
			// 就是将两行
			// _buffer.WriteString('<a>xxx</a>')
			// _buffer.WriteString('<a>xxx</a>')
			// 简化为一行
			// _buffer.WriteString('<a>xxx</a><a>xxx</a>')
			htmlBUF += buf
		} else {
			// 本行为Golang
			// 检查是否有html需要输出
			if htmlBUF != "" {
				// 输出html
				_buffer.WriteString(*buffer + ".WriteString(\"" + htmlBUF + "\")\n")
				// 清空缓存
				htmlBUF = ""
			}
			_buffer.WriteString(buf + "\n")
		}
		// 检查是否已经读取完毕
		if ok {
			break
		}
	}
	return _buffer.String()
}

type readLine struct {
	// 用于记录每行的起始位置
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
	// 行的结束位置
	r.end = r.start + n
	// 提取一行
	line := string(str[r.start:r.end])
	// 下一行的开始位置
	r.start = r.end + 1

	return line, false
}
