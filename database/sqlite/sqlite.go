package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/a6910438/go-logger"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	db *sql.DB
}

func NewSqlite(dbFilePath string) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		return nil, err
	}
	return &Sqlite{db: db}, nil
}

func (s *Sqlite) GetLastBlockHeight() (int, error) {
	var height int
	row := s.db.QueryRow("SELECT HEIGHT FROM BLOCK_HEIGHT")
	err := row.Scan(&height)
	if err != nil {
		_, err := s.db.Exec("INSERT INTO BLOCK_HEIGHT(id, HEIGHT)  values(1, 9150000)")
		if err != nil {
			logger.Errorf("Sqlite GetLastBlockHeight Init:", err)
			return 0, err
		}
	}
	return height, nil
}

func (s *Sqlite) UpdateBlockHeight(height int) (int, error) {
	sql := "update BLOCK_HEIGHT set HEIGHT=%d where id= 1"
	_, err := s.db.Exec(fmt.Sprintf(sql, height))
	if err != nil {
		logger.Errorf("Sqlite GetLastBlockHeight Init:", err)
		return 0, err
	}
	return 0, nil
}
