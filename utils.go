package chart

import (
	"fmt"
	"hash/fnv"
	"math/rand"
	"time"

	"github.com/jinzhu/now"
	"github.com/moisespsena-go/error-wrap"
)

// range specification, note that min <= max
type IntRange struct {
	min, max int
}

// get next random value within the interval including min and max
func (ir *IntRange) NextRandom(r *rand.Rand) int {
	return r.Intn(ir.max-ir.min+1) + ir.min
}

func Random(min, max int, cb func(i, v int) bool) {
	r := rand.New(rand.NewSource(55))
	ir := IntRange{min, max}
	i := 0
	for ; cb(i, ir.NextRandom(r)); i++ {
	}
}

/*
date format 2015-01-23
*/
func ParseDateRange(start, end string) (startDate, endDate time.Time, err error) {
	startDate, err = now.Parse(start)
	if err != nil {
		err = errwrap.Wrap(err, "Parse start date %q", start)
		return
	}

	endDate, err = now.Parse(end)

	if err != nil {
		err = errwrap.Wrap(err, "Parse end date %q", end)
		return
	}

	if endDate.UnixNano() < startDate.UnixNano() {
		endDate = now.EndOfDay()
	} else {
		endDate = endDate.AddDate(0, 0, 1)
	}
	return
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func hashString(s string) string {
	return fmt.Sprintf("%x", hash(s))
}
