package models

type Bucket struct {
	Count int     `json:"count"`
	Key   string  `json:"key"`
	Start int64   `json:"start"`
	End   int64   `json:"end"`
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Avg   float64 `json:"avg"`
}

type Buckets []Bucket
