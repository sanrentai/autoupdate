package autoupdate

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type AutoUpdate struct {
	Url      string // 下载路径
	Softname string // 软件名
	CurrVer  string // 当前版本
	CurrName string // 当前运行路径及文件名
}

func (au *AutoUpdate) Update() error {

	if strings.HasSuffix(au.CurrName, "update.exe") { //开始更新

		time.Sleep(3 * time.Second)

		if au.copyFile() {

			cmd := exec.Command("cmd", "/c", "start", filepath.Dir(au.CurrName)+"\\"+au.Softname+".exe")
			cmd.Start()
			cmd.Wait()
			os.Exit(0)

		}

	} else { //检测更新

		resp, err := http.Get(au.Url + "/" + au.Softname + ".txt?num=" + fmt.Sprintf("%d", rand.Intn(1000)))
		if err != nil {
			return err
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		newVer := string(body)
		if newVer != au.CurrVer {
			fmt.Println("远端版本号:", newVer)
			fmt.Println("本地版本号:", au.CurrVer)
			if err := au.getNewVer(); err != nil {
				return err
			} else {
				cmd := exec.Command("cmd", "/c", "start", filepath.Dir(au.CurrName)+"\\update.exe")
				cmd.Start()
				cmd.Wait()
				os.Exit(0)
			}

		} else {

			os.Remove("update.exe")

		}

	}

	return nil

}

func (au *AutoUpdate) NeedUpdate() bool {
	if strings.HasSuffix(au.CurrName, "update.exe") { //开始更新
		return true
	}
	os.Remove("update.exe")

	resp, err := http.Get(au.Url + "/" + au.Softname + ".txt?num=" + fmt.Sprintf("%d", rand.Intn(1000)))
	if err != nil {
		return false
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	newVer := string(body)
	return newVer != au.CurrVer
}

func (au *AutoUpdate) getNewVer() error {

	var fsize int64
	buf := make([]byte, 32*1024)
	var written int64

	client := http.Client{Timeout: 900 * time.Second}

	resp, err := client.Get(au.Url + "/" + au.Softname + ".exe?num=" + fmt.Sprintf("%d", rand.Intn(1000)))

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.Status == "200 OK" {
		fmt.Println("开始下载")
		fsize, err = strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 32)
		newFile, err := os.Create("update.exe")

		if err != nil {
			return err
		}

		defer newFile.Close()

		for {

			//读取bytes
			nr, er := resp.Body.Read(buf)
			if nr > 0 {
				//写入bytes
				nw, ew := newFile.Write(buf[0:nr])
				//数据长度大于0
				if nw > 0 {
					written += int64(nw)
				}
				//写入出错
				if ew != nil {
					err = ew
					break
				}
				//读取是数据长度不等于写入的数据长度
				if nr != nw {
					err = io.ErrShortWrite
					break
				}
			}
			if er != nil {
				if er != io.EOF {
					err = er
				}
				break
			}
			go func(fsize, written int64) {

				i := int(written * 100 / fsize)
				m := i / 5

				fmt.Fprintf(os.Stdout, "%f%% [%s]\r", float32(written*100)/float32(fsize), getS(m, "#")+getS(20-m, "="))

			}(fsize, written)

		}

		// _, err = io.Copy(newFile, resp.Body)

		return nil

	} else {
		fmt.Println("更新程序出错")
		return errors.New(resp.Status)

	}

}

func (au *AutoUpdate) copyFile() bool {

	var fsize int64
	buf := make([]byte, 32*1024)
	var written int64

	source_open, err := os.Open(au.CurrName)
	info, _ := source_open.Stat()
	fsize = info.Size()

	if err != nil {
		return false
	}

	defer source_open.Close()

	dest_open, err := os.OpenFile(filepath.Dir(au.CurrName)+"\\"+au.Softname+".exe", os.O_CREATE|os.O_WRONLY, 644)

	if err != nil {

		return false

	}

	defer dest_open.Close()

	//进行数据拷贝

	for {

		//读取bytes
		nr, er := source_open.Read(buf)
		if nr > 0 {
			//写入bytes
			nw, ew := dest_open.Write(buf[0:nr])
			//数据长度大于0
			if nw > 0 {
				written += int64(nw)
			}
			//写入出错
			if ew != nil {
				err = ew
				break
			}
			//读取是数据长度不等于写入的数据长度
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
		go func(fsize, written int64) {

			i := int(written * 100 / fsize)
			// 打印进度条
			fmt.Fprintf(os.Stdout, "%f%% [%s]\r", float32(written*100)/float32(fsize), getS(i, "#")+getS(100-i, " "))

		}(fsize, written)

	}

	// _, copy_err := io.Copy(dest_open, source_open)

	if err != nil {

		return false

	} else {

		return true

	}

}

func getS(n int, char string) (s string) {
	for i := 1; i <= n; i++ {
		s += char
	}
	return
}
