package mysql

import (
	"database/sql"
	"flag"
	"strings"

	"github.com/go-sql-driver/mysql"
)

var user, password, addr string

func init() {
	flag.StringVar(&user, "user", "root", "mysql user")
	flag.StringVar(&password, "password", "", "mysql password")
	flag.StringVar(&addr, "addr", "127.0.0.1:3306", "mysql addr")
}

type Option int

const (
	EnableMultiStatements Option = iota
)

func Open(options ...Option) (*sql.DB, error) {
	cfg := mysql.NewConfig()
	cfg.User = user
	cfg.Passwd = password
	cfg.Net = "tcp"
	cfg.Addr = addr
	cfg.InterpolateParams = true
	for _, opt := range options {
		switch opt {
		case EnableMultiStatements:
			cfg.MultiStatements = true
		}
	}
	var err error
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(0)
	return db, err
}

const MWDBPrefix = "mwdb"

func LoadSchemaNames(db *sql.DB) (map[string]struct{}, error) {
	rows, err := db.Query("select schema_name from information_schema.schemata")
	if err != nil {
		return nil, err
	}

	schemata := map[string]struct{}{}
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			_ = rows.Close()
			return nil, err
		}
		if !strings.HasPrefix(name, MWDBPrefix) {
			continue
		}
		schemata[name] = struct{}{}
	}
	return schemata, nil
}
