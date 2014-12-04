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
--------

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

源码规范
--------

作为`http.HandleFunc`的函数，必须返回值只有一个`string`（因为多了就不符合`http.HandleFunc`的要求了）

然后里面必须创建一个类型为`*bytes.Buffer`的变量，变量名称需要和设置的一样（默认是`_buffer`）。这个缓冲器是用于存放要输出的HTML代码的

然后作为`http.HandleFunc`的函数的返回值需要是`_buffer.String()`，如果你自定义了缓冲器变量名称，需要相应的改变

负责处理部分html的函数，推荐写成上面例子里的样式（接收一个缓冲器变量），就像：

```
func test1(_buffer *bytes.Buffer) {
    <a>test1</a>
}
```

因为会转换成：

```
func test1(_buffer *bytes.Buffer) {
    _buffer.WriteString("<a>test1</a>")
}
```

这个样子

当然你要是写成函数里新建一个缓冲器变量然后再`return`也是没问题的

不过那样就需要在使用函数的地方再加一个输出了（相当麻烦）

所以就乖乖按照标准写吧

然后还有一点要注意

就是分隔符输出变量那里，比如`{{t}}`

这里的分隔符只有在HTML中才会被转换，就像`<a>{{t}}</a>`

写在别的地方是不会转换的，比如：

```
<a>
if i==0 {
    {{t1}}
}else{
    {{t2}}
}
</a>
```

必须写成

```
if i==0 {
    <a>{{t1}}</a>
}else{
    <a>{{t2}}</a>
}
```

这样才可以

自动转换
--------

每次更新完模板后都要手动执行一次转换工具一定是超麻烦的

所以需要自动的工具

这里推荐[Unknwon](https://github.com/Unknwon)的[BRA](https://github.com/Unknwon/bra/)（名字好像有点奇怪？），当然其他相似功能的工具也是可以的（如果有更好的还拜托推荐一下）

具体使用方法这里就不说了，自己去看吧

下载
----

`go get github.com/Bluek404/gohtml`

OR

[![Gobuild Download](http://gobuild.io/badge/github.com/Bluek404/gohtml/downloads.svg)](http://gobuild.io/github.com/Bluek404/gohtml)
