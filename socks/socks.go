package socks

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	socks5 "github.com/nodauf/Go-RouterSocks/go-socks5"
	router "github.com/nodauf/Go-RouterSocks/router"
)

var serverSocks5 *socks5.Server

func StartSocks(ip string, port int) {
	address := ip + ":" + strconv.Itoa(port)
	errorMsg := make(chan error)
	go listenAndAccept(address, errorMsg)
	status := <-errorMsg
	if status != nil {
		log.Fatalln(status)
	}
	log.Println("[*] Server socks server on " + address)

}

func listenAndAccept(address string, status chan error) {
	var err error
	serverSocks5, err = socks5.New(&socks5.Config{})
	if err != nil {
		status <- err
	}
	ln, err := net.Listen("tcp", address)
	if err != nil {
		status <- err
	}
	status <- nil
	for {
		conn, err := ln.Accept()
		//log.Println("Got a client")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Errors accepting!")
		}
		//log.Println("Passing off to socks5")
		go func() {
			//firstBytes, secondBytes, dest, err := serverSocks5.GetDest(conn)
			firstBytes, secondBytes, dest, request, err := serverSocks5.GetDest(conn)
			if err != nil {
				// The errors are print to stdout by go-socks5
				//log.Println(err)
				conn.Close()
				return
			}
			_, remoteSocks := router.GetRoute(dest)
			if remoteSocks != "" {
				err := connectToSocks(firstBytes, secondBytes, conn, remoteSocks)
				// If the socks server is no longer available, we have the error conneciton refused
				if err != nil && strings.Contains(err.Error(), "connection refused") {
					// The route is no longer valid and we delete it
					fmt.Println("Remote socks " + remoteSocks + " server no longer available. Got connection refused")
					//router.DeleteRoutes(network)
				}
			} else {
				fmt.Println("\n[-] Unkown route for " + dest + " using direct connection without proxy")
				serverSocks5.Handle(request, conn)
			}
		}()
	}
}

func connectToSocks(firstBytes []byte, secondBytes []byte, src net.Conn, remoteSocks string) error {

	var proxy net.Conn
	//log.Println("Connecting to remote socks")
	proxy, err := net.Dial("tcp", remoteSocks)
	if err != nil {
		log.Println(err)
		return err
	}
	defer src.Close()
	defer proxy.Close()
	// Send first request
	proxy.Write(firstBytes)
	// Empty the buffer
	buf := make([]byte, 100)
	proxy.Read(buf)
	// Send second request
	proxy.Write(secondBytes)

	chanToRemote := streamCopy(proxy, src)
	chanToStdout := streamCopy(src, proxy)
	select {
	case <-chanToStdout:
		//log.Println("Remote connection is closed")
	case <-chanToRemote:
		//log.Println("Local program is terminated")
	}
	return nil
}

// Performs copy operation between streams: os and tcp streams
func streamCopy(src io.Reader, dst io.Writer) <-chan int {
	buf := make([]byte, 1024)
	syncChannel := make(chan int)
	go func() {
		defer func() {
			if con, ok := dst.(net.Conn); ok {
				con.Close()
				//log.Printf("Connection from %v is closed\n", con.RemoteAddr())
			}
			syncChannel <- 0 // Notify that processing is finished
		}()
		for {

			var nBytes int
			var err error
			nBytes, err = src.Read(buf)

			if err != nil {
				if err != io.EOF {
					//log.Printf("Read error: %s\n", err)
				}
				break
			}
			_, err = dst.Write(buf[0:nBytes])
			if err != nil {
				//log.Fatalf("Write error: %s\n", err)
			}
		}
	}()
	return syncChannel
}
