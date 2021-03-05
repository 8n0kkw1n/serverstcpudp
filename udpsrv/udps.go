package udpsrv

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	conf "servers/config"
)

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

// NewUDPServerStart - Start new UDP server
func NewUDPServerStart(cf conf.UDPConfig) error {
	s, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", cf.Host, cf.Port))
	if err != nil {
		return fmt.Errorf("Failed resolve udp address: %v", err)
	}

	connection, err := net.ListenUDP("udp4", s)
	if err != nil {
		return fmt.Errorf("Failed to listen: %v", err)
	}

	defer connection.Close()
	buffer := make([]byte, 1024)
	rand.Seed(time.Now().Unix())

	for {
		n, addr, err := connection.ReadFromUDP(buffer)
		fmt.Print("-> ", string(buffer[0:n-1]))

		if strings.TrimSpace(string(buffer[0:n])) == "STOP" {
			fmt.Println("Exiting UDP server!")
			return nil
		}

		data := []byte(strconv.Itoa(random(1, 1001)))
		fmt.Printf("data: %s\n", string(data))
		_, err = connection.WriteToUDP(data, addr)
		if err != nil {
			return fmt.Errorf("Failed to write: %v", err)
		}

		ch := make(chan os.Signal)
		signal.Notify(ch, os.Interrupt)

		go func() {
			select {
			case sig := <-ch:
				fmt.Printf("Got %s signal. Aborting...\n", sig)
				os.Exit(1)
			}
		}()
	}
}
