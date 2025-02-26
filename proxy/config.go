package main

import (
	"github.com/BurntSushi/toml"
	"github.com/je4/bremote/common"
	"log"
)

type Config struct {
	Logfile         string
	Loglevel        string
	InstanceName    string
	CertPEM         string
	KeyPEM          string
	CaPEM           string
	TLSAddr         string
	WebRoot         string
	KVDBFile        string
	NTPHost         string
	RuntimeInterval common.Duration
}

func LoadConfig(filepath string) Config {
	var conf Config
	_, err := toml.DecodeFile(filepath, &conf)
	if err != nil {
		log.Fatalln("Error on loading config: ", err)
	}
	return conf
}
