package main

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/p2ptunnel/p2ptunnel/pkg/httplogger"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"io"
	"log"
	"net"
	"strconv"
)

var (
	revLookup   map[string]string
	forwardPort int
)

func agent(ctx *cli.Context) error {
	conf, err := readConf(ctx.GlobalString("conf"))
	if err != nil {
		return err
	}
	if len(ctx.Args()) != 1 {
		return errors.New("Please provide forwarding port number")
	}

	verbose = ctx.GlobalBool("verbose")
	if verbose {
		logger = httplogger.New(nil)
	}

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
	defer cancel()

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
		if err := stream.Reset(); err != nil {
			log.Printf("while reset stream: %v", err)
		}
		return
	}
	/*
		var requestSize = make([]byte, 2)
		// Read the incoming packet's size as a binary value.
		_, err := stream.Read(requestSize)
		if err != nil {
			stream.Close()
			return
		}

		// Decode the incoming packet's size from binary.
		size := binary.LittleEndian.Uint16(requestSize)
	*/

	// TODO: use persistent connection
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

	defer func() {
		if err = conn.Close(); err != nil {
			fmt.Printf("close forwarding connection got: %v\n", err)
		}
	}()

	//fmt.Printf("got %d byte, length: %d\n", requestSize, size)

	// Read in the packet until completion.
	var buffer = make([]byte, 1024)
	defer func() {
		// free memory explicitly
		buffer = nil
	}()

	if verbose {
		logger.Reset()
		fmt.Println("--->")
	}

	for {
		readSize, err := stream.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("(EOF)")
				return
			}
			stream.Close()
			fmt.Printf("read error: %v", err)
			return
		}

		if verbose {
			logger.Print(buffer)
		}

		writeSize, err := conn.Write(buffer[:readSize])
		if err != nil {
			fmt.Println(err)
			return
		}
		if writeSize != readSize {
			fmt.Printf("forward expect write %d bytes, actually %d bytes\n", readSize, writeSize)
			return
		}

		if readSize < 1024 {
			break
		}
	}

	if verbose {
		logger.Reset()
		fmt.Println("<---")
	}

	readSize := 0
	for readSize < len(buffer) {
		readSize, err = conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("reach EOF of read reply")
			} else {
				fmt.Printf("failed to read reply: %s\n", err)
			}
			return
		}
		if readSize == 0 {
			fmt.Println("read empty data")
			break
		}
		if verbose {
			logger.Print(buffer[:readSize])
		}
		writeSize, err := stream.Write(buffer[:readSize])
		if err != nil {
			fmt.Printf("failed to write reply to peer: %s\n", err)
			return
		}
		if writeSize != readSize {
			fmt.Printf("reply expect write %d bytes, actually %d bytes\n", readSize, writeSize)
			return
		}
	}
}
