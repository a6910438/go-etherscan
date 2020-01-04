package bolt

import (
	"etherscan-go/cmd/app/config"
	"etherscan-go/types"
	"github.com/a6910438/go-logger"
	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
)

type Bolt struct {
	db *storm.DB
}

func NewBolt() (*Bolt, error) {
	db, err := storm.Open("bolt.db")
	if err != nil {
		return nil, err
	}
	return &Bolt{db: db}, nil
}

func (b *Bolt) CreateTable() error {
	//也可以在这里对表做插入操作
	err := b.db.Save(&types.Block{ID: 1, Height: config.Cfg.Block.Height})
	//更新数据库失败
	if err != nil {
		logger.Errorf("Bolt CreateTable err", err)
		return err
	}
	return nil
}

func (b *Bolt) Update(block types.Block) error {
	err := b.db.Update(&block)
	//更新数据库失败
	if err != nil {
		logger.Errorf("Bolt Update err", err)
		return err
	}
	return nil
}

func (b *Bolt) Select() (types.Block, error) {
	var block types.Block
	err := b.db.Select(q.Eq("ID", 1)).First(&block)
	//更新数据库失败
	if err != nil {
		b.CreateTable()
		if err := b.db.Select(q.Eq("ID", 1)).First(&block); err != nil {
			logger.Errorf("Bolt Select err", err)
			return block, err
		}
	}
	return block, nil
}
