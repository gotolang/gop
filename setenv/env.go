package main

import (
	"log"
	"os"
	"runtime"

	"golang.org/x/sys/windows/registry"
)

func main() {
	log.Println("set windows env")
	os.Setenv("P", "Golang") //session environment

	if runtime.GOOS == "windows" {
		//HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment
		k, err := registry.OpenKey(registry.LOCAL_MACHINE, "SYSTEM\\CurrentControlSet\\Control\\Session Manager\\Environment", registry.QUERY_VALUE)
		if err != nil {
			log.Println("registry open SYSTEM\\CurrentControlSet\\Control\\Session Manager\\Environment ")
		}
		defer k.Close()

		err = k.SetStringValue("ORACLE_HOME", "C:\\oracle\\product\\10.2")
		if err != nil {
			log.Println("set value error")
		}
	} else if runtime.GOOS == "darwin" { // macos

	} else { //linux

	}

}
