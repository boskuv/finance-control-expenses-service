package expenses_service

import "time"

type Expense struct {
	ID         uint64
	Amount     float64
	Timestamp  time.Time
	Comment    string
	CategoryID uint8
	UserID     uint64
}
