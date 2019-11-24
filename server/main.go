package server

import (
	"github.com/dashotv/flame"
	"github.com/dashotv/mercury"

	"github.com/nats-io/nats.go"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type Server struct {
	URL     string
	merc    *mercury.Mercury
	channel chan *flame.Response
	client  *flame.Client
}

func New(URL string) (*Server, error) {
	var err error
	s := &Server{URL: URL}

	s.merc, err = mercury.New("blaze", nats.DefaultURL)
	if err != nil {
		return nil, err
	}

	s.channel = make(chan *flame.Response, 5)
	if err := s.merc.Sender("blaze", s.channel); err != nil {
		return nil, err
	}

	s.client = flame.NewClient(URL)

	return s, nil
}

func (s *Server) Start() error {
	logrus.Info("starting blaze...")

	c := cron.New(cron.WithSeconds())
	if _, err := c.AddFunc("* * * * * *", s.Sender); err != nil {
		return err
	}
	c.Start()

	for {
		select {}
	}

	return nil
}

func (s *Server) Sender() {
	logrus.WithField("prefix", "genma").Info("sending message")

	resp, err := s.client.List()
	if err != nil {
		logrus.Errorf("flame list error: %s", err)
		return
	}

	s.channel <- resp
}
