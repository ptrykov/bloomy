package server

import (
	"bytes"
	"fmt"
	"net"

	"encoding/binary"
	"github.com/gansidui/gotcp"
)

const (
	apiVersionBytes = 4
	apiCodeBytes    = 1
)

// API_VERSION 4 bytes

// API_CODE 1 byte
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
	ApiCode        byte
	CollectionName string
	Optional1      string
}

func (p *BloomyPacket) Serialize() []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf[0:], p.ApiVersion)
	buf = append(buf, p.ApiCode)
	buf = append(buf, p.CollectionName...)
	buf = append(buf, paramEnd...)
	buf = append(buf, p.Optional1...)
	buf = append(buf, paramEnd...)
	buf = append(buf, endTag...)
	return buf
}

func NewBloomyPacket(data []byte) *BloomyPacket {
	apiVersion := data[:4]
	apiCode := data[4]
	collectionEnd := bytes.Index(data[5:], paramEnd)
	collection := data[5:collectionEnd]
	optional1End := bytes.Index(data[collectionEnd+3:], paramEnd)
	optional1 := data[collectionEnd+3 : optional1End]
	return &BloomyPacket{
		ApiVersion:     binary.LittleEndian.Uint32(apiVersion),
		ApiCode:        apiCode,
		CollectionName: string(collection),
		Optional1:      string(optional1),
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

		readLengh, err := conn.Read(data)

		if err != nil { //EOF, or worse
			return nil, err
		}

		if readLengh == 0 { // Connection maybe closed by the client
			return nil, gotcp.ErrConnClosing
		} else {
			fullBuf.Write(data[:readLengh])

			index := bytes.Index(fullBuf.Bytes(), endTag)
			if index > -1 {
				command := fullBuf.Next(index)
				fullBuf.Next(2) // skipping endTag size
				//fmt.Println(string(command))
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
		c.AsyncWritePacket(NewBloomyPacketOut(commandType, []byte(command)), 0)
	}

	return true
}

func (this *BloomyCallback) OnClose(c *gotcp.Conn) {
	fmt.Println("OnClose:", c.GetExtraData())
}
