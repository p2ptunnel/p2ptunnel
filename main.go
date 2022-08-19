package main

import (
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/p2ptunnel/p2ptunnel/pkg/httplogger"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
)

const (
	defaultConf          = "./conf/p2ptunnel.yml"
	defaultAgentPort     = 8011
	defaultConnectorPort = 8012
)

var (
	verbose bool
	logger  *httplogger.HTTPLogger
)

func main() {
	app := cli.NewApp()

	app.Usage = "p2p tunnel"
	app.Email = "dwebfan@gmail.com"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "conf, c",
			Usage: "config file path",
			Value: defaultConf,
		},
		cli.BoolFlag{
			Name:  "verbose, v",
			Usage: "print more debug message",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:      "init",
			Usage:     "user friendly name of agent or connection",
			ArgsUsage: "[name]",
			Aliases:   []string{"i"},
			Action:    initConf,
		},
		{
			Name:      "add",
			Usage:     "add peer name and its ID",
			ArgsUsage: "[peer name] [peer id]",
			Aliases:   []string{"a"},
			Action:    addPeer,
		},
		{
			Name:      "remove",
			Usage:     "remove peer name and its ID",
			ArgsUsage: "[peer name] [peer id]",
			Aliases:   []string{"a"},
			Action:    removePeer,
		},
		{
			Name:      "agent",
			Usage:     "start p2p tunnel agent service",
			Action:    agent,
			ArgsUsage: "[forward port]",
		},
		{
			Name:   "connector",
			Usage:  "start p2p tunnel connector service",
			Action: connector,
			Flags: []cli.Flag{
				cli.UintFlag{
					Name:  "port, p",
					Usage: "connector's listening port",
					Value: defaultConnectorPort,
				},
			},
		},
	}
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func initConf(ctx *cli.Context) error {
	if len(ctx.Args()) != 1 {
		return errors.New("Please provide one user friendly name for the agent or connector")
	}
	// Create New Libp2p Node
	host, err := libp2p.New()
	if err != nil {
		return err
	}

	// Get Node's Private Key
	keyBytes, err := crypto.MarshalPrivateKey(host.Peerstore().PrivKey(host.ID()))
	if err != nil {
		return err
	}

	// Setup an initial default command.
	conf := &Config{
		Name:       ctx.Args()[0],
		ID:         host.ID().Pretty(),
		PrivateKey: string(keyBytes),
	}

	configFile := ctx.GlobalString("conf")
	err = os.MkdirAll(filepath.Dir(configFile), os.ModePerm)
	if err != nil {
		return err
	}

	f, err := os.Create(configFile)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write out config to file.
	err = yaml.NewEncoder(f).Encode(conf)
	if err != nil {
		return err
	}

	// Print config creation message to user
	fmt.Printf("Initialized new config at %s\n", configFile)
	fmt.Printf("Please remember your ID: %s\n", conf.ID)
	return nil
}

func addPeer(ctx *cli.Context) error {
	if len(ctx.Args()) != 2 {
		return errors.New("Please provide both peer name and peer ID")
	}

	conf := &Config{}
	configFile := ctx.GlobalString("conf")
	f, err := os.OpenFile(configFile, os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	err = yaml.NewDecoder(f).Decode(conf)
	if err != nil {
		return err
	}

	p, ok := conf.Peers[ctx.Args()[0]]
	if ok {
		return errors.Errorf("Peer %s has been added with ID %s", ctx.Args()[0], p.ID)
	}
	conf.Peers[ctx.Args()[0]] = Peer{ID: ctx.Args()[1]}

	// Write out config to file.
	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	err = yaml.NewEncoder(f).Encode(conf)
	if err != nil {
		return err
	}

	fmt.Printf("%s - %s has been saved in config file: %s\n", ctx.Args()[0], ctx.Args()[1], configFile)
	return nil
}

func removePeer(ctx *cli.Context) error {
	if len(ctx.Args()) != 2 {
		return errors.New("Please provide both peer name and peer ID")
	}

	configFile := ctx.GlobalString("conf")
	conf, err := readConf(configFile)

	p, ok := conf.Peers[ctx.Args()[0]]
	if !ok {
		return errors.Errorf("Peer %s is not in config file %s", ctx.Args()[0], configFile)
	}
	if p.ID != ctx.Args()[1] {
		return errors.Errorf("To-be-removed Peer ID %s is different from input %s", p.ID, ctx.Args()[1])
	}

	delete(conf.Peers, ctx.Args()[0])

	f, err := os.OpenFile(configFile, os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write out config to file.
	err = yaml.NewEncoder(f).Encode(conf)
	if err != nil {
		return err
	}

	fmt.Printf("%s has been removed from config file %s\n", ctx.Args()[0], configFile)
	return nil
}
