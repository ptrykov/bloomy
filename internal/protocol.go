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
}

func (p *BloomyPacket) Serialize() []byte {
	buf := []byte{0}
	return buf
}

func NewBloomyPacket(data []byte) *BloomyPacket {
	apiVersion := binary.LittleEndian.Uint32(data[:4])
	// fmt.Println("api version: ", apiVersion)

	apiCode := binary.LittleEndian.Uint32(data[4:8])
	// fmt.Println("api code: ", apiCode)

	tailBuf := data[8:]
	collectionEnd := bytes.Index(tailBuf, paramEnd)
	collection := string(tailBuf[:collectionEnd])
	// fmt.Println("collection: ", collection)

	tailBuf = tailBuf[collectionEnd+3:] // Skipping paramEnd
	optional1End := bytes.Index(tailBuf, paramEnd)
	optional1 := string(tailBuf[:optional1End])

	return &BloomyPacket{
		ApiVersion:     apiVersion,
		ApiCode:        apiCode,
		CollectionName: collection,
		Optional1:      optional1,
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

type BloomyCallback struct {
}

func (this *BloomyCallback) OnConnect(c *gotcp.Conn) bool {
	addr := c.GetRawConn().RemoteAddr()
	c.PutExtraData(addr)
	fmt.Println("OnConnect:", addr)
	c.AsyncWritePacket(NewBloomyPacketOut(0, []byte("Welcome to bloomy")), 0)
	return true
}

func (this *BloomyCallback) OnMessage(c *gotcp.Conn, p gotcp.Packet) bool {
	packet := p.(*BloomyPacket)
	command := packet.CollectionName
	commandType := packet.ApiCode
	fmt.Println("OnMessage: ", command, " ", commandType)

	switch commandType {
	case 1:
		c.AsyncWritePacket(NewBloomyPacketOut(1, []byte{0}), 0)
	case 2:
		c.AsyncWritePacket(NewBloomyPacketOut(2, []byte("2")), 0)
	case 3:
		return false
	case 4:
		return true
	case 5:
		return true
	case 6:
		return true
	default:
		c.AsyncWritePacket(NewBloomyPacketOut(3, []byte(command)), 0)
	}

	return true
}

func (this *BloomyCallback) OnClose(c *gotcp.Conn) {
	fmt.Println("OnClose:", c.GetExtraData())
}
