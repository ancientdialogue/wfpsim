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

var (
	attackFrames         [][]int
	attackReleases       = [][]int{{17}, {12}, {20, 29}, {29}}
	multitargetTapFrames [][]int
)

const normalHitNum = 4

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackReleases[0][0], 32) // N1 -> Walk
	attackFrames[0][action.ActionAttack] = 22
	attackFrames[0][action.ActionAim] = 22

	attackFrames[1] = frames.InitNormalCancelSlice(attackReleases[1][0], 34) // N2 -> CA
	attackFrames[1][action.ActionAttack] = 24
	attackFrames[1][action.ActionWalk] = 31

	attackFrames[2] = frames.InitNormalCancelSlice(attackReleases[2][1], 86) // N3 -> Walk
	attackFrames[2][action.ActionAttack] = 39
	attackFrames[2][action.ActionAim] = 81

	attackFrames[3] = frames.InitNormalCancelSlice(attackReleases[3][0], 66) // N4 -> Walk
	attackFrames[3][action.ActionAttack] = 59
	attackFrames[3][action.ActionAim] = 500 // TODO: this action is illegal; need better way to handle it

	multitargetTapFrames = make([][]int, 1)
	multitargetTapFrames[0] = frames.InitAbilSlice(10)
	multitargetTapFrames[0][action.ActionAttack] = 10
	multitargetTapFrames[0][action.ActionAim] = 10
}

// Normal attack damage queue generator
// relatively standard with no major differences versus other bow characters
// Has "travel" parameter, used to set the number of frames that the arrow is in the air (default = 10)
func (c *char) Attack(p map[string]int) (action.Info, error) {
	if c.nightsoulState.HasBlessing() {
		return c.MultitargetFireTap(p)
	}
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  attacks.AttackTagNormal,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypePierce,
		Element:    attributes.Physical,
		Durability: 25,
	}

	for i, mult := range attack[c.NormalCounter] {
		ai.Mult = mult[c.TalentLvlAttack()]
		c.Core.QueueAttack(
			ai,
			combat.NewBoxHit(
				c.Core.Combat.Player(),
				c.Core.Combat.PrimaryTarget(),
				geometry.Point{Y: -0.5},
				0.1,
				1,
			),
			attackReleases[c.NormalCounter][i],
			attackReleases[c.NormalCounter][i]+travel,
		)
	}

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackReleases[c.NormalCounter][len(attackReleases[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}, nil
}

func (c *char) MultitargetFireTap(_ map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Multitarget Fire (Press)",
		AttackTag:  attacks.AttackTagNormal,
		ICDTag:     attacks.ICDTagChascaShot,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       skillMultitarget[c.TalentLvlSkill()],
	}

	ap := combat.NewBoxHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 5, 5)
	c.Core.QueueAttack(ai, ap, 0, 0, c.particleCB)

	return action.Info{
		Frames:          frames.NewAbilFunc(multitargetTapFrames[0]),
		AnimationLength: multitargetTapFrames[0][action.InvalidAction],
		CanQueueAfter:   multitargetTapFrames[0][action.ActionBurst],
		State:           action.NormalAttackState,
	}, nil
}
