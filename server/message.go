package server

import "github.com/dashotv/flame"

type Message struct {
	Response flame.Response
	Sender   string
	Sent     string
}
