package citlali

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

const (
	itzpapaInterval           = 60 // looking at footage, seems like both attack and NS consumption intervals are the same
	obsidianTzitzimitlHitmark = 22
	particleICDKey            = "citlali-particle-icd"
	itzpapaKey                = "itzpapa-key"
	frostFallAbil             = "Frostfall Storm DMG"
)

var (
	skillFrames []int
)

func init() {
	skillFrames = frames.InitAbilSlice(40) // E -> Swap
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Obsidian Tzitzimitl DMG",
		AttackTag:      attacks.AttackTagElementalArt,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagNone,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Cryo,
		Durability:     25,
		Mult:           skill[c.TalentLvlSkill()],
	}
	c.QueueCharTask(func() {
		c.SetCD(action.ActionSkill, 16*60)
		c.addShield()
	}, 1)

	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 6), obsidianTzitzimitlHitmark, obsidianTzitzimitlHitmark, c.particleCB)

	if c.nightsoulState.HasBlessing() {
		c.nightsoulState.GeneratePoints(24)
	} else {
		c.nightsoulState.EnterBlessing(24)
	}

	c.itzpapaSrc = c.Core.F
	c.ActivateItzpapa(c.Core.F)

	if c.Base.Cons >= 1 {
		c.numStellarBlades = 10
		c.c2()
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionBurst],
		State:           action.SkillState,
	}, nil
}

func (c *char) ActivateItzpapa(src int) {
	// try to activate Itzpapa each time Citlali gains NS points to avoid event subscribtion
	if c.nightsoulState.Points() >= 50 {
		// if it's activation or REactivation
		if !c.StatusIsActive(itzpapaKey) || src != c.itzpapaSrc {
			// this status is active only when Itzpapa is in "attack mode"
			c.AddStatus(itzpapaKey, -1, false)
			c.QueueCharTask(c.ItzpapaHit(src), itzpapaInterval)
		}
	}
}

func (c *char) ItzpapaHit(src int) func() {
	return func() {
		if src != c.itzpapaSrc {
			return
		}
		if !c.StatusIsActive(itzpapaKey) {
			return
		}
		if c.nightsoulState.Points() == 0 {
			c.nightsoulState.ExitBlessing()
			c.DeleteStatus(itzpapaKey)
			c.numStellarBlades = 0 // C1
			return
		}
		ai := combat.AttackInfo{
			ActorIndex:     c.Index,
			Abil:           frostFallAbil,
			AttackTag:      attacks.AttackTagElementalArt,
			AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
			ICDTag:         attacks.ICDTagElementalArtHold,
			ICDGroup:       attacks.ICDGroupCitlaliSkill,
			StrikeType:     attacks.StrikeTypeDefault,
			Element:        attributes.Cryo,
			Durability:     25,
			Mult:           frostfall[c.TalentLvlSkill()],
			FlatDmg:        c.a4DmgSkill(),
		}
		c.nightsoulState.ConsumePoints(8)
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 6), 0, 0)
		c.QueueCharTask(c.ItzpapaHit(src), itzpapaInterval)
		c.c4Skull()
	}
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 0.2*60, false)
	c.Core.QueueParticle(c.Base.Key.String(), 5, attributes.Cryo, c.ParticleDelay)
}
