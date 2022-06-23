package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"sync"
	"sync/atomic"
	"time"

	"github.com/robert-nix/mysql-loadgen/internal/mysql"
	"github.com/schollz/progressbar/v3"
)

func main() {
	var concurrency int
	var dbCount int
	flag.IntVar(&concurrency, "concurrency", 32, "number of concurrent database imports to run")
	flag.IntVar(&dbCount, "databases", 90000, "total number of mwdb* databases to create")
	flag.Parse()

	db, err := mysql.Open(mysql.EnableMultiStatements)
	if err != nil {
		panic(err)
	}

	schemata, err := mysql.LoadSchemaNames(db)
	if err != nil {
		panic(err)
	}

	if _, ok := schemata["shared"]; !ok {
		fmt.Printf("creating shared db ...")
		sharedbSQL, err := loadSQLFixture("fixtures/shared.sql")
		if err != nil {
			panic(err)
		}
		err = createDatabase(db, "shared", sharedbSQL)
		if err != nil {
			panic(err)
		}
		fmt.Printf(" done\n")
	}

	mediawikiSQL, err := loadSQLFixture("fixtures/mediawiki.sql")
	if err != nil {
		panic(err)
	}

	var loaded int64

	nameChan := make(chan string, concurrency)
	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		go func() {
			for name := range nameChan {
				if name == "" {
					break
				}
				err = createDatabase(db, name, mediawikiSQL)
				if err != nil {
					panic(err)
				}
				atomic.AddInt64(&loaded, 1)
			}
			wg.Done()
		}()
	}

	bar := progressbar.Default(int64(dbCount), "creating mwdbs")
	go func() {
		ticker := time.NewTicker(time.Second)
		for range ticker.C {
			bar.Set64(atomic.LoadInt64(&loaded))
		}
	}()

	for i := 0; i < dbCount; i++ {
		name := fmt.Sprintf("mwdb%010d", i)
		if _, ok := schemata[name]; ok {
			atomic.AddInt64(&loaded, 1)
			continue
		}
		nameChan <- name
	}
	close(nameChan)
	wg.Wait()

	bar.Set64(int64(dbCount))
	fmt.Println("finished")
}

func loadSQLFixture(filename string) (string, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func createDatabase(db *sql.DB, schema, initSQL string) error {
	ctx := context.Background()
	conn, err := db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.ExecContext(ctx, "CREATE DATABASE `"+schema+"`")
	if err != nil {
		return err
	}

	_, err = conn.ExecContext(ctx, "USE `"+schema+"`")
	if err != nil {
		return err
	}

	_, err = conn.ExecContext(ctx, initSQL)
	if err != nil {
		return err
	}
	return nil
}
