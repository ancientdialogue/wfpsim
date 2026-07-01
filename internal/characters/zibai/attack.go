package zibai

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
	attackHits            = []int{1, 1, 2, 1}
	attackHitmarks        = [][]int{{10}, {10}, {10, 20}, {26}}
	attackHitlagHaltFrame = [][]float64{{0.00}, {0.04}, {0.04, 0.00}, {0.04}}
	attackDefHalt         = [][]bool{{false}, {true}, {true, false}, {true}}
	attackHitboxes        = [][]float64{{1.5, 3.8}, {2}, {1, 1.5}, {1.7}}
	attackOffsets         = []float64{0, 0.8, 0.5, 1.8}
	attackFanAngles       = []float64{360, 180, 360, 360}
)

// 242
const normalHitNum = 4

func init() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 30)
	attackFrames[0][action.ActionAttack] = 24

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 30)
	attackFrames[1][action.ActionAttack] = 22

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 28)
	attackFrames[2][action.ActionAttack] = 32

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 34)
	attackFrames[3][action.ActionCharge] = 500 // TODO: this action is illegal; need better way to handle it
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(skillKey) {
		return c.skillAttack()
	}
	for i := range attackHits[c.NormalCounter] {
		ai := info.AttackInfo{
			ActorIndex:         c.Index(),
			Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
			AttackTag:          attacks.AttackTagNormal,
			ICDTag:             attacks.ICDTagNormalAttack,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeSlash,
			Element:            attributes.Physical,
			Durability:         25,
			Mult:               attack[c.NormalCounter][i][c.TalentLvlAttack()],
			HitlagFactor:       0.01,
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter][i] * 60,
			CanBeDefenseHalted: attackDefHalt[c.NormalCounter][i],
		}
		if c.NormalCounter == 1 || c.NormalCounter == 4 {
			ai.StrikeType = attacks.StrikeTypeSlash
		}
		ap := combat.NewCircleHitOnTargetFanAngle(
			c.Core.Combat.Player(),
			info.Point{Y: attackOffsets[c.NormalCounter]},
			attackHitboxes[c.NormalCounter][0],
			attackFanAngles[c.NormalCounter],
		)
		if c.NormalCounter == 0 || c.NormalCounter == 2 {
			ap = combat.NewBoxHitOnTarget(
				c.Core.Combat.Player(),
				info.Point{Y: attackOffsets[c.NormalCounter]},
				attackHitboxes[c.NormalCounter][0],
				attackHitboxes[c.NormalCounter][1],
			)
		}
		// the multihit part generates no hitlag so this is fine
		c.Core.QueueAttack(ai, ap, attackHitmarks[c.NormalCounter][i], attackHitmarks[c.NormalCounter][i])
	}

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackFrames[c.NormalCounter][action.ActionDash],
		State:           action.NormalAttackState,
	}, nil
}

func (c *char) skillAttack() (action.Info, error) {
	for i := range attackHits[c.NormalCounter] {
		ai := info.AttackInfo{
			ActorIndex:         c.Index(),
			Abil:               fmt.Sprintf("Normal %v (Skill)", c.NormalCounter),
			AttackTag:          attacks.AttackTagNormal,
			ICDTag:             attacks.ICDTagNormalAttack,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeSlash,
			Element:            attributes.Geo,
			IgnoreInfusion:     true,
			Durability:         25,
			Mult:               skillAttack[c.NormalCounter][i][c.TalentLvlAttack()],
			UseDef:             true,
			HitlagFactor:       0.01,
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter][i] * 60,
			CanBeDefenseHalted: attackDefHalt[c.NormalCounter][i],
		}

		ap := combat.NewCircleHitOnTargetFanAngle(
			c.Core.Combat.Player(),
			info.Point{Y: attackOffsets[c.NormalCounter]},
			attackHitboxes[c.NormalCounter][0],
			attackFanAngles[c.NormalCounter],
		)
		if c.NormalCounter == 0 || c.NormalCounter == 2 {
			ap = combat.NewBoxHitOnTarget(
				c.Core.Combat.Player(),
				info.Point{Y: attackOffsets[c.NormalCounter]},
				attackHitboxes[c.NormalCounter][0],
				attackHitboxes[c.NormalCounter][1],
			)
		}
		if c.NormalCounter == 3 && c.Core.Player.GetMoonsignLevel() >= 2 {
			ai.Mult = skillLastAttackBonus[c.TalentLvlAttack()]
			c.Core.QueueAttack(ai, ap, attackHitmarks[c.NormalCounter][i], attackHitmarks[c.NormalCounter][i], c.particleCB, c.radianceCB)
			ai.Abil += lunarCrystallizeAbil
			ai.AttackTag = attacks.AttackTagDirectLunarCrystallize
			ai.Durability = 0
			ai.HitlagHaltFrames = 0
			ai.CanBeDefenseHalted = false
			ai.IgnoreDefPercent = 1
			ai.Mult *= c.c4N4Bonus()
			c.Core.QueueAttack(ai, ap, attackHitmarks[c.NormalCounter][i], attackHitmarks[c.NormalCounter][i], c.particleCB, c.radianceCB)
		} else {
			c.Core.QueueAttack(ai, ap, attackHitmarks[c.NormalCounter][i], attackHitmarks[c.NormalCounter][i], c.particleCB, c.radianceCB)
		}
	}

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackFrames[c.NormalCounter][action.ActionDash],
		State:           action.NormalAttackState,
	}, nil
}
