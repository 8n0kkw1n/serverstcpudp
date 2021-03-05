package tcpsrv

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"

	conf "servers/config"
)

// NewTCPClientStart - connect new TCP client
func NewTCPClientStart(cf conf.TCPConfig) error {
	c, err := net.Dial("tcp", fmt.Sprintf("%s:%d", cf.Host, cf.Port))
	if err != nil {
		return fmt.Errorf("Failed to listen: %v", err)
	}

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		fmt.Fprintf(c, text+"\n")

		message, _ := bufio.NewReader(c).ReadString('\n')
		fmt.Print("->: " + message)

		if strings.TrimSpace(string(text)) == "STOP" {
			fmt.Println("TCP client exiting...")
			return nil
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
