package klee

import (
	"fmt"
	"math"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	attackFrames          [][]int
	attackFramesWithLag   [][]int
	attackHitmarks        = []int{15, 21, 34}
	attackPoiseDMG        = []float64{65, 65, 130}
	attackHitmarksWithLag []int
	attackRadius          = []float64{1, 1, 1.5}
)

const normalHitNum = 3

func init() {
	attackHitmarksWithLag = make([]int, len(attackHitmarks))
	copy(attackHitmarksWithLag, attackHitmarks)
	for i := range attackHitmarksWithLag {
		attackHitmarksWithLag[i] += 9
	}
	attackFrames = make([][]int, normalHitNum)

	// N1 -> x
	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 34)
	attackFrames[0][action.ActionAttack] = 23
	attackFrames[0][action.ActionCharge] = 23
	attackFrames[0][action.ActionSkill] = 6
	attackFrames[0][action.ActionBurst] = 6
	attackFrames[0][action.ActionDash] = 6
	attackFrames[0][action.ActionJump] = 6

	// N2 -> x
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 41)
	attackFrames[1][action.ActionAttack] = 27
	attackFrames[1][action.ActionCharge] = 25
	attackFrames[1][action.ActionSkill] = 2
	attackFrames[1][action.ActionBurst] = 2
	attackFrames[1][action.ActionDash] = 2
	attackFrames[1][action.ActionJump] = 2

	// N3 -> x
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 72)
	attackFrames[2][action.ActionAttack] = 46
	attackFrames[2][action.ActionCharge] = 47

	// N1 -> x (9f lag)
	attackFramesWithLag = make([][]int, len(attackFrames))
	for i := range attackFrames {
		attackFramesWithLag[i] = make([]int, len(attackFrames[i]))
		copy(attackFramesWithLag[i], attackFrames[i])
	}
	add9FrameLag(attackFramesWithLag[0])
}

// klee has 9f lag on normals when using after skill/dash/burst
func add9FrameLag(frames []int) {
	for i := range frames {
		switch action.Action(i) {
		case action.ActionBurst,
			action.ActionDash,
			action.ActionJump,
			action.ActionSkill:
		default:
			frames[i] += 9
		}
	}
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  attacks.AttackTagNormal,
		ICDTag:     attacks.ICDTagKleeFireDamage,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		PoiseDMG:   attackPoiseDMG[c.NormalCounter],
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}

	done := false
	tryPerformAttack := func() {
		if done {
			return
		}
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(
				c.Core.Combat.Player(),
				c.Core.Combat.PrimaryTarget(),
				nil,
				attackRadius[c.NormalCounter],
			),
			0,
			travel,
			c.makeA1CB(),
		)
		c.c1(travel)
		done = true
	}

	// Hexerei: Secret Rite
	// When she performs the third Normal Attack in the sequence, an Explosive Spark
	// will be consumed to unleash an additional attack equivalent to Boom-Boom Strike.
	if c.IsHexerei && c.Core.Player.GetHexereiCount() > 1 && c.StatusIsActive(a1SparkKey) {
		if c.NormalCounter == 2 || c.Base.Cons == 6 && c.NormalCounter < 2 && c.Core.Rand.Float64() < 0.4 {
			c.queueCoordinatedCharge()
		}
	}

	defer func() {
		c.AdvanceNormalIndex()
		c.savedNormalCounter = c.NormalCounter
	}()

	adjustedFrames := attackFrames
	adjustedHitmarks := attackHitmarks
	switch c.Core.Player.CurrentState() {
	case action.DashState,
		action.SkillState,
		action.BurstState:
		adjustedFrames = attackFramesWithLag
		adjustedHitmarks = attackHitmarksWithLag
	}

	canQueueAfter := math.MaxInt
	for _, f := range adjustedFrames[c.NormalCounter] {
		if f < canQueueAfter {
			canQueueAfter = f
		}
	}
	actionInfo := action.Info{
		Frames:          frames.NewAttackFunc(c.Character, adjustedFrames),
		AnimationLength: adjustedFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   canQueueAfter,
		State:           action.NormalAttackState,
		OnRemoved: func(next action.AnimationState) {
			switch next {
			case action.SkillState,
				action.BurstState,
				action.DashState,
				action.JumpState:
				tryPerformAttack()
			}
		},
	}
	actionInfo.QueueAction(tryPerformAttack, adjustedHitmarks[c.NormalCounter])
	return actionInfo, nil
}

func (c *char) queueCoordinatedCharge() {
	delay := 30
	travel := 10

	c.Core.Tasks.Add(func() {
		ai := c.getChargeAttackInfo()
		snap := c.applySpark(&ai)
		ai.Abil += " (Coordinated)"
		c.Core.QueueAttackWithSnap(
			ai,
			snap,
			combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 3),
			travel,
			c.makeA4CB(),
		)
	}, delay)

	c.Core.Log.NewEvent("coordinated CA triggered", glog.LogCharacterEvent, c.Index()).
		Write("expected hit", c.Core.F+delay+travel)

	c.c1(delay - travel)
}
