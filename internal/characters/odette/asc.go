package odette

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

const (
	stellarBonusKey  = "odette-stellar-bonus"
	radianceSwirlKey = "radiance-stellar-swirl"
	a1Key            = "odette-a1"
)

func (c *char) a1Init() {
	if c.Base.Ascension < 1 {
		return
	}

	c.AddReactBonusMod(character.ReactBonusMod{
		Base: modifier.NewBase(a1Key+"-buff", -1),
		Amount: func(ai info.AttackInfo) float64 {
			switch ai.AttackTag {
			case
				attacks.AttackTagDirectStellarConduct,
				attacks.AttackTagDirectStellarSwirl,
				attacks.AttackTagReactionStellarSwirl:
			default:
				return 0
			}
			if !c.StatusIsActive(danceDoubleKey) {
				return 0
			}
			return float64(c.a1StacksSelf) * 0.15
		},
	})

	for _, char := range c.Core.Player.Chars() {
		if char.Index() == c.Index() {
			continue
		}
		char.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBase(a1Key+"-buff", -1),
			Amount: func(ai info.AttackInfo) float64 {
				switch ai.AttackTag {
				case
					attacks.AttackTagDirectStellarConduct,
					attacks.AttackTagDirectStellarSwirl,
					attacks.AttackTagReactionStellarSwirl:
				default:
					return 0
				}
				if !c.StatusIsActive(danceDoubleKey) {
					return 0
				}
				return float64(c.a1StacksOthers) * 0.15
			},
		})
	}

	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...any) {
		prev := args[0].(int)
		next := args[1].(int)
		if prev == c.Index() {
			src := c.Core.F
			c.a1Src = src
			c.Core.Tasks.Add(func() { c.a1Ticker(src) }, 60)
		} else if next == c.Index() {
			// cancel the a1Ticker
			c.a1Src = -1
		}
	}, a1Key)
}

func (c *char) a1Ticker(src int) {
	// don't need to check asc because it's only called by a1Init()
	if c.a1Src != src {
		return
	}

	if c.a1StacksSelf == 0 {
		return
	}

	// TODO: This check isn't needed because we cancel the task when we swap back to Odette
	if c.Core.Player.Active() == c.Index() {
		return
	}

	if !c.StatusIsActive(danceDoubleKey) {
		return
	}

	stacks := min(c.c1a1Remove(), c.a1StacksSelf)
	c.a1StacksSelf -= stacks * c.c6a1ReduceMod()
	c.a1StacksOthers = min(c.a1StacksOthers+stacks, c.c1a1Stacks())

	c.Core.Tasks.Add(func() { c.a1Ticker(src) }, 60)
}

func (c *char) a1OnDanceSummon() {
	if c.Base.Ascension < 1 {
		return
	}

	c.a1StacksSelf = c.c1a1Stacks()
	c.a1StacksOthers = 0
}

func (c *char) a4StellarGlimmerMult() float64 {
	if c.Base.Ascension < 1 {
		return 1
	}

	scaling := max(c.TotalAtk()-1000, 0)
	buff := min(scaling/100*0.015, 0.3)

	return 1.0 + buff
}

func (c *char) stellarInit() {
	reactable.EnableStellarConduct(c.Core)
	reactable.EnableStellarSwirl(c.Core)

	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) {
		atk := args[1].(*info.AttackEvent)

		switch atk.Info.AttackTag {
		case attacks.AttackTagDirectStellarConduct:
		case attacks.AttackTagDirectStellarSwirl:
		default:
			return
		}

		bonus := min(c.TotalAtk()/100.0*0.007, 0.14)

		if c.Core.Flags.LogDebug {
			c.Core.Log.NewEvent("odette adding stellar base damage", glog.LogCharacterEvent, c.Index()).Write("bonus", bonus)
		}

		atk.Info.BaseDmgBonus += bonus
	}, stellarBonusKey)

	c.Core.Events.Subscribe(event.OnLunarReactionAttack, func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.AttackTag != attacks.AttackTagReactionStellarSwirl {
			return
		}

		bonus := min(c.TotalAtk()/100.0*0.007, 0.14)

		if c.Core.Flags.LogDebug {
			c.Core.Log.NewEvent("odette adding stellar base damage", glog.LogCharacterEvent, c.Index()).Write("bonus", bonus)
		}

		atk.Info.BaseDmgBonus += bonus
	}, stellarBonusKey+"-reaction")

	c.Core.Events.Subscribe(event.OnStellarSwirl, func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}

		c.AddStatus(radianceSwirlKey, 8*60, false)
	}, stellarBonusKey)
}
