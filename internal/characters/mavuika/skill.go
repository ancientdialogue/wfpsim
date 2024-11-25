package mavuika

import (
	"errors"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var (
	skillFrames       []int
	skillRecastFrames []int
)

const (
	skillHitmark   = 20
	particleICDKey = "mavuika-particle-icd"
)

func init() {
	skillFrames = frames.InitAbilSlice(39) // E -> N1
	skillFrames[action.ActionSkill] = 30
	skillFrames[action.ActionBurst] = 29
	skillFrames[action.ActionDash] = 26
	skillFrames[action.ActionJump] = 28
	skillFrames[action.ActionSwap] = 25

	skillRecastFrames = frames.InitAbilSlice(74) // E -> N1
	skillRecastFrames[action.ActionSkill] = 45
	skillRecastFrames[action.ActionBurst] = 45
	skillRecastFrames[action.ActionDash] = 45
	skillRecastFrames[action.ActionJump] = 49
	skillRecastFrames[action.ActionSwap] = 44
}

func (c *char) nightsoulPointReduceFunc(src int) func() {
	return func() {
		if c.nightsoulSrc != src {
			return
		}
		val := 1.0
		if c.armamentState == bike {
			val += 0.8
			if c.Core.Player.CurrentState() == action.ChargeAttackState {
				val += 0.4
			}
		}
		c.reduceNightsoulPoints(val)
		c.Core.Tasks.Add(c.nightsoulPointReduceFunc(src), 12)
	}
}

func (c *char) reduceNightsoulPoints(val float64) {
	val *= c.nightsoulConsumptionMul()
	if val == 0 {
		return
	}
	c.nightsoulState.ConsumePoints(val)

	// don't exit nightsoul while in NA/Plunge/Charge of Flamestride
	if c.armamentState == bike {
		switch c.Core.Player.CurrentState() {
		case action.NormalAttackState, action.PlungeAttackState, action.ChargeAttackState:
			return
		}
	}

	if c.nightsoulState.Points() < 0.001 {
		c.exitNightsoul()
	}
}

func (c *char) exitNightsoul() {
	if !c.nightsoulState.HasBlessing() {
		return
	}
	c.nightsoulState.ExitBlessing()
	c.nightsoulState.ClearPoints()
	c.nightsoulSrc = -1
	c.NormalHitNum = normalHitNum
	c.NormalCounter = 0
	c.c2OnNightsoulExit()
}
func (c *char) enterNightsoulOrRegenerate(points float64) {
	if !c.nightsoulState.HasBlessing() {
		c.nightsoulState.EnterBlessing(points)
		c.nightsoulSrc = c.Core.F
		c.Core.Tasks.Add(c.nightsoulPointReduceFunc(c.nightsoulSrc), 12)
		c.c2OnNightsoulEnter()
		return
	}
	c.nightsoulState.GeneratePoints(points)
}
func (c *char) Skill(p map[string]int) (action.Info, error) {
	h := p["hold"]

	if c.nightsoulState.HasBlessing() {
		if h > 0 {
			return action.Info{}, errors.New("cannot hold E while in Nightsoul Blessing")
		}
		ai := combat.AttackInfo{
			ActorIndex:     c.Index,
			Abil:           "The Named Moment",
			AttackTag:      attacks.AttackTagElementalArt,
			ICDTag:         attacks.ICDTagNone,
			AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
			ICDGroup:       attacks.ICDGroupDefault,
			StrikeType:     attacks.StrikeTypePierce,
			Element:        attributes.Pyro,
			Durability:     25,
			Mult:           skill[c.TalentLvlSkill()],
		}
		ap := combat.NewCircleHitOnTarget(
			c.Core.Combat.Player(),
			geometry.Point{Y: 1.0},
			6,
		)
		c.Core.QueueAttack(ai, ap, skillHitmark, skillHitmark, c.particleCB)
		c.enterNightsoulOrRegenerate(c.nightsoulState.MaxPoints)

		switch c.armamentState {
		case ring:
			c.enterBike()
		default:
			c.exitBike()
		}
		c.SetCDWithDelay(action.ActionSkill, 15*60, 18)
		return action.Info{
			Frames:          frames.NewAbilFunc(skillFrames),
			AnimationLength: skillFrames[action.InvalidAction],
			CanQueueAfter:   skillFrames[action.ActionSwap],
			State:           action.SkillState,
		}, nil
	}
	c.enterNightsoulOrRegenerate(c.nightsoulState.MaxPoints)
	if h > 0 {
		return c.skillHold(), nil
	}
	return c.skillPress(), nil
}

func (c *char) enterBike() {
	c.armamentState = bike
	c.NormalHitNum = bikeHitNum
	c.c6Bike()
}

func (c *char) exitBike() {
	c.armamentState = ring
	c.NormalHitNum = normalHitNum
	c.ringSrc = c.Core.F

	c.QueueCharTask(c.skillRing(c.ringSrc), 120)
	c.c2Ring()
}

func (c *char) skillHold() action.Info {
	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "The Named Moment (Flamestrider)",
		AttackTag:      attacks.AttackTagElementalArt,
		ICDTag:         attacks.ICDTagNone,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeBlunt,
		PoiseDMG:       75,
		Element:        attributes.Pyro,
		Durability:     25,
		Mult:           skill[c.TalentLvlSkill()],
	}
	ap := combat.NewCircleHitOnTarget(
		c.Core.Combat.Player(),
		geometry.Point{Y: 1.0},
		6,
	)
	c.Core.QueueAttack(ai, ap, skillHitmark, skillHitmark, c.particleCB)
	c.enterBike()
	c.SetCDWithDelay(action.ActionSkill, 15*60, 18)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionSwap],
		State:           action.SkillState,
	}
}

func (c *char) skillPress() action.Info {
	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "The Named Moment",
		AttackTag:      attacks.AttackTagElementalArt,
		ICDTag:         attacks.ICDTagNone,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypePierce,
		Element:        attributes.Pyro,
		Durability:     25,
		Mult:           skill[c.TalentLvlSkill()],
	}
	ap := combat.NewCircleHitOnTarget(
		c.Core.Combat.Player(),
		geometry.Point{Y: 1.0},
		6,
	)
	c.Core.QueueAttack(ai, ap, skillHitmark, skillHitmark, c.particleCB)
	c.exitBike()
	c.SetCDWithDelay(action.ActionSkill, 15*60, 18)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionSwap],
		State:           action.SkillState,
	}
}

func (c *char) skillRing(src int) func() {
	return func() {
		if c.ringSrc != src {
			return
		}
		if c.armamentState != ring {
			return
		}
		if !c.nightsoulState.HasBlessing() {
			return
		}
		ai := combat.AttackInfo{
			ActorIndex:     c.Index,
			Abil:           "Rings of Searing Radiance",
			AttackTag:      attacks.AttackTagElementalArt,
			ICDTag:         attacks.ICDTagNone,
			AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
			ICDGroup:       attacks.ICDGroupDefault,
			StrikeType:     attacks.StrikeTypePierce,
			Element:        attributes.Pyro,
			Durability:     25,
			Mult:           skillRing[c.TalentLvlSkill()],
		}
		ap := combat.NewCircleHitOnTarget(
			c.Core.Combat.Player(),
			geometry.Point{Y: 1.0},
			6,
		)
		c.Core.QueueAttack(ai, ap, 0, 0, c.c6RingCB())
		c.reduceNightsoulPoints(3)
		c.QueueCharTask(c.skillRing(src), 120)
	}
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 0.5*60, true)
	c.Core.QueueParticle(c.Base.Key.String(), 5, attributes.Pyro, c.ParticleDelay)
}
