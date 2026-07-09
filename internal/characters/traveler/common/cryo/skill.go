package cryo

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	skillFrames       [][]int
	skillTickHitmarks = []int{5, 8}
)

const (
	skillTapHitmark     = 24
	skillFirstTickDelay = 60
	skillInterval       = 3 * 60
	skillSscICD         = 0.05 * 60
	particleICDKey      = "travelercryo-particle-icd"
	skillKey            = "travelercryo-e"
	skillICDKey         = "travelercryo-e-icd"
	skillStacksMax      = 8
)

func init() {
	skillFrames = make([][]int, 2)

	// Male
	// Tap
	skillFrames[0] = frames.InitAbilSlice(49) // E -> N1
	skillFrames[0][action.ActionDash] = 31
	skillFrames[0][action.ActionJump] = 31
	skillFrames[0][action.ActionSwap] = 48

	// Female
	// Tap
	skillFrames[1] = frames.InitAbilSlice(49) // E -> N1
	skillFrames[1][action.ActionDash] = 31
	skillFrames[1][action.ActionJump] = 31
	skillFrames[1][action.ActionSwap] = 48
}

func (c *Traveler) Skill(p map[string]int) (action.Info, error) {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	c.skillTravel = travel

	ai := info.AttackInfo{
		ActorIndex:     c.Index(),
		Abil:           "Pierce of the Ice Fog",
		AttackTag:      attacks.AttackTagElementalArt,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagTravelerHoldDMG,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Cryo,
		Durability:     25,
		Mult:           skill[c.TalentLvlSkill()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 3.0),
		skillTapHitmark,
		skillTapHitmark,
		c.particleCB,
	)

	src := c.Core.F
	c.skillSrc = src
	c.QueueCharTask(func() { c.skillTicker(src) }, skillFirstTickDelay)
	c.AddStatus(skillKey, c.skillDur+skillFirstTickDelay, false)
	c.SetCD(action.ActionSkill, 15*60)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames[c.gender]),
		AnimationLength: skillFrames[c.gender][action.InvalidAction],
		CanQueueAfter:   skillFrames[c.gender][action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *Traveler) skillTicker(src int) {
	if c.skillSrc != src {
		return
	}

	if !c.StatusIsActive(skillKey) {
		return
	}

	if c.getRadiance() != radianceStellarConduct {
		c.queueCrystal(skillTickHitmarks...)
	}

	c.QueueCharTask(func() { c.skillTicker(src) }, skillInterval)
}

func (c *Traveler) queueCrystal(hitmarks ...int) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Frostpierce Star",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagElementalArt,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       crystal[c.TalentLvlSkill()],
	}

	for _, delay := range hitmarks {
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 3.0),
			delay,
			delay+c.skillTravel,
			c.crystalCB,
			c.c2CB,
		)
	}
}

func (c *Traveler) naCaPlungeCB(a info.AttackCB) {
	if a.Target.Type() != info.TargettableEnemy {
		return
	}

	if c.getRadiance() != radianceStellarConduct {
		return
	}

	if !c.StatusIsActive(skillKey) {
		return
	}

	if c.StatusIsActive(skillICDKey) {
		return
	}

	c.AddStatus(skillICDKey, skillSscICD, true)
	c.queueCrystal(3)
}

func (c *Traveler) crystalCB(a info.AttackCB) {
	if a.Target.Type() != info.TargettableEnemy {
		return
	}

	c.flostglowStacks = min(c.flostglowStacks+1, skillStacksMax)
}

func (c *Traveler) particleCB(a info.AttackCB) {
	if a.Target.Type() != info.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 0.5*60, true)

	c.Core.QueueParticle(c.Base.Key.String(), 3.0, attributes.Cryo, c.ParticleDelay)
}
