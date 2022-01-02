package transaction

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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
	if len(dto.SenderIDs) > 0 {
		filter.senderID = newInt64Filter(dto.SenderIDs...)
	}
	if len(dto.ReceiverIDs) > 0 {
		filter.receiverID = newInt64Filter(dto.ReceiverIDs...)
	}

	if dto.Amount != nil {
		filter.amount = newFloatRangeFilter(dto.Amount.From, dto.Amount.To)
	}

	if dto.Timestamp != nil {
		filter.timestamp = newDateRangeFilter(dto.Timestamp.From, dto.Timestamp.To)
	}

	if len(dto.Types) > 0 {
		filter.transactionType = newStringFilter(dto.Types...)
	}
	return filter
}

func (s transactionFilter) Empty() bool {
	return s.senderID == nil && s.receiverID == nil && s.amount == nil && s.timestamp == nil && s.transactionType == nil
}

func (s transactionFilter) BuildQuery(limit, offset int) string {
	var query strings.Builder
	query.WriteString("SELECT * FROM transaction WHERE")
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
	return fmt.Sprintf("%s %s ORDER BY ID DESC LIMIT %d OFFSET %d;", query.String(), strings.Join(filters, " AND "), limit, offset)
}

func (s *transactionFilter) Bind(r *http.Request) error {
	return nil
}

type stringFilter []string

func (f stringFilter) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte(`""`)) {
		return nil
	}
	var v []string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	f = v
	return nil
}

func (f stringFilter) MarshalJSON() ([]byte, error) {
	return json.Marshal(f)
}

func newStringFilter(values ...string) stringFilter {
	return values
}

func (f stringFilter) Empty() bool {
	return f == nil || len(f) == 0
}

func (f stringFilter) Build(fieldName string) string {
	vals := make([]string, len(f))
	for i, v := range f {
		vals[i] = fmt.Sprintf("'%s'", v)
	}
	return fmt.Sprintf("%s IN (%s)", fieldName, strings.Join(vals, ", "))
}

type int64Filter []int64

func (f int64Filter) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte(`""`)) {
		return nil
	}
	var v []int64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	f = v
	return nil
}

func (f int64Filter) MarshalJSON() ([]byte, error) {
	return json.Marshal(f)
}

func newInt64Filter(values ...int64) int64Filter {
	return values
}

func (f int64Filter) Empty() bool {
	return f == nil || len(f) == 0
}

func (f int64Filter) Build(fieldName string) string {
	values := make([]string, len(f))
	for i, v := range f {
		values[i] = fmt.Sprintf("%d", v)
	}
	return fmt.Sprintf("%s IN (%s)", fieldName, strings.Join(values, ", "))
}

type floatRangeFilter struct {
	From float64 `json:"from"`
	To   float64 `json:"to"`
}

func newFloatRangeFilter(from, to float64) *floatRangeFilter {
	return &floatRangeFilter{
		From: from,
		To:   to,
	}
}

func (f *floatRangeFilter) Empty() bool {
	return f == nil || f.From == 0 && f.To == 0
}

func (f floatRangeFilter) Build(fieldName string) string {
	return fmt.Sprintf("%s BETWEEN %f AND %f", fieldName, f.From, f.To)
}

type dateRangeFilter struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

func newDateRangeFilter(from, to time.Time) *dateRangeFilter {
	return &dateRangeFilter{
		From: from,
		To:   to,
	}
}

func (f *dateRangeFilter) Empty() bool {
	return f == nil || f.From.IsZero() && f.To.IsZero()
}

func (f dateRangeFilter) Build(fieldName string) string {
	return fmt.Sprintf("%s BETWEEN '%s' AND '%s'", fieldName, f.From.Format("2006-01-02 15:04:05"), f.To.Format("2006-01-02 15:04:05"))
}
