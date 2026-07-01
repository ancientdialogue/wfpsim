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
	attackHitmarks        = [][]int{{12}, {8}, {5, 15, 24}, {21}, {11, 18}}
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

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 20)

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 17)

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][2], 27)
	attackFrames[2][action.ActionCharge] = 24

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 20)
	attackFrames[3][action.ActionCharge] = 20

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][1], 70)
	attackFrames[4][action.ActionCharge] = 500 // TODO: this action is illegal; need better way to handle it
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
