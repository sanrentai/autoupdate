# autoupdate

软件自动更新(只适用于 windows 平台)

软件部署目录下需要包括两个文件，如 demo.exe, update.exe，其中 demo.exe 为程序文件，update.exe 为升级程序。

以下为 update.exe 的样例代码:

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
