package chasca

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

var multitargetFrames [][]int

var c6MultitargetFrames [][]int

var c6FatalRoundsFrames [][]int

var aimedHitmarks = []int{15, 86}

var multitargetHitmarks = []int{3, 6, 9, 12, 15, 18}

var firstBulletLoadFrames = 33

var additionalBulletLoadFrames = 20

var c6firstBulletLoadFrames = 20

var c6AdditionalBulletLoadFrames = 14

var fatalBulletLoadFrames = 2

func init() {
	aimedFrames = make([][]int, 2)

	// Aimed Shot
	aimedFrames[0] = frames.InitAbilSlice(26)
	aimedFrames[0][action.ActionDash] = aimedHitmarks[0]
	aimedFrames[0][action.ActionJump] = aimedHitmarks[0]

	// Fully-Charged Aimed Shot
	aimedFrames[1] = frames.InitAbilSlice(96)
	aimedFrames[1][action.ActionDash] = aimedHitmarks[1]
	aimedFrames[1][action.ActionJump] = aimedHitmarks[1]

	multitargetFrames = make([][]int, 6)
	multitargetFrames[0] = frames.InitAbilSlice(firstBulletLoadFrames)
	multitargetFrames[0][action.ActionAim] = multitargetHitmarks[0]
	multitargetFrames[1] = frames.InitAbilSlice(additionalBulletLoadFrames + firstBulletLoadFrames)
	multitargetFrames[1][action.ActionAim] = multitargetHitmarks[1]
	multitargetFrames[2] = frames.InitAbilSlice(additionalBulletLoadFrames*2 + firstBulletLoadFrames)
	multitargetFrames[2][action.ActionAim] = multitargetHitmarks[2]
	multitargetFrames[3] = frames.InitAbilSlice(additionalBulletLoadFrames*3 + firstBulletLoadFrames)
	multitargetFrames[3][action.ActionAim] = multitargetHitmarks[3]
	multitargetFrames[4] = frames.InitAbilSlice(additionalBulletLoadFrames*4 + firstBulletLoadFrames)
	multitargetFrames[4][action.ActionAim] = multitargetHitmarks[4]
	multitargetFrames[5] = frames.InitAbilSlice(additionalBulletLoadFrames*5 + firstBulletLoadFrames)
	multitargetFrames[5][action.ActionAim] = multitargetHitmarks[5]

	c6MultitargetFrames = make([][]int, 6)
	c6MultitargetFrames[0] = frames.InitAbilSlice(c6firstBulletLoadFrames)
	c6MultitargetFrames[0][action.ActionAim] = multitargetHitmarks[0]
	c6MultitargetFrames[1] = frames.InitAbilSlice(c6AdditionalBulletLoadFrames + c6firstBulletLoadFrames)
	c6MultitargetFrames[1][action.ActionAim] = multitargetHitmarks[1]
	c6MultitargetFrames[2] = frames.InitAbilSlice(c6AdditionalBulletLoadFrames*2 + c6firstBulletLoadFrames)
	c6MultitargetFrames[2][action.ActionAim] = multitargetHitmarks[2]
	c6MultitargetFrames[3] = frames.InitAbilSlice(c6AdditionalBulletLoadFrames*3 + c6firstBulletLoadFrames)
	c6MultitargetFrames[3][action.ActionAim] = multitargetHitmarks[3]
	c6MultitargetFrames[4] = frames.InitAbilSlice(c6AdditionalBulletLoadFrames*4 + c6firstBulletLoadFrames)
	c6MultitargetFrames[4][action.ActionAim] = multitargetHitmarks[4]
	c6MultitargetFrames[5] = frames.InitAbilSlice(c6AdditionalBulletLoadFrames*5 + c6firstBulletLoadFrames)
	c6MultitargetFrames[5][action.ActionAim] = multitargetHitmarks[5]

	c6FatalRoundsFrames = make([][]int, 6)
	c6FatalRoundsFrames[0] = frames.InitAbilSlice(fatalBulletLoadFrames * 1)
	c6FatalRoundsFrames[0][action.ActionAim] = multitargetHitmarks[0]
	c6FatalRoundsFrames[1] = frames.InitAbilSlice(fatalBulletLoadFrames * 2)
	c6FatalRoundsFrames[1][action.ActionAim] = multitargetHitmarks[1]
	c6FatalRoundsFrames[2] = frames.InitAbilSlice(fatalBulletLoadFrames * 3)
	c6FatalRoundsFrames[2][action.ActionAim] = multitargetHitmarks[2]
	c6FatalRoundsFrames[3] = frames.InitAbilSlice(fatalBulletLoadFrames * 4)
	c6FatalRoundsFrames[3][action.ActionAim] = multitargetHitmarks[3]
	c6FatalRoundsFrames[4] = frames.InitAbilSlice(fatalBulletLoadFrames * 5)
	c6FatalRoundsFrames[4][action.ActionAim] = multitargetHitmarks[4]
	c6FatalRoundsFrames[5] = frames.InitAbilSlice(fatalBulletLoadFrames * 6)
	c6FatalRoundsFrames[5][action.ActionAim] = multitargetHitmarks[5]
}

func (c *char) Aimed(p map[string]int) (action.Info, error) {
	if c.nightsoulState.HasBlessing() {
		return c.MultitargetFireHold(p)
	}
	hold, ok := p["hold"]
	if !ok {
		hold = attacks.AimParamLv1
	}
	switch hold {
	case attacks.AimParamPhys:
	case attacks.AimParamLv1:
	default:
		return action.Info{}, fmt.Errorf("invalid hold param supplied, got %v", hold)
	}
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	weakspot := p["weakspot"]

	ai := combat.AttackInfo{
		ActorIndex:           c.Index,
		Abil:                 "Fully-Charged Aimed Shot",
		AttackTag:            attacks.AttackTagExtra,
		ICDTag:               attacks.ICDTagNone,
		ICDGroup:             attacks.ICDGroupDefault,
		StrikeType:           attacks.StrikeTypePierce,
		Element:              attributes.Anemo,
		Durability:           25,
		Mult:                 aim[c.TalentLvlAttack()],
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
		aimedHitmarks[hold],
		aimedHitmarks[hold]+travel,
	)

	return action.Info{
		Frames:          frames.NewAbilFunc(aimedFrames[hold]),
		AnimationLength: aimedFrames[hold][action.InvalidAction],
		CanQueueAfter:   aimedHitmarks[hold],
		State:           action.AimState,
	}, nil
}

func (c *char) MultitargetFireHold(p map[string]int) (action.Info, error) {
	hold, ok := p["hold"]
	if !ok {
		hold = 6 // Default 6 bullet
	}

	if hold < 1 || hold > 6 {
		return action.Info{}, fmt.Errorf("invalid hold param supplied, got %v", hold)
	}
	c.loadShadowhuntShells(hold)
	for _, element := range c.shadowhuntShells {
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Shining Shadowhunt Shell",
			AttackTag:  attacks.AttackTagExtra,
			ICDTag:     attacks.ICDTagChascaShot,
			ICDGroup:   attacks.ICDGroupChascaShot,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    element,
			Durability: 25,
			Mult:       skillShiningShadowhunt[c.TalentLvlSkill()],
		}
		if element != attributes.Anemo && c.Base.Cons >= 2 && !c.StatusIsActive(c2icd) {
			c2ai := combat.AttackInfo{
				ActorIndex: c.Index,
				Abil:       "C2 Shining Shadowhunt Shell",
				AttackTag:  attacks.AttackTagExtra,
				ICDTag:     attacks.ICDTagChascaShot,
				ICDGroup:   attacks.ICDGroupChascaShot,
				StrikeType: attacks.StrikeTypeDefault,
				Element:    element,
				Durability: 25,
				Mult:       400 / 100,
			}
			ap := combat.NewBoxHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 5, 5)
			c.Core.QueueAttack(c2ai, ap, 0, 2)

			c.AddStatus(c2icd, -1, false)
		}
		if element == attributes.Anemo {
			ai.Abil = "Shadowhunt Shell"
			ai.Element = attributes.Anemo
			ai.Mult = skillShadowhunt[c.TalentLvlSkill()]
		}
		ap := combat.NewBoxHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 1, 1)

		c.Core.QueueAttack(ai, ap, 0, firstBulletLoadFrames+additionalBulletLoadFrames*(hold-1), c.particleCB)
		c.c6()
	}

	c.DeleteStatus(c2icd)
	if c.Base.Cons < 6 {
		return action.Info{
			Frames:          frames.NewAbilFunc(multitargetFrames[hold-1]),
			AnimationLength: multitargetFrames[hold-1][action.InvalidAction],
			CanQueueAfter:   multitargetFrames[hold-1][action.ActionSkill],
			State:           action.ChargeAttackState,
		}, nil
	}
	return action.Info{
		Frames:          frames.NewAbilFunc(c6MultitargetFrames[hold-1]),
		AnimationLength: c6MultitargetFrames[hold-1][action.InvalidAction],
		CanQueueAfter:   c6MultitargetFrames[hold-1][action.ActionBurst],
		State:           action.ChargeAttackState,
	}, nil
}
