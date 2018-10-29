package main

import (
	"flag"
	"fmt"
	sc "github.com/BraindeadBZH/science_cursor_api"
	"os"
)

func main() {
	config := flag.String("config", "/etc/science_cursor/api.conf.json", "Path to the configuration file")
	flag.Parse()

	err := sc.Run(*config)
	if err != nil {
		fmt.Println("Server exited with an error:", err.Error())
		os.Exit(-1)
	} else {
		os.Exit(0)
	}
}
