package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/burntsushi/toml"
)

type apps struct {
	Apps map[string]app
}

type app struct {
	Chinese string
	Zip     string
	Dir     string
	Ini     string
	Exe     string
	Desktop bool
}

func decodeTOML(tomlfile string) (map[string]app, error) {
	var applications apps

	_, err := toml.DecodeFile(tomlfile, &applications)
	if err != nil {
		return nil, err
	}

	// fmt.Println(applications.Apps)
	return applications.Apps, nil
}
func listApps() ([]byte, error) {

	fmt.Println("=====================start=======================")

	apps, err := decodeTOML("apps.toml")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	appsWithJSON, err := json.Marshal(apps)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	fmt.Println(string(appsWithJSON))
	fmt.Println("====================end=========================")
	return appsWithJSON, nil
}

func showApps(appWithJSON []byte) (map[string]app, error) {
	var apps map[string]app
	err := json.Unmarshal(appWithJSON, &apps)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// fmt.Println(apps)
	return apps, nil
}

func main() {
	appsWithJSON, err := listApps()
	if err != nil {
		log.Println(err)
		return
	}
	appsWithMaps, err := showApps(appsWithJSON)
	if err != nil {
		log.Println(err)
		return
	}
	for key := range appsWithMaps {
		// fmt.Println(key, value)
		fmt.Println(appsWithMaps[key].Chinese)
	}

}
