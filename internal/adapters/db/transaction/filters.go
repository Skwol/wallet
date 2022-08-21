package transaction

import (
	"fmt"
	"strings"
	"time"

	"github.com/skwol/wallet/internal/domain/transaction"
)

type transactionFilter struct {
	senderID        int64Filter
	receiverID      int64Filter
	amount          *floatRangeFilter
	timestamp       *dateRangeFilter
	transactionType stringFilter
}

func newTransactionFilter(dto *transaction.FilterTransactionsDTO) transactionFilter {
	var filter transactionFilter
	if dto == nil {
		return filter
	}

	filter.senderID = newInt64Filter(dto.SenderIDs...)
	filter.receiverID = newInt64Filter(dto.ReceiverIDs...)
	filter.amount = newFloatRangeFilter(dto.Amount.From, dto.Amount.To)
	filter.timestamp = newDateRangeFilter(dto.Timestamp.From, dto.Timestamp.To)
	filter.transactionType = newStringFilter(dto.Types...)

	return filter
}

func (s transactionFilter) Empty() bool {
	return s.senderID == nil && s.receiverID == nil && s.amount == nil && s.timestamp == nil && s.transactionType == nil
}

func (s transactionFilter) BuildQuery(limit, offset int) string {
	var filters []string
	if !s.senderID.Empty() {
		filters = append(filters, s.senderID.Build("sender_id"))
	}
	if !s.receiverID.Empty() {
		filters = append(filters, s.receiverID.Build("receiver_id"))
	}
	if !s.amount.Empty() {
		filters = append(filters, s.amount.Build("amount"))
	}
	if !s.timestamp.Empty() {
		filters = append(filters, s.timestamp.Build("date"))
	}
	if !s.transactionType.Empty() {
		filters = append(filters, s.transactionType.Build("tran_type"))
	}
	var filter string
	if len(filters) > 0 {
		filter = fmt.Sprintf("WHERE %s ", strings.Join(filters, " AND "))
	}
	return fmt.Sprintf("SELECT id, sender_id, receiver_id, amount, date, tran_type FROM transaction %sORDER BY id ASC LIMIT %d OFFSET %d;", filter, limit, offset)
}

type stringFilter []string

func newStringFilter(values ...string) stringFilter {
	if len(values) == 0 {
		return nil
	}
	return values
}

func (f stringFilter) Empty() bool {
	return len(f) == 0
}

func (f stringFilter) Build(fieldName string) string {
	vals := make([]string, len(f))
	for i, v := range f {
		vals[i] = fmt.Sprintf("'%s'", v)
	}
	return fmt.Sprintf("%s IN (%s)", fieldName, strings.Join(vals, ", "))
}

type int64Filter []int64

func newInt64Filter(values ...int64) int64Filter {
	if len(values) == 0 {
		return nil
	}
	return values
}

func (f int64Filter) Empty() bool {
	return len(f) == 0
}

func (f int64Filter) Build(fieldName string) string {
	values := make([]string, len(f))
	for i, v := range f {
		values[i] = fmt.Sprintf("%d", v)
	}
	return fmt.Sprintf("%s IN (%s)", fieldName, strings.Join(values, ", "))
}

type floatRangeFilter struct {
	From float64
	To   float64
}

func newFloatRangeFilter(from, to float64) *floatRangeFilter {
	if from == 0 && to == 0 {
		return nil
	}
	return &floatRangeFilter{
		From: from,
		To:   to,
	}
}

func (f *floatRangeFilter) Empty() bool {
	return f == nil || f.From == 0 && f.To == 0
}

func (f floatRangeFilter) Build(fieldName string) string {
	if f.From > 0 && f.To == 0 {
		return fmt.Sprintf("%s > %f", fieldName, f.From)
	} else if f.From == 0 && f.To > 0 {
		return fmt.Sprintf("%s < %f", fieldName, f.To)
	}
	return fmt.Sprintf("%s BETWEEN %f AND %f", fieldName, f.From, f.To)
}

type dateRangeFilter struct {
	From time.Time
	To   time.Time
}

func newDateRangeFilter(from, to time.Time) *dateRangeFilter {
	if from.IsZero() && to.IsZero() {
		return nil
	}
	return &dateRangeFilter{
		From: from,
		To:   to,
	}
}

func (f *dateRangeFilter) Empty() bool {
	return f == nil || f.From.IsZero() && f.To.IsZero()
}

func (f dateRangeFilter) Build(fieldName string) string {
	if !f.From.IsZero() && f.To.IsZero() {
		return fmt.Sprintf("%s > '%s'", fieldName, f.From.Format("2006-01-02 15:04:05"))
	} else if f.From.IsZero() && !f.To.IsZero() {
		return fmt.Sprintf("%s < '%s'", fieldName, f.To.Format("2006-01-02 15:04:05"))
	}
	return fmt.Sprintf("%s BETWEEN '%s' AND '%s'", fieldName, f.From.Format("2006-01-02 15:04:05"), f.To.Format("2006-01-02 15:04:05"))
}
