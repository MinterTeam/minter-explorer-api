package transaction

import (
	"encoding/base64"
	"encoding/hex"
	"github.com/MinterTeam/minter-explorer-api/v2/coins"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/minter-explorer-api/v2/transaction/data_resources"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/MinterTeam/minter-go-sdk/v2/transaction"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"reflect"
	"strconv"
	"time"
)

type Resource struct {
	Txn       uint64                 `json:"txn"`
	Hash      string                 `json:"hash"`
	Nonce     uint64                 `json:"nonce"`
	Block     uint64                 `json:"height"`
	Timestamp string                 `json:"timestamp"`
	GasCoin   resource.Interface     `json:"gas_coin"`
	Gas       string                 `json:"gas"`
	GasPrice  uint64                 `json:"gas_price"`
	Fee       string                 `json:"fee"`
	Type      uint8                  `json:"type"`
	Payload   string                 `json:"payload"`
	From      string                 `json:"from"`
	Data      resource.ItemInterface `json:"data"`
	RawTx     string                 `json:"raw_tx"`
}

func (Resource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	tx := model.(models.Transaction)

	return Resource{
		Txn:       tx.ID,
		Hash:      tx.GetHash(),
		Nonce:     tx.Nonce,
		Block:     tx.BlockID,
		Timestamp: tx.CreatedAt.Format(time.RFC3339),
		Gas:       strconv.FormatUint(tx.Gas, 10),
		GasPrice:  tx.GasPrice,
		Fee:       helpers.Fee2Bip(tx.GetFee()),
		GasCoin:   new(coins.IdResource).Transform(*tx.GasCoin),
		Type:      tx.Type,
		Payload:   base64.StdEncoding.EncodeToString(tx.Payload[:]),
		From:      tx.FromAddress.GetAddress(),
		Data:      TransformTxData(tx),
		RawTx:     hex.EncodeToString(tx.RawTx),
	}
}

type TransformTxConfig struct {
	Model    resource.ItemInterface
	Resource resource.Interface
}

var transformConfig = map[uint8]TransformTxConfig{
	uint8(transaction.TypeSend):                   {Model: new(api_pb.SendData), Resource: data_resources.Send{}},
	uint8(transaction.TypeSellCoin):               {Model: new(api_pb.SellCoinData), Resource: data_resources.SellCoin{}},
	uint8(transaction.TypeSellAllCoin):            {Model: new(api_pb.SellAllCoinData), Resource: data_resources.SellAllCoin{}},
	uint8(transaction.TypeBuyCoin):                {Model: new(api_pb.BuyCoinData), Resource: data_resources.BuyCoin{}},
	uint8(transaction.TypeCreateCoin):             {Model: new(api_pb.CreateCoinData), Resource: data_resources.CreateCoin{}},
	uint8(transaction.TypeDeclareCandidacy):       {Model: new(api_pb.DeclareCandidacyData), Resource: data_resources.DeclareCandidacy{}},
	uint8(transaction.TypeDelegate):               {Model: new(api_pb.DelegateData), Resource: data_resources.Delegate{}},
	uint8(transaction.TypeUnbond):                 {Model: new(api_pb.UnbondData), Resource: data_resources.Unbond{}},
	uint8(transaction.TypeRedeemCheck):            {Model: new(api_pb.RedeemCheckData), Resource: data_resources.RedeemCheck{}},
	uint8(transaction.TypeCreateMultisig):         {Model: new(api_pb.CreateMultisigData), Resource: data_resources.CreateMultisig{}},
	uint8(transaction.TypeMultisend):              {Model: new(api_pb.MultiSendData), Resource: data_resources.Multisend{}},
	uint8(transaction.TypeEditCandidate):          {Model: new(api_pb.EditCandidateData), Resource: data_resources.EditCandidate{}},
	uint8(transaction.TypeSetCandidateOnline):     {Model: new(api_pb.SetCandidateOnData), Resource: data_resources.SetCandidate{}},
	uint8(transaction.TypeSetCandidateOffline):    {Model: new(api_pb.SetCandidateOffData), Resource: data_resources.SetCandidate{}},
	uint8(transaction.TypeSetHaltBlock):           {Model: new(api_pb.SetHaltBlockData), Resource: data_resources.SetHaltBlock{}},
	uint8(transaction.TypeRecreateCoin):           {Model: new(api_pb.RecreateCoinData), Resource: data_resources.RecreateCoin{}},
	uint8(transaction.TypeEditCoinOwner):          {Model: new(api_pb.EditCoinOwnerData), Resource: data_resources.EditCoinOwner{}},
	uint8(transaction.TypeEditMultisig):           {Model: new(api_pb.EditMultisigData), Resource: data_resources.EditMultisigData{}},
	uint8(transaction.TypePriceVote):              {Model: new(api_pb.PriceVoteData), Resource: data_resources.PriceVote{}},
	uint8(transaction.TypeEditCandidatePublicKey): {Model: new(api_pb.EditCandidatePublicKeyData), Resource: data_resources.EditCandidatePublicKey{}},
}

func TransformTxData(tx models.Transaction) resource.Interface {
	config := transformConfig[tx.Type]

	val := reflect.New(reflect.TypeOf(config.Model).Elem()).Interface().(proto.Message)
	err := protojson.Unmarshal(tx.Data, val)
	helpers.CheckErr(err)

	return config.Resource.Transform(val, tx)
}
