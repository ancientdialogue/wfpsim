package vodyanitsa

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c1Key = "vodyanitsa-c1"
	c2Key = "vodyanitsa-c2"
	c4Key = "vodyanitsa-c4"
	c6Key = "vodyanitsa-c6"
	c4Dur = 6 * 60
)

func (c *char) c1Init() {
	if c.Base.Cons < 1 {
		return
	}

	c.c1Buff = make([]float64, attributes.EndStatType)
}

func (c *char) c1OnHeal() {
	if c.Base.Cons < 1 {
		return
	}

	for _, char := range c.Core.Player.Chars() {
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(c1Key, 3*60),
			AffectedStat: attributes.ATK,
			Amount: func() []float64 {
				c.c1Buff[attributes.ATK] = c.MaxHP() * 0.007
				return c.c1Buff
			},
		})
	}
}

func (c *char) c2Init() {
	if c.Base.Cons < 2 {
		return
	}

	c.c2Buff = make([]float64, attributes.EndStatType)

	c.Core.Events.Subscribe(event.OnLunarReactionAttack, func(args ...any) {
		atk, ok := args[1].(*info.AttackEvent)
		if !ok {
			return
		}

		if c.c6c2CheckActive() && atk.Info.ActorIndex != c.Core.Player.Active() {
			return
		}

		if atk.Info.AttackTag != attacks.AttackTagReactionStellarSwirl {
			return
		}

		char := c.Core.Player.Chars()[atk.Info.ActorIndex]
		if !char.StatusIsActive(c2Key) {
			return
		}

		if !c.recentSSW() {
			return
		}

		atk.Snapshot.Stats[attributes.CD] += 0.6
	}, c2Key)
}

func (c *char) c2OnSkillAttack() {
	if c.Base.Cons < 2 {
		return
	}

	for _, char := range c.Core.Player.Chars() {
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase(c2Key, 5*60),
			Amount: func(atk *info.AttackEvent, _ info.Target) []float64 {
				if c.c6c2CheckActive() && atk.Info.ActorIndex != c.Core.Player.Active() {
					return nil
				}
				if c.recentSSW() {
					if atk.Info.AttackTag == attacks.AttackTagDirectStellarSwirl {
						c.c2Buff[attributes.CD] = 0.6
						return c.c2Buff
					}
				} else {
					if atk.Info.Element == attributes.Cryo || atk.Info.Element == attributes.Hydro {
						c.c2Buff[attributes.CD] = 0.5
						return c.c2Buff
					}
				}
				return nil
			},
		})
	}
}

func (c *char) c2SkillDurBonus() int {
	if c.Base.Cons < 2 {
		return 0
	}
	return 9 * 60
}

func (c *char) c4Init() {
	if c.Base.Cons < 4 {
		return
	}
	c.c4Buff = make([]float64, attributes.EndStatType)
	c.c4Stacks = NewRingQueue[int](3)
}

func (c *char) c4BeforeHeal() float64 {
	if c.Base.Cons < 4 {
		return 0
	}

	if c.Core.Player.ActiveChar().CurrentHPRatio() < 0.4 {
		return 0.5
	}

	c.c4Stacks.PushOverwrite(c.TimePassed)
	filter := func(src int) bool {
		return c.TimePassed <= src+c4Dur
	}

	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag(c4Key, c4Dur),
		AffectedStat: attributes.HPP,
		Amount: func() []float64 {
			c.c4Buff[attributes.HPP] = float64(c.c4Stacks.Count(filter)) * 0.2
			return c.c4Buff
		},
	})

	return 0
}

func (c *char) c6Init() {
	if c.Base.Cons < 6 {
		return
	}

	c.c6Buff = make([]float64, attributes.EndStatType)
	c.c6Buff[attributes.CryoP] = 0.5
	c.c6Buff[attributes.HydroP] = 0.5

	c.Core.Events.Subscribe(event.OnApplyAttack, func(args ...any) {
		atk := args[0].(*info.AttackEvent)

		if !c.StatusIsActive(skillKey) {
			return
		}
		// do not apply elevation to Reaction damage here because the elevation is already applied at the contributor level
		if atk.Info.AttackTag == attacks.AttackTagDirectStellarSwirl {
			atk.Info.Elevation += 0.25
		}
	}, c6Key)

	c.Core.Events.Subscribe(event.OnLunarReactionAttack, func(args ...any) {
		atk := args[1].(*info.AttackEvent)

		if !c.StatusIsActive(skillKey) {
			return
		}
		if atk.Info.AttackTag == attacks.AttackTagReactionStellarSwirl {
			atk.Info.Elevation += 0.25
		}
	}, c6Key)
}

func (c *char) c6c2CheckActive() bool {
	if c.Base.Cons < 6 {
		return true
	}
	return false
}

func (c *char) c6OnSkill() {
	if c.Base.Cons < 6 {
		return
	}

	for _, char := range c.Core.Player.Chars() {
		char.AddStatMod(character.StatMod{
			Base: modifier.NewBase(c6Key, 16*60+c.c2SkillDurBonus()),
			Amount: func() []float64 {
				return c.c6Buff
			},
		})
	}
}
