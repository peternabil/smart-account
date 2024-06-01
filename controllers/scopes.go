package controllers

import (
	"net/http"
	"strconv"
	"time"
)

func getPaginationArgs(r *http.Request) (int, int) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	if page <= 0 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(q.Get("page_size"))
	switch {
	case pageSize > 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}
	return page, pageSize
}

func getDates(r *http.Request) (time.Time, time.Time, error, error) {
	q := r.URL.Query()
	startDatestr := q.Get("start_date")
	startDate, sError := time.Parse("2006-01-02T15:04:05Z", startDatestr)
	if sError != nil {
		return time.Now(), time.Now(), sError, nil
	}
	endDatestr := q.Get("end_date")
	endDate, sError := time.Parse("2006-01-02T15:04:05Z", endDatestr)
	if sError != nil {
		return time.Now(), time.Now(), nil, sError
	}
	return startDate, endDate, nil, nil
}
