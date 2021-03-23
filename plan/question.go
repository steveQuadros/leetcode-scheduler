package plan

import (
	"time"
)

type Question struct {
	Link, Title, Difficulty string
	InitialDate *time.Time
	Reviews []*Review
	CompletedAt *time.Time
	Intervals []int
}

type Review struct {
	// the date this should be reviewed
	Date time.Time
	// if the review was completed
	Completed bool
}

func (q *Question) Schedule(initialDate time.Time, intervals []int) {
	q.InitialDate = &initialDate
	q.Reviews = nil // reset
	for _, interval := range intervals {
		q.Reviews = append(q.Reviews, &Review{Date: initialDate.AddDate(0, 0, interval-1)})
	}
}

func (q *Question) Do(intervals []int) {
	q.DoAt(time.Now(), intervals)
}

func (q *Question) DoAt(doneDate time.Time, intervals []int) {
	q.CompletedAt = &doneDate
	q.Schedule(doneDate, intervals)
}