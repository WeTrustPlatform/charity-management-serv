package db

// EventName is either Staked or Unstaked
type EventName string

const (
	Staked   EventName = "Staked"
	Unstaked EventName = "Unstaked"
)

// EventInput is the aggregation of all the input parameters
// of two events Staked and Unstaked
type EventInput struct {
	User   string `json:"user"`
	Amount string `json:"amount"`
	Total  string `json:"total"`
	Data   string `json:"data"`
}

// TimeLockedStakingEvent emitted when users interact with the staking V0 contract
// https://github.com/WeTrustPlatform/staking-dapp/blob/master/packages/truffle/contracts/TimeLockedStaking.sol
// Cross referenced https://github.com/ethereum/go-ethereum/blob/master/core/types/log.go
// Use camelcase to be consistent with Web3 response
type TimeLockedStakingEvent struct {
	Address          string     `json:"address"`
	BlockHash        string     `json:"blockHash"`
	BlockNumber      uint64     `json:"blockNumber"`
	Event            EventName  `json:"event" sql:"type:ENUM('Staked', 'Unstaked')"`
	ID               string     `json:"id"`
	LogIndex         uint       `json:"logIndex"`
	ReturnValues     EventInput `json:"returnValues"`
	TransactionHash  string     `json:"transactionHash"`
	TransactionIndex uint       `json:"transactionIndex"`
}
