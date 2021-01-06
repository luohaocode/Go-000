package main

import (
	"fmt"
	"sync"
	"time"
)

type Bucket struct {
	timestamp time.Time
	count     int
}

type Window struct {
	buckets        []*Bucket
	BucketDuration time.Duration
	WinDuration    time.Duration
	mu             sync.RWMutex
}

func NewWindow(size int, bucketDuration time.Duration, winDuration time.Duration) *Window {
	return &Window{buckets: make([]*Bucket, 0, size), BucketDuration: bucketDuration, WinDuration: winDuration}
}

type SlidingWindowCounter struct {
	window      *Window
	maxRequests int
}

func NewSlidingWindowCounter(bucketDuration time.Duration, winDuration time.Duration, maxRequests int) *SlidingWindowCounter {
	return &SlidingWindowCounter{
		window:      NewWindow(int(winDuration/bucketDuration), bucketDuration, winDuration),
		maxRequests: maxRequests,
	}
}

func (c *SlidingWindowCounter) Run() error {

	if c.Sum() >= c.maxRequests {
		return fmt.Errorf("too many request")
	}

	c.Increment()

	return nil
}

func (c *SlidingWindowCounter) getCurrentBucket() *Bucket {
	now := time.Now()
	var lastSlot *Bucket
	if len(c.window.buckets) == 0 {
		lastSlot = &Bucket{timestamp: now, count: 0}
		c.window.buckets = append(c.window.buckets, lastSlot)
	}

	lastSlot = c.window.buckets[len(c.window.buckets)-1]
	if lastSlot.timestamp.Add(c.window.BucketDuration).Before(now) || lastSlot.timestamp.Add(c.window.BucketDuration).Equal(now) {
		lastSlot = &Bucket{timestamp: now, count: 0}
		c.window.buckets = append(c.window.buckets, lastSlot)
	}

	return lastSlot
}

func (c *SlidingWindowCounter) RemoveOldBucket() {
	now := time.Now()
	timeoutOffset := -1
	for i, b := range c.window.buckets {
		if b.timestamp.Add(c.window.WinDuration).After(now) || b.timestamp.Add(c.window.WinDuration).Equal(now) {
			break
		}
		timeoutOffset = i
	}
	if timeoutOffset > -1 {
		c.window.buckets = c.window.buckets[timeoutOffset+1:]
	}
}

func (c *SlidingWindowCounter) Increment() {
	c.window.mu.Lock()
	defer c.window.mu.Unlock()

	b := c.getCurrentBucket()
	b.count++

	c.RemoveOldBucket()
}

func (c *SlidingWindowCounter) Sum() int {
	c.window.mu.RLock()
	defer c.window.mu.RUnlock()
	count := 0
	for _, b := range c.window.buckets {
		count += b.count
	}
	return count
}

func (c *SlidingWindowCounter) Avg() int {
	return c.Sum() / (int(c.window.WinDuration / c.window.BucketDuration))
}

func main() {
	counter := NewSlidingWindowCounter(10*time.Millisecond, time.Second, 10)

	for i := 0; i < 20; i++ {
		go func() {
			counter.Run()
		}()
	}

	time.Sleep(time.Second)

	fmt.Println(counter.Sum())
}
