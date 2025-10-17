package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/haykm86/foreverstor/p2p"
)

func makeServer(listenAddr string, nodes ...string) *FileServer {
	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    listenAddr,
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		// TODO onPeer func
	}

	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)

	fileServerOpts := FileServerOpts{
		EncKey:            newEncryptionKey(),
		StorageRoot:       listenAddr + "_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
		BootstrapNodes:    nodes,
	}

	s := NewFileServer(fileServerOpts)
	tcpTransport.OnPeer = s.OnPeer
	return s
}

func main() {
	s1 := makeServer(":33033", "")
	s2 := makeServer(":44044", ":33033")
	s3 := makeServer(":55055", ":44044", ":33033")

	go func() {
		log.Fatal(s1.Start())
	}()

	time.Sleep(2 * time.Second)
	go func() {
		log.Fatal(s2.Start())
	}()

	time.Sleep(2 * time.Second)
	go s3.Start()
	time.Sleep(2 * time.Second)
	key := "coolpicture.jpg"
	data := bytes.NewReader([]byte("my big data file here!"))
	s3.Store(key, data)
	s3.store.Delete(s3.ID, key)

	time.Sleep(5 * time.Second)

	r, err := s3.Get("coolpicture.jpg")
	if err != nil {
		log.Fatal(err)
	}

	b, err := io.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nIt's main here for you: ", string(b))
	select {}
}
