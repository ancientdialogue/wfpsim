package emilie

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

const burstDelay = 110
const burstDur = 2.8 * 60 // 168

const burstTargettingICD = 0.7 * 60
const burstTargettingICDKey = "emilie-burst-targetting-icd"

var burstFrames []int

func init() {
	burstFrames = frames.InitAbilSlice(118)
	burstFrames[action.ActionAttack] = 114
	burstFrames[action.ActionAim] = 114
	burstFrames[action.ActionSkill] = 117
	burstFrames[action.ActionDash] = 117
	burstFrames[action.ActionSwap] = 115
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	c.QueueCharTask(func() {
		var level = max(c.lumidouceLvl, 1)
		// remove current lumi
		c.removeLumi(c.lumidouceSrc)()
		c.QueueCharTask(c.removeLumi(c.Core.F), burstDur+1)
		c.queueLumiSpawn(burstDur+1, skillLumiFirstTick, level)

		player := c.Core.Combat.Player()
		c.lumidoucePos = geometry.CalcOffsetPoint(player.Pos(), geometry.Point{Y: 1.5}, player.Direction())
		c.lumidouceSrc = c.Core.F
		c.lumidouceLvl = 3

		for i := 0; i < burstDur; i += 18 {
			c.Core.Tasks.Add(c.burstTick, i)
		}
		c.QueueCharTask(c.scentCheck(c.Core.F), 0)
	}, burstDelay)
	c.ConsumeEnergy(6)
	c.SetCD(action.ActionBurst, 13.5*60)
	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *char) burstTick() {
	burstArea := combat.NewCircleHitOnTarget(c.lumidoucePos, nil, 8)
	enemy := c.Core.Combat.RandomEnemyWithinArea(
		burstArea,
		func(e combat.Enemy) bool {
			return !e.StatusIsActive(burstTargettingICDKey)
		},
	)
	var pos geometry.Point
	if enemy != nil {
		pos = enemy.Pos()
		enemy.AddStatus(burstTargettingICDKey, burstTargettingICD, true) // same enemy can't be targeted again for 0.7s
	} else {
		pos = geometry.CalcRandomPointFromCenter(burstArea.Shape.Pos(), 0.5, 9.5, c.Core.Rand)
	}
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Lumidouce Case (Level 3)",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypePierce,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       burstDMG[c.TalentLvlBurst()],
	}

	// deal dmg after a certain delay
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(pos, nil, 2.5), 8, 8)
}
