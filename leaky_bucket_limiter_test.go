package limiter

import (
	"testing"
	"time"
)

func TestLeakyBucketLimiter(t *testing.T) {
	type args struct {
		peakLevel       int
		currentVelocity int
	}
	testCases := []struct {
		name string
		args args
		want *LeakyBucketLimiter
	}{
		{
			name: "60",
			args: args{
				peakLevel:       60,
				currentVelocity: 10,
			},
			want: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l := NewLeakyBucketLimiter(tc.args.peakLevel, tc.args.currentVelocity)
			successCount := 0
			for i := 0; i < tc.args.peakLevel; i++ {
				if l.TryAcquire() {
					successCount++
				}
			}
			if successCount != tc.args.peakLevel {
				t.Errorf("NewLeakyBucketLimiter() got = %v, want %v", successCount, tc.args.peakLevel)
				return
			}

			successCount = 0
			for i := 0; i < tc.args.peakLevel; i++ {
				if l.TryAcquire() {
					successCount++
				}
				time.Sleep(time.Second / 10)
			}
			if successCount != tc.args.peakLevel-tc.args.currentVelocity {
				t.Errorf("NewLeakyBucketLimiter() got = %v, want %v", successCount, tc.args.peakLevel-tc.args.currentVelocity)
				return
			}
		})
	}
}
