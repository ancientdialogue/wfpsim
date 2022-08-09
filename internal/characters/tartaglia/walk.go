package tartaglia

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
)

var walkFrames []int

func init() {
	walkFrames = frames.InitAbilSlice(1)
	walkFrames[action.ActionSkill] = 4
}

func (c *char) Walk(p map[string]int) action.ActionInfo {
	f, ok := p["f"]
	if !ok {
		f = 1
	}
	return action.ActionInfo{
		Frames: func(next action.Action) int {
			if f < walkFrames[next] {
				return walkFrames[next]
			}
			return f
		},
		AnimationLength: walkFrames[action.ActionSkill],
		CanQueueAfter:   f,
		State:           action.WalkState,
	}
}
