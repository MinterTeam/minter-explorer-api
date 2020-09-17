package transaction

import (
	"encoding/json"
	"github.com/MinterTeam/minter-explorer-api/coins"
	txModels "github.com/MinterTeam/minter-explorer-api/transaction/data_models"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/MinterTeam/minter-go-sdk/v2/transaction"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
)

type Service struct {
	coinRepository *coins.Repository
}

func NewService(coinRepository *coins.Repository) *Service {
	return &Service{coinRepository}
}

func (s *Service) PrepareTransactionsModel(txs []models.Transaction) ([]models.Transaction, error) {
	for key, tx := range txs {
		prepared, err := s.PrepareTransactionModel(&tx)
		if err != nil {
			return nil, err
		}

		txs[key] = *prepared
	}

	return txs, nil
}

func (s *Service) PrepareTransactionModel(tx *models.Transaction) (*models.Transaction, error) {
	if tx.Type == uint8(transaction.TypeRedeemCheck) {
		data := new(api_pb.RedeemCheckData)

		err := json.Unmarshal(tx.Data, data)
		if err != nil {
			return nil, err
		}

		checkData, err := s.transformCheckDataToModel(data.RawCheck)
		if err != nil {
			return nil, err
		}

		tx.IData = txModels.Check{
			RawCheck: data.RawCheck,
			Proof:    data.Proof,
			Check:    *checkData,
		}
	}

	return tx, nil
}

func (s Service) transformCheckDataToModel(raw string) (*txModels.CheckData, error) {
	data, err := transaction.DecodeCheckBase64(raw)
	if err != nil {
		return nil, err
	}

	sender, err := data.Sender()
	if err != nil {
		return nil, err
	}

	coin, err := s.coinRepository.FindByID(uint(data.Coin))
	if err != nil {
		return nil, err
	}

	gasCoin, err := s.coinRepository.FindByID(uint(data.GasCoin))
	if err != nil {
		return nil, err
	}

	return &txModels.CheckData{
		Coin:     coin,
		GasCoin:  gasCoin,
		Nonce:    data.Nonce,
		Value:    data.Value,
		Sender:   sender,
		DueBlock: data.DueBlock,
	}, nil
}
