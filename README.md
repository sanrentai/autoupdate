# autoupdate
软件自动更新，用于windows
在服务器上保存两个文件 比如  bbb.txt,bbb.exe   txt为可执行文件,exe为需要更新执行程序

```golang
package main

import (
	"os"
	"path/filepath"

	"github.com/sanrentai/autoupdate"
)

func main() {
	mainUpdate()
}

func mainUpdate() {
	f, _ := filepath.Abs(os.Args[0])

	au := autoupdate.AutoUpdate{
		Url:      "服务器地址",
		Softname: "bbb",   // 软件名
		CurrVer:  "0.0.1", // 当前版本
		CurrName: f,       // 当前运行路径及文件名
	}

	if au.NeedUpdate() {
		err := au.Update()
		if err != nil {
			panic(err)
		}
	}

}

```
