package limiter

import (
	"testing"
	"time"
)

func TestFixedWindowLimiter(t *testing.T) {
	type args struct {
		limit  int
		window time.Duration
	}

	testCases := []struct {
		name string
		args args
		want *FixedWindowLimiter
	}{
		{
			name: "100_second",
			args: args{
				limit:  100,
				window: time.Second,
			},
			want: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l := NewFixedWindowLimiter(tc.args.limit, tc.args.window)
			successCount := 0
			for i := 0; i < tc.args.limit*2; i++ {
				if l.TryAcquire() {
					successCount++
				}
			}
			if successCount != tc.args.limit {
				t.Errorf("NewFixedWindowLimiter() = %v, want %v", successCount, tc.args.limit)
			}

			time.Sleep(time.Second)
			successCount = 0
			for i := 0; i < tc.args.limit*2; i++ {
				if l.TryAcquire() {
					successCount++
				}
			}
			if successCount != tc.args.limit {
				t.Errorf("NewFixedWindowLimiter() = %v, want %v", successCount, tc.args.limit)
			}
		})
	}
}
