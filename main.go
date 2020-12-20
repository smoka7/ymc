package main

import (
	"flag"
	"log"
	"strconv"
	"ymc/src/ymc"
)

func main() {
	var host, port, password *string
	var err error
	host = flag.String("h", "localhost", "host of mpd server")
	port = flag.String("p", "6601", "port of mpd server")
	password = flag.String("P", "", "password of mpd server")
	flag.Parse()
	command := flag.Arg(0)
	m := ymc.Server{
		Host:     *host,
		Port:     *port,
		Password: *password,
	}
	conn := m.Connect()
	defer conn.Close()
	switch command {
	case "play":
		err = conn.Pause(false)
		if err != nil {
			log.Println(err)
		}
	case "pause":
		err = conn.Pause(true)
		if err != nil {
			log.Println(err)
		}
	case "next":
		err = conn.Next()
		if err != nil {
			log.Println(err)
		}
	case "prev":
		err = conn.Previous()
		if err != nil {
			log.Println(err)
		}
	case "toggle":
		status := ymc.GetStatus(conn)
		if status["state"] == "play" {
			err = conn.Pause(true)
		} else {
			err = conn.Pause(false)
		}
		if err != nil {
			log.Println(err)
		}
	case "random":
		status := ymc.GetStatus(conn)
		if status["random"] == "0" {
			err = conn.Random(true)
		} else {
			err = conn.Random(false)
		}
		if err != nil {
			log.Println(err)
		}
	case "repeat":
		status := ymc.GetStatus(conn)
		if status["repeat"] == "0" {
			err = conn.Repeat(true)
		} else {
			err = conn.Repeat(false)
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
			err = conn.SetVolume(volume)
			if err != nil {
				log.Println(err)
			}
		}

	default:
		break
	}
	ymc.PrintStatus(conn)
}
