package types

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type Transaction struct {
	Nonce        uint64
	To           common.Address
	Amount       *big.Int
	GasLimit     uint64
	GasPrice     *big.Int
	Data         []byte
	TokenAddress common.Address
}
