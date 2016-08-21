package models

type Sample struct {
	Id     int     `json:"id"`
	Key    string  `json:"key"`
	Labels string  `json:"labels,omitempty"`
	Time   int64   `json:"time"`
	Value  float64 `json:"value"`
}

type Bucket struct {
	Count int     `json:"count"`
	Key   string  `json:"key"`
	Start int64   `json:"start"`
	End   int64   `json:"end"`
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Avg   float64 `json:"avg"`
}

type Samples []Sample
type Buckets []Bucket
