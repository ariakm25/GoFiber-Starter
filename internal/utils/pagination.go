package utils

import (
	"math"

	"gorm.io/gorm"
)

type PaginateParam struct {
	DB      *gorm.DB
	Page    int
	Limit   int
	ShowSQL bool
}

type Paginator struct {
	TotalItems int         `json:"total_items"`
	TotalPage  int         `json:"total_page"`
	Items      interface{} `json:"items"`
	Offset     int         `json:"offset"`
	Limit      int         `json:"limit"`
	Page       int         `json:"page"`
	NextPage   int         `json:"next_page"`
	PrevPage   int         `json:"prev_page"`
}

func Paginate(p *PaginateParam, result interface{}) *Paginator {
	db := p.DB.Session(&gorm.Session{})

	if p.ShowSQL {
		db = db.Debug()
	}

	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit == 0 {
		p.Limit = 10
	}

	done := make(chan bool, 1)
	var paginator Paginator
	var count int64
	var offset int

	go countRecords(db, result, done, &count)

	if p.Page == 1 {
		offset = 0
	} else {
		offset = (p.Page - 1) * p.Limit
	}

	db.Limit(p.Limit).Offset(offset).Find(result)
	<-done

	paginator.TotalItems = int(count)
	paginator.Items = result
	paginator.Page = p.Page

	paginator.Offset = offset
	paginator.Limit = p.Limit
	paginator.TotalPage = int(math.Ceil(float64(count) / float64(p.Limit)))

	if p.Page > 1 {
		paginator.PrevPage = p.Page - 1
	} else {
		paginator.PrevPage = p.Page
	}

	if p.Page == paginator.TotalPage {
		paginator.NextPage = p.Page
	} else {
		paginator.NextPage = p.Page + 1
	}
	return &paginator
}

func countRecords(db *gorm.DB, anyType interface{}, done chan bool, count *int64) {
	db.Model(anyType).Count(count)
	done <- true
}
