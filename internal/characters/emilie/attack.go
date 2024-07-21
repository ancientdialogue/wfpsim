package emilie

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
	attackFrames   [][]int
	attackHitmarks = []int{15, 12, 21, 27}

	attackOffsets = [][]float64{{-0.1, 0.1, 0.5, -2.5}, {-0.1, 0.1, 0.8, -2.5}}

	attackHitlagHaltFrame = []float64{0.01, 0.01, 0.02, 0.02}
	attackHitboxes        = [][][]float64{{{1.5, 2.8}, {1.7}, {1.9}, {5, 6}}, {{1.5, 3}, {2.3}, {2.2}, {6, 7}}}
	attackStrikeType      = []attacks.StrikeType{attacks.StrikeTypeSpear, attacks.StrikeTypeSlash, attacks.StrikeTypeSlash, attacks.StrikeTypeSlash}
)

const normalHitNum = 4

func init() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 34) // N1 -> CA/Walk
	attackFrames[0][action.ActionAttack] = 31

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 28) // N2 -> CA/Walk
	attackFrames[1][action.ActionAttack] = 23

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 48) // N3 -> Walk
	attackFrames[2][action.ActionAttack] = 36
	attackFrames[2][action.ActionCharge] = 45

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 58) // N4 -> CA/Walk
	attackFrames[3][action.ActionAttack] = 53
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:          attacks.AttackTagNormal,
		ICDTag:             attacks.ICDTagNormalAttack,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attackStrikeType[c.NormalCounter],
		Element:            attributes.Physical,
		Durability:         25,
		Mult:               attack[c.NormalCounter][c.TalentLvlAttack()],
		HitlagFactor:       0.01,
		HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter] * 60,
		CanBeDefenseHalted: true,
	}

	c.QueueCharTask(
		func() {
			var c6cb combat.AttackCBFunc
			var ap combat.AttackPattern
			c6Index := 0
			// TODO: Check if DMG bonus still applies if c6 runs out between start of NA and the hit

			// if c.Base.Cons >= 6 && c.StatusIsActive(c6Key) {
			// 	c6Index = 1
			// 	ai.Element = attributes.Dendro
			// 	ai.IgnoreInfusion = true
			// 	ai.FlatDmg = 2.5 * c.getTotalAtk()
			// 	c6cb = c.c6cb
			// }
			switch c.NormalCounter {
			case 0, 3:
				ap = combat.NewBoxHitOnTarget(
					c.Core.Combat.Player(),
					geometry.Point{Y: attackOffsets[c6Index][c.NormalCounter]},
					attackHitboxes[c6Index][c.NormalCounter][0],
					attackHitboxes[c6Index][c.NormalCounter][1],
				)
			case 1, 2:
				ap = combat.NewCircleHitOnTarget(
					c.Core.Combat.Player(),
					geometry.Point{Y: attackOffsets[c6Index][c.NormalCounter]},
					attackHitboxes[c6Index][c.NormalCounter][0],
				)
			}
			c.Core.QueueAttack(ai, ap, 0, 0, c6cb)
		}, attackHitmarks[c.NormalCounter],
	)

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}, nil
}
