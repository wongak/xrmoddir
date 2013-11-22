// X Rebirth Mod Directory
//
// Server
package main

import (
	"flag"
	"log"
)

var cfgFile = flag.String("c", "cfg/xrmoddir.cfg.json", "Path to config file.")

func main() {
	flag.Parse()
	cfg, err := loadConfig(*cfgFile)
	if err != nil {
		log.Fatal(err)
	}
	_, err = connectDb(cfg)
	if err != nil {
		log.Fatal(err)
	}
}
