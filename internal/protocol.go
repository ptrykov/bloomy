package server

import (
	"bytes"
	"fmt"
	"net"

	"encoding/binary"
	"github.com/gansidui/gotcp"
)

// API_VERSION 4 bytes

// API_CODE 4 bytes
// 1 create filter (name, size)
// 2 load filter (name)
// 3 delete filter (name)
// 4 insert (name, element)
// 5 test (name, element)
// 6 remove (name, element)

// Collection String

// Optional 1 String

var (
	endTag   = []byte("\r\n") // Telnet command's end tag
	paramEnd = []byte(" \t ")
)

// Packet
type BloomyPacket struct {
	ApiVersion     uint32
	ApiCode        uint32
	CollectionName string
	Optional1      string
	Optional2      string
}

func (p *BloomyPacket) Serialize() []byte {
	buf := []byte{0}
	return buf
}

func parseParam(data []byte, end []byte) (param string, count int) {
	index := bytes.Index(data, end)
	if index > -1 {
		return string(data[:index]), index
	}
	return "", 0
}

func processParams(data []byte, end []byte) (collection string, optional1 string, optional2 string) {
	var index, cSize, o1Size = 0, 0, 0

	collection, cSize = parseParam(data, end)
	index += cSize + 3 // Skip paramEnd

	optional1, o1Size = parseParam(data[index:], end)
	index += o1Size + 3 // Skip paramEnd

	optional2, _ = parseParam(data[index:], end)

	return
}

func NewBloomyPacket(data []byte) *BloomyPacket {
	apiVersion := binary.LittleEndian.Uint32(data[:4])
	// fmt.Println("api version: ", apiVersion)

	apiCode := binary.LittleEndian.Uint32(data[4:8])
	// fmt.Println("api code: ", apiCode)

	collection, param1, param2 := processParams(data[8:], paramEnd)
	fmt.Println("params collection ", collection, "param1 ", param1, "param2", param2)
	return &BloomyPacket{
		ApiVersion:     apiVersion,
		ApiCode:        apiCode,
		CollectionName: collection,
		Optional1:      param1,
		Optional2:      param2,
	}
}

// Packet Out
type BloomyPacketOut struct {
	Type byte
	Data []byte
}

func (p *BloomyPacketOut) Serialize() []byte {
	buf := make([]byte, 1)
	buf[0] = p.Type
	buf = append(buf, p.Data...)
	buf = append(buf, endTag...)
	return buf
}

func NewBloomyPacketOut(answerType byte, data []byte) *BloomyPacketOut {
	return &BloomyPacketOut{
		Type: answerType,
		Data: data,
	}
}

type BloomyProtocol struct {
}

func (this *BloomyProtocol) ReadPacket(conn *net.TCPConn) (gotcp.Packet, error) {
	fullBuf := bytes.NewBuffer([]byte{})
	for {
		data := make([]byte, 1024)

		readLen, err := conn.Read(data)
		if err != nil { //EOF, or worse
			return nil, err
		}

		if readLen == 0 { // Connection maybe closed by the client
			return nil, gotcp.ErrConnClosing
		} else {
			fullBuf.Write(data[:readLen])

			index := bytes.Index(fullBuf.Bytes(), endTag)
			if index > -1 {
				command := fullBuf.Next(index)
				fullBuf.Next(2) // skipping endTag size
				fmt.Println("command")
				fmt.Println(string(command))
				return NewBloomyPacket(command), nil
			}
		}
	}
}
