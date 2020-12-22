package main

import (
	"flag"
	"log"
	"strconv"
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
		if err != nil {
			log.Println(err)
		}
	case "pause":
		err = m.Client.Pause(true)
		if err != nil {
			log.Println(err)
		}
	case "next":
		err = m.Client.Next()
		if err != nil {
			log.Println(err)
		}
	case "prev":
		err = m.Client.Previous()
		if err != nil {
			log.Println(err)
		}
	case "toggle":
		status := m.GetStatus()
		if status["state"] == "play" {
			err = m.Client.Pause(true)
		} else {
			err = m.Client.Pause(false)
		}
		if err != nil {
			log.Println(err)
		}
	case "random":
		status := m.GetStatus()
		if status["random"] == "0" {
			err = m.Client.Random(true)
		} else {
			err = m.Client.Random(false)
		}
		if err != nil {
			log.Println(err)
		}
	case "repeat":
		status := m.GetStatus()
		if status["repeat"] == "0" {
			err = m.Client.Repeat(true)
		} else {
			err = m.Client.Repeat(false)
		}
		if err != nil {
			log.Println(err)
		}
	case "volume":
		volume, err := strconv.Atoi(flag.Arg(1))
		if err != nil {
			log.Println(err)
		}
		if volume > 100 || volume < 0 {
			log.Println("volume range is between [0-100]")
		} else {
			err = m.Client.SetVolume(volume)
			if err != nil {
				log.Println(err)
			}
		}
	default:
		log.Println("arg isn't valid")
	}
	m.PrintStatus()
}
