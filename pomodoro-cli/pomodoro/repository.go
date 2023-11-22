package pomodoro

import "time"

type Repository interface {
	Create(i Interval) (int64, error)
	Update(i Interval) error
	ByID(id int64) (Interval, error)
	Last() (Interval, error)
	Breaks(n int) ([]Interval, error)
	CategorySummary(day time.Time, filter string) (time.Duration, error)
}
