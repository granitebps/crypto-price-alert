package types

import "time"

type Alert struct {
	Email     string    `json:"email"`
	Ticker    string    `json:"ticker"`
	Pair      string    `json:"pair"`
	Price     uint      `json:"price"`
	EmailedAt time.Time `json:"emailed_at"`
}

type IndodaxResponse map[string]map[string]interface{}
