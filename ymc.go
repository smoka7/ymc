package main

import (
	"flag"
	"log"
	"os"
	"ymc/src/ymc"
)

func main() {
	var host, port, password, config *string
	var err error
	host = flag.String("h", "localhost", "host of mpd server")
	config = flag.String("c", "", "path of config file")
	port = flag.String("p", "6601", "port of mpd server")
	password = flag.String("P", "", "password of mpd server")
	flag.Parse()
	m := ymc.Connection{
		Host:     *host,
		Port:     *port,
		Password: *password,
	}
	if *config != "" {
		m.ParseConfig(*config)
	}
	m.Connect()
	defer m.Client.Close()
	switch flag.Arg(0) {
	case "":
		break
	case "play":
		err = m.Client.Pause(false)
		ymc.CheckErr(err)
	case "pause":
		err = m.Client.Pause(true)
		ymc.CheckErr(err)
	case "next":
		err = m.Client.Next()
		ymc.CheckErr(err)
	case "prev":
		err = m.Client.Previous()
		ymc.CheckErr(err)
	case "toggle":
		m.Toggle()
	case "random":
		m.Random()
	case "repeat":
		m.Repeat()
	case "volume":
		m.ChangeVolume()
	case "update":
		m.UpdateDatabase()
	case "stop":
		m.Stop()
	case "shuffle":
		m.Shuffle()
	case "seek":
		m.Seek()
	case "single":
		m.Single()
	case "crossfade":
		m.Crossfade()
	case "outputs":
		m.ListOutputs()
		os.Exit(0)
	case "disable":
		m.DisableOutput()
		os.Exit(0)
	case "enable":
		m.EnableOutput()
		os.Exit(0)
	default:
		log.Println("arg isn't valid")
	}
	m.PrintStatus()
}
