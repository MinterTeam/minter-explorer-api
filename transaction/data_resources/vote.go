package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
)

type VoteCommission struct {
	PubKey                  string `json:"pub_key"`
	Height                  uint64 `json:"height"`
	Coin                    Coin   `json:"coin"`
	PayloadByte             string `json:"payload_byte"`
	Send                    string `json:"send"`
	BuyBancor               string `json:"buy_bancor"`
	SellBancor              string `json:"sell_bancor"`
	SellAllBancor           string `json:"sell_all_bancor"`
	BuyPool                 string `json:"buy_pool"`
	SellPool                string `json:"sell_pool"`
	SellAllPool             string `json:"sell_all_pool"`
	CreateTicker3           string `json:"create_ticker3"`
	CreateTicker4           string `json:"create_ticker4"`
	CreateTicker5           string `json:"create_ticker5"`
	CreateTicker6           string `json:"create_ticker6"`
	CreateTicker7_10        string `json:"create_ticker7_10"`
	CreateCoin              string `json:"create_coin"`
	CreateToken             string `json:"create_token"`
	RecreateCoin            string `json:"recreate_coin"`
	RecreateToken           string `json:"recreate_token"`
	DeclareCandidacy        string `json:"declare_candidacy"`
	Delegate                string `json:"delegate"`
	Unbond                  string `json:"unbond"`
	RedeemCheck             string `json:"redeem_check"`
	SetCandidateOn          string `json:"set_candidate_on"`
	SetCandidateOff         string `json:"set_candidate_off"`
	CreateMultisig          string `json:"create_multisig"`
	MultisendBase           string `json:"multisend_base"`
	MultisendDelta          string `json:"multisend_delta"`
	EditCandidate           string `json:"edit_candidate"`
	SetHaltBlock            string `json:"set_halt_block"`
	EditTickerOwner         string `json:"edit_ticker_owner"`
	EditMultisig            string `json:"edit_multisig"`
	PriceVote               string `json:"price_vote"`
	EditCandidatePublicKey  string `json:"edit_candidate_public_key"`
	CreateSwapPool          string `json:"create_swap_pool"`
	AddLiquidity            string `json:"add_liquidity"`
	RemoveLiquidity         string `json:"remove_liquidity"`
	EditCandidateCommission string `json:"edit_candidate_commission"`
	MoveStake               string `json:"move_stake"`
	MintToken               string `json:"mint_token"`
	BurnToken               string `json:"burn_token"`
	VoteCommission          string `json:"vote_commission"`
	VoteUpdate              string `json:"vote_updated"`
}

func (VoteCommission) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*api_pb.VoteCommissionData)

	return VoteCommission{
		PubKey:                  data.PubKey,
		Height:                  data.Height,
		Coin:                    new(Coin).Transform(data.Coin),
		PayloadByte:             data.PayloadByte,
		Send:                    data.Send,
		BuyBancor:               data.BuyBancor,
		SellBancor:              data.SellBancor,
		SellAllBancor:           data.SellAllBancor,
		BuyPool:                 data.BuyPool,
		SellPool:                data.SellPool,
		SellAllPool:             data.SellAllPool,
		CreateTicker3:           data.CreateTicker3,
		CreateTicker4:           data.CreateTicker4,
		CreateTicker5:           data.CreateTicker5,
		CreateTicker6:           data.CreateTicker6,
		CreateTicker7_10:        data.CreateTicker7_10,
		CreateCoin:              data.CreateCoin,
		CreateToken:             data.CreateToken,
		RecreateCoin:            data.RecreateCoin,
		RecreateToken:           data.RecreateToken,
		DeclareCandidacy:        data.DeclareCandidacy,
		Delegate:                data.Delegate,
		Unbond:                  data.Unbond,
		RedeemCheck:             data.RedeemCheck,
		SetCandidateOn:          data.SetCandidateOn,
		SetCandidateOff:         data.SetCandidateOff,
		CreateMultisig:          data.CreateMultisig,
		MultisendBase:           data.MultisendBase,
		MultisendDelta:          data.MultisendDelta,
		EditCandidate:           data.EditCandidate,
		SetHaltBlock:            data.SetHaltBlock,
		EditTickerOwner:         data.EditTickerOwner,
		EditMultisig:            data.EditMultisig,
		PriceVote:               data.PriceVote,
		EditCandidatePublicKey:  data.EditCandidatePublicKey,
		CreateSwapPool:          data.CreateSwapPool,
		AddLiquidity:            data.AddLiquidity,
		RemoveLiquidity:         data.RemoveLiquidity,
		EditCandidateCommission: data.EditCandidateCommission,
		MoveStake:               data.MoveStake,
		MintToken:               data.MintToken,
		BurnToken:               data.BurnToken,
		VoteCommission:          data.VoteCommission,
		VoteUpdate:              data.VoteUpdate,
	}
}

type VoteUpdate struct {
	PubKey  string `json:"pub_key"`
	Height  uint64 `json:"height"`
	Version string `json:"version"`
}

func (VoteUpdate) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*api_pb.VoteUpdateData)

	return VoteUpdate{
		PubKey:  data.PubKey,
		Height:  data.Height,
		Version: data.Version,
	}
}
