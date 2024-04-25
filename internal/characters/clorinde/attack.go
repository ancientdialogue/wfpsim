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
	attackHitmarks        = [][]int{{17}, {16}, {24, 32}, {37, 38, 39}, {30}}
	attackHitlagHaltFrame = [][]float64{{0.03}, {0.03}, {0.03, 0.03}, {0.02, 0.02, 0.02}, {0.03}}
	attackHitlagFactor    = [][]float64{{0.01}, {0.01}, {0.01, 0.01}, {0.05, 0.05, 0.05}, {0.05}}
	attackDefHalt         = [][]bool{{true}, {true}, {true, true}, {true, true, true}, {true}}
	attackHitboxes        = [][]float64{{1.7}, {1.9}, {2.1, 2.1}, {2, 3.5}, {2.5}} // n4 is a box
	attackOffsets         = []float64{1.1, 1.3, 1.2, 1.3, 1.4}
)

const normalHitNum = 5

func init() {
	attackFrames = make([][]int, normalHitNum)

	//TODO: these are all chiori frames

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 22) // N1 -> CA
	attackFrames[0][action.ActionAttack] = 19                                // N1 -> N2 rerecorded

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 33) // N2 -> CA
	attackFrames[1][action.ActionAttack] = 22                                // N2 -> N3 rerecorded

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 42) // N3 -> CA
	attackFrames[2][action.ActionAttack] = 41                                // N3 -> N4

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 42) // N4 -> CA
	attackFrames[3][action.ActionAttack] = 41                                // N4 -> N4

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 59) // N5 -> N1
	attackFrames[4][action.ActionCharge] = 500                               // TODO: this action is illegal; need better way to handle it
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
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
