package vodyanitsa

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
	skillHitmark   = 32
	particleICDKey = "vodyanitsa-particle-icd"
	skillKey       = "vodyanitsa-skill"
)

func init() {
	skillFrames = frames.InitAbilSlice(56)
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Rechitativ: Sonorous Dawn",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagElementalArt,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
		UseHP:      true,
	}

	ap := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 3)

	c.Core.QueueAttack(ai, ap, 0, 0, c.particleCB)

	c.c2OnSkillAttack()

	c.skillSrc = c.Core.F
	c.Core.Tasks.Add(func() { c.skillHealTask(c.skillSrc) }, skillHitmark+14+3*60)
	c.Core.Tasks.Add(func() { c.skillAttackTask(c.skillSrc) }, skillHitmark+14+1.5*60)

	c.AddStatus(skillKey, 16*60+c.c2SkillDurBonus(), false)

	c.SetCD(action.ActionSkill, 16*60)

	c.a4OnSkill()
	c.c6OnSkill()

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillHitmark,
		State:           action.SkillState,
	}, nil
}

func (c *char) skillAttackTask(src int) {
	if !c.StatusIsActive(skillKey) {
		return
	}

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Horn of Spring's Call",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagElementalArt,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
		UseHP:      true,
	}

	ap := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 3)

	c.Core.QueueAttack(ai, ap, 0, 0, c.particleCB)

	c.c2OnSkillAttack()
	c.Core.Tasks.Add(func() { c.skillAttackTask(src) }, 3*60)
}

func (c *char) skillHealTask(src int) {
	if !c.StatusIsActive(skillKey) {
		return
	}

	maxhp := c.MaxHP()

	c.c4BeforeHeal()

	c.Core.Player.Heal(info.HealInfo{
		Caller:  c.Index(),
		Target:  c.Core.Player.Active(),
		Message: "Song of Ages Past",
		Src:     skillHealRatio[c.TalentLvlSkill()]*maxhp + skillHealFlat[c.TalentLvlSkill()],
		Bonus:   c.Stat(attributes.Heal),
	})
	c.c1OnHeal()

	c.Core.Tasks.Add(func() { c.skillHealTask(src) }, 1.5*60)
}

func (c *char) particleCB(ac info.AttackCB) {
	if ac.Target.Type() != info.TargettableEnemy {
		return
	}

	if c.StatusIsActive(particleICDKey) {
		return
	}

	c.AddStatus(particleICDKey, 0.2*60, true)
	c.Core.QueueParticle(c.Base.Key.String(), 1, attributes.Hydro, c.ParticleDelay)
}
