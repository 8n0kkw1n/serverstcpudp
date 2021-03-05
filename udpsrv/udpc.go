package udpsrv

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	conf "servers/config"
	"strings"
)

// NewUDPClientStart - connect new TCP client
func NewUDPClientStart(cf conf.UDPConfig) error {
	s, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", cf.Host, cf.Port))
	c, err := net.DialUDP("udp4", nil, s)
	if err != nil {
		return fmt.Errorf("Failed to listen: %v", err)
	}

	fmt.Printf("The UDP server is %s\n", c.RemoteAddr().String())
	defer c.Close()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		data := []byte(text + "\n")
		_, err = c.Write(data)
		if strings.TrimSpace(string(data)) == "STOP" {
			fmt.Println("Exiting UDP client!")
			return nil
		}

		if err != nil {
			return fmt.Errorf("Failed to written: %v", err)
		}

		buffer := make([]byte, 1024)
		n, _, err := c.ReadFromUDP(buffer)
		if err != nil {
			return fmt.Errorf("Failed to listen: %v", err)
		}
		fmt.Printf("Reply: %s\n", string(buffer[0:n]))

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
