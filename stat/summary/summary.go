package summary

import (
	"sync"
	"time"
)

type Summary interface {
	Add(value int64)
	Value() (value, count int64)
	Reset()
}

type bucket struct {
	next  *bucket
	value int64
	count int64
}

func (b *bucket) Add(value int64) {
	b.value += value
	b.count++
}

func (b *bucket) Reset() {
	b.value = 0
	b.count = 0
}

type summary struct {
	buckets         []bucket
	cur             *bucket
	lastTime        int64
	widthPerWindows int64
	lock            sync.RWMutex
}

func New(window time.Duration, windowCount int) Summary {
	buckets := make([]bucket, windowCount, windowCount)
	for i := 1; i < len(buckets); i++ {
		buckets[i-1].next = &buckets[i]
	}
	buckets[len(buckets)-1].next = &buckets[0]
	widthPerWindows := window.Nanoseconds() / int64(windowCount)
	return &summary{
		buckets:         buckets,
		cur:             &buckets[0],
		lastTime:        time.Now().UnixNano(),
		widthPerWindows: widthPerWindows,
	}
}

func (s *summary) Add(value int64) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.lastBucket().Add(value)
}

func (s *summary) Value() (value, count int64) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.lastBucket()
	for _, b := range s.buckets {
		value += b.value
		count += b.count
	}
	return value, count
}

func (s *summary) Reset() {
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, b := range s.buckets {
		b.Reset()
	}
}

// 获取当前应该写哪个 bucket, 并把期间未使用的 bucket reset
func (s *summary) lastBucket() *bucket {
	now := time.Now().UnixNano()
	i := (now - s.lastTime) / s.widthPerWindows
	if i > int64(len(s.buckets)) {
		i = int64(len(s.buckets))
	}
	if i > 0 {
		s.lastTime = now
	}
	for ;i > 0; i-- {
		s.cur = s.cur.next
		s.cur.Reset()
	}
	return s.cur
}