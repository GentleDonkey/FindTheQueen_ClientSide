package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

var (
	connected     bool
	connectedSync sync.Mutex
)

func main() {
	fmt.Println("Client started...")
	for {
		connectedSync.Lock()
		alreadyConnected := connected
		connectedSync.Unlock()
		if !alreadyConnected {
			conn, err := net.Dial("tcp", "127.0.0.1:7621")
			if err != nil {
				fmt.Println(err.Error())
				time.Sleep(time.Duration(5) * time.Second)
				continue
			}
			fmt.Println(conn.RemoteAddr().String() + ": connected")
			connectedSync.Lock()
			connected = true
			connectedSync.Unlock()
			go receiveData(conn)
			authentication(conn)
			sendNumber(conn)
		}
		time.Sleep(time.Duration(5) * time.Second)
	}
}

func authentication(conn net.Conn){
	fmt.Println("Please input your user name")
	reader := bufio.NewReader(os.Stdin)
	user, _ := reader.ReadString('\n')
	_, err := fmt.Fprintf(conn, user)
	if err != nil {
		fmt.Println(conn.RemoteAddr().String() + ": end sending data")
	}
	fmt.Println("Please input your password")
	pwd, _ := reader.ReadString('\n')
	_, err = fmt.Fprintf(conn, pwd)
	if err != nil {
		fmt.Println(conn.RemoteAddr().String() + ": end sending data")
	}
	fmt.Println("Please wait...")
}
func sendNumber(conn net.Conn){
	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		_, err := fmt.Fprintf(conn, text)
		if err != nil {
			fmt.Println(conn.RemoteAddr().String() + ": end sending data")
		}
		fmt.Println("Please wait...")
	}
}

func receiveData(conn net.Conn) {
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(conn.RemoteAddr().String() + ": disconnected")
			conn.Close()
			connectedSync.Lock()
			connected = false
			connectedSync.Unlock()
			fmt.Println(conn.RemoteAddr().String() + ": end receiving data")
			return
		}
		fmt.Print(conn.RemoteAddr().String() + ": " + message)
	}
}