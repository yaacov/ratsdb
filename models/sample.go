package models

type Sample struct {
	Id     int     `json:"id"`
	Key    string  `json:"key"`
	Labels string  `json:"labels,omitempty"`
	Time   int64   `json:"time"`
	Value  float64 `json:"value"`
}

type Samples []Sample
