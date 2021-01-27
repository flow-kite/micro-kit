package lrubucket

import (
	"sync"
	"time"
)

type Bucket struct {
	locker *sync.Mutex
	bulk   map[time.Time]map[string]bool //

	timeoutTicker *time.Ticker
}
