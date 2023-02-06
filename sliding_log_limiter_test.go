package limiter

import (
	"testing"
	"time"
)

func TestSlidingLogLimiter(t *testing.T) {
	type args struct {
		smallWindow int
		strategies  []*SlidingLogLimiterStrategy
	}
	testCases := []struct {
		name string
		args args
		want *SlidingLogLimiter
	}{
		{
			name: "60_5seconds",
			args: args{
				smallWindow: int(time.Second),
				strategies: []*SlidingLogLimiterStrategy{
					NewSlidingLogLimiterStrategy(10, time.Minute),
					NewSlidingLogLimiterStrategy(100, time.Hour),
				},
			},
			want: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewNewSlidingLogLimiter(time.Duration(tc.args.smallWindow), tc.args.strategies...)
			if err != nil {
				return
			}
		})
	}
}
