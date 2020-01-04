package transfer

import (
	"context"
	"crypto/ecdsa"
	"etherscan-go/token"
	types2 "etherscan-go/types"
	"fmt"
	"github.com/a6910438/go-logger"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/shopspring/decimal"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type dbBaser interface {
	GetAssests() ([]*types2.UserAssest, error)
	GetConfigByWithdrawNumber() (*types2.Config, error)
}

type Transfer struct {
	db     dbBaser
	client *ethclient.Client
}

func NewTransfer(db dbBaser, client *ethclient.Client) (*Transfer, error) {
	return &Transfer{
		db:     db,
		client: client,
	}, nil
}

// 每天凌晨执行一次
func (t *Transfer) Start() {
	ests, err := t.db.GetAssests()
	if err != nil {
		logger.Errorf("Transfer GetAssests err", err)
	}

	t.GetTokenBalance("0xf8251ab6210606c3e06f8429b90ba946b4c5c94c", "0x4b446871737a443f5f4052d0278ef3183bb2271a")
	t.TransferERC20("", common.HexToAddress("0xddf3d8fe538c249e592c89ed91fa103dd3d2a9af"), common.HexToAddress("0x4b446871737a443f5f4052d0278ef3183bb2271a"), decimal.New(1, 0))

	for _, e := range ests {
		config, err := t.db.GetConfigByWithdrawNumber()
		if err != nil {
			logger.Errorf("Transfer GetConfigByWithdrawNumber err", err)
		}
		condition, _ := decimal.NewFromString(config.CValue)
		// ETH余额
		balance, err := t.client.BalanceAt(context.Background(), common.HexToAddress(e.Address), nil)
		// 余额是否大于提币数量 大于直接走归集逻辑
		if condition.Cmp(decimal.NewFromBigInt(balance, 0)) <= 0 {

		}
	}
}

func (t *Transfer) GetTokenBalance(address, contractAddress string) (*big.Int, error) {
	instance, err := token.NewToken(common.HexToAddress(contractAddress), t.client)
	if err != nil {
		logger.Errorf("Transfer NewToken err", err)
		return nil, nil
	}
	bal, err := instance.BalanceOf(&bind.CallOpts{}, common.HexToAddress(address))
	if err != nil {
		logger.Errorf("Transfer BalanceOf err", err)
		return nil, nil
	}
	return bal, nil
}

func (t *Transfer) TransferETH(hexKey string, toAddress common.Address, value *big.Int) error {
	// 导入私钥
	privateKey, publicKeyECDSA := t.ImportPrivateKey(hexKey)

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := t.client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		logger.Errorf("Transfer PendingNonceAt err", err)
		return err
	}

	gasLimit := uint64(21000) // in units
	gasPrice, err := t.client.SuggestGasPrice(context.Background())
	if err != nil {
		logger.Errorf("Transfer SuggestGasPrice err", err)
		return err
	}

	var data []byte
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	chainID, err := t.client.NetworkID(context.Background())
	if err != nil {
		logger.Errorf("Transfer NetworkID err", err)
		return err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		logger.Errorf("Transfer SignTx err", err)
		return err
	}

	err = t.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		logger.Errorf("Transfer SendTransaction err", err)
		return err
	}

	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())
	return nil
}

func (t *Transfer) TransferERC20(hexKey string, toAddress common.Address, tokenAddress common.Address, number decimal.Decimal) (string, error) {
	// 导入私钥
	privateKey, publicKeyECDSA := t.ImportPrivateKey(hexKey)

	// 获取公钥地址
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 获取Nonce值
	nonce, err := t.client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		logger.Errorf("Transfer PendingNonceAt err", err)
		return "", err
	}

	// 获取Gas费用
	gasPrice, err := t.client.SuggestGasPrice(context.Background())
	if err != nil {
		logger.Errorf("Transfer SuggestGasPrice err", err)
		return "", err
	}
	//toAddress := common.HexToAddress("0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d")
	//tokenAddress := common.HexToAddress("0x28b149020d2152179873ec60bed6bf7cd705775d")

	// 获取Gas Limit
	data, amount := t.GasLimitData(toAddress, number)
	gasLimit, err := t.client.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &toAddress,
		Data: data,
	})
	if err != nil {
		logger.Errorf("Transfer EstimateGas err", err)
		return "", err
	}

	// 发送交易
	hash, err := t.SignAndSendTx(privateKey, types2.Transaction{Nonce: nonce, GasPrice: gasPrice, GasLimit: gasLimit, TokenAddress: tokenAddress, Amount: amount, Data: data})
	if err != nil {
		logger.Errorf("Transfer SignAndSendTx err", err)
		return "", err
	}
	return hash, nil
}

func (t *Transfer) ImportPrivateKey(hexkey string) (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	privateKey, err := crypto.HexToECDSA(hexkey)

	if err != nil {
		logger.Errorf("Transfer Casting Private Key err", err)
		return nil, nil
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	fmt.Println(crypto.PubkeyToAddress(*publicKeyECDSA).Hex())
	if !ok {
		logger.Errorf("Transfer Casting Public Key err", err)
		return nil, nil
	}
	return privateKey, publicKeyECDSA
}

func (t *Transfer) GasLimitData(toAddress common.Address, number decimal.Decimal) ([]byte, *big.Int) {
	// 转账方法名称
	transferFnSignature := []byte("transfer(address,uint256)")
	methodID := crypto.Keccak256(transferFnSignature)[:4]

	// 合约地址
	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)

	// 计算转出金额并转字节
	unit, _ := decimal.NewFromString("1000000000000000000")
	amount := number.Mul(unit)
	n := new(big.Int)
	n, ok := n.SetString(amount.String(), 10)
	if !ok {
		logger.Errorf("SetString err", ok)
	}
	paddedAmount := common.LeftPadBytes(n.Bytes(), 32)

	// 数据拼接到一起
	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)
	return data, n
}

func (t *Transfer) SignAndSendTx(privateKey *ecdsa.PrivateKey, tran types2.Transaction) (string, error) {
	// 创建一个交易结构体
	tx := types.NewTransaction(tran.Nonce, tran.TokenAddress, tran.Amount, tran.GasLimit, tran.GasPrice, tran.Data)
	fmt.Println(tran.GasPrice)
	fmt.Println(tran.GasLimit)
	fmt.Println(tran.Amount)
	fmt.Println(common.BytesToHash(tran.Data).String())

	chainID, err := t.client.NetworkID(context.Background())
	if err != nil {
		logger.Errorf("Transfer NetworkID err", err)
		return "", err
	}

	// 打包签名交易哈希
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		logger.Errorf("Transfer NewEIP155Signer err", err)
		return "", nil
	}

	// 发送交易
	err = t.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		logger.Errorf("Transfer SendTransaction err", err)
		return "", nil
	}

	return signedTx.Hash().Hex(), nil
}
