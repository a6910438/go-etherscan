package main

import (
	"etherscan-go/cmd/app/config"
	"etherscan-go/cmd/app/watch"
	"etherscan-go/database/mysql"
	"etherscan-go/database/sqlite"
	"etherscan-go/log"
	"flag"
	"fmt"
	"github.com/a6910438/go-logger"
	"github.com/judwhite/go-svc/svc"
)

type program struct {
	watch *watch.Watcher //监听
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

	sqlite, err := sqlite.NewSqlite("D:/etherscan-go/build/deposit.db")
	if err != nil {
		logger.Errorf("init sqlite error: %v", err)
		return err
	}

	watch, err := watch.NewWatcher(mysql, sqlite, config.Cfg.Eth.Node)
	if err != nil {
		logger.Errorf("NewWatcher error: %v", err)
		return err
	}

	p.watch = watch

	logger.Info("program inited")
	return nil
}

func (p *program) Start() error {

	logger.Info("program start")

	p.watch.Start()
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
