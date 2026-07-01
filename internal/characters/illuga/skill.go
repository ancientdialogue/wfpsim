package illuga

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var skillFrames []int

const (
	skillHitmark   = 21
	particleICDKey = "illuga-particle-icd"
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
		Element:    attributes.Geo,
		Durability: 25,
		UseDef:     true,
		Mult:       skillTapDef[c.TalentLvlSkill()],
		FlatDmg:    skillTapEM[c.TalentLvlSkill()] * c.Stat(attributes.EM),
	}

	c.Core.QueueAttack(
		ai,
		combat.NewBoxHitOnTarget(c.Core.Combat.Player(), nil, 4, 4.1),
		skillHitmark,
		skillHitmark,
		c.baseParticleCB,
	)

	c.SetCDWithDelay(action.ActionSkill, 360, skillHitmark-2)
	c.a1OnSkillBurst()

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
		Element:    attributes.Geo,
		Durability: 25,
		UseDef:     true,
		Mult:       skillHoldDef[c.TalentLvlSkill()],
		FlatDmg:    skillHoldEM[c.TalentLvlSkill()] * c.Stat(attributes.EM),
	}

	c.Core.QueueAttack(
		ai,
		combat.NewBoxHitOnTarget(c.Core.Combat.Player(), nil, 4, 4.1),
		hitmark,
		hitmark,
		c.baseParticleCB,
	)

	c.SetCDWithDelay(action.ActionSkill, 15*60, hitmark-2)
	c.a1OnSkillBurst()

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
	c.AddStatus(particleICDKey, 0.3*60, true)
	particles := 4.0
	if c.Core.Rand.Float64() < 0.5 {
		particles = 5
	}
	c.Core.QueueParticle(c.Base.Key.String(), particles, attributes.Geo, c.ParticleDelay)
}
