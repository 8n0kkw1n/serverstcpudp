package main

import (
	"fmt"
	"log"
	"os"

	conf "servers/config"
	tcpsrv "servers/tcpsrv"
	udpsrv "servers/udpsrv"
)

func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Println("Please provide protocol and mode.")
		return
	}

	cnf, err := conf.NewConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Can't read the config: %s", err)
	}

	switch {
	case args[1] == "TCP":
		if args[2] == "server" {
			if err := tcpsrv.NewTCPServerStart(cnf.TCP); err != nil {
				log.Fatalf("Couldn't launch TCP server: %s", err)
			}
		} else if args[2] == "client" {
			if err := tcpsrv.NewTCPClientStart(cnf.TCP); err != nil {
				log.Fatalf("Couldn't connect TCP server: %s", err)
			}
		} else {
			log.Println("Undefinted")
		}
	case args[1] == "UDP":
		if args[2] == "server" {
			if err := udpsrv.NewUDPServerStart(cnf.UDP); err != nil {
				log.Fatalf("Couldn't launch UDP server: %s", err)
			}
		} else if args[2] == "client" {
			if err := udpsrv.NewUDPClientStart(cnf.UDP); err != nil {
				log.Fatalf("Couldn't connect UDP server: %s", err)
			}
		} else {
			log.Println("Undefinted")
		}

	default:
		fmt.Println("No information available for that day.")
	}
}
