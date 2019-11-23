package main

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/dashotv/flame"
	"github.com/dashotv/mercury"

	"github.com/nats-io/nats.go"
)

type TorrentsByIndex []flame.Torrent

func (a TorrentsByIndex) Len() int           { return len(a) }
func (a TorrentsByIndex) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a TorrentsByIndex) Less(i, j int) bool { return a[i].Queue < a[j].Queue }

func main() {
	m, err := mercury.New("blaze", nats.DefaultURL)
	if err != nil {
		panic(err)
	}

	fmt.Println("starting receiver...")
	channel := make(chan *flame.Response, 5)
	m.Receiver("blaze", channel)

	for {
		select {
		case r := <-channel:
			//fmt.Printf("received: %#v\n", r)
			logrus.Infof("received message")
			sort.Sort(TorrentsByIndex(r.Torrents))
			for _, t := range r.Torrents {
				logrus.Infof("%3.0f %6.2f%% %10.2fmb %8.8s %s\n", t.Queue, t.Progress, t.SizeMb(), t.State, t.Name)
			}
		case <-time.After(30 * time.Second):
			fmt.Println("timeout")
			os.Exit(0)
		}
	}
}
