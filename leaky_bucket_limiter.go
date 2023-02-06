// 漏桶是模拟一个漏水的桶，请求相当于往桶里倒水，处理请求的速度相当于水漏出的速度。
// 主要用于请求处理速率较为稳定的服务，需要使用生产者消费者模式把请求放到一个队列里，让消费者以一个较为稳定的速率处理。

package limiter

import (
	"sync"
	"time"
)

type LeakyBucketLimiter struct {
	peakLevel       int        // 最高水位
	currentLevel    int        // 当前水位
	currentVelocity int        // 水流速度/秒
	lasTime         time.Time  // 上次放水时间
	mutex           sync.Mutex // 避免并发问题
}

func NewLeakyBucketLimiter(peakLevel, currentVelocity int) Limiter {
	return &LeakyBucketLimiter{
		peakLevel:       peakLevel,
		currentVelocity: currentVelocity,
	}
}

func (l *LeakyBucketLimiter) TryAcquire() bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	// 尝试放水
	now := time.Now()
	// 距离上次放水时间
	interval := now.Sub(l.lasTime)
	if interval >= time.Second {
		// 当前水位-距离上次放水的时间(秒)*水流速度
		l.currentLevel = maxInt(0, l.currentLevel-int(interval/time.Second)*l.currentVelocity)
		l.lasTime = now
	}

	// 若到达最高水位，请求失败
	if l.currentLevel > l.peakLevel {
		return false
	}
	// 若没有到达最高水位，当前水位+1，请求成功
	l.currentLevel++
	return true
}

func maxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}
