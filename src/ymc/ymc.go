package ymc

import (
	"bufio"
	"flag"
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

var err error

func CheckErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func (c *Connection) ParseConfig(path string) {
	home := os.Getenv("HOME")
	file, err := os.Open(home + "/" + path)
	CheckErr(err)
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
	address := c.Host + ":" + c.Port
	if c.Password == "" {
		c.Client, err = mpd.Dial("tcp", address)
	} else {
		c.Client, err = mpd.DialAuthenticated("tcp", address, c.Password)
	}
	CheckErr(err)
	fmt.Println("connected to", address)
}
func (c *Connection) GetStatus() (status map[string]string) {
	status, err := c.Client.Status()
	CheckErr(err)
	return
}
func (c *Connection) PrintStatus() {
	status := c.GetStatus()
	song, err := c.Client.CurrentSong()
	CheckErr(err)
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
func (c *Connection) Toggle() {
	status := c.GetStatus()
	if status["state"] == "play" {
		err = c.Client.Pause(true)
		CheckErr(err)
		return
	}
	err = c.Client.Pause(false)
	CheckErr(err)
}
func (c *Connection) Repeat() {
	status := c.GetStatus()
	if status["repeat"] == "0" {
		err = c.Client.Repeat(true)
		CheckErr(err)
		return
	}
	err = c.Client.Repeat(false)
	CheckErr(err)
}
func (c *Connection) Random() {
	status := c.GetStatus()
	if status["random"] == "0" {
		err = c.Client.Random(true)
		CheckErr(err)
		return
	}
	err = c.Client.Random(false)
	CheckErr(err)
}
func (c *Connection) ChangeVolume() {
	volume, err := strconv.Atoi(flag.Arg(1))
	CheckErr(err)
	if volume > 100 || volume < 0 {
		fmt.Println("volume range is between [0-100]")
		return
	}
	err = c.Client.SetVolume(volume)
	CheckErr(err)
}
func (c *Connection) UpdateDatabase() {
	dir := flag.Arg(1)
	jobId, err := c.Client.Update(dir)
	CheckErr(err)
	fmt.Println("updating Database #" + string(jobId))
}
