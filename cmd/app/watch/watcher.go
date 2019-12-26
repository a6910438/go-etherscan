package watch

import (
	"etherscan-go/types"
	"fmt"
	"github.com/a6910438/go-logger"
	"github.com/ethereum/go-ethereum/params"
	"github.com/onrik/ethrpc"
	"github.com/shopspring/decimal"
	"math/big"
	"time"
)

type dbBaser interface {
	Add(deposit types.Deposit) (id int64, err error)
	AddRecord(record types.Record, tableName string) (id int64, err error)
	GetDepositByHash(hash string) (*types.Deposit, error)
	GetCoinByAddress(contractAddress string) (*types.Coin, error)
	GetAssestByAddress(to string) (*types.UserAssest, error)
	UpdateAssestById(assest *types.UserAssest) (err error)
}

type sqlite interface {
	GetLastBlockHeight() (int, error)
	UpdateBlockHeight(height int) (int, error)
}

type Watcher struct {
	db   dbBaser
	sql  sqlite
	node string
}

func NewWatcher(db dbBaser, sql sqlite, node string) (*Watcher, error) {
	return &Watcher{
		db:   db,
		sql:  sql,
		node: node,
	}, nil
}

func (w *Watcher) Start() {
	client := ethrpc.New(w.node)

	height, err := w.sql.GetLastBlockHeight()
	if err != nil {
		logger.Errorf("Watcher Init Sqlite3 err", err)
		return
	}

	for {
		bestHeight, err := client.EthBlockNumber()
		if err != nil {
			logger.Errorf("Watcher EthBlockNumber err", err)
			time.Sleep(5 * time.Second)
			continue
		}
		for i := height; i <= bestHeight; i++ {
			count, err := client.EthGetBlockTransactionCountByNumber(i)
			if err != nil {
				logger.Errorf("Watcher EthGetBlockTransactionCountByNumber err", err)
				time.Sleep(5 * time.Second)
				continue
			}
			logger.Info(fmt.Sprintf("Start to synchronize %d blocks, the number of transactions is %d", i, count))
			for j := 0; j < count; j++ {
				transaction, err := client.EthGetTransactionByBlockNumberAndIndex(i, j)
				if err != nil {
					logger.Errorf("Watcher EthGetTransactionByBlockNumberAndIndex err", err)
					time.Sleep(5 * time.Second)
					continue
				}
				_, err = w.db.GetDepositByHash(transaction.Hash)
				if err == nil {
					// 数据库有重复充值记录
					continue
				}

				coin, err := w.db.GetCoinByAddress(transaction.To)
				if err != nil {
					// 合约地址不是数据库存放地址
					continue
				}

				txRe, err := client.EthGetTransactionReceipt(transaction.Hash)
				if err != nil {
					logger.Errorf("Watcher EthGetTransactionReceipt err", err)
					time.Sleep(5 * time.Second)
					continue
				}

				if txRe != nil && txRe.Status == "0x1" && len(txRe.Logs) == 1 && len(txRe.Logs[0].Topics) == 3 {
					to := "0x" + txRe.Logs[0].Topics[2][26:]
					amount := decimal.NewFromBigInt(hexToBigInt(txRe.Logs[0].Data), 0).Div(decimal.NewFromBigInt(big.NewInt(params.Ether), 0))
					fee := decimal.NewFromBigInt(big.NewInt(int64(txRe.GasUsed)).Mul(big.NewInt(int64(txRe.GasUsed)), &transaction.GasPrice), 0).Div(decimal.NewFromBigInt(big.NewInt(params.Ether), 0))
					fmt.Print(fee)
					assest, err := w.db.GetAssestByAddress(to)
					if err != nil {
						time.Sleep(5 * time.Second)
						continue
					}
					assest.AvaBalance = assest.AvaBalance.Add(amount)
					w.db.UpdateAssestById(assest)

					time := time.Now().Unix()
					w.db.Add(types.Deposit{Amount: amount, UserId: assest.UserId, CoinId: coin.Id, Receice: to, Fee: fee, Send: transaction.From, Hash: transaction.Hash, CreateTime: time})

				}
			}
			// 更新块高
			w.sql.UpdateBlockHeight(i)

		}
	}
}

func hexToBigInt(hex string) *big.Int {
	n := new(big.Int)
	n, _ = n.SetString(hex[2:], 16)

	return n
}
