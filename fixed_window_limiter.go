// 每开启一个新的窗口，在窗口时间大小内，可以通过窗口请求上限个请求。
// 该算法主要是会存在临界问题，如果流量都集中在两个窗口的交界处，那么突发流量会是设置上限的两倍。
package limiter

import (
	"sync"
	"time"
)

type FixedWindowLimiter struct {
	limit    int           // 窗口请求上限
	window   time.Duration // 窗口时间大小
	counter  int           // 计数器
	lastTime time.Time     // 上一次请求的时间
	mutex    sync.Mutex    // 避免并发问题
}

func NewFixedWindowLimiter(limit int, window time.Duration) Limiter {
	return &FixedWindowLimiter{
		limit:  limit,
		window: window,
	}
}

func (l *FixedWindowLimiter) TryAcquire() bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	now := time.Now()
	// 如果当前窗口失效，计数器清0，开启新的窗口
	if now.Sub(l.lastTime) > l.window {
		l.counter = 0
		l.lastTime = now
	}

	// 若到达窗口请求上限，请求失败
	if l.counter >= l.limit {
		return false
	}

	// 若没到窗口请求上限，计数器+1，请求成功
	l.counter++

	return true
}
