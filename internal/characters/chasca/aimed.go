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

var aimedHitmarks = []int{15, 86}

var multitargetHitmarks = []int{3, 6, 9, 12, 15, 18}

var firstBulletLoad = 30

var additionalBulletLoad = 20

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
	multitargetFrames[0] = frames.InitAbilSlice(firstBulletLoad)
	multitargetFrames[0][action.ActionDash] = multitargetHitmarks[0]
	multitargetFrames[1] = frames.InitAbilSlice(additionalBulletLoad + firstBulletLoad)
	multitargetFrames[1][action.ActionAim] = multitargetHitmarks[1]
	multitargetFrames[2] = frames.InitAbilSlice(additionalBulletLoad*2 + firstBulletLoad)
	multitargetFrames[2][action.ActionAim] = multitargetHitmarks[2]
	multitargetFrames[3] = frames.InitAbilSlice(additionalBulletLoad*3 + firstBulletLoad)
	multitargetFrames[3][action.ActionAim] = multitargetHitmarks[3]
	multitargetFrames[4] = frames.InitAbilSlice(additionalBulletLoad*4 + firstBulletLoad)
	multitargetFrames[4][action.ActionAim] = multitargetHitmarks[4]
	multitargetFrames[5] = frames.InitAbilSlice(additionalBulletLoad*5 + firstBulletLoad)
	multitargetFrames[5][action.ActionAim] = multitargetHitmarks[5]

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

	for i, element := range c.shadowhuntShells {
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Shining Shadowhunt Shell",
			AttackTag:  attacks.AttackTagExtra,
			ICDTag:     attacks.ICDTagChascaShot,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    element,
			Durability: 25,
			Mult:       skillShiningShadowhunt[c.TalentLvlSkill()],
		}
		if element != attributes.Anemo && !c.StatusIsActive(c2icd) {
			c2ai := combat.AttackInfo{
				ActorIndex: c.Index,
				Abil:       "Shining Shadowhunt Shell",
				AttackTag:  attacks.AttackTagExtra,
				ICDTag:     attacks.ICDTagChascaShot,
				ICDGroup:   attacks.ICDGroupDefault,
				StrikeType: attacks.StrikeTypeDefault,
				Element:    element,
				Durability: 25,
				Mult:       skillShiningShadowhunt[c.TalentLvlSkill()],
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
		c.Core.QueueAttack(ai, ap, 0, 2*i, c.particleCB)
	}

	c.DeleteStatus(c2icd)

	return action.Info{
		Frames:          frames.NewAbilFunc(multitargetFrames[hold-1]),
		AnimationLength: multitargetFrames[hold-1][action.InvalidAction],
		CanQueueAfter:   multitargetFrames[hold-1][action.ActionBurst],
		State:           action.NormalAttackState,
	}, nil
}
