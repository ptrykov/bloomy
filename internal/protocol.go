package protocol

import (
	"bytes"
	"fmt"
	"net"
	"strings"

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
	copy(buf, p.ApiVersion)
	buf = append(buf, p.ApiCode)
	buf = append(buf, p.CollectionName...)
	buf = append(buf, p.paramEnd...)
	buf = append(buf, p.Optional1...)
	buf = append(buf, p.paramEnd...)
	buf = append(buf, endTag...)
	return buf
}

func NewBloomyPacket(data []byte) *BloomyPacket {
	apiVersion := data[:4]
	apiCode := data[4:5]
	collectionEnd := bytes.Index(data[5:], paramEnd)
	collection := data[5:collectionEnd]
	optional1End := bytes.Index(data[collectionEnd+3:], paramEnd)
	optional1 := data[collectionEnd+3 : optional1End]
	return &BloomyPacket{
		ApiVersion:     apiVersion,
		ApiCode:        apiCode,
		CollectionName: collectionEnd,
		Optional1:      optional1,
	}
}

// Packet Out
type BloomyPacketOut struct {
	Type uint32
	Data []byte
}

func (p *BloomyPacketOut) Serialize() []byte {
	buf := make([]byte, 4)
	copy(buf, p.Type)
	buf = append(buf, p.Data...)
	buf = append(buf, endTag...)
	return buf
}

func NewBloomyPacketOut(answerType uint32, data []byte) *BloomyPacket {
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
				return NewBloomyPacket(command)
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
	packet := p.(*BlooomyPacket)
	command := packet.GetData()
	commandType := packet.ApiCode()

	switch commandType {
	case 1:
		c.AsyncWritePacket(NewBloomyPacketOut(1, []byte(true)), 0)
	case 2:
		c.AsyncWritePacket(NewBloomyPacketOut(2, []byte("2")), 0)
	case 3:
		return false
	case 4:
		return
	case 5:
		return
	case 6:
		return
	default:
		c.AsyncWritePacket(NewBloomyPacketOut(commandType, []byte("unknow command")), 0)
	}

	return true
}

func (this *BloomyCallback) OnClose(c *gotcp.Conn) {
	fmt.Println("OnClose:", c.GetExtraData())
}
