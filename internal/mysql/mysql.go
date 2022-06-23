package mysql

import (
	"context"
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

func Use(ctx context.Context, db *sql.DB, schema string) (*sql.Conn, error) {
	conn, err := db.Conn(ctx)
	if err != nil {
		return nil, err
	}
	_, err = conn.ExecContext(ctx, "USE `"+schema+"`")
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	return conn, nil
}

func FetchRevision(ctx context.Context, conn *sql.Conn, namespace int, title string) error {
	var revID, revPage, revMinorEdit, revDeleted, revLen, revParentID, revCommentCID, revActor, pageNamespace, pageID, pageLatest, pageIsRedirect, pageLen int
	var revTimestamp, revSha1, revCommentText, revCommentData, revUser, revUserText, pageTitle, userName []byte

	return conn.QueryRowContext(ctx, "SELECT  rev_id,rev_page,rev_timestamp,rev_minor_edit,rev_deleted,rev_len,rev_parent_id,rev_sha1,comment_rev_comment.comment_text AS `rev_comment_text`,comment_rev_comment.comment_data AS `rev_comment_data`,comment_rev_comment.comment_id AS `rev_comment_cid`,actor_rev_user.actor_user AS `rev_user`,actor_rev_user.actor_name AS `rev_user_text`,temp_rev_user.revactor_actor AS `rev_actor`,page_namespace,page_title,page_id,page_latest,page_is_redirect,page_len,user_name  FROM `revision` JOIN `revision_comment_temp` `temp_rev_comment` ON ((temp_rev_comment.revcomment_rev = rev_id)) JOIN `comment` `comment_rev_comment` ON ((comment_rev_comment.comment_id = temp_rev_comment.revcomment_comment_id)) JOIN `revision_actor_temp` `temp_rev_user` ON ((temp_rev_user.revactor_rev = rev_id)) JOIN `shared`.`actor` `actor_rev_user` ON ((actor_rev_user.actor_id = temp_rev_user.revactor_actor)) JOIN `page` ON ((page_id = rev_page)) LEFT JOIN `shared`.`user` ON ((actor_rev_user.actor_user != 0) AND (user_id = actor_rev_user.actor_user))   WHERE page_namespace = ? AND page_title = ? AND (rev_id=page_latest)  LIMIT 1", 0, title).Scan(
		&revID, &revPage, &revTimestamp, &revMinorEdit,
		&revDeleted, &revLen, &revParentID, &revSha1,
		&revCommentText, &revCommentData, &revCommentCID,
		&revUser, &revUserText, &revActor,
		&pageNamespace, &pageTitle, &pageID,
		&pageLatest, &pageIsRedirect, &pageLen,
		&userName,
	)
}
