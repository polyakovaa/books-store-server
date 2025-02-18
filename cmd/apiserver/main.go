package main

import (
	"flag"
	"log"

	_ "net/http/pprof"

	"github.com/BurntSushi/toml"
	"github.com/polyakovaa/standartserver3/internal/app/apiserver"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file (.toml or .env file)")
}

func main() {
	flag.Parse()
	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Println("can not find path to config, app will use default confs:", err)
	}

	//server instance
	s := apiserver.New(config)

	//server start
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}

}
