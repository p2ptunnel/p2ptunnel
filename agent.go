package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"net"
	"strconv"
)

var (
	revLookup   map[string]string
	forwardPort int
	forwardBuf  []byte
)

func agent(ctx *cli.Context) error {
	conf, err := readConf(ctx.GlobalString("conf"))
	if err != nil {
		return err
	}
	if len(ctx.Args()) != 1 {
		return errors.New("Please provide forwarding port number")
	}

	forwardBuf = make([]byte, 1000000)
	forwardPort, err = strconv.Atoi(ctx.Args()[0])
	if err != nil {
		return err
	}

	// Setup reverse lookup hash map for authentication.
	revLookup = make(map[string]string, len(conf.Peers))
	// Setup Peer Table for Quick Packet --> Dest ID lookup
	peerTable := make(map[string]peer.ID)
	for name, p := range conf.Peers {
		revLookup[p.ID] = name
		peerTable[name], err = peer.Decode(p.ID)
		if err != nil {
			return err
		}
	}

	// Setup System Context
	cctx, cancel := context.WithCancel(context.Background())

	fmt.Println("[+] Creating LibP2P Node")

	// Create P2P Node
	fmt.Printf("My ID: %s\n", conf.ID)
	host, _, err := CreateNode(
		cctx,
		conf.PrivateKey,
		ctx.Uint("port"),
		streamHandlerAgent,
	)
	if err != nil {
		return err
	}

	// Register the application to listen for SIGINT/SIGTERM
	go signalExit(cancel, host)

	<-cctx.Done()
	return nil
}

func streamHandlerAgent(stream network.Stream) {
	// If the remote node ID isn't in the list of known nodes don't respond.
	fmt.Printf("stream handle: %s from %+v\n", stream.Conn().RemotePeer().Pretty(), revLookup)
	if _, ok := revLookup[stream.Conn().RemotePeer().Pretty()]; !ok {
		fmt.Println("not found, need reset")
		stream.Reset()
		return
	}
	var requestSize = make([]byte, 2)
	// Read the incoming packet's size as a binary value.
	_, err := stream.Read(requestSize)
	if err != nil {
		stream.Close()
		return
	}

	// Decode the incoming packet's size from binary.
	size := binary.LittleEndian.Uint16(requestSize)

	//fmt.Printf("got %d byte, length: %d\n", requestSize, size)

	// Read in the packet until completion.
	var request = make([]byte, size)
	var plen uint16 = 0
	for plen < size {
		tmp, err := stream.Read(request[plen:size])
		plen += uint16(tmp)
		if err != nil {
			stream.Close()
			fmt.Printf("read error: %v", err)
			return
		}
		//fmt.Printf("read %d byte, %s\n", tmp, string(request))
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "localhost:"+strconv.Itoa(forwardPort))
	if err != nil {
		fmt.Printf("resolve local service tcp:%d : %s", forwardPort, err)
		return
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	for err != nil {
		fmt.Println(err)
		return
	}
	l, err := conn.Write(request[:size])
	if err != nil {
		fmt.Println(err)
		return
	}
	if l != int(size) {
		fmt.Printf("forward expect write %d bytes, actually %d bytes\n", size, l)
		return
	}

	for {
		rl, err := conn.Read(forwardBuf)
		if err != nil {
			fmt.Println(err)
			return
		}
		if rl == 0 {
			fmt.Println("reach EOF")
			break
		}
		fmt.Printf("read feedback %s\n", string(forwardBuf[:rl]))
		wl, err := stream.Write(forwardBuf[:rl])
		if err != nil {
			fmt.Println(err)
			continue
		}
		if wl != rl {
			fmt.Printf("reply expect write %d bytes, actually %d bytes\n", l, wl)
			return
		}
		if rl < len(forwardBuf) {
			// reach end
			break
		}
	}
	stream.Write(nil)
	if err = conn.Close(); err != nil {
		fmt.Printf("close forwarding connection got: %v\n", err)
	}
}
