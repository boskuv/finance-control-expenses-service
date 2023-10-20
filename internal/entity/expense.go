package entity

import (
	"time"

	"github.com/google/uuid"
)

type Expense struct {
	ID         uint64
	Amount     float64
	Timestamp  time.Time
	Comment    string
	CategoryID uint8
	UserID     uuid.UUID
}
