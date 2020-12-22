package ymc

import (
	"bufio"
	"fmt"
	"github.com/fhs/gompd/mpd"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Connection struct {
	Host     string
	Port     string
	Password string
	Client   *mpd.Client
}

func (c *Connection) ParseConfig(path string) {
	home := os.Getenv("HOME")
	file, err := os.Open(home + "/" + path)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	parameters := []string{"host", "port", "password"}
	for scanner.Scan() {
		for _, parameter := range parameters {
			value := scanner.Text()
			match, _ := regexp.MatchString(parameter+"=", value)
			if match == true {
				value = strings.TrimPrefix(value, parameter+"=")
				switch value {
				case "host":
					c.Host = value
				case "port":
					c.Port = value
				case "password":
					c.Password = value
				}
			}
		}
	}
}
func (c *Connection) Connect() {
	var err error
	address := c.Host + ":" + c.Port
	if c.Password == "" {
		c.Client, err = mpd.Dial("tcp", address)
	} else {
		c.Client, err = mpd.DialAuthenticated("tcp", address, c.Password)
	}
	if err != nil {
		log.Fatalln("connection error", err)
	}
	fmt.Println("connected to", address)
	return
}
func (c *Connection) GetStatus() (status map[string]string) {
	status, err := c.Client.Status()
	if err != nil {
		log.Fatalln(err)
	}
	return
}
func (c *Connection) PrintStatus() {
	status := c.GetStatus()
	song, err := c.Client.CurrentSong()
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
