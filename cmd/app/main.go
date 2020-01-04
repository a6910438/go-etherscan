package main

import (
	"etherscan-go/cmd/app/config"
	transfer2 "etherscan-go/cmd/app/transfer"
	"etherscan-go/cmd/app/watch"
	"etherscan-go/database/bolt"
	"etherscan-go/database/mysql"
	"etherscan-go/log"
	"flag"
	"fmt"
	"github.com/a6910438/go-logger"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/judwhite/go-svc/svc"
)

type program struct {
	watch    *watch.Watcher //监听
	transfer *transfer2.Transfer
}

func (p *program) Init(env svc.Environment) error {
	//初始化日志
	if err := log.Init(config.Cfg.Log.Path, config.Cfg.Log.File, config.Cfg.Log.Level); err != nil {
		fmt.Println("init log error:", err)
		return err
	}
	//初始化数据库
	mysql, err := mysql.NewMysql(config.Cfg.Db)
	if err != nil {
		logger.Errorf("init db error: %v", err)
		return err
	}

	bolt, err := bolt.NewBolt()
	if err != nil {
		logger.Errorf("init boltdb error: %v", err)
		return err
	}

	watch, err := watch.NewWatcher(mysql, bolt, config.Cfg.Eth.Node)
	if err != nil {
		logger.Errorf("NewWatcher error: %v", err)
		return err
	}

	client, err := ethclient.Dial(config.Cfg.Eth.Node)
	if err != nil {
		logger.Errorf("Dial error: %v", err)
		return err
	}

	transfer, err := transfer2.NewTransfer(mysql, client)

	p.watch = watch
	p.transfer = transfer

	logger.Info("program inited")
	return nil
}

func (p *program) Start() error {

	logger.Info("program start")

	go p.transfer.Start()
	go p.watch.Start()

	return nil
}

func (p *program) Stop() error {

	logger.Info("program stopped")
	return nil
}

func main() {
	cfg := flag.String("C", "config.json", "configuration file")
	flag.Parse()

	if err := config.Init(*cfg); err != nil {
		fmt.Println("init config error:", err.Error())
		return
	}

	app := &program{}
	if err := svc.Run(app); err != nil {
		logger.Println(err)
	}
}
