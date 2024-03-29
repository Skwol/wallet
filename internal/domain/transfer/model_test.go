package transfer

import (
	"reflect"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/skwol/wallet/pkg/clock"
)

func Test_createTransfer(t *testing.T) {
	clk := clock.NewFake(time.Date(2021, 10, 10, 10, 0, 0, 0, time.UTC))
	type args struct {
		dto *CreateTransferDTO
	}
	tests := []struct {
		name    string
		args    args
		want    *Transfer
		wantErr error
	}{
		{
			name:    "test missing sender",
			args:    args{dto: &CreateTransferDTO{Amount: 100, Receiver: WalletDTO{ID: 1}}},
			want:    nil,
			wantErr: errors.New("missing sender or receiver"),
		},
		{
			name:    "test missing receiver",
			args:    args{dto: &CreateTransferDTO{Amount: 100, Sender: WalletDTO{ID: 1}}},
			want:    nil,
			wantErr: errors.New("missing sender or receiver"),
		},
		{
			name:    "test missing amount",
			args:    args{dto: &CreateTransferDTO{Receiver: WalletDTO{ID: 1}, Sender: WalletDTO{ID: 2}}},
			want:    nil,
			wantErr: errors.New("amount should be greater then 0"),
		},
		{
			name:    "test negative amount",
			args:    args{dto: &CreateTransferDTO{Amount: -1, Receiver: WalletDTO{ID: 1}, Sender: WalletDTO{ID: 2}}},
			want:    nil,
			wantErr: errors.New("amount should be greater then 0"),
		},
		{
			name:    "test same sender and receiver",
			args:    args{dto: &CreateTransferDTO{Amount: 100, Sender: WalletDTO{ID: 1}, Receiver: WalletDTO{ID: 1}}},
			want:    nil,
			wantErr: errors.New("transfer can not be performed when sender and receiver is the same wallet"),
		},
		{
			name:    "test receiver does not have enough money",
			args:    args{dto: &CreateTransferDTO{Amount: 100, Sender: WalletDTO{ID: 1}, Receiver: WalletDTO{ID: 2}}},
			want:    nil,
			wantErr: errors.New("sender does not have enough 'money' for transfer"),
		},
		{
			name:    "test ok",
			args:    args{dto: &CreateTransferDTO{Amount: 100, Sender: WalletDTO{ID: 1, Balance: 150}, Receiver: WalletDTO{ID: 2, Balance: 50}}},
			want:    &Transfer{Amount: 100, Timestamp: clk.Now(), Sender: Wallet{ID: 1, Balance: 50}, Receiver: Wallet{ID: 2, Balance: 150}},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createTransfer(tt.args.dto, clk.Now())
			if tt.wantErr != nil {
				if err == nil || tt.wantErr.Error() != err.Error() {
					t.Errorf("createTransfer() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createTransfer() = %v, want %v", got, tt.want)
			}
		})
	}
}
