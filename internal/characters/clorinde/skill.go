package clorinde

import (
	"math"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var (
	skillFrames            []int
	skillDashNoBOLFrames   []int
	skillDashLowBOLFrames  []int
	skillDashFullBOLFrames []int
	c6Stacks               int
)

const (
	skillStateKey  = "clorinde-night-watch"
	tolerance      = 0.0000001
	skillCD        = 16 * 60
	particleICDKey = "clorinde-particle-icd"

	// TODO: all hit marks
	skillDashNoBOLHitmark   = 24
	skillDashLowBOLHitmark  = 24
	skillDashFullBOLHitmark = 24
)

func init() {
	// TODO: all these frames are gusses
	skillFrames = frames.InitAbilSlice(32)
	skillFrames[action.ActionSkill] = 27

	skillDashNoBOLFrames = frames.InitAbilSlice(28)
	skillDashNoBOLFrames[action.ActionAttack] = 23
	skillDashLowBOLFrames = frames.InitAbilSlice(28)
	skillDashLowBOLFrames[action.ActionAttack] = 23
	skillDashFullBOLFrames = frames.InitAbilSlice(28)
	skillDashFullBOLFrames[action.ActionAttack] = 23
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	// first press activates skill state
	// sequential presses pew pew stuff
	if c.StatusIsActive(skillStateKey) {
		return c.skillDash(p)
	}
	c6Stacks = 6
	c.QueueCharTask(c.c6skill, 0)
	c.AddStatus(skillStateKey, int(60*skillStateDuration[0]), true)

	c.SetCD(action.ActionSkill, skillCD)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionSkill],
		State:           action.SkillState,
	}, nil
}

func (c *char) skillDash(p map[string]int) (action.Info, error) {
	// depending on BOL lvl it does either 1 hit or 3 hit
	ratio := c.currentHPDebtRatio()
	switch {
	case ratio >= 1:
		if c.Base.Cons >= 6 && c6Stacks > 0 {
			c.c6()
			c6Stacks -= 1
		}
		return c.skillDashFullBOL(p)
	case math.Abs(ratio) < tolerance:
		return c.skillDashNoBOL(p)
	default:
		return c.skillDashRegular(p)
	}
}

func (c *char) gainBOLOnAttack() {
	c.ModifyHPDebtByRatio(skillBOLGain[c.TalentLvlSkill()])
}

func (c *char) skillDashNoBOL(_ map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Skill Dash (No BOL)",
		AttackTag:      attacks.AttackTagNormal,
		ICDTag:         attacks.ICDTagNormalAttack,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeSlash,
		Element:        attributes.Electro,
		Durability:     25,
		Mult:           skillLungeNoBOL[c.TalentLvlSkill()],
		IgnoreInfusion: true,
	}
	// TODO: what's the size of this??
	ap := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 0.6)
	// TODO: assume no snapshotting on this
	c.Core.QueueAttack(ai, ap, skillDashNoBOLHitmark, skillDashNoBOLHitmark, c.particleCB)
	// TODO: no idea if this counts as a normal attack state or not. pretend it does for now
	return action.Info{
		Frames:          frames.NewAbilFunc(skillDashNoBOLFrames),
		AnimationLength: skillDashNoBOLFrames[action.InvalidAction],
		CanQueueAfter:   skillDashNoBOLFrames[action.ActionAttack], //TODO: fastest cancel?
		State:           action.SkillState,
	}, nil
}

func (c *char) skillDashFullBOL(_ map[string]int) (action.Info, error) {
	for i := 0; i < 3; i++ {
		ai := combat.AttackInfo{
			ActorIndex:     c.Index,
			Abil:           "Skill Dash (Full BOL)",
			AttackTag:      attacks.AttackTagNormal,
			ICDTag:         attacks.ICDTagNormalAttack,
			ICDGroup:       attacks.ICDGroupDefault,
			StrikeType:     attacks.StrikeTypeSlash,
			Element:        attributes.Electro,
			Durability:     25,
			Mult:           skillLungeFullBOL[c.TalentLvlSkill()],
			IgnoreInfusion: true,
		}
		// TODO: what's the size of this??
		ap := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 0.8)
		// TODO: assume no snapshotting on this
		c.Core.QueueAttack(ai, ap, skillDashFullBOLHitmark, skillDashFullBOLHitmark, c.particleCB)
	}

	// TODO: timing on this heal?
	c.skillHeal(skillLungeFullBOLHeal[0], "skill (>= 100%)")

	// TODO: no idea if this counts as a normal attack state or not. pretend it does for now
	return action.Info{
		Frames:          frames.NewAbilFunc(skillDashFullBOLFrames),
		AnimationLength: skillDashFullBOLFrames[action.InvalidAction],
		CanQueueAfter:   skillDashFullBOLFrames[action.ActionAttack], //TODO: fastest cancel?
		State:           action.SkillState,
	}, nil
}

func (c *char) skillDashRegular(_ map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Skill Dash (< 100% BOL)",
		AttackTag:      attacks.AttackTagNormal,
		ICDTag:         attacks.ICDTagNormalAttack,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeSlash,
		Element:        attributes.Electro,
		Durability:     25,
		Mult:           skillLungeLowBOL[c.TalentLvlSkill()],
		IgnoreInfusion: true,
	}
	// TODO: what's the size of this??
	ap := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 0.8)
	// TODO: assume no snapshotting on this
	c.Core.QueueAttack(ai, ap, skillDashLowBOLHitmark, skillDashLowBOLHitmark, c.particleCB)

	// TODO: timing on this heal?
	c.skillHeal(skillLungeLowBOLHeal[0], "skill (< 100%)")

	// TODO: no idea if this counts as a normal attack state or not. pretend it does for now
	return action.Info{
		Frames:          frames.NewAbilFunc(skillDashLowBOLFrames),
		AnimationLength: skillDashLowBOLFrames[action.InvalidAction],
		CanQueueAfter:   skillDashLowBOLFrames[action.ActionAttack], //TODO: fastest cancel?
		State:           action.SkillState,
	}, nil
}

func (c *char) skillHeal(bolMult float64, msg string) {
	amt := c.CurrentHPDebt() * bolMult
	c.heal(&info.HealInfo{
		Caller:  c.Index,
		Target:  c.Index,
		Message: msg,
		Src:     amt,
		Bonus:   c.Stat(attributes.Heal), // TODO: confirms that it scales with healing %
	})
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 2*60, true)

	c.Core.QueueParticle(c.Base.Key.String(), 1, attributes.Electro, c.ParticleDelay)
}
