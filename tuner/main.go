package main

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/prometheus/procfs"
	"github.com/robert-nix/mysql-loadgen/internal/mysql"
)

func main() {
	t := &tuner{
		option:      "table-definition-cache",
		min:         1000,
		max:         501000,
		step:        100000,
		concurrency: 32,
		sampleTime:  90 * time.Second,
	}

	t.printHeader()
	err := t.tune()
	if err != nil {
		panic(err)
	}
}

type mysqlInstance struct {
	cmd         *exec.Cmd
	ready, done chan struct{}

	proc procfs.Proc
}

func startMysql(args ...string) (*mysqlInstance, error) {
	bashArgs := append([]string{"./scripts/run-mysql.sh"}, args...)
	cmd := exec.Command("bash", bashArgs...)
	instance := &mysqlInstance{
		cmd:   cmd,
		ready: make(chan struct{}, 1),
		done:  make(chan struct{}),
	}
	cmd.Stderr = instance
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	go func() {
		cmd.Wait()
		close(instance.done)
		close(instance.ready)
	}()
	instance.proc, err = procfs.NewProc(cmd.Process.Pid)
	if err != nil {
		// only happens when /proc isn't available
		panic(err)
	}
	<-instance.ready
	return instance, err
}

func (i *mysqlInstance) Write(p []byte) (n int, err error) {
	if bytes.Contains(p, []byte("[ERROR]")) {
		fmt.Printf("%s", string(p))
	}
	if bytes.Contains(p, []byte("ready for connections")) {
		i.ready <- struct{}{}
	}
	return len(p), nil
}

func (i *mysqlInstance) close() {
	_ = i.cmd.Process.Signal(syscall.SIGTERM)
	<-i.done
}

func (i *mysqlInstance) stat() (procfs.ProcStat, error) {
	return i.proc.Stat()
}

type tuner struct {
	option         string
	min, max, step int
	concurrency    int

	sampleTime time.Duration

	queries, errors int64
}

func (t *tuner) tune() error {
	for v := t.min; v <= t.max; v += t.step {
		err := t.sample(v)
		if err != nil {
			return err
		}
	}
	return nil
}

var schemata map[string]struct{}

func (t *tuner) printHeader() {
	fmt.Printf("%s\truntimeMS\tittimeUS\tqueries\terrors\trssPages\tutimeTicks\topenedTables\n", t.option)
}

func (t *tuner) sample(v int) error {
	inst, err := startMysql(fmt.Sprintf("--%s=%d", t.option, v))
	if err != nil {
		return err
	}
	defer inst.close()

	db, err := mysql.Open()
	if err != nil {
		return err
	}
	defer db.Close()

	if schemata == nil {
		schemata, err = mysql.LoadSchemaNames(db)
		if err != nil {
			return err
		}
	}

	start := time.Now()
	done := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < t.concurrency; i++ {
		wg.Add(1)
		go t.sampleThread(db, done, &wg, i)
	}

	tick := time.NewTicker(100 * time.Millisecond)
	last := start
	var lastUTime uint
	{
		procStat, _ := inst.stat()
		lastUTime = procStat.UTime
	}
	lastOpenedTables, _ := fetchGlobalStatusVar(db, "Opened_tables")
	for {
		now := <-tick.C
		totalDur := now.Sub(start)
		itDur := now.Sub(last)
		last = now

		qs := atomic.SwapInt64(&t.queries, 0)
		es := atomic.SwapInt64(&t.errors, 0)
		procStat, _ := inst.stat()
		dUTime := procStat.UTime - lastUTime
		lastUTime = procStat.UTime
		openedTables, _ := fetchGlobalStatusVar(db, "Opened_tables")
		dOpenedTables := openedTables - lastOpenedTables
		lastOpenedTables = openedTables
		fmt.Printf("%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\n", v, totalDur.Milliseconds(), itDur.Microseconds(), qs, es, procStat.RSS, dUTime, dOpenedTables)
		if totalDur >= t.sampleTime {
			close(done)
			break
		}
	}

	wg.Wait()
	return nil
}

func fetchGlobalStatusVar(db *sql.DB, name string) (int, error) {
	var retName string
	var value int
	err := db.QueryRow("SHOW GLOBAL STATUS LIKE ?", name).Scan(&retName, &value)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func (t *tuner) sampleThread(db *sql.DB, done chan struct{}, wg *sync.WaitGroup, n int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano() ^ int64(n)))
	skipFirstN := r.Intn(len(schemata))
	for {
		for schema := range schemata {
			select {
			case <-done:
				wg.Done()
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
				atomic.AddInt64(&t.errors, 1)
			} else {
				atomic.AddInt64(&t.queries, 1)
			}
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
