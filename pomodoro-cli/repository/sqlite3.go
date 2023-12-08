//go:build !inmemory && !containers

package repository

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zawlinnnaing/go-clis/pomodoro-cli/pomodoro"
	"sync"
	"time"
)

const (
	createIntervalTable string = `
	CREATE TABLE IF NOT EXISTS "interval" (
    "id" INTEGER,
    "start_time" DATETIME NOT NULL,
    "planned_duration" INTEGER DEFAULT 0,
    "actual_duration" INTEGER DEFAULT 0,
    "category" TEXT NOT NULL,
    "state" INTEGER DEFAULT 1,
    PRIMARY KEY("id")
);`
)

type dbRepo struct {
	db *sql.DB
	sync.RWMutex
}

func (repo *dbRepo) Create(interval pomodoro.Interval) (int64, error) {
	repo.Lock()
	defer repo.Unlock()
	insertStatement, err := repo.db.Prepare("INSERT INTO interval VALUES(NULL,?,?,?,?,?)")
	if err != nil {
		return 0, err
	}
	defer insertStatement.Close()
	res, err := insertStatement.Exec(interval.StartTime, interval.PlannedDuration, interval.ActualDuration, interval.Category, interval.State)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (repo *dbRepo) Update(interval pomodoro.Interval) error {
	repo.Lock()
	defer repo.Unlock()
	updateStatement, err := repo.db.Prepare("UPDATE interval SET start_time=?, actual_duration=?, state=? WHERE id=?")
	if err != nil {
		return err
	}
	res, err := updateStatement.Exec(interval.StartTime, interval.ActualDuration, interval.State, interval.ID)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	return err
}

func (repo *dbRepo) ByID(id int64) (pomodoro.Interval, error) {
	repo.RLock()
	defer repo.RUnlock()
	row := repo.db.QueryRow("SELECT * from interval WHERE id=?", id)
	interval := pomodoro.Interval{}
	err := row.Scan(&interval.ID, &interval.StartTime, &interval.PlannedDuration, &interval.ActualDuration, &interval.Category, &interval.State)
	return interval, err
}

func (repo *dbRepo) Last() (pomodoro.Interval, error) {
	repo.RLock()
	defer repo.RUnlock()
	interval := pomodoro.Interval{}
	err := repo.db.QueryRow("SELECT * from interval ORDER BY id desc LIMIT 1").Scan(
		&interval.ID, &interval.StartTime, &interval.PlannedDuration, &interval.ActualDuration, &interval.Category, &interval.State)
	if errors.Is(err, sql.ErrNoRows) {
		return interval, pomodoro.ErrNoIntervals
	}
	return interval, err
}

func (repo *dbRepo) Breaks(n int) ([]pomodoro.Interval, error) {
	repo.RLock()
	defer repo.RUnlock()

	rows, err := repo.db.Query("SELECT * from interval WHERE CATEGORY LIKE '%Break' ORDER BY id desc LIMIT ?", n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var data []pomodoro.Interval
	for rows.Next() {
		interval := pomodoro.Interval{}
		err := rows.Scan(&interval.ID, &interval.StartTime, &interval.PlannedDuration, &interval.ActualDuration, &interval.Category, &interval.State)
		if err != nil {
			return nil, err
		}
		data = append(data, interval)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return data, nil
}

func (repo *dbRepo) CategorySummary(day time.Time, filter string) (time.Duration, error) {
	repo.RLock()
	defer repo.RUnlock()
	dbStatement := `SELECT sum(actual_duration) from interval WHERE category like ? AND strftime('%Y-%m-%d',start_time, 'localtime')=strftime('%Y-%m-%d', ? , 'localtime');`
	var dbResult sql.NullInt64
	var totalDuration time.Duration
	err := repo.db.QueryRow(dbStatement, filter, day).Scan(&dbResult)
	if err != nil {
		return totalDuration, err
	}
	if dbResult.Valid {
		totalDuration = time.Duration(dbResult.Int64)
	}
	return totalDuration, nil
}

func NewSQLite3Repo(dbFile string) (*dbRepo, error) {
	if dbFile == "" {
		return nil, errors.New("DB file not provided")
	}
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetMaxOpenConns(1)

	if err = db.Ping(); err != nil {
		return nil, err
	}

	if _, err := db.Exec(createIntervalTable); err != nil {
		return nil, err
	}

	return &dbRepo{db: db}, err
}
