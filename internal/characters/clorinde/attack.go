package clorinde

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
	attackFrames [][]int
	// TODO: these are made up hitmarks
	attackHitmarks        = [][]int{{17}, {12}, {27, 31}, {16, 22, 23}, {18}}
	attackHitlagHaltFrame = [][]float64{{0.03}, {0.03}, {0.03, 0.03}, {0.02, 0.02, 0.02}, {0.03}}
	attackHitlagFactor    = [][]float64{{0.01}, {0.01}, {0.01, 0.01}, {0.05, 0.05, 0.05}, {0.05}}
	attackDefHalt         = [][]bool{{true}, {true}, {true, true}, {true, true, true}, {true}}
	attackHitboxes        = [][]float64{{1.7}, {1.9}, {2.1, 2.1}, {2, 3.5}, {2.5}} // n4 is a box
	attackOffsets         = []float64{1.1, 1.3, 1.2, 1.3, 1.4}

	skillAttackFrames []int
)

const (
	skillAttackHitmark = 8 //  TODO: travel time??
)

const normalHitNum = 5

func init() {
	attackFrames = make([][]int, normalHitNum)

	//TODO: these frames are just basically random guesses

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 20)
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 15)
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 42)
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][2], 34)
	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][0], 56)

	skillAttackFrames = frames.InitAbilSlice(18) // TODO: this is a rough estimate
	skillAttackFrames[action.ActionSkill] = 13
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(skillStateKey) {
		return c.skillAttack(p)
	}

	for i, mult := range attack[c.NormalCounter] {
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
			AttackTag:          attacks.AttackTagNormal,
			ICDTag:             attacks.ICDTagNormalAttack,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeSlash,
			Element:            attributes.Physical,
			Durability:         25,
			Mult:               mult[c.TalentLvlAttack()],
			HitlagFactor:       attackHitlagFactor[c.NormalCounter][i],
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter][i] * 60,
			CanBeDefenseHalted: attackDefHalt[c.NormalCounter][i],
		}

		ap := combat.NewCircleHitOnTarget(
			c.Core.Combat.Player(),
			geometry.Point{Y: attackOffsets[c.NormalCounter]},
			attackHitboxes[c.NormalCounter][0],
		)
		if c.NormalCounter == 3 {
			ap = combat.NewBoxHitOnTarget(
				c.Core.Combat.Player(),
				geometry.Point{Y: attackOffsets[c.NormalCounter]},
				attackHitboxes[c.NormalCounter][0],
				attackHitboxes[c.NormalCounter][1],
			)
		}

		c.Core.QueueAttack(ai, ap, attackHitmarks[c.NormalCounter][i], attackHitmarks[c.NormalCounter][i])
	}

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}, nil
}

func (c *char) skillAttack(_ map[string]int) (action.Info, error) {
	// TODO: not sure if we need a counter here
	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Piercing Shot",
		AttackTag:      attacks.AttackTagNormal,
		ICDTag:         attacks.ICDTagNormalAttack,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeSlash,
		Element:        attributes.Electro,
		Durability:     25,
		Mult:           skillEnhancedNA[c.TalentLvlSkill()],
		IgnoreInfusion: true,
	}
	t := c.Core.Combat.PrimaryTarget()
	var ap combat.AttackPattern
	if c.currentHPDebtRatio() < 1 {
		// TODO: assume this is just a big rectangle center on target
		ap = combat.NewBoxHitOnTarget(t, nil, 2, 14)
	} else {
		ai.Abil = "Normal Shot"
		ai.Mult = skillNA[c.TalentLvlSkill()]
		// TODO: how big is this??
		ap = combat.NewCircleHitOnTarget(t, nil, 2)
	}
	// TODO: assume no snapshotting on this
	c.Core.QueueAttack(ai, ap, skillAttackHitmark, skillAttackHitmark)

	// TODO: timing on this?
	c.gainBOLOnAttack()

	return action.Info{
		Frames:          frames.NewAbilFunc(skillAttackFrames),
		AnimationLength: skillAttackFrames[action.InvalidAction],
		CanQueueAfter:   skillAttackFrames[action.ActionSkill], //TODO: fastest cancel?
		State:           action.NormalAttackState,
	}, nil
}
