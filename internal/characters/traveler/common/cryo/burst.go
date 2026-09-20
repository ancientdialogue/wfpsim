package cryo

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var burstFrames [][]int

var burstTickHitmarks = []int{0, 0 + 33, 0 + 33 + 2 + 3, 0 + 33 + 2, 0 + 33 + 2 + 3 + 4}

const burstSpawnFrame = 36

func init() {
	burstFrames = make([][]int, 2)

	// Male, assuming same as female for now
	burstFrames[0] = frames.InitAbilSlice(75)
	burstFrames[0][action.ActionSkill] = 74 // Q -> E
	burstFrames[0][action.ActionJump] = 74  // Q -> J
	burstFrames[0][action.ActionSwap] = 73  // Q -> Swap

	// Female
	burstFrames[1] = frames.InitAbilSlice(75)
	burstFrames[1][action.ActionSkill] = 74 // Q -> E
	burstFrames[1][action.ActionJump] = 74  // Q -> J
	burstFrames[1][action.ActionSwap] = 73  // Q -> Swap
}

func (c *Traveler) Burst(p map[string]int) (action.Info, error) {
	c.SetCD(action.ActionBurst, 15*60)
	c.ConsumeEnergy(7)

	attack := func() {
		ai := info.AttackInfo{
			ActorIndex: c.Index(),
			Abil:       "Ice Javelin",
			AttackTag:  attacks.AttackTagElementalBurst,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Cryo,
			Durability: 25,
			Mult:       burst[c.TalentLvlBurst()] + flowGlowBonus[c.TalentLvlBurst()]*float64(c.flostglowStacks),
		}

		switch c.getRadiance() {
		case radianceStellarConduct:
			ai.Abil += stellarConductText
			ai.AttackTag = attacks.AttackTagDirectStellarConduct
			ai.Durability = 0
			ai.Mult = burstSSC[c.TalentLvlBurst()] + flowGlowBonusSSC[c.TalentLvlBurst()]*float64(c.flostglowStacks)
			ai.IgnoreDefPercent = 1
		case radianceStellarSwirl:
			ai.Abil += stellarSwirlText
			ai.AttackTag = attacks.AttackTagDirectStellarSwirl
			ai.Durability = 0
			ai.Mult = burstSSw[c.TalentLvlBurst()] + flowGlowBonusSSw[c.TalentLvlBurst()]*float64(c.flostglowStacks)
			ai.IgnoreDefPercent = 1
		default:
		}

		hits := 3
		if c.flostglowStacks == skillStacksMax {
			hits += 2
		}

		for _, delay := range burstTickHitmarks[:hits] {
			c.Core.QueueAttack(
				ai,
				combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 4.5),
				delay-burstSpawnFrame,
				delay-burstSpawnFrame,
			)
		}
		c.c6OnBurst(c.flostglowStacks)
		c.flostglowStacks = 0
	}

	c.QueueCharTask(attack, burstSpawnFrame)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames[c.gender]),
		AnimationLength: burstFrames[c.gender][action.InvalidAction],
		CanQueueAfter:   burstFrames[c.gender][action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}
