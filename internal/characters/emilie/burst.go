package emilie

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var (
	burstFrames []int
)

const (
	burstHitmark  = 104
	burstDuration = 168
	burstInterval = 18
	burstMarkKey  = "emilie-burst-mark"
)

func init() {
	burstFrames = frames.InitAbilSlice(127)
	burstFrames[action.ActionAttack] = 102
	burstFrames[action.ActionSkill] = 102
	burstFrames[action.ActionDash] = 103
	burstFrames[action.ActionJump] = 103
	burstFrames[action.ActionSwap] = 93
}

func (c *char) Burst(_ map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Aromatic Explication",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       burstDMG[c.TalentLvlBurst()],
	}

	oldLvl := c.lumidouceLvl

	c.QueueCharTask(func() {
		burstArea := combat.NewCircleHitOnTarget(c.lumidoucePos, nil, 12.5)

		if c.lumidouceSrc != -1 {
			c.removeLumi(c.lumidouceSrc)()
		}

		for i := 21; i <= burstDuration; i += burstInterval {
			enemy := c.Core.Combat.RandomEnemyWithinArea(
				burstArea,
				func(e combat.Enemy) bool {
					return !e.StatusIsActive(burstMarkKey)
				},
			)
			var pos geometry.Point
			if enemy != nil {
				pos = enemy.Pos()
				enemy.AddStatus(burstMarkKey, 0.7*60, false)
			} else {
				pos = geometry.CalcRandomPointFromCenter(burstArea.Shape.Pos(), 0.5, 12.5, c.Core.Rand)
			}

			ap := combat.NewCircleHitOnTarget(pos, geometry.Point{Y: 4.5}, 2.5)
			c.Core.QueueAttack(
				ai,
				ap,
				0,
				0,
			)
		}
	}, 21)

	c.QueueCharTask(func() {
		c.lumidouceLvl = oldLvl

		c.lumidouceSrc = c.Core.F

		c.genScents()

		c.Core.Tasks.Add(c.lumiTick(c.Core.F), skillLumiFirstTick)
		c.Core.Tasks.Add(c.removeLumi(c.Core.F), 22*60)
	}, burstDuration)

	c.ConsumeEnergy(12)
	c.SetCD(action.ActionBurst, 13.5*60)
	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap],
		State:           action.BurstState,
	}, nil
}
