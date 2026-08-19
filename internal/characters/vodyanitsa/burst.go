package vodyanitsa

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var burstFrames []int

const (
	burstHitmark = 91
)

func init() {
	burstFrames = frames.InitAbilSlice(110) // Q -> Dash
	burstFrames[action.ActionAttack] = 108
	burstFrames[action.ActionSkill] = 108
	burstFrames[action.ActionJump] = 109
	burstFrames[action.ActionSwap] = 107
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	mult := burst[c.TalentLvlBurst()]
	if c.StatusIsActive(skillKey) {
		mult += burstBonus[c.TalentLvlBurst()]
	}

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Sink With Thee",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeSlash,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       mult,
		UseHP:      true,
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 3),
		burstHitmark,
		burstHitmark,
	)

	c.ConsumeEnergy(4)
	c.SetCD(action.ActionBurst, 15*60)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}
