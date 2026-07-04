package sandrone

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	burstFrames   []int
	burstHitmarks = []int{123, 123 + 12, 123 + 12 + 12}
	rayHitmark    = 123 + 12 + 12 + 30
)

func init() {
	burstFrames = frames.InitAbilSlice(128) // Q -> CA
	burstFrames[action.ActionAttack] = 101  // Q -> N1
	burstFrames[action.ActionSkill] = 100   // Q -> E
	burstFrames[action.ActionDash] = 103    // Q -> D
	burstFrames[action.ActionJump] = 103    // Q -> J
	burstFrames[action.ActionWalk] = 105    // Q -> Swap
	burstFrames[action.ActionSwap] = 102    // Q -> Swap
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Bombardment",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}

	ap := combat.NewBoxHitOnTarget(
		c.Core.Combat.Player(),
		info.Point{Y: -5},
		14,
		12,
	)

	for _, delay := range burstHitmarks {
		c.QueueCharTask(func() { c.Core.QueueAttack(ai, ap, 0, 0) }, delay)
	}

	c.QueueCharTask(func() {
		aiRay := info.AttackInfo{
			ActorIndex: c.Index(),
			Abil:       "Convective Inhibition Ray",
			AttackTag:  attacks.AttackTagElementalBurst,
			ICDTag:     attacks.ICDTagElementalBurst,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Cryo,
			Durability: 25,
			Mult:       burstRay[c.TalentLvlSkill()],
		}

		switch c.getRadiance() {
		case radianceStellarConduct:
			aiRay.Abil += stellarConductText
			aiRay.AttackTag = attacks.AttackTagDirectStellarConduct
			aiRay.Durability = 0
			aiRay.Mult = burstRaySSC[c.TalentLvlSkill()] * c.a1OnBurstRayStellar()
			aiRay.IgnoreDefPercent = 1
		case radianceStellarSwirl:
			aiRay.Abil += stellarSwirlText
			aiRay.AttackTag = attacks.AttackTagDirectStellarSwirl
			aiRay.Durability = 0
			aiRay.Mult = burstRaySSw[c.TalentLvlSkill()] * c.a1OnBurstRayStellar()
			aiRay.IgnoreDefPercent = 1
		}

		c.Core.QueueAttack(aiRay, ap, 0, 0)
	}, rayHitmark)

	c.ConsumeEnergy(7)
	c.SetCDWithDelay(action.ActionBurst, 15*60, 0)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSkill], // earliest cancel
		State:           action.BurstState,
	}, nil
}
