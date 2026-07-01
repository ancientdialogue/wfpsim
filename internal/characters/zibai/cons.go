package zibai

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c1Key = "zibai-c1"
	c2Key = "zibai-c2"
	c4Key = "zibai-c4"
	c6Key = "zibai-c6"
)

func (c *char) c1Init() {
	if c.Base.Cons < 1 {
		return
	}

	c.AddReactBonusMod(character.ReactBonusMod{
		Base: modifier.NewBase(c1Key+"-buff", -1),
		Amount: func(ai info.AttackInfo) float64 {
			if ai.ActorIndex != c.Index() {
				return 0
			}

			if !c.StatusIsActive(c1Key) {
				return 0
			}

			if ai.Abil != skillAbil2 {
				return 0
			}
			if c.Core.Flags.LogDebug {
				c.Core.Log.NewEvent("Adding C1 react bonus", glog.LogCharacterEvent, c.Index())
			}

			c.QueueCharTask(func() { c.DeleteStatus(c1Key) }, 1)
			return 2.20
		},
	})
}

func (c *char) c1MaxSkillsPerSkill() int {
	if c.Base.Cons < 1 {
		return 4
	}
	return 5
}

func (c *char) c1OnSkill() {
	if c.Base.Cons < 1 {
		return
	}
	c.addRadiance(100)
	c.AddStatus(c1Key, 15*60, true)
}

func (c *char) c2Init() {
	if c.Base.Cons < 2 {
		return
	}

	for _, char := range c.Core.Player.Chars() {
		char.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBase(c2Key+"-buff", -1),
			Amount: func(ai info.AttackInfo) float64 {
				if !c.StatusIsActive(skillKey) {
					return 0
				}

				if c.Core.Flags.LogDebug {
					c.Core.Log.NewEvent("Adding C2 react bonus", glog.LogCharacterEvent, c.Index())
				}

				switch ai.AttackTag {
				case attacks.AttackTagReactionLunarCrystallize:
				case attacks.AttackTagDirectLunarCrystallize:
				default:
					return 0
				}
				return 0.3
			},
		})
	}
}

func (c *char) c2A4Mult() float64 {
	if c.Base.Cons < 2 {
		return 0.6
	}

	if c.Core.Player.GetMoonsignLevel() < 2 {
		return 0.6
	}

	return 5.5
}

func (c *char) c4ResetNormalCount() {
	if c.Base.Cons < 4 {
		c.Character.ResetNormalCounter()
		return
	}

	if !c.StatusIsActive(skillKey) {
		c.Character.ResetNormalCounter()
		return
	}
}

func (c *char) c4SkillCB(ac info.AttackCB) {
	if c.Base.Cons < 4 {
		return
	}

	if _, ok := ac.Target.(*enemy.Enemy); !ok {
		return
	}
	c.AddStatus(c4Key, 30*60, false)
}

func (c *char) c4N4Bonus() float64 {
	if c.Base.Cons < 4 {
		return 1.0
	}

	if !c.StatusIsActive(c4Key) {
		return 1.0
	}
	return 2.5
}

func (c *char) c6Init() {
	if c.Base.Cons < 6 {
		return
	}
	hook := func(aeInd int, tag attacks.AttackTag) func(args ...any) {
		return func(args ...any) {
			atk := args[aeInd].(*info.AttackEvent)
			if atk.Info.AttackTag != tag {
				return
			}

			if !c.StatusIsActive(c6Key) {
				return
			}

			if atk.Info.ActorIndex != c.Index() {
				return
			}

			if c.Core.Flags.LogDebug {
				c.Core.Log.NewEvent("Adding zibai c6 lunar crystallize elevation", glog.LogCharacterEvent, c.Index()).Write("amt", c.c6Elev)
			}
			atk.Info.Elevation += c.c6Elev
		}
	}

	c.Core.Events.Subscribe(event.OnApplyAttack, hook(0, attacks.AttackTagDirectLunarCrystallize), c6Key+"-direct")
	c.Core.Events.Subscribe(event.OnLunarReactionAttack, hook(1, attacks.AttackTagReactionLunarCrystallize), c6Key+"-reaction")
}

func (c *char) c6RadianceEff() float64 {
	if c.Base.Cons < 6 {
		return 1.0
	}
	return 1.5
}

func (c *char) c6ConsumeRadiance() {
	if c.Base.Cons < 6 {
		c.consumeRadiance(70)
		return
	}

	c.c6Elev = max((c.radiance-70)*0.016, 0)
	c.AddStatus(c6Key, 3*60, true)
	c.consumeRadiance(c.radiance)
}
