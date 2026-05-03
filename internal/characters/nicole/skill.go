package nicole

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var skillFrames []int

const (
	skillHitmark   = 30
	particleICDKey = "nicole-particle-icd"
	skillBuffKey   = "grace-of-kenosis"
)

func init() {
	skillFrames = frames.InitAbilSlice(75)
	skillFrames[action.ActionAttack] = 40
	skillFrames[action.ActionCharge] = 69
	skillFrames[action.ActionSkill] = 68
	skillFrames[action.ActionBurst] = 34
	skillFrames[action.ActionDash] = 37
	skillFrames[action.ActionJump] = 35
	skillFrames[action.ActionSwap] = 74
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Revelation: Uncreated Light",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagKleeFireDamage,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}
	ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 4)

	c.Core.QueueAttack(ai, ap, skillHitmark, skillHitmark, c.particleCB)

	for _, char := range c.Core.Player.Chars() {
		char.DeleteStatMod(c.graceOfKenosis.ModKey)
	}

	c.QueueCharTask(func() {
		for _, char := range c.Core.Player.Chars() {
			char.AddStatMod(c.graceOfKenosis)
		}
		c.addShield()
	}, skillHitmark-1)

	c.a1OnSkill()

	c.SetCDWithDelay(action.ActionSkill, 16*60, 4)
	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionBurst], // earliest cancel is before skillHitmark
		State:           action.SkillState,
	}, nil
}

func (c *char) skillInit() {
	m := make([]float64, attributes.EndStatType)
	c.graceOfKenosis = character.StatMod{
		Base:         modifier.NewBase(skillBuffKey, 20*60),
		AffectedStat: attributes.ATK,
		Extra:        true,
		Amount: func() ([]float64, bool) {
			atk := c.SelectStat(true, attributes.BaseATK, attributes.ATKP, attributes.ATK).TotalATK()
			m[attributes.ATK] = min(atk*skillBuff[c.TalentLvlSkill()], skillBuffCap[c.TalentLvlSkill()]) + c.c2SkillBuff()
			return m, true
		},
	}
}

func (c *char) particleCB(a info.AttackCB) {
	if a.Target.Type() != info.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 0.5*60, true)
	c.Core.QueueParticle(c.Base.Key.String(), 5, attributes.Pyro, c.ParticleDelay)
}
