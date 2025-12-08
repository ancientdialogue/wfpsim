package zibai

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var burstFrames []int

var burstHitmarks = []int{101, 110}

func init() {
	burstFrames = frames.InitAbilSlice(139) // Q -> N1/E
	burstFrames[action.ActionDash] = 123    // Q -> D
	burstFrames[action.ActionJump] = 123    // Q -> J
	burstFrames[action.ActionSwap] = 138    // Q -> Swap
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	// deal damage when created
	for i, mult := range burst {
		ai := info.AttackInfo{
			ActorIndex: c.Index(),
			Abil:       fmt.Sprintf("Burst %v", i),
			AttackTag:  attacks.AttackTagElementalBurst,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeBlunt,
			PoiseDMG:   500,
			Element:    attributes.Geo,
			Durability: 25,
			Mult:       mult[c.TalentLvlBurst()],
		}
		if i == 1 {
			ai.AttackTag = attacks.AttackTagDirectLunarCrystallize
			ai.Durability = 0
			ai.IgnoreDefPercent = 1
		}

		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(c.Core.Combat.Player(), info.Point{Y: 5}, 7),
			burstHitmarks[i],
			burstHitmarks[i],
		)
	}

	c.SetCD(action.ActionBurst, 15*60)
	c.ConsumeEnergy(7)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}, nil
}
