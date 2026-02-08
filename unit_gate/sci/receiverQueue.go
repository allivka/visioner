package sci

import (
	"context"
	"time"
)

func MaintainChannel(ctx context.Context, flow chan float64, target int, cycle time.Duration) {
	go func(ctx context.Context) {
		timer := time.NewTimer(cycle)
		for {
			select {

			case <-ctx.Done():
				return
			case <-timer.C:
				timer = time.NewTimer(cycle)
				if len(flow) > target {
					<-flow
				}
			}
		}
	}(ctx)
}
