// 与漏桶算法的相反，令牌桶会不断地把令牌添加到桶里，而请求会从桶中获取令牌，只有拥有令牌地请求才能被接受。
// 因为桶中可以提前保留一些令牌，所以它允许一定地突发流量通过。

package limiter

import (
	"sync"
	"time"
)

type TokenBucketLimiter struct {
	capacity      int        // 容量
	currentTokens int        // 令牌数量
	rate          int        // 发放令牌速率/秒
	lastTime      time.Time  // 上次发放令牌时间
	mutex         sync.Mutex // 避免并发问题
}

func NewTokenBucketLimiter(capacity, rate int) Limiter {
	return &TokenBucketLimiter{
		capacity: capacity,
		rate:     rate,
		lastTime: time.Now(),
	}
}

func (l *TokenBucketLimiter) TryAcquire() bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	// 尝试发放令牌
	now := time.Now()
	// 距离上次发放令牌的时间
	interval := now.Sub(l.lastTime)
	if interval >= time.Second {
		// 当前令牌数量+距离上次发放令牌的时间(秒)*发放令牌速率
		l.currentTokens = MinInt(l.capacity, l.currentTokens+int(interval/time.Second)*l.rate)
		l.lastTime = now
	}

	// 如果没有令牌，请求失败
	if l.currentTokens == 0 {
		return false
	}
	// 如果有令牌，当前令牌-1，请求成功
	l.currentTokens--
	return true
}

func MinInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}
