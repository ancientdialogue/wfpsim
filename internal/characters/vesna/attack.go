package vesna

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	attackFrames          [][]int
	attackHitmarks        = [][]int{{14}, {6}, {9, 9 + 12}, {16}, {12}, {24}}
	attackHitlagHaltFrame = [][]float64{{0.04}, {0.04}, {0.04, 0.04}, {0.04}, {0.04}, {0.04}}
	attackDefHalt         = [][]bool{{true}, {true}, {true, true}, {true}, {true}, {true}}
	attackHitboxes        = [][]float64{{1.7}, {2}, {1, 1.5}, {1.7}, {1.7}, {1.7}}
	attackOffsets         = []float64{1.8, 0.8, 0.5, 1.8, 1.8, 1.8}
	attackFanAngles       = []float64{360, 180, 360, 360, 360, 360}
)

const normalHitNum = 6

func init() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 27)
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 23)
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 35)
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 27)
	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][0], 22)
	attackFrames[5] = frames.InitNormalCancelSlice(attackHitmarks[5][0], 50)
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(c6Key) {
		return c.c6Attack()
	}

	for i, hitmark := range attackHitmarks[c.NormalCounter] {
		ai := info.AttackInfo{
			ActorIndex:         c.Index(),
			Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
			AttackTag:          attacks.AttackTagNormal,
			ICDTag:             attacks.ICDTagNormalAttack,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeSlash,
			Element:            attributes.Physical,
			Durability:         25,
			Mult:               attack[c.NormalCounter][c.TalentLvlAttack()],
			HitlagFactor:       0.01,
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter][i] * 60,
			CanBeDefenseHalted: attackDefHalt[c.NormalCounter][i],
		}
		if c.NormalCounter == 1 || c.NormalCounter == 4 {
			ai.StrikeType = attacks.StrikeTypeSlash
		}

		if c.StatusIsActive(skillKey) {
			ai.Element = attributes.Anemo
			ai.IgnoreInfusion = true
		}

		ap := combat.NewCircleHitOnTargetFanAngle(
			c.Core.Combat.Player(),
			info.Point{Y: attackOffsets[c.NormalCounter]},
			attackHitboxes[c.NormalCounter][0],
			attackFanAngles[c.NormalCounter],
		)
		if c.NormalCounter == 2 {
			ap = combat.NewBoxHitOnTarget(
				c.Core.Combat.Player(),
				info.Point{Y: attackOffsets[c.NormalCounter]},
				attackHitboxes[c.NormalCounter][0],
				attackHitboxes[c.NormalCounter][1],
			)
		}
		c.Core.QueueAttack(ai, ap, hitmark, hitmark)
	}

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackFrames[c.NormalCounter][action.ActionDash],
		State:           action.NormalAttackState,
	}, nil
}
