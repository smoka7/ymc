package ymc

import (
	"fmt"
	"github.com/fhs/gompd/mpd"
	"log"
	"math"
	"strconv"
)

type Server struct {
	Host     string
	Port     string
	Password string
}

func (s *Server) Connect() (conn *mpd.Client) {
	var err error
	address := s.Host + ":" + s.Port
	if s.Password == "" {
		conn, err = mpd.Dial("tcp", address)
	} else {
		conn, err = mpd.DialAuthenticated("tcp", address, s.Password)
	}
	if err != nil {
		log.Fatalln("connection error", err)
	}
	log.Println("connected to", address)
	return
}
func GetStatus(conn *mpd.Client) (status map[string]string) {
	status, err := conn.Status()
	if err != nil {
		log.Fatalln(err)
	}
	return
}
func PrintStatus(conn *mpd.Client) {
	status := GetStatus(conn)
	song, err := conn.CurrentSong()
	if err != nil {
		log.Fatalln(err)
	}
	if status["state"] == "play" || status["state"] == "pause" {
		el, _ := strconv.ParseFloat(status["elapsed"], 64)
		du, _ := strconv.ParseFloat(status["duration"], 64)
		m := math.Round(du / float64(60))
		s := math.Round(math.Mod(du, float64(60)))
		percentage := int(el / du * 100)
		fmt.Println(song["Artist"], "-", song["Title"])
		fmt.Printf("state:%s %g:%g %d%% \n", status["state"], m, s, percentage)
	}
	fmt.Println("volume:", status["volume"], "random:", status["random"], "repeat:", status["repeat"], "consume:", status["consume"])
}
