package limiter

import (
	"testing"
	"time"
)

func TestSlidingWindowLimiter(t *testing.T) {
	type args struct {
		limit       int
		window      int
		smallWindow int
	}
	testCases := []struct {
		name string
		args args
		want *SlidingWindowLimiter
	}{
		{
			name: "60_5seconds",
			args: args{
				limit:       60,
				window:      int(time.Second) * 5,
				smallWindow: int(time.Second),
			},
			want: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l, err := NewSlidingWindowLimiter(tc.args.limit, time.Duration(tc.args.window), time.Duration(tc.args.smallWindow))
			if err != nil {
				t.Errorf("NewSlidingWindowLimiter() error = %v", err)
				return
			}

			successCount := 0
			for i := 0; i < tc.args.limit/2; i++ {
				if l.TryAcquire() {
					successCount++
				}
			}
			if successCount != tc.args.limit/2 {
				t.Errorf("NewSlidingWindowLimiter() got = %v, want %v", successCount, tc.args.limit/2)
				return
			}

			time.Sleep(time.Second * 2)
			successCount = 0
			for i := 0; i < tc.args.limit/2; i++ {
				if l.TryAcquire() {
					successCount++
				}
			}
			if successCount != tc.args.limit/2 {
				t.Errorf("NewSlidingWindowLimiter() got = %v, want %v", successCount, tc.args.limit-tc.args.limit/2)
				return
			}

			time.Sleep(time.Second * 3)
			successCount = 0
			for i := 0; i < tc.args.limit/2; i++ {
				if l.TryAcquire() {
					successCount++
				}
			}
			if successCount != tc.args.limit/2 {
				t.Errorf("NewSlidingWindowLimiter() got = %v, want %v", successCount, tc.args.limit/2)
				return
			}
		})
	}
}
