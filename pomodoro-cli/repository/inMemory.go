//go:build inmemory
// +build inmemory

package repository

import (
	"fmt"
	"github.com/zawlinnnaing/go-clis/pomodoro-cli/pomodoro"
	"sync"
)

type InMemoryRepo struct {
	sync.RWMutex
	intervals []pomodoro.Interval
}

func (repo *InMemoryRepo) Create(i pomodoro.Interval) (int64, error) {
	repo.Lock()
	defer repo.Unlock()

	i.ID = int64(len(repo.intervals)) + 1
	repo.intervals = append(repo.intervals, i)
	return i.ID, nil
}

func (repo *InMemoryRepo) Update(interval pomodoro.Interval) error {
	repo.Lock()
	defer repo.Unlock()
	if interval.ID == 0 || interval.ID > int64(len(repo.intervals)) {
		return fmt.Errorf("%w: %d", pomodoro.ErrInvalidID, interval.ID)
	}
	repo.intervals[interval.ID-1] = interval
	return nil
}

func (repo *InMemoryRepo) ByID(id int64) (pomodoro.Interval, error) {
	repo.RLock()
	defer repo.RUnlock()
	interval := pomodoro.Interval{}
	if id == 0 || id > int64(len(repo.intervals)) {
		return interval, fmt.Errorf("%w: %d", pomodoro.ErrInvalidID, id)
	}
	return repo.intervals[id-1], nil
}

func (repo *InMemoryRepo) Last() (pomodoro.Interval, error) {
	repo.RLock()
	defer repo.RUnlock()
	interval := pomodoro.Interval{}
	if len(repo.intervals) == 0 {
		return interval, pomodoro.ErrNoIntervals
	}
	return repo.intervals[len(repo.intervals)-1], nil
}

func (repo *InMemoryRepo) Breaks(n int) ([]pomodoro.Interval, error) {
	repo.RLock()
	defer repo.RUnlock()
	var data []pomodoro.Interval
	for k := len(repo.intervals) - 1; k >= 0; k-- {
		interval := &repo.intervals[k]
		if interval.Category == pomodoro.CategoryPomodoro {
			continue
		}
		data = append(data, *interval)
		if len(data) == n {
			return data, nil
		}
	}
	return data, nil
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{intervals: []pomodoro.Interval{}}
}
