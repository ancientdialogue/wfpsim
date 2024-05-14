package sethos

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var skillFrames []int

func init() {
	skillFrames = frames.InitAbilSlice(28)
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Ancient Rite: Thunderous Roar of Sand",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}

	snap := c.Snapshot(&ai)
	ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 4.5)

	c.Core.QueueAttackWithSnap(ai, snap, ap, 14, c.makeParticleCB(), c.makeEnergyRegenCB())

	c.SetCDWithDelay(action.ActionSkill, 8*60, 10)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionSwap], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) makeParticleCB() combat.AttackCBFunc {
	done := false
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		if done {
			return
		}
		done = true
		c.Core.QueueParticle(c.Base.Key.String(), 2, attributes.Electro, c.ParticleDelay)
	}
}
func (c *char) makeEnergyRegenCB() combat.AttackCBFunc {
	done := false
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}

		if done {
			return
		}

		// assuming that the skill can only do electro reactions, since the skill applies electro
		// therefore we aren't checking for the list of electro reactions
		// Aggravate/Spread don't set the Reacted flag
		if !a.AttackEvent.Reacted && !a.AttackEvent.Info.Catalyzed {
			return
		}
		done = true
		c.AddEnergy("sethos-skill", skillEnergyRegen[c.TalentLvlSkill()])
		c.c2AddStack()
	}
}
