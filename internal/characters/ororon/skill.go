package ororon

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var skillFrames []int

const skillHitmark = 30

func init() {
	skillFrames = frames.InitAbilSlice(52) // E -> D
	skillFrames[action.ActionAttack] = 29  // E -> N1
	skillFrames[action.ActionAim] = 30     // E -> CA
	skillFrames[action.ActionBurst] = 32   // E -> Q
	skillFrames[action.ActionJump] = 51    // E -> J
	skillFrames[action.ActionSwap] = 50    // E -> Swap
}

const particleICDKey = "ororon-particle-icd"

func (c *char) Skill(p map[string]int) (action.Info, error) {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}

	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Night's Sling",
		AttackTag:      attacks.AttackTagElementalArt,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagElementalArt,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Electro,
		Durability:     25,
		Mult:           skill[c.TalentLvlSkill()],
	}

	enemies := []targets.TargetKey{c.Core.Combat.PrimaryTarget().Key()}
	maxHits := 3 + c.c1ExtraBounce()
	for i := 0; len(enemies) < maxHits && i < c.Core.Combat.EnemyCount(); i++ {
		newKey := c.Core.Combat.Enemies()[i].Key()
		if newKey == c.Core.Combat.PrimaryTarget().Key() {
			continue
		}
		enemies = append(enemies, newKey)
	}
	for i, e := range enemies {
		c.Core.QueueAttack(
			ai,
			combat.NewSingleTargetHit(e),
			skillHitmark,
			skillHitmark+travel*(i+1),
			c.particleCB,
			c.makeA4cb(),
			c.makeC1cb(),
		)
	}

	c.SetCDWithDelay(action.ActionSkill, 15*60, 7)
	c.a1OnSkill()
	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionAttack], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 6*60, true)
	c.Core.QueueParticle(c.Base.Key.String(), 3, attributes.Electro, c.ParticleDelay)
}
