package chart

import (
	"time"
)

type Record struct {
	Total string
	Date  time.Time
}

type OutputChart struct {
	Label       string
	Description string
	Result      []Record
}
