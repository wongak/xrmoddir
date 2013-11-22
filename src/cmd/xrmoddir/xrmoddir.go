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
	_, err := loadConfig(*cfgFile)
	if err != nil {
		log.Fatal(err)
	}
}
