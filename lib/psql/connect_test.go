package psql

import (
	"context"
	"fmt"
	"management_users/lib/cfg"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestConnect(t *testing.T) {

	config := cfg.Config{
		SqlType: "postgres",
		SqlHost: "192.168.20.114",
		SqlPort: "5400",
		SqlUser: "postgres",
		SqlPass: "postgres",
		SqlDb:   "storage",
	}

	pool, err := InitPgxPool(&config)
	if err != nil {
		t.Fatalf("Unable set psql connection: %s", err.Error())
	}

	type benchmark struct {
		Tests     int
		Successes int
		Time      int64
		Max       int64
		Min       int64
		Avg       int64
		mux       *sync.Mutex
	}

	t.Run("test_ping", func(t *testing.T) {
		wg := sync.WaitGroup{}
		startCh := make(chan bool)

		b := benchmark{
			Tests: 1000,
			mux:   new(sync.Mutex),
		}

		for i := 1; i <= b.Tests; i++ {
			wg.Add(1)

			go func(start chan bool, b *benchmark) {
				defer wg.Done()
				<-start

				now := time.Now()
				conn, err := pool.Acquire(context.Background())
				if err != nil {
					return
				}
				defer conn.Release()

				err = conn.Ping(context.Background())
				if err != nil {
					return
				}

				difference := time.Now().Sub(now).Nanoseconds()
				b.mux.Lock()
				b.Successes += 1
				b.Time += difference
				if b.Max < difference {
					b.Max = difference
				}
				if b.Min > difference || b.Min == 0 {
					b.Min = difference
				}
				b.mux.Unlock()
			}(startCh, &b)
		}

		close(startCh)
		wg.Wait()
		b.Avg = b.Time / int64(b.Tests)

		var lines []string
		lines = append(lines, fmt.Sprintf("\nResults of %d tests:", b.Tests))
		lines = append(lines, fmt.Sprintf("	Successes: %d", b.Successes))
		lines = append(lines, fmt.Sprintf("	Min: %.2fms", float64(b.Min)/1000000))
		lines = append(lines, fmt.Sprintf("	Max: %.2fms", float64(b.Max)/1000000))
		lines = append(lines, fmt.Sprintf("	Avg: %.2fms", float64(b.Avg)/1000000))
		t.Log(strings.Join(lines, "\n"))
	})
}
