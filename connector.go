package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

func connector(ctx *cli.Context) error {
	conf, err := readConf(ctx.GlobalString("conf"))
	if err != nil {
		return err
	}

	switch len(conf.Peers) {
	case 0:
		return errors.New("Remote agent ID is not found, please add firstly")
	case 1:
	default:
		return errors.New("Multiple agents are added. Only support 1 for now")
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
	host, dht, err := CreateNode(
		cctx,
		conf.PrivateKey,
		uint(forwardPort),
		streamHandlerConnector,
	)
	if err != nil {
		return err
	}

	// Setup P2P Discovery
	go Discover(cctx, host, dht, peerTable)
	go prettyDiscovery(cctx, host, peerTable)

	// Register the application to listen for SIGINT/SIGTERM
	go signalExit(cancel, host)

	localAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", ctx.Uint("port")))
	if err != nil {
		return err
	}
	l, err := net.ListenTCP("tcp", localAddr)
	if err != nil {
		return err
	}
	defer l.Close()
	for {
		select {
		case <-cctx.Done():
			return cctx.Err()
		default:
			// Wait for a connection.
			if err := l.SetDeadline(time.Now().Add(time.Second)); err != nil {
				// this is an error for registering timeout with SetDeadline()
				return err
			}
			conn, err := l.Accept()
			if err != nil {
				// if it is due to out timeout expiration we will continue
				if os.IsTimeout(err) {
					continue
				}
				return err
			}
			// Handle the connection in a new goroutine.
			// The loop then returns to accepting, so that
			// multiple connections may be served concurrently.
			go func(ctx context.Context, c net.Conn) {
				// Echo all incoming data.
				buf := make([]byte, 1000000)
				size, err := c.Read(buf)
				if err != nil {
					fmt.Println(err)
				}
				body := make([]byte, 2)
				binary.LittleEndian.PutUint16(body, uint16(size))
				body = append(body, buf[:size]...)

				err = sendToRemote(ctx, host, peerTable, body, c)
				if err != nil {
					fmt.Println(err)
				}
				// Shut down the connection.
				c.Close()
			}(cctx, conn)
		}
	}

	return nil
}

func sendToRemote(ctx context.Context, node host.Host, peerTable map[string]peer.ID, body []byte, local io.Writer) error {
	fmt.Printf("remote table: %+v\n", peerTable)
retry:
	for name, id := range peerTable {
		fmt.Printf("discover peer: %s - %v\n", name, id)
		stream, err := node.NewStream(ctx, id, Protocol)
		if err != nil {
			if strings.HasPrefix(err.Error(), "failed to dial") ||
				strings.HasPrefix(err.Error(), "no addresses") {
				// Attempt to connect to peers slowly when they aren't found.
				fmt.Println(err)
				time.Sleep(5 * time.Second)
				goto retry
			} else {
				return err
			}
		}
		fmt.Printf("[+] Connection to %s Successful. Network Ready.\n", name)
		l, err := stream.Write(body)
		if err != nil {
			return err
		}
		if l <= 0 {
			fmt.Printf("write wrong size payload: %d\n", l)
		}

		/*
			var requestSize = make([]byte, 2)
			// Read the incoming packet's size as a binary value.
			_, err = stream.Read(requestSize)
			if err != nil {
				stream.Close()
				return errors.Wrapf(err, "read reply")
			}
		*/

		// Decode the incoming packet's size from binary.
		//size := binary.LittleEndian.Uint16(requestSize)
		size := 100000

		reply := make([]byte, size)
		r, err := stream.Read(reply)
		if err != nil {
			return err
		}
		if r == 0 {
			return errors.Errorf("read empty data")
		}
		fmt.Printf("read reply %d: %s\n", r, string(reply))
		w, err := local.Write(reply)
		if err != nil {
			return err
		}
		if int(size) != w {
			return errors.Errorf("expect reply write %d, actual write %d\n", size, w)
		}
		err = stream.Close()
		if err != nil {
			return err
		}

		// only support one peer
		break
	}
	return nil
}

func streamHandlerConnector(stream network.Stream) {
	// If the remote node ID isn't in the list of known nodes don't respond.
	if _, ok := revLookup[stream.Conn().RemotePeer().Pretty()]; !ok {
		stream.Reset()
		return
	}
	var packet = make([]byte, 1000000)
	var packetSize = make([]byte, 2)
	for {
		// Read the incoming packet's size as a binary value.
		_, err := stream.Read(packetSize)
		if err != nil {
			stream.Close()
			return
		}

		// Decode the incoming packet's size from binary.
		size := binary.LittleEndian.Uint16(packetSize)

		// Read in the packet until completion.
		var plen uint16 = 0
		for plen < size {
			tmp, err := stream.Read(packet[plen:size])
			plen += uint16(tmp)
			if err != nil {
				stream.Close()
				return
			}
		}
		fmt.Printf("read %d from %s: %s\n", size, stream.ID(), string(packet[:size]))
	}
}
