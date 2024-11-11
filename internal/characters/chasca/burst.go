package chasca

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var (
	burstFrames           []int
	burstSkillStateFrames []int
)

const (
	burstCD = 15 * 60
)

func init() {
	burstFrames = frames.InitAbilSlice(128) // Q - N1/E/Jump/Walk
	burstFrames[action.ActionDash] = 127
	burstFrames[action.ActionSwap] = 127

	burstSkillStateFrames = frames.InitAbilSlice(128) // Q - Jump/Walk
	burstSkillStateFrames[action.ActionAttack] = 127
	burstSkillStateFrames[action.ActionSkill] = 127
	burstSkillStateFrames[action.ActionDash] = 127
	burstSkillStateFrames[action.ActionSwap] = 127
}
func (c *char) Burst(p map[string]int) (action.Info, error) {
	c.DeleteStatus(c4icd)
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Galesplitting Soulseeker Shell",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagChascaBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       burstDMG[c.TalentLvlBurst()],
	}

	ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5)
	c.Core.QueueAttack(ai, ap, 0, 1)
	c.SetCD(action.ActionBurst, burstCD)
	c.ConsumeEnergy(1)
	c.BurstConversion()
	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *char) BurstConversion() (action.Info, error) {
	numConversions := len(c.conversionElements) * 2
	for i := 0; i < 6; i++ {
		element := attributes.Anemo
		if i < numConversions && len(c.conversionElements) > 0 {
			randomIndex := c.Core.Rand.Intn(len(c.conversionElements))
			element = c.conversionElements[randomIndex]
		}

		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Soulseeker Shell",
			AttackTag:  attacks.AttackTagElementalBurst,
			ICDTag:     attacks.ICDTagChascaBurst,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    element,
			Durability: 25,
			Mult:       burstSoulseeker[c.TalentLvlBurst()],
		}
		if element != attributes.Anemo {
			if c.Base.Cons >= 4 && !c.StatusIsActive(c4icd) {
				c4ai := combat.AttackInfo{
					ActorIndex: c.Index,
					Abil:       "C4 Radiant Soulseeker Shells",
					AttackTag:  attacks.AttackTagElementalBurst,
					ICDTag:     attacks.ICDTagChascaBurst,
					ICDGroup:   attacks.ICDGroupDefault,
					StrikeType: attacks.StrikeTypeDefault,
					Element:    element,
					Durability: 25,
					Mult:       400 / 100,
				}
				ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5)
				c.Core.QueueAttack(c4ai, ap, 5, 5)
				c.AddStatus(c4icd, -1, false)
			}
			c.c4energy()
			ai.Abil = "Radiant Soulseeker Shell"
			ai.Mult = burstRadiantSoulseeker[c.TalentLvlBurst()]
		}

		ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5)
		c.Core.QueueAttack(ai, ap, i*5, i*5+5)
	}
	return action.Info{
		Frames:          frames.NewAbilFunc(burstSkillStateFrames),
		AnimationLength: burstSkillStateFrames[action.InvalidAction],
		CanQueueAfter:   burstSkillStateFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}
