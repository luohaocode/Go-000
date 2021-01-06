package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Bucket struct {
	timestamp time.Time
	count     int
}

type Window struct {
	buckets []*Bucket
	mu      sync.RWMutex
}

func NewWindow(size int) *Window {
	return &Window{buckets: make([]*Bucket, 0, size)}
}

func countRequest(buckets []*Bucket) (count int) {
	for _, b := range buckets {
		count += b.count
	}
	return
}

type SlidingWindowCounter struct {
	SlotDuration time.Duration
	WinDuration  time.Duration
	window       *Window
	maxRequests  int
}

func NewSlidingWindowCounter(slotDuration time.Duration, winDuration time.Duration, maxRequests int) *SlidingWindowCounter {
	return &SlidingWindowCounter{
		SlotDuration: slotDuration,
		WinDuration:  winDuration,
		window:       NewWindow(int(winDuration / slotDuration)),
		maxRequests:  maxRequests,
	}
}

func (l *SlidingWindowCounter) Run() error {
	l.window.mu.Lock()
	defer l.window.mu.Unlock()

	now := time.Now()
	timeoutOffset := -1
	for i, ts := range l.window.buckets {
		if ts.timestamp.Add(l.WinDuration).After(now) {
			break
		}
		timeoutOffset = i
	}
	if timeoutOffset > -1 {
		l.window.buckets = l.window.buckets[timeoutOffset+1:]
	}

	if countRequest(l.window.buckets) >= l.maxRequests {
		return fmt.Errorf("too many request")
	}

	var lastSlot *Bucket
	if len(l.window.buckets) > 0 {
		lastSlot = l.window.buckets[len(l.window.buckets)-1]
		if lastSlot.timestamp.Add(l.SlotDuration).Before(now) {
			lastSlot = &Bucket{timestamp: now, count: 1}
			l.window.buckets = append(l.window.buckets, lastSlot)
		} else {
			lastSlot.count++
		}
	} else {
		lastSlot = &Bucket{timestamp: now, count: 1}
		l.window.buckets = append(l.window.buckets, lastSlot)
	}

	return nil
}

func main() {
	counter := NewSlidingWindowCounter(10*time.Millisecond, time.Second, 10)

	var count int32
	for i := 0; i < 20; i++ {
		go func() {
			if err := counter.Run(); err == nil {
				atomic.AddInt32(&count, 1)
			}
		}()
	}

	time.Sleep(time.Second)

	fmt.Println(count)
}
