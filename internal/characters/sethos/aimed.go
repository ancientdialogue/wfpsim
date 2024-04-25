package sethos

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var aimedFrames [][]int
var aimedA1Frames []int

var aimedHitmarks = []int{15, 86, 375}
var startCharge = aimedHitmarks[0]

func init() {
	// outside of E status
	aimedFrames = make([][]int, 3)

	// Aimed Shot
	aimedFrames[0] = frames.InitAbilSlice(23)
	aimedFrames[0][action.ActionDash] = aimedHitmarks[0]
	aimedFrames[0][action.ActionJump] = aimedHitmarks[0]

	// Fully-Charged Aimed Shot
	aimedFrames[1] = frames.InitAbilSlice(94)
	aimedFrames[1][action.ActionDash] = aimedHitmarks[1]
	aimedFrames[1][action.ActionJump] = aimedHitmarks[1]

	// Fully-Charged Aimed Shot (Crowfeather)
	aimedA1Frames = frames.InitAbilSlice(383)
	aimedA1Frames[action.ActionDash] = aimedHitmarks[2]
	aimedA1Frames[action.ActionJump] = aimedHitmarks[2]
}

func (c *char) Aimed(p map[string]int) (action.Info, error) {
	hold, ok := p["hold"]
	if !ok {
		// is this a good default? it's gonna take 6s to do without energy
		hold = attacks.AimParamLv2
	}
	switch hold {
	case attacks.AimParamPhys:
	case attacks.AimParamLv1:
	case attacks.AimParamLv2:
		return c.ShadowPierce(p)
	default:
		return action.Info{}, fmt.Errorf("invalid hold param supplied, got %v", hold)
	}

	skip, energy := c.a1Calc()
	if skip > aimedHitmarks[hold]-startCharge {
		skip = aimedHitmarks[hold] - startCharge
	}

	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	weakspot := p["weakspot"]

	c.QueueCharTask(func() {
		ai := combat.AttackInfo{
			ActorIndex:           c.Index,
			Abil:                 "Fully-Charged Aimed Shot",
			AttackTag:            attacks.AttackTagExtra,
			ICDTag:               attacks.ICDTagNone,
			ICDGroup:             attacks.ICDGroupDefault,
			StrikeType:           attacks.StrikeTypePierce,
			Element:              attributes.Dendro,
			Durability:           25,
			Mult:                 fullaim[c.TalentLvlAttack()],
			HitWeakPoint:         weakspot == 1,
			HitlagHaltFrames:     0.12 * 60,
			HitlagFactor:         0.01,
			HitlagOnHeadshotOnly: true,
			IsDeployable:         true,
		}
		if hold < attacks.AimParamLv1 {
			ai.Abil = "Aimed Shot"
			ai.Element = attributes.Physical
			ai.Mult = aim[c.TalentLvlAttack()]
		}

		c.Core.QueueAttack(
			ai,
			combat.NewBoxHit(
				c.Core.Combat.Player(),
				c.Core.Combat.PrimaryTarget(),
				geometry.Point{Y: -0.5},
				0.1,
				1,
			),
			0,
			travel,
		)
		if hold == attacks.AimParamLv1 {
			c.Energy -= energy * 0.5
		}
	}, aimedHitmarks[hold]-skip)

	return action.Info{
		Frames:          func(next action.Action) int { return aimedFrames[hold][next] - skip },
		AnimationLength: aimedFrames[hold][action.InvalidAction] - skip,
		CanQueueAfter:   aimedHitmarks[hold] - skip,
		State:           action.AimState,
	}, nil
}

func (c *char) ShadowPierce(p map[string]int) (action.Info, error) {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	weakspot := p["weakspot"]

	skip, energy := c.a1Calc()
	if skip > aimedHitmarks[2]-startCharge {
		skip = aimedHitmarks[2] - startCharge
	}
	hitHaltFrames := 0.0
	if weakspot == 1 {
		hitHaltFrames = 0.12 * 60
	}

	c.QueueCharTask(func() {
		ai := combat.AttackInfo{
			ActorIndex:           c.Index,
			Abil:                 "Shadow Piercing Arrow",
			AttackTag:            attacks.AttackTagExtra,
			ICDTag:               attacks.ICDTagNone,
			ICDGroup:             attacks.ICDGroupDefault,
			StrikeType:           attacks.StrikeTypePierce,
			Element:              attributes.Electro,
			Durability:           25,
			Mult:                 shadowpierceAtk[c.TalentLvlAttack()],
			HitWeakPoint:         weakspot == 1,
			HitlagHaltFrames:     hitHaltFrames,
			HitlagFactor:         0.01,
			HitlagOnHeadshotOnly: true,
			IsDeployable:         true,
			FlatDmg:              shadowpierceEM[c.TalentLvlAttack()] * c.Stat(attributes.EM),
		}

		deltaPos := c.Core.Combat.Player().Pos().Sub(c.Core.Combat.PrimaryTarget().Pos())
		dist := deltaPos.Magnitude()

		// simulate piercing. Extends from player to 15 units behind primary target
		ap := combat.NewBoxHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			geometry.Point{Y: -dist},
			0.1,
			15+dist,
		)

		c.Core.QueueAttack(
			ai,
			ap,
			0,
			travel,
		)
		c.Energy -= energy
	}, aimedHitmarks[2]-skip)

	return action.Info{
		Frames:          func(next action.Action) int { return aimedFrames[2][next] - skip },
		AnimationLength: aimedFrames[2][action.InvalidAction] - skip,
		CanQueueAfter:   aimedHitmarks[2] - skip,
		State:           action.AimState,
	}, nil
}
