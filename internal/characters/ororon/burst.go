package ororon

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var burstFrames []int

const burstHitmark = 51 // Initial Hit

func init() {
	burstFrames = frames.InitAbilSlice(80) // Q -> CA
	burstFrames[action.ActionAttack] = 78  // Q -> N1
	burstFrames[action.ActionSkill] = 57   // Q -> E
	burstFrames[action.ActionDash] = 58    // Q -> D
	burstFrames[action.ActionJump] = 58    // Q -> J
	burstFrames[action.ActionSwap] = 56    // Q -> Swap
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	// first zap has no icd and hits everyone
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Ritual (Initial)",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Electro,
		Durability: 0,
		Mult:       burstRitual[c.TalentLvlBurst()],
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5),
		burstHitmark,
		burstHitmark,
		c.makeC2cb(),
	)

	ai = combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Soundwave Collision (Tick)",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       burstSonic[c.TalentLvlBurst()],
	}

	firstTick := 111 // first tick at 111
	burstInterval := int(45 * c.c4BurstInterval())
	for i := 0; i <= 9*60; i += burstInterval {
		progress := i + firstTick
		c.QueueCharTask(func() {
			c.Core.QueueAttack(
				ai,
				combat.NewCircleHitFanAngle(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 10, 30),
				0,
				0,
				c.makeC2cb(),
			)
		}, progress)
	}
	c.c2OnBurst()
	c.c6OnBurst()
	c.SetCDWithDelay(action.ActionBurst, 15*60, 0)
	c.ConsumeEnergy(0)
	c.c4EnergyRestore()
	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}
