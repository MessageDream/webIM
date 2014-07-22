package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	//"time"
)

func Log(v ...interface{}) {
	fmt.Println(v...)
}

func test(err error, mesg string) {
	if err != nil {
		Log("CLIENT: ERROR: ", mesg)
		os.Exit(-1)
	}
}

func read(con net.Conn) (map[string]interface{}, error) {
	buf := make([]byte, 4096)
	data := make(map[string]interface{})
	n, err := con.Read(buf)
	if err == io.EOF {
		con.Close()
		wait.Done()
		data["err"] = "remote server closed!"
		return data, err

	}
	if err != nil {
		con.Close()
		wait.Done()
		data["err"] = "Error in reading!"
		return data, err
	}

	json.Unmarshal(buf[0:n], &data)
	return data, nil
}

func clientsender(cn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	for {
		input, _ := reader.ReadBytes('\n')
		if string(input) == "/quit\n" {
			cn.Write([]byte("/quit"))
			break
		}
		datalen := len(input) - 1
		input = append(intToBytes(datalen), input[0:len(input)-1]...)
		cn.Write(input[0:len(input)])
	}
	wait.Done()
}

func clientreceiver(name string, cn net.Conn) {
	for {
		data, err := read(cn)
		if err != nil {
			fmt.Printf("remote error %s\n", data["err"])
			break
		} else {
			user := data["User"].(map[string]interface{})
			rname := user["Name"].(string)
			content := data["Content"].(string)
			eventType := data["Type"].(float64)
			if eventType == 0 || eventType == 1 {
				if rname == name {
					rname = "You"
				}
				var typestr string
				if eventType == 0 {
					typestr = "joined"
				} else {
					typestr = "left"
				}
				fmt.Printf("%s %s the chat room\n", rname, typestr)

			} else if rname != name {
				fmt.Printf("%s> %s\n", rname, content)
			}
		}
	}
	wait.Done()
}

func intToBytes(nNum int) []byte {
	bytesRet := make([]byte, 4)
	bytesRet[0] = (byte)((nNum >> 24) & 0xFF)
	bytesRet[1] = (byte)((nNum >> 16) & 0xFF)
	bytesRet[2] = (byte)((nNum >> 8) & 0xFF)
	bytesRet[3] = (byte)(nNum & 0xFF)
	return bytesRet
}

func intToByte(i int) []byte {
	abyte0 := make([]byte, 4)
	abyte0[0] = (byte)(0xff & i)
	abyte0[1] = (byte)((0xff00 & i) >> 8)
	abyte0[2] = (byte)((0xff0000 & i) >> 16)
	abyte0[3] = (byte)((0xff000000 & i) >> 24)
	return abyte0
}

var wait sync.WaitGroup

func main() {
	Log("client start ")

	destination := "127.0.0.1:3001"
	Log("connecting to ", destination)
	cn, err := net.Dial("tcp", destination)
	test(err, "dialing")
	defer cn.Close()
	Log("server connected ")

	// get the user name
	fmt.Print("Please give you name: ")
	reader := bufio.NewReader(os.Stdin)
	name, _ := reader.ReadBytes('\n')

	form := make(map[string]string, 2)
	form["uname"] = string(name[0 : len(name)-1])
	fmt.Print("Please give the room ID: ")
	roomid, _ := reader.ReadBytes('\n')
	form["roomid"] = string(roomid[0 : len(roomid)-1])
	data, err := json.Marshal(form)
	test(err, "json.Marshal")
	datalen := len(data)
	bytelen := intToBytes(datalen)
	data = append(bytelen, data...)

	cn.Write(data)

	wait = sync.WaitGroup{}

	Log("start receiver")
	wait.Add(1)
	go clientreceiver(string(name), cn)
	Log("start sender")
	wait.Add(1)
	go clientsender(cn)

	wait.Wait()
	Log("client stoped")
}
