package pkg

import (
	"sync"
	"time"
)

type TokenBuket struct {
	cap int64
	tokens  int64
	lock sync.Mutex
	rate int64   //每秒加入令牌速率
	lastTime  int64   // 上一次加入的时间
}

func NewTokenBuket(cap int64, rate int64) *TokenBuket {
	if cap <= 0 || rate <= 0 {
		panic("bucket cap or rate can't be negative")
	}
	bucket := &TokenBuket{cap: cap, tokens:cap, rate:rate}

	return  bucket
}

func (bucket *TokenBuket)addToken()  {
	bucket.lock.Lock()
	defer bucket.lock.Unlock()
	if bucket.tokens + bucket.rate <= bucket.cap {
		bucket.tokens  += bucket.rate
	} else {
		bucket.tokens = bucket.cap
	}
}

func (bucket *TokenBuket) IsAccept()  bool {
	bucket.lock.Lock()
	defer bucket.lock.Unlock()

	now := time.Now().Unix()
	bucket.tokens = bucket.tokens + (now - bucket.lastTime) * bucket.rate

	if bucket.tokens > bucket.cap {
		bucket.tokens = bucket.cap
	}
	bucket.lastTime = now

	if bucket.tokens > 0 {
		bucket.tokens--
		return true
	}
	return false
}
