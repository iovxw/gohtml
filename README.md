GoHTML
======

```
 ██████╗         ██╗  ██╗████████╗███╗   ███╗██╗
██╔════╝  █████╗ ██║  ██║╚══██╔══╝████╗ ████║██║
██║  ███╗██╔══██╗███████║   ██║   ██╔████╔██║██║
██║   ██║██║  ██║██╔══██║   ██║   ██║╚██╔╝██║██║
╚██████╔╝╚█████╔╝██║  ██║   ██║   ██║ ╚═╝ ██║███████╗
 ╚═════╝  ╚════╝ ╚═╝  ╚═╝   ╚═╝   ╚═╝     ╚═╝╚══════╝
```

Golang HTML模板，非正式项目

基本就是写着玩的

可以将Golang和HTML写在一起，然后用本工具转换成Golang文件

基本功能就是将`<a>`转换成`buffer.WriteString("<a>")`

例子
----

**转换前：**

```
package tpl

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

	print:=func() {
		<a>text</a>
	}
	print()

	test1(_buffer)
	test2(_buffer, t)

	return _buffer.String()
}

func test1(_buffer *bytes.Buffer) {
	<a>test1</a>
}

func test2(_buffer *bytes.Buffer, t string) {
	<a>test{{t}}</a>
}
```

**转换后：**

```
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

	print:=func() {
		_buffer.WriteString("<a>text</a>")
	}
	print()

	test1(_buffer)
	test2(_buffer,t)

	return _buffer.String()
}

func test1(_buffer *bytes.Buffer) {
	_buffer.WriteString("<a>test1</a>")
}

func test2(_buffer *bytes.Buffer, t string) {
	_buffer.WriteString("<a>test")
	_buffer.WriteString(t)
	_buffer.WriteString("</a>")
}
```

然后就可以直接调用里面的函数作为`http.HandleFunc`

使用帮助
-------

```
$ go run GoHTML.go

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
  方括号[]为选填项目，尖括号<>为必填项目
```

下载
----

[![Gobuild Download](http://gobuild.io/badge/github.com/Bluek404/GoHTML/downloads.svg)](http://gobuild.io/github.com/Bluek404/GoHTML)
