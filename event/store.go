package event

import "time"

type Store interface {
	Save(event Event) error
	GetByType(eventType string) ([]Event, error)
	GetAll() ([]Event, error)
	GetByTimeRange(start, end time.Time) ([]Event, error)
	Count() int
	GetLatest(n int) ([]Event, error)
}
