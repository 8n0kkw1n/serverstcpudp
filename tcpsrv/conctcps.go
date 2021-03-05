package tcpsrv

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	conf "servers/config"
	"strings"
	"time"
)

func handleTCPConnect(c net.Conn) error {
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			return fmt.Errorf("Failed reader: %v", err)
		}

		temp := strings.TrimSpace(string(netData))
		if temp == "STOP" {
			fmt.Println("Client exiting TCP server!")
			break
		}

		fmt.Println("-> ", temp)
		t := time.Now()
		myTime := t.Format(time.RFC3339) + "\n"
		c.Write([]byte(myTime))
	}

	c.Close()

	return nil
}

// NewTCPServerStart - Start new TCP server
func NewTCPServerStart(cf conf.TCPConfig) error {
	l, err := net.Listen("tcp4", fmt.Sprintf("%s:%d", cf.Host, cf.Port))
	if err != nil {
		return fmt.Errorf("Failed to listen: %v", err)
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			return fmt.Errorf("Failed accept: %v", err)
		}

		go handleTCPConnect(c)

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
