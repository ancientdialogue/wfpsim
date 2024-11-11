package chasca

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var c2icd = "chasca-c2-icd"
var c4energy = "chasca-c4-energy"
var c4icd = "chasca-c4-icd"

const fatalRoundsKey = "chasca-fatal-rounds"
const fatalRoundsIcd = 60 * 3
const fatalRoundsCDBuff = 1.2
const fatalBulletLoad = 6

func (c *char) c2stacks() int {
	if c.Base.Cons >= 2 {
		return 1
	}
	return 0
}

func (c *char) c4energy() {
	if c.Base.Cons < 4 {
		return
	}
	c.AddEnergy(c4energy, 1.5)
}

func (c *char) c6() (action.Info, error) {
	if c.Base.Cons < 6 {
		return action.Info{}, fmt.Errorf("no C6")
	}
	if !c.StatusIsActive(fatalRoundsKey) {
		c.loadShadowhuntShells(fatalBulletLoad) // assuming it always load 6 bullets for now
		for _, element := range c.shadowhuntShells {
			c6ai := combat.AttackInfo{
				ActorIndex: c.Index,
				Abil:       "Fatal Rounds: Shining Shadowhunt Shell",
				AttackTag:  attacks.AttackTagExtra,
				ICDTag:     attacks.ICDTagChascaShot,
				ICDGroup:   attacks.ICDGroupDefault,
				StrikeType: attacks.StrikeTypeDefault,
				Element:    element,
				Durability: 25,
				Mult:       skillShiningShadowhunt[c.TalentLvlSkill()],
			}
			ap := combat.NewBoxHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 5, 5)
			c.Core.QueueAttack(c6ai, ap, 0, 2)
		}

		c6DMG := make([]float64, attributes.EndStatType)
		c6DMG[attributes.CD] = fatalRoundsCDBuff
		c.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase("chasca-c6-fatal-rounds-cdmg-bonus", -1),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				if atk.Info.Abil != "Fatal Rounds: Shining Shadowhunt Shell" {
					return nil, false
				}
				return c6DMG, true
			},
		})
	}

	c.AddStatus(fatalRoundsKey, fatalRoundsIcd, false)
	return action.Info{
		Frames:          frames.NewAbilFunc(c6FatalRoundsFrames[fatalBulletLoad-1]),
		AnimationLength: multitargetFrames[fatalBulletLoad-1][action.InvalidAction],
		CanQueueAfter:   multitargetFrames[fatalBulletLoad-1][action.ActionBurst],
		State:           action.NormalAttackState,
	}, nil
}
