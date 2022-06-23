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
					conn, err := db.Conn(context.Background())
					if err != nil {
						atomic.AddInt64(&errors, 1)
						continue
					}
					err = execFetch(conn, r, schema)
					conn.Close()
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

func execFetch(conn *sql.Conn, r *rand.Rand, schema string) error {
	title := pageTitles[r.Intn(len(pageTitles))]

	var revID, revPage, revMinorEdit, revDeleted, revLen, revParentID, revCommentCID, revActor, pageNamespace, pageID, pageLatest, pageIsRedirect, pageLen int
	var revTimestamp, revSha1, revCommentText, revCommentData, revUser, revUserText, pageTitle, userName []byte
	ctx := context.Background()
	_, err := conn.ExecContext(ctx, "use `"+schema+"`")
	if err != nil {
		log.Printf("err changing DB: %v", err)
		return err
	}
	err = conn.QueryRowContext(ctx, "SELECT  rev_id,rev_page,rev_timestamp,rev_minor_edit,rev_deleted,rev_len,rev_parent_id,rev_sha1,comment_rev_comment.comment_text AS `rev_comment_text`,comment_rev_comment.comment_data AS `rev_comment_data`,comment_rev_comment.comment_id AS `rev_comment_cid`,actor_rev_user.actor_user AS `rev_user`,actor_rev_user.actor_name AS `rev_user_text`,temp_rev_user.revactor_actor AS `rev_actor`,page_namespace,page_title,page_id,page_latest,page_is_redirect,page_len,user_name  FROM `revision` JOIN `revision_comment_temp` `temp_rev_comment` ON ((temp_rev_comment.revcomment_rev = rev_id)) JOIN `comment` `comment_rev_comment` ON ((comment_rev_comment.comment_id = temp_rev_comment.revcomment_comment_id)) JOIN `revision_actor_temp` `temp_rev_user` ON ((temp_rev_user.revactor_rev = rev_id)) JOIN `shared`.`actor` `actor_rev_user` ON ((actor_rev_user.actor_id = temp_rev_user.revactor_actor)) JOIN `page` ON ((page_id = rev_page)) LEFT JOIN `shared`.`user` ON ((actor_rev_user.actor_user != 0) AND (user_id = actor_rev_user.actor_user))   WHERE page_namespace = ? AND page_title = ? AND (rev_id=page_latest)  LIMIT 1", 0, title).Scan(
		&revID, &revPage, &revTimestamp, &revMinorEdit,
		&revDeleted, &revLen, &revParentID, &revSha1,
		&revCommentText, &revCommentData, &revCommentCID,
		&revUser, &revUserText, &revActor,
		&pageNamespace, &pageTitle, &pageID,
		&pageLatest, &pageIsRedirect, &pageLen,
		&userName,
	)
	// ErrNoRows indicates data inconsistency, so it's intentional to log it here
	if err != nil {
		log.Printf("err querying revision: %v", err)
		return err
	}
	return nil
}
