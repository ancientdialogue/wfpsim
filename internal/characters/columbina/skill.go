package columbina

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

var skillFrames []int

const (
	skillHitmark   = 24
	particleICDKey = "columbina-particle-icd"
	skillKey       = "columbina-skill"
	gravityKey     = "columbina-gravity"
	gravityMax     = 60
	LCInd          = 0
	// LBInd          = 1
	// LCrInd         = 2
)

func init() {
	skillFrames = frames.InitAbilSlice(61)
	skillFrames[action.ActionDash] = 29
	skillFrames[action.ActionJump] = 29
}

func (c *char) skillInit() {
	c.Core.Events.Subscribe(event.OnLunarCharged, func(args ...any) bool {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return false
		}
		if !c.StatusIsActive(skillKey) {
			return false
		}
		c.AddStatus(gravityKey, 2*60, false)
		if !c.gravityTask {
			c.gravityAccum()
		}
		return false
	}, "columbina-gravity-lc")

	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...any) bool {
		atk := args[1].(*info.AttackEvent)
		if !attackTagIsLunar(atk.Info.AttackTag) {
			return false
		}
		if !c.StatusIsActive(skillKey) {
			return false
		}
		c.AddStatus(gravityKey, 2*60, false)
		if !c.gravityTask {
			c.gravityAccum()
		}
		return false
	}, "columbina-gravity-on-dmg")
}

func attackTagIsLunar(atkTag attacks.AttackTag) bool {
	return attacks.LunarReactionStartDelim < atkTag && atkTag < attacks.LunarReactionEndDelim || attacks.DirectLunarReactionStartDelim < atkTag && atkTag < attacks.DirectLunarReactionEndDelim
}

func (c *char) gravityAccum() {
	if !c.StatusIsActive(gravityKey) {
		c.gravityTask = false
		return
	}

	if !c.StatusIsActive(skillKey) {
		c.gravityTask = false
		return
	}
	c.gravityTask = true
	c.gravity[c.gravityLastReaction] += 1 // 10 gravity per 1s
	if c.totalGravity() >= gravityMax {
		c.gravityTick()
	}
	c.QueueCharTask(c.gravityAccum, 6)
}

func (c *char) totalGravity() float64 {
	sum := 0.0
	for _, g := range c.gravity {
		sum += g
	}
	return sum
}

func (c *char) clearGravity() {
	for i := range c.gravity {
		c.gravity[i] = 0
	}
}

func (c *char) gravityTick() {
	maxReaction := 0
	maxGravity := 0.0
	for i, g := range c.gravity {
		if g > maxGravity {
			maxGravity = g
			maxReaction = i
		}
	}
	c.clearGravity()

	var mult float64
	var atkTag attacks.AttackTag
	var elem attributes.Element
	var abil string
	switch maxReaction {
	case LCInd:
		mult = skillLC[c.TalentLvlSkill()]
		atkTag = attacks.AttackTagDirectLunarCharged
		elem = attributes.Electro
		abil = "Skill Gravity (Lunar-Charged)"
	default:
		return
	}

	ai := info.AttackInfo{
		ActorIndex:       c.Index(),
		Abil:             abil,
		AttackTag:        atkTag,
		ICDTag:           attacks.ICDTagNone,
		StrikeType:       attacks.StrikeTypeDefault,
		Element:          elem,
		Durability:       0,
		UseHP:            true,
		IgnoreDefPercent: 1.0,
		Mult:             mult,
	}

	ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 6)
	c.a1OGravityTick()

	c.Core.QueueAttack(ai, ap, 1, 1)
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	c.QueueCharTask(func() {
		c.skillSrc = c.Core.F
		ai := info.AttackInfo{
			ActorIndex: c.Index(),
			Abil:       "Skill",
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Hydro,
			Durability: 25,
			UseHP:      true,
			Mult:       skill[c.TalentLvlSkill()],
		}
		ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 6)
		c.Core.QueueAttack(ai, ap, 0, 0)
		if !c.StatusIsActive(skillKey) {
			c.clearGravity()
		}
		c.AddStatus(skillKey, 25*60+1, true)
		c.QueueCharTask(c.skillTickTask(c.skillSrc), 126)
		c.SetCDWithDelay(action.ActionSkill, 17*60, 0)
	}, skillHitmark)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillHitmark,
		State:           action.SkillState,
	}, nil
}

func (c *char) particleCB(a info.AttackCB) {
	if a.Target.Type() != info.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 3*60, false)
	if c.Core.Rand.Float64() < 0.67 {
		c.Core.QueueParticle(c.Base.Key.String(), 1, attributes.Hydro, c.ParticleDelay)
	} else {
		c.Core.QueueParticle(c.Base.Key.String(), 2, attributes.Hydro, c.ParticleDelay)
	}
}

// Helper function that handles damage, healing, and particle components of every tick of her E
func (c *char) skillTick() {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Skill (DoT)",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		UseHP:      true,
		Mult:       skillDoT[c.TalentLvlSkill()],
	}
	ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 6)
	c.Core.QueueAttack(ai, ap, 0, 0, c.particleCB)
}

// Handles repeating skill damage ticks. Split into a separate function as you can only have 1 jellyfish on field at once
// Skill snapshots, so inputs into the function are the originating snapshot
func (c *char) skillTickTask(src int) func() {
	return func() {
		// Basically stops "old" casts of E from working, and also stops further ticks from that source
		if c.skillSrc > src {
			return
		}

		if !c.StatusIsActive(skillKey) {
			return
		}

		c.skillTick()

		c.Core.Tasks.Add(c.skillTickTask(src), 120)
	}
}
