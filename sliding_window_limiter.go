// 滑动窗口类似于固定窗口，它只是把大窗口切分成多个小窗口，每次向右移动一个小窗口，它可以避免两倍的突发流量。
// 固定窗口可以说是滑动窗口的一种特殊情况，只要滑动窗口里面的小窗口和大窗口大小一样。
// 窗口算法都有一个问题，当流量达到上限，后面的请求都会被拒绝。
package limiter

import (
	"fmt"
	"sync"
	"time"
)

type SlidingWindowLimiter struct {
	limit        int           // 窗口请求上线
	window       int64         // 窗口时间大小
	smallWindow  int64         // 小窗口时间大小
	smallWindows int64         // 小窗口数量
	counters     map[int64]int // 小窗口计数器
	mutex        sync.Mutex    // 避免并发问题
}

func NewSlidingWindowLimiter(limit int, window, smallWindow time.Duration) (Limiter, error) {
	// 窗口时间必须能够被小窗口时间整除
	if window%smallWindow != 0 {
		return nil, fmt.Errorf("window cannot be split by integers")
	}

	return &SlidingWindowLimiter{
		limit:        limit,
		window:       int64(window),
		smallWindow:  int64(smallWindow),
		smallWindows: int64(window / smallWindow),
		counters:     make(map[int64]int),
	}, nil
}

func (l *SlidingWindowLimiter) TryAcquire() bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	// 获取当前小窗口值
	currentSmallWindow := time.Now().UnixNano() / l.smallWindow * l.smallWindow
	// 获取起始小窗口值
	startSmallWindow := currentSmallWindow - l.smallWindow*(l.smallWindows-1)

	// 计算当前窗口的请求总数
	count := 0
	for smallWindow, counter := range l.counters {
		if smallWindow < startSmallWindow {
			delete(l.counters, smallWindow)
		} else {
			count += counter
		}
	}

	// 若到达窗口请求上限，请求失败
	if count > l.limit {
		return false
	}

	// 若没到窗口请求上限，当前小窗口计数器+1，请求成功
	l.counters[currentSmallWindow]++

	return true
}
