package socket

import (
	"bufio"
	"encoding/binary"
	//"errors"
	"io"
	"net"
	"sync"

	"github.com/MessageDream/webIM/modules/log"
)

type Server struct {
	scanners       []*bufio.Scanner
	listeners      []*net.TCPListener
	connections    []net.Conn
	wait           sync.WaitGroup
	parser         func([]byte) string
	lastError      error
	OnConnected    func(conn net.Conn)
	OnDisconnected func(conn net.Conn)
	OnMessage      func(msg string, conn net.Conn)
}

//NewServer returns a new Server
func NewServer() *Server {
	server := new(Server)

	return server
}

//Sets the parser
func (self *Server) SetParser(function func([]byte) string) {
	self.parser = function
}

// //Sets  format
// func (self *Server) SetFormat(format Format) {
// 	self.format = format
// }

// //Sets the handler, this halder with receive every syslog entry
// func (self *Server) SetHandler(handler string) {
// 	self.handler = handler
// }

//Configure the server for listen on an UDP addr
func (self *Server) ListenUDP(addr string) error {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}

	connection, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}

	self.connections = append(self.connections, connection)
	return nil
}

//Configure the server for listen on an unix socket
func (self *Server) ListenUnixgram(addr string) error {
	unixAddr, err := net.ResolveUnixAddr("unixgram", addr)
	if err != nil {
		return err
	}

	connection, err := net.ListenUnixgram("unixgram", unixAddr)
	if err != nil {
		return err
	}

	self.connections = append(self.connections, connection)
	return nil
}

//Configure the server for listen on a TCP addr
func (self *Server) ListenTCP(addr string) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}

	self.listeners = append(self.listeners, listener)
	return nil
}

//Starts the server, all the go routines goes to live
func (self *Server) Boot() error {

	for _, listerner := range self.listeners {
		self.goAcceptConnection(listerner)
	}

	for _, connection := range self.connections {
		self.goScanConnection(connection)
	}

	return nil
}

func (self *Server) goAcceptConnection(listerner *net.TCPListener) {
	self.wait.Add(1)
	go func(listerner *net.TCPListener) {
		for {
			connection, err := listerner.Accept()
			if err != nil {
				log.Error("Accept connection error %v", err)
				continue
			}
			log.Warn("connection %s is connected", connection.RemoteAddr())
			if self.OnConnected != nil {
				self.OnConnected(connection)
			}
			self.goReadConnection(connection)
		}

		self.wait.Done()
	}(listerner)
}

func (self *Server) goScanConnection(connection net.Conn) {
	scanner := bufio.NewScanner(connection)
	self.scanners = append(self.scanners, scanner)
	self.wait.Add(1)
	go self.scan(scanner)
}

func (self *Server) scan(scanner *bufio.Scanner) {
	for scanner.Scan() {
		if self.parser != nil {
			self.parser([]byte(scanner.Text()))
		}
	}

	self.wait.Done()
}

func (self *Server) goReadConnection(connection net.Conn) {
	self.wait.Add(1)
	go self.read(connection)
}

func (self *Server) read(connection net.Conn) {

	// close connection
	defer func() {
		if self.OnDisconnected != nil {
			self.OnDisconnected(connection)
		}
		if connection != nil {
			connection.Close()
		}
		self.wait.Done()
	}()

	isHeadLoaded := false
	bodyLen := 0
	reader := bufio.NewReader(connection)

Out:
	for {
		if !isHeadLoaded {
			headLenSl := make([]byte, 4)
			log.Info("Ready for reading package head")

			readedHeadLen := 0

			for readedHeadLen < 4 {
				len, err := reader.Read(headLenSl)
				if err == io.EOF {
					log.Warn("Client exit: %s\n", connection.RemoteAddr())
					return
				}
				if err != nil {
					log.Warn("Read error: %s\n", err)
					continue
				}
				readedHeadLen += len
			}
			bodyLen = int(binary.BigEndian.Uint32(headLenSl))
			log.Info("Package head read success, body length is : %d", bodyLen)
			isHeadLoaded = true
		}

		if isHeadLoaded {
			log.Info("Ready for reading package body")
			bodySl := make([]byte, bodyLen)

			readedBodyLen := 0

			for readedBodyLen < bodyLen {
				len, err := reader.Read(bodySl)
				if err == io.EOF {
					log.Warn("Client exit: %s\n", connection.RemoteAddr())
					return
				}
				if err != nil {
					log.Warn("Read body error: %s\n", err)
					break Out
				}
				readedBodyLen += len
			}
			log.Info("Read body success")

			var msg string = ""
			if self.parser != nil {
				msg = self.parser(bodySl)
			} else {
				msg = string(bodySl)
			}

			if self.OnConnected != nil {
				self.OnMessage(msg, connection)
			}
			isHeadLoaded = false
		}
	}

}

// func (self *Server) getParserRFC3164(line []byte) *rfc3164.Parser {
// 	parser := rfc3164.NewParser(line)

// 	return parser
// }

// func (self *Server) getParserRFC5424(line []byte) *rfc5424.Parser {
// 	parser := rfc5424.NewParser(line)

// 	return parser
// }

//Returns the last error
func (self *Server) GetLastError() error {
	return self.lastError
}

//Kill the server
func (self *Server) Kill() error {
	for _, connection := range self.connections {
		err := connection.Close()
		if err != nil {
			return err
		}
	}

	for _, listener := range self.listeners {
		err := listener.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

//Waits until the server stops
func (self *Server) Wait() {
	self.wait.Wait()
}
