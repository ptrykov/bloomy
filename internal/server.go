package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gansidui/gotcp"
	filter "github.com/ptrykov/bloomy/pkg"
	bf "github.com/ptrykov/bloomy/pkg/bloom_filters"
)

type Server struct {
	// config
	// map string -> filter
	config  *ServerConfig
	filters map[string]filter.Filter
}

func NewServer(cfg *ServerConfig) *Server {
	return &Server{
		config:  cfg,
		filters: make(map[string]filter.Filter),
	}
}

func (s *Server) CreateFilter(name string, size uint) (filter.Filter, error) {
	if _, ok := s.filters[name]; ok != true {
		s.filters[name] = bf.NewCounting(size)
	}
	return s.filters[name], nil
}

func (s *Server) loadFilters() {

}

func (s *Server) DeleteFilter(name string) {
	delete(s.filters, name)
}

func (s *Server) Add(name string, value *[]byte) {
	s.filters[name].Add(value)
}

func (s *Server) Test(name string, value *[]byte) bool {
	return s.filters[name].Test(value)
}

func (s *Server) Remove(name string, value *[]byte) bool {
	return s.filters[name].Remove(value)
}

func (s *Server) Run() {
	s.CreateFilter("users", 5)
	val := []byte("simone")
	s.Add("users", &val)
	val = []byte("pavel")
	s.Add("users", &val)
	filter := bf.Load(s.filters["users"].Dump())
	fmt.Println("include pavel?", filter.Test(&val))
	val = []byte("simone")
	fmt.Println("include simone?", filter.Test(&val))
	val = []byte("sebastian")
	fmt.Println("include sebastian?", filter.Test(&val))

	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":"+s.config.Port)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	config := &gotcp.Config{
		PacketReceiveChanLimit: s.config.Channels,
		PacketSendChanLimit:    s.config.Channels,
	}

	srv := gotcp.NewServer(config, s, &BloomyProtocol{})

	go srv.Start(listener, time.Second)
	fmt.Println("listening:", listener.Addr())

	// catches system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)

	srv.Stop()
}

func (sc *Server) OnConnect(c *gotcp.Conn) bool {
	addr := c.GetRawConn().RemoteAddr()
	c.PutExtraData(addr)
	fmt.Println("OnConnect:", addr)
	c.AsyncWritePacket(NewBloomyPacketOut(0, []byte("Welcome to bloomy")), 0)
	return true
}

func (sc *Server) OnMessage(c *gotcp.Conn, p gotcp.Packet) bool {
	packet := p.(*BloomyPacket)
	command := packet.CollectionName
	commandType := packet.ApiCode
	optional1Bytes := []byte(packet.Optional1)
	switch commandType {
	case 1:
		size, err := strconv.ParseUint(packet.Optional1, 10, 32)
		checkError(err)
		fmt.Println("size ", size)
		sc.CreateFilter(packet.CollectionName, uint(size))
		c.AsyncWritePacket(NewBloomyPacketOut(1, []byte{0}), 0)
	case 2:
		c.AsyncWritePacket(NewBloomyPacketOut(2, []byte("-1")), 0)
	case 3:
		c.AsyncWritePacket(NewBloomyPacketOut(3, []byte("-1")), 0)
	case 4:
		sc.Add(packet.CollectionName, &optional1Bytes)
		c.AsyncWritePacket(NewBloomyPacketOut(4, []byte{0}), 0)
	case 5:
		if sc.Test(packet.CollectionName, &optional1Bytes) {
			c.AsyncWritePacket(NewBloomyPacketOut(5, []byte("1")), 0)
		} else {
			c.AsyncWritePacket(NewBloomyPacketOut(5, []byte("0")), 0)
		}
	case 6:
		if sc.Remove(packet.CollectionName, &optional1Bytes) {
			c.AsyncWritePacket(NewBloomyPacketOut(6, []byte("1")), 0)
		} else {
			c.AsyncWritePacket(NewBloomyPacketOut(6, []byte("0")), 0)
		}
	default:
		c.AsyncWritePacket(NewBloomyPacketOut(3, []byte(command)), 0)
	}

	return true
}

func (sc *Server) OnClose(c *gotcp.Conn) {
	fmt.Println("OnClose:", c.GetExtraData())
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
