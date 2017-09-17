package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tenmozes/sanchez/pkg/pump"
	"github.com/tenmozes/sanchez/pkg/server"
	"github.com/tenmozes/sanchez/pkg/spirit"
	"gopkg.in/yaml.v2"
	"periph.io/x/periph/host"
)

var pumpsPath = flag.String("p", "pumps.yml", "available pumps config with connected spirits")
var cocktailsPath = flag.String("c", "cocktails.yml", "list of cocktails")

type pumpConfig struct {
	Spirit  string `yaml:"spirit"`
	PumpPin int    `yaml:"pump"`
}

const port = 8080

func main() {
	flag.Parse()
	host.Init()
	pumps := []pumpConfig{}
	if err := readYaml(*pumpsPath, &pumps); err != nil {
		log.Fatalf("can't read pumps config %v", err)
	}
	cocktails := []spirit.Cocktail{}
	if err := readYaml(*cocktailsPath, &cocktails); err != nil {
		log.Fatalf("can't read pumps config %v", err)
	}
	workers := make([]pump.Worker, 0, len(pumps))
	for _, p := range pumps {
		workers = append(workers, pump.NewPump(p.Spirit, p.PumpPin))
	}
	spiritService, err := spirit.NewService(workers, cocktails)
	if err != nil {
		log.Fatalf("can't initialize spirit service %v", err)
	}
	srv := server.NewServer(spiritService)
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		srv.Run(port)
	}()
	log.Printf("got signal %v", <-s)
	srv.Close()
}

func readYaml(path string, config interface{}) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, config)
}
