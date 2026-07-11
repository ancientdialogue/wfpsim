package alyosha

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var skillFrames []int

const (
	skillHitmark   = 21
	particleICDKey = "alyosha-particle-icd"
	skillMarkKey   = "alyohsa-hunters-mark"
	skillBuffKey   = "alyohsa-hunters-precision"
)

func init() {
	skillFrames = frames.InitAbilSlice(46)
	skillFrames[action.ActionDash] = 28
	skillFrames[action.ActionJump] = 28
	skillFrames[action.ActionBurst] = 28
	skillFrames[action.ActionSwap] = 45
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	hold, ok := p["hold"]
	if ok && hold > 0 {
		return c.skillHold(hold)
	}

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Skill",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       skillTap[c.TalentLvlSkill()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewBoxHitOnTarget(c.Core.Combat.Player(), nil, 4, 4.1),
		skillHitmark,
		skillHitmark,
		c.baseParticleCB,
		c.triggerSkillMarkCB(true),
	)

	c.SetCDWithDelay(action.ActionSkill, 360, skillHitmark-2)

	return action.Info{
		Frames:          func(next action.Action) int { return skillFrames[next] },
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) skillHold(hold int) (action.Info, error) {
	hitmark := skillHitmark + hold + 15

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Skill (Hold)",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       skillHold[c.TalentLvlSkill()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewBoxHitOnTarget(c.Core.Combat.Player(), nil, 4, 4.1),
		hitmark,
		hitmark,
		c.baseParticleCB,
		c.triggerSkillMarkCB(true),
	)

	c.SetCDWithDelay(action.ActionSkill, 15*60, hitmark-2)

	return action.Info{
		Frames:          func(next action.Action) int { return skillFrames[next] + hold },
		AnimationLength: skillFrames[action.InvalidAction] + hold,
		CanQueueAfter:   skillFrames[action.ActionDash] + hold, // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) baseParticleCB(a info.AttackCB) {
	if a.Target.Type() != info.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 1*60, true)

	c.Core.QueueParticle(c.Base.Key.String(), 5, attributes.Electro, c.ParticleDelay)
}

func (c *char) skillInit() {
	m := make([]float64, attributes.EndStatType)
	for _, char := range c.Core.Player.Chars() {
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase(skillBuffKey+"-atkp", -1),
			AffectedStat: attributes.ATKP,
			Amount: func() []float64 {
				if !c.StatusIsActive(skillBuffKey) {
					return nil
				}
				m[attributes.ATKP] = skillBuff[c.TalentLvlSkill()] * float64(c.skillStacks)
				return m
			},
		})
	}
}

func (c *char) triggerSkillMarkCB(canApply bool) info.AttackCBFunc {
	return func(ac info.AttackCB) {
		if ac.Target.Type() != info.TargettableEnemy {
			return
		}

		e, ok := ac.Target.(*enemy.Enemy)
		if !ok {
			return
		}

		if e.StatusIsActive(skillMarkKey) {
			e.DeleteStatus(skillMarkKey)
			if !c.StatusIsActive(skillBuffKey) {
				c.skillStacks = 0
			}

			c.skillStacks = min(c.skillStacks+1, c.c6MaxSkillStacks())
			c.AddStatus(skillBuffKey, 15*60, true)
			return
		}

		if canApply {
			e.AddStatus(skillMarkKey, 15*60, true)
		}
	}
}
