package cli

import (
	"flag"
	"fmt"
	"runtime"

	"github.com/Gunyoung-Kim/blockchain/explorer"
	"github.com/Gunyoung-Kim/blockchain/rest"
)

func usage() {
	fmt.Printf("Please use the following flags:\n\n")
	fmt.Printf("-port: 	Set the port of the server\n")
	fmt.Printf("-mode: 	Choose between 'html' and 'rest' or 'both'\n")
	runtime.Goexit() // for execute defer in main
}

//Start CLI
func Start() {
	port := flag.Int("port", 4000, "Set Port of this server ")
	mode := flag.String("mode", "rest", "Choose between 'html' and 'rest' or 'both'")

	flag.Parse()

	switch *mode {
	case "html":
		explorer.Start(*port)
	case "rest":
		rest.Start(*port)
	case "both":
		go rest.Start(*port)
		explorer.Start(*port + 1)
	default:
		usage()
	}

	flag.Parse()
}
