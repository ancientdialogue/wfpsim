package sethos

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
)

var burstFrames []int

const (
	burstStart   = 57
	burstBuffKey = "sethosburst"
)

func init() {
	burstFrames = frames.InitAbilSlice(82) // Q -> N1/E
	burstFrames[action.ActionDash] = 59    // Q -> D
	burstFrames[action.ActionJump] = 60    // Q -> J
	burstFrames[action.ActionSwap] = 66    // Q -> Swap
}

// Sets Xiao's burst damage state
func (c *char) Burst(p map[string]int) (action.Info, error) {
	c.QueueCharTask(func() {
		c.AddStatus(burstBuffKey, 8*60, true)
	}, burstStart)

	c.SetCDWithDelay(action.ActionBurst, 15*60, 29)
	c.ConsumeEnergy(36)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}, nil
}
