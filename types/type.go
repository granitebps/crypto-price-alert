package types

import "time"

type Alert struct {
	Email     string    `json:"email"`
	Ticker    string    `json:"ticker"`
	Pair      string    `json:"pair"`
	Price     uint      `json:"price"`
	EmailedAt time.Time `json:"emailed_at"`
}

type Alias map[string]map[string]string

type IndodaxResponse map[string]map[string]interface{}

type CoingeckoResponse map[string]map[string]float64
