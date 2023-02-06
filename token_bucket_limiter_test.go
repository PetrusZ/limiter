package limiter

import (
	"testing"
	"time"
)

func TestTokenBucketLimiter(t *testing.T) {
	type args struct {
		capacity int
		rate     int
	}

	testCases := []struct {
		name string
		args args
		want *TokenBucketLimiter
	}{
		{
			name: "60",
			args: args{
				capacity: 60,
				rate:     10,
			},
			want: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l := NewTokenBucketLimiter(tc.args.capacity, tc.args.rate)
			time.Sleep(time.Second)
			successCount := 0
			for i := 0; i < tc.args.capacity; i++ {
				if l.TryAcquire() {
					successCount++
				}
			}
			if successCount != tc.args.rate {
				t.Errorf("NewTokenBucketLimiter() got = %v, want %v", successCount, tc.args.rate)
				return
			}

			successCount = 0
			for i := 0; i < tc.args.capacity; i++ {
				if l.TryAcquire() {
					successCount++
				}
				time.Sleep(time.Second / 10)
			}
			if successCount != tc.args.capacity-tc.args.rate {
				t.Errorf("NewTokenBucketLimiter() got = %v, want %v", successCount, tc.args.capacity-tc.args.rate)
				return
			}
		})
	}
}
