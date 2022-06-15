package model

import (
	"time"
)

type Auction struct {
	UUID  string    `msgpack:"uuid" bson:"uuid"`
	Start time.Time `msgpack:"start" bson:"start"`
}

type AuctionResponse struct {
	Count int64     `json:"count"`
	From  time.Time `json:"from"`
	To    time.Time `json:"to"`
}
