package transfer

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rodolfoalvesg/api-banking/api/domain/entities/transfers"
	"github.com/stretchr/testify/assert"
)

var ErrorTransfersNoListed = errors.New("transfers not listed")

//TestCreateTransfer, teste de caso de uso para criar e registrar transferência
func TestCreateTransfer(t *testing.T) {
	t.Parallel()

	tFake := transfers.Transfer{
		AccountOriginID:      uuid.New().String(),
		AccountDestinationID: uuid.New().String(),
		Amount:               2550,
		CreatedAt:            time.Now().UTC(),
	}

	type TestCase struct {
		name         string
		transferMock transfers.TransferMock
		transferFake transfers.Transfer
		want         uuid.UUID
	}

	transferFakeID := uuid.New()

	testCase := []TestCase{
		{
			name: "transfer created successfully",
			transferMock: transfers.TransferMock{
				OnSaveTransfer: func(transfers.Transfer) (uuid.UUID, error) {
					return transferFakeID, nil
				},
			},
			transferFake: tFake,
			want:         transferFakeID,
		},
		{
			name: "transfer not created",
			transferMock: transfers.TransferMock{
				OnSaveTransfer: func(transfers.Transfer) (uuid.UUID, error) {
					return uuid.UUID{}, errors.New("Conflict")
				},
			},
			transferFake: tFake,
			want:         uuid.UUID{},
		},
	}

	for _, tc := range testCase {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			usecase := NewUsecaseTransfers(transfers.TransferMock{
				OnSaveTransfer: tt.transferMock.OnSaveTransfer,
			})

			got, err := usecase.CreateTransfer(context.Background(), tt.transferFake)

			if got != tt.want && err != nil {
				t.Errorf("Error %s", err)
			}

		})
	}
}

//TestShowTransfers, teste de caso de uso para exibição de transferências de um usuário
func TestShowTransfers(t *testing.T) {
	t.Parallel()

	type TestCaseList struct {
		name         string
		transferMock transfers.TransferMock
		accId        string
		wantErr      error
		want         []transfers.Transfer
	}

	myListTransfers := []transfers.Transfer{}

	testCase := []TestCaseList{
		{
			name: "transfers listed",
			transferMock: transfers.TransferMock{
				OnListAllTransfer: func(string) ([]transfers.Transfer, error) {
					return myListTransfers, nil
				},
			},
			accId:   uuid.New().String(),
			wantErr: nil,
			want:    []transfers.Transfer{},
		},
		{
			name: "transfers not listed",
			transferMock: transfers.TransferMock{
				OnListAllTransfer: func(string) ([]transfers.Transfer, error) {
					return []transfers.Transfer{}, ErrorTransfersNoListed
				},
			},
			accId:   uuid.New().String(),
			wantErr: ErrorTransfersNoListed,
			want:    []transfers.Transfer{nil},
		},
	}

	for _, tc := range testCase {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			usecase := NewUsecaseTransfers(transfers.TransferMock{
				OnListAllTransfer: tt.transferMock.OnListAllTransfer,
			})

			listedAllTransfers, err := usecase.ShowTransfers(context.Background(), tt.accId)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.ElementsMatch(t, tt.want, listedAllTransfers)

		})
	}

}
