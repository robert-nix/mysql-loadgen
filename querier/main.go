package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/robert-nix/mysql-loadgen/internal/mysql"
)

var schemata map[string]struct{}

func main() {
	flag.Parse()

	db, err := mysql.Open()
	if err != nil {
		panic(err)
	}

	schemata, err = mysql.LoadSchemaNames(db)
	if err != nil {
		panic(err)
	}

	for i := 16; i <= 512; i *= 2 {
		measureQPS(db, i)
	}
}

const sampleTime = 15 * time.Second

func measureQPS(db *sql.DB, concurrency int) {
	var queries int64
	var errors int64
	closed := make(chan struct{})
	start := time.Now()
	for i := 0; i < concurrency; i++ {
		i := i
		go func() {
			r := rand.New(rand.NewSource(time.Now().UnixNano() ^ int64(i)))
			skipFirstN := r.Intn(len(schemata))
			for {
				for schema := range schemata {
					select {
					case <-closed:
						return
					default:
					}
					if skipFirstN > 0 {
						skipFirstN--
						continue
					}
					if r.Intn(2) == 0 {
						continue
					}
					err := execFetch(db, r, schema)
					if err != nil {
						atomic.AddInt64(&errors, 1)
					} else {
						atomic.AddInt64(&queries, 1)
					}
				}
			}
		}()
	}

	reportTick := time.NewTicker(time.Second)
	lastT := time.Now()
	for t := range reportTick.C {
		dur := t.Sub(lastT)
		lastT = t
		qs := atomic.LoadInt64(&queries)
		es := atomic.LoadInt64(&errors)
		atomic.StoreInt64(&queries, 0)
		atomic.StoreInt64(&errors, 0)
		qps := float64(qs) / dur.Seconds()
		eps := float64(es) / dur.Seconds()
		runDur := time.Since(start)
		fmt.Printf("%d\t%d\t%d\t%d\t%d\t%f\t%f\n", concurrency, runDur.Microseconds(), dur.Microseconds(), qs, es, qps, eps)
		if runDur > sampleTime {
			close(closed)
			break
		}
	}
}

var pageTitles = []string{
	"Page1",
	"Page2",
	"Page3",
	"Page4",
	"Page5",
	"Page6",
	"Page7",
}

func execFetch(db *sql.DB, r *rand.Rand, schema string) error {
	ctx := context.Background()
	conn, err := mysql.Use(ctx, db, schema)
	if err != nil {
		log.Printf("err changing DB: %v", err)
		return err
	}

	err = mysql.FetchRevision(ctx, conn, 0, pageTitles[r.Intn(len(pageTitles))])
	_ = conn.Close()
	// ErrNoRows indicates data inconsistency, so it's intentional to log it here
	if err != nil {
		log.Printf("err querying revision: %v", err)
		return err
	}
	return nil
}
