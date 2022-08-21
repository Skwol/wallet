package transaction

import "time"

type DTO struct {
	ID         int64
	SenderID   int64
	ReceiverID int64
	Amount     float64
	Timestamp  time.Time
	Type       TranType
}

type FilterTransactionsDTO struct {
	SenderIDs   []int64
	ReceiverIDs []int64
	Amount      FloatRangeFilter
	Timestamp   DateRangeFilter
	Types       []string
}

type FloatRangeFilter struct {
	From float64
	To   float64
}

type DateRangeFilter struct {
	From time.Time
	To   time.Time
}
