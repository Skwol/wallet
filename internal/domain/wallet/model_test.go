package wallet

import (
	"reflect"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/skwol/wallet/pkg/clock"
)

func TestWallet_Update(t *testing.T) {
	clk := clock.NewFake(time.Date(2021, 10, 10, 10, 0, 0, 0, time.UTC))
	type fields struct {
		ID      int64
		Name    string
		Balance float64
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
			name:    "test balance less then zero",
			fields:  fields{ID: 1},
			args:    args{wallet: &UpdateWalletDTO{CreateWalletDTO: CreateWalletDTO{Balance: -1}}},
			want:    nil,
			wantErr: errors.New("balance can not be less then 0"),
		},
		{
			name:    "test balance should be updated",
			fields:  fields{ID: 1, Balance: 10},
			args:    args{wallet: &UpdateWalletDTO{CreateWalletDTO: CreateWalletDTO{Balance: 10}}},
			want:    nil,
			wantErr: errors.New("balance should be updated"),
		},
		{
			name:   "test OK set balance to zero",
			fields: fields{ID: 1, Balance: 1},
			args:   args{wallet: &UpdateWalletDTO{CreateWalletDTO: CreateWalletDTO{Balance: 0}}},
			want: &Wallet{ID: 1, Balance: 0, TransactionsToApply: []Transaction{{
				SenderID: 1, ReceiverID: 1, Amount: 1, Timestamp: clk.Now(), Type: TranTypeWithdraw,
			}}},
			wantErr: nil,
		},
		{
			name:   "test OK deposit",
			fields: fields{ID: 1, Balance: 1},
			args:   args{wallet: &UpdateWalletDTO{CreateWalletDTO: CreateWalletDTO{Balance: 20}}},
			want: &Wallet{ID: 1, Balance: 20, TransactionsToApply: []Transaction{{
				SenderID: 1, ReceiverID: 1, Amount: 19, Timestamp: clk.Now(), Type: TranTypeDeposit,
			}}},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Wallet{
				ID:      tt.fields.ID,
				Name:    tt.fields.Name,
				Balance: tt.fields.Balance,
			}
			got, err := w.Update(tt.args.wallet, clk.Now())
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

func Test_newWallet(t *testing.T) {
	clk := clock.NewFake(time.Date(2021, 10, 10, 10, 0, 0, 0, time.UTC))
	type args struct {
		dto *CreateWalletDTO
	}
	tests := []struct {
		name    string
		args    args
		want    *Wallet
		wantErr error
	}{
		{
			name:    "test balance less then 0",
			args:    args{&CreateWalletDTO{Balance: -1}},
			want:    nil,
			wantErr: errors.New("balance can not be less then zero"),
		},
		{
			name:    "test ok",
			args:    args{&CreateWalletDTO{Balance: 0, Name: "test name"}},
			want:    &Wallet{Name: "test name", Balance: 0},
			wantErr: nil,
		},
		{
			name:    "test ok with balance",
			args:    args{&CreateWalletDTO{Balance: 1, Name: "test name"}},
			want:    &Wallet{Name: "test name", Balance: 1, TransactionsToApply: []Transaction{{Amount: 1, Timestamp: clk.Now(), Type: TranTypeDeposit}}},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newWallet(tt.args.dto, clk.Now())

			if tt.wantErr != nil {
				if err == nil || tt.wantErr.Error() != err.Error() {
					t.Errorf("newWallet() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newWallet() = %v, want %v", got, tt.want)
			}
		})
	}
}
