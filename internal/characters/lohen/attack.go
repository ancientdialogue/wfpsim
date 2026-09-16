package lohen

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
	attackHitmarks        = [][]int{{8}, {16}, {13, 13 + 11, 13 + 11 + 11}, {10}, {18, 18 + 20}}
	attackHitlagHaltFrame = [][]float64{{0.03}, {0.03}, {0, 0, 0.03}, {0.03}, {0, 0.09}}
	attackHitboxes        = [][][]float64{
		{{1.2, 3}},
		{{1.6, 3.3}},
		{{1.2, 3.3}, {1.2, 3.3}, {1.2, 3.3}},
		{{1.6, 3.3}},
		{{1.6, 3.3}, {1.2, 3.3}},
	}
)

const (
	normalHitNum = 5
)

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 26)
	attackFrames[0][action.ActionAttack] = 12
	attackFrames[0][action.ActionCharge] = 12

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 39)
	attackFrames[1][action.ActionAttack] = 23
	attackFrames[1][action.ActionCharge] = 18

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][2], 55)
	attackFrames[2][action.ActionAttack] = 37
	attackFrames[2][action.ActionCharge] = 36

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 46)
	attackFrames[3][action.ActionAttack] = 30
	attackFrames[3][action.ActionCharge] = 16

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][1], 73)
	attackFrames[4][action.ActionAttack] = 69
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(skillKey) {
		return c.skillAttack(), nil
	}
	for i, mult := range attack[c.NormalCounter] {
		ai := info.AttackInfo{
			ActorIndex:         c.Index(),
			Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
			Mult:               mult[c.TalentLvlAttack()],
			AttackTag:          attacks.AttackTagNormal,
			ICDTag:             attacks.ICDTagNormalAttack,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeSpear,
			Element:            attributes.Physical,
			Durability:         25,
			HitlagFactor:       0.01,
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter][i] * 60,
			CanBeDefenseHalted: true,
		}

		ap := combat.NewBoxHitOnTarget(
			c.Core.Combat.Player(),
			nil,
			attackHitboxes[c.NormalCounter][i][0],
			attackHitboxes[c.NormalCounter][i][1],
		)
		c.QueueCharTask(func() {
			c.Core.QueueAttack(ai, ap, 0, 0)
		}, attackHitmarks[c.NormalCounter][i])
	}

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}, nil
}

func (c *char) skillAttack() action.Info {
	for i, mult := range skillAttack[c.NormalCounter] {
		ai := info.AttackInfo{
			ActorIndex:         c.Index(),
			Abil:               fmt.Sprintf("Normal %v (Masterstroke)", c.NormalCounter),
			Mult:               mult[c.TalentLvlSkill()],
			AttackTag:          attacks.AttackTagNormal,
			ICDTag:             attacks.ICDTagNormalAttack,
			ICDGroup:           attacks.ICDGroupLohenSkillAttack,
			StrikeType:         attacks.StrikeTypeSpear,
			Element:            attributes.Cryo,
			Durability:         25,
			HitlagFactor:       0.01,
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter][i] * 60,
			CanBeDefenseHalted: true,
			IgnoreInfusion:     true,
		}

		ap := combat.NewBoxHitOnTarget(
			c.Core.Combat.Player(),
			nil,
			attackHitboxes[c.NormalCounter][i][0],
			attackHitboxes[c.NormalCounter][i][1],
		)
		c.QueueCharTask(func() {
			c.Core.QueueAttack(ai, ap, 0, 0, c.particleCB, c.joyCB, c.c2MakeCB())
		}, attackHitmarks[c.NormalCounter][i])
	}

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}
}
