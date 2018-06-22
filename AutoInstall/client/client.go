package main

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	ole "github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"github.com/manifoldco/promptui"
)

func checkErrAtMainFunc(err error) {
	if err != nil {
		log.Println(err)
		fmt.Scanln()
	}
}

func openINI(path string) {
	cmd := exec.Command("cmd", "/C", "start", path)
	err := cmd.Start()
	if err != nil {
		return
	}
}

func genShortcut(goos string, arch string) error {

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED|ole.COINIT_SPEED_OVER_MEMORY)
	defer ole.CoUninitialize()

	oleShellObject, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		return err
	}
	defer oleShellObject.Release()
	wshell, err := oleShellObject.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return err
	}
	defer wshell.Release()
	cs, err := oleutil.CallMethod(wshell, "CreateShortcut", "dst")
	idispatch := cs.ToIDispatch()
	oleutil.PutProperty(idispatch, "TargetPath", "src")
	oleutil.CallMethod(idispatch, "Save")
	return nil
}

func unzipLocalfile(unzip2Dir string, file *os.File, goos string) (bool, error) {

	// fileInfo, err := file.Stat()
	// if err != nil {
	// 	return false, err
	// }

	// var unzipDir string
	// if goos == "windows" {
	// 	unzipDir = unzip2Dir + fileInfo.Name() + "/"
	// } else {
	// 	unzipDir = unzip2Dir + fileInfo.Name() + "\\"
	// }
	rc, err := zip.OpenReader(file.Name())
	if err != nil {
		return false, err
	}
	defer rc.Close()

	for i, f := range rc.Reader.File {
		fmt.Println(i, f)
		frc, err := f.Open()
		if err != nil {
			return false, err
		}
		defer frc.Close()

		fpath := filepath.Join(unzip2Dir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(fpath), f.Mode())
			lf, err := os.Create(fpath)
			if err != nil {
				return false, err
			}
			defer lf.Close()

			_, err = io.Copy(lf, frc)
			if err != nil {
				return false, err
			}
		}

	}

	return true, nil
}

func download(url string, dir string, app string, goos string) (*os.File, error) {
	resp, err := http.Get(url + app)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var file *os.File
	if goos == "windows" {
		file, err = os.Create(dir + "\\" + app)
	} else {
		file, err = os.Create(dir + "/" + app)
	}

	if err != nil {
		return nil, err
	}
	// defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func showApps(apps []string) (index int, result string, err error) {

	// templates := promptui.SelectTemplates{
	// 	Active:   `👉  {{ .Title | cyan | bold }}`,
	// 	Inactive: `   {{ .Title | cyan }}`,
	// 	Selected: `{{ "✔" | green | bold }} {{ "Recipe" | bold }}: {{ .Title | cyan }}`,
	// }

	list := promptui.Select{
		Label: "请选择要安装的程序",
		Items: apps,
		// Templates: &templates,
	}
	index, result, err = list.Run()
	if err != nil {
		return 0, "", err
	}
	return index, result, nil
}

func listApps(url string) (apps []string, err error) {

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	con, err := ioutil.ReadAll(resp.Body)

	// var rb []byte
	// _, err = resp.Body.Read(rb)

	if err != nil {
		return nil, err
	}
	if len(con) == 0 {
		return nil, errors.New("程序列表为空")
	}

	apps = strings.Split(string(con), ";")

	return apps, nil
}

func areYouReady(url string) (index int, result string, err error) {
	propmt := promptui.Select{
		Label: "是否从 " + url + " 安装程序？",
		Items: []string{"是", "否"},
	}
	index, result, err = propmt.Run()
	if err != nil {
		return 0, "", err
	}
	return index, result, nil

}

func main() {

	operationSystem := runtime.GOOS
	architecture := runtime.GOARCH
	// url := "http://172.42.1.221:9090"
	url := "http://localhost:9090"
	url4ListApps := url + "/applist"
	url4Download := url + "/download?app="
	var download2Dir string
	var unzip2Dir string
	if operationSystem == "windows" {
		download2Dir = "c:\\"
		unzip2Dir = "c:\\"
	} else {
		download2Dir = "/Users/damao/"
		unzip2Dir = "/Users/damao/"
	}

	_, result, err := areYouReady(url)
	checkErrAtMainFunc(err)
	if result == "否" {
		os.Exit(0)
	}

	var apps []string
	// apps = make([]string, 20)
	apps, err = listApps(url4ListApps)
	checkErrAtMainFunc(err)
	if apps == nil {
		os.Exit(0)
	}

	_, result, err = showApps(apps)
	checkErrAtMainFunc(err)

	file, err := download(url4Download, download2Dir, result, operationSystem)
	checkErrAtMainFunc(err)
	defer file.Close()

	_, err = unzipLocalfile(unzip2Dir, file, operationSystem)
	checkErrAtMainFunc(err)

	err = genShortcut(operationSystem, architecture)

}
