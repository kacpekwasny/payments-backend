package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	cmt "github.com/kacpekwasny/commontools"
)

var CONFIG_MAP = LoadConfig()

func LoadConfig() map[string]string {
	var conf = map[string]string{}
	f, err := os.Open(os.Args[1])
	cmt.Pc(err)
	defer f.Close()
	bytes, err := ioutil.ReadAll(f)
	cmt.Pc(err)
	err = json.Unmarshal(bytes, &conf)
	cmt.Pc(err)
	checkConfigKeys(conf)
	return conf
}

func checkConfigKeys(m map[string]string) {
	// panic when a key is lacking
	var required_keys = []string{"listen address"}
	for _, k := range required_keys {
		if _, ok := m[k]; !ok {
			cmt.Pc(errors.New("missing key: " + k))
		}
	}
}
