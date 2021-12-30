package wallet

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/skwol/wallet/internal/domain/transaction"
)

func TestWallet_Update(t *testing.T) {
	timeNow = func() time.Time {
		return time.Date(2021, 10, 10, 10, 0, 0, 0, time.UTC)
	}
	type fields struct {
		ID        int64
		Name      string
		AccountID int64
		Balance   float64
	}
	type args struct {
		wallet *UpdateWalletDTO
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Wallet
		wantErr error
	}{
		{
			name:    "test another account id",
			fields:  fields{ID: 1, AccountID: 1},
			args:    args{wallet: &UpdateWalletDTO{AccountID: 2}},
			want:    nil,
			wantErr: fmt.Errorf("wrong account"),
		},
		{
			name:    "test balance less then zero",
			fields:  fields{ID: 1, AccountID: 1},
			args:    args{wallet: &UpdateWalletDTO{AccountID: 1, Balance: -1}},
			want:    nil,
			wantErr: fmt.Errorf("balance can not be less then 0"),
		},
		{
			name:    "test balance should be updated",
			fields:  fields{ID: 1, AccountID: 1, Balance: 10},
			args:    args{wallet: &UpdateWalletDTO{AccountID: 1, Balance: 10}},
			want:    nil,
			wantErr: fmt.Errorf("balance should be updated"),
		},
		{
			name:   "test OK set balance to zero",
			fields: fields{ID: 1, AccountID: 1, Balance: 1},
			args:   args{wallet: &UpdateWalletDTO{AccountID: 1, Balance: 0}},
			want: &Wallet{ID: 1, AccountID: 1, Balance: 0, TransactionsToApply: []transaction.Transaction{{
				SenderID: 1, ReceiverID: 1, Amount: 1, Timestamp: timeNow(), Type: transaction.TranTypeWithdraw,
			}}},
			wantErr: nil,
		},
		{
			name:   "test OK deposit",
			fields: fields{ID: 1, AccountID: 1, Balance: 1},
			args:   args{wallet: &UpdateWalletDTO{AccountID: 1, Balance: 20}},
			want: &Wallet{ID: 1, AccountID: 1, Balance: 20, TransactionsToApply: []transaction.Transaction{{
				SenderID: 1, ReceiverID: 1, Amount: 19, Timestamp: timeNow(), Type: transaction.TranTypeDeposit,
			}}},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Wallet{
				ID:        tt.fields.ID,
				Name:      tt.fields.Name,
				AccountID: tt.fields.AccountID,
				Balance:   tt.fields.Balance,
			}
			got, err := w.Update(tt.args.wallet)
			if tt.wantErr != nil {
				if err == nil || tt.wantErr.Error() != err.Error() {
					t.Errorf("Wallet.Update() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Wallet.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}
