package main

import (
	"fmt"
	"net/http"
	"os"
)

import (
	"./certs"
	"./hashes"
	"./libraries"
	"./obfuscation"
	"./traffic"
)

func main() {
	if len(os.Args) < 3 {
		help()
		os.Exit(1)
	}
	switch mode := os.Args[1]; mode {
	case "certs":
		certs.Routes()
	case "hashes":
		hashes.Routes()
	case "libraries":
		libraries.Routes()
	case "obfuscation":
		obfuscation.Routes()
	case "traffic":
		traffic.Routes()
	default:
		os.Exit(0)
	}
	port := fmt.Sprintf(":%v", os.Args[2])
	fmt.Println("Listening in", port)
	http.ListenAndServe(port, nil)
}

func help() {
	fmt.Println("Usage:")
	fmt.Println("\tcontrol <mode> <port>")
}
