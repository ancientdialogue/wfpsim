package sandrone

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var skillFrames []int

const skillParticleICDKey = "sandrone-particle-icd"

func init() {
	skillFrames = frames.InitAbilSlice(600) // need to full exit E before actual walk occurs
	skillFrames[action.ActionAttack] = 10
	skillFrames[action.ActionCharge] = 46
	skillFrames[action.ActionSkill] = 46 // TODO
	skillFrames[action.ActionBurst] = 6
	skillFrames[action.ActionDash] = 26
	skillFrames[action.ActionSwap] = 28
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}

	c.QueueCharTask(func() {
		ai := info.AttackInfo{
			ActorIndex: c.Index(),
			Abil:       "Prism Shot",
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Cryo,
			Durability: 25,
			Mult:       skill[c.TalentLvlSkill()],
		}

		ap := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 0.5)
		c.Core.QueueAttack(ai, ap, 5, 5+travel, c.particleCB)

		switch c.getRadiance() {
		case radianceStellarConduct:
			ai.Abil += stellarConductText
			ai.AttackTag = attacks.AttackTagDirectStellarConduct
			ai.Durability = 0
			ai.Mult = skillSSC[c.TalentLvlSkill()]
			ai.IgnoreDefPercent = 1
		case radianceStellarSwirl:
			ai.Abil += stellarSwirlText
			ai.AttackTag = attacks.AttackTagDirectStellarSwirl
			ai.Durability = 0
			ai.Mult = skillSSw[c.TalentLvlSkill()]
			ai.IgnoreDefPercent = 1
		}

		ai.Mult *= c.a1OnSkill()
		c.Core.QueueAttack(ai, ap, 5, 5+travel, c.particleCB)

		c.reduceDecode(100)
	}, 15)

	c.SetCD(action.ActionSkill, 4*60)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionBurst], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) particleCB(a info.AttackCB) {
	if a.Target.Type() != info.TargettableEnemy {
		return
	}
	if c.StatusIsActive(skillParticleICDKey) {
		return
	}
	c.AddStatus(skillParticleICDKey, 2.5*60, true)
	c.Core.QueueParticle(c.Base.Key.String(), 1, attributes.Cryo, c.ParticleDelay)
}
