package chart

import (
	"fmt"

	"github.com/vjeantet/jodaTime"

	"github.com/moisespsena-go/aorm"
)

type ChartInfo interface {
	Label() string
	Description() string
}

type Finder interface {
	Find(ctx *Context) (res []Record, err error)
}

type FinderFunc func(ctx *Context) (res []Record, err error)

func (f FinderFunc) Find(ctx *Context) (res []Record, err error) {
	return f(ctx)
}

type QueryFinder struct {
	DB        *aorm.DB
	TableName string
	FieldName string
}

func (q *QueryFinder) Find(ctx *Context) (res []Record, err error) {
	DB := q.DB

	if q.TableName != "" {
		DB = DB.Table(q.TableName)
	}

	fieldName := "created_at"

	if q.FieldName != "" {
		fieldName = q.FieldName
	}

	type Tress struct {
		Date  string
		Total string
	}

	var ress []Tress

	if err = DB.Where(aorm.IQ(fieldName+" BETWEEN ? AND ?"), ctx.StartDate, ctx.EndDate).
		Select(fmt.Sprint("cast(date(", fieldName, ") as text) as date, count(*) as total")).
		Group(fmt.Sprint("date(", fieldName, ")")).
		Order(fmt.Sprint("date(", fieldName, ")")).
		Scan(&ress).Error; err == nil {
		res = make([]Record, len(ress))
		for i, r := range ress {
			res[i].Total = r.Total
			res[i].Date, _ = jodaTime.Parse("YYYY-MM-dd", r.Date)
		}
	}
	return
}

type FinderInfo interface {
	ChartInfo
	Finder
}

type finderInfo struct {
	Finder
	label       string
	description string
}

func (i *finderInfo) Label() string {
	return i.label
}

func (i *finderInfo) Description() string {
	return i.description
}

func NewFinderInfo(finder Finder, label, description string) FinderInfo {
	return &finderInfo{finder, label, description}
}
