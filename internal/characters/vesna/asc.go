package vesna

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

const (
	a1Key           = "vesna-a1"
	a4Key           = "vesna-a4"
	stellarBonusKey = "vesna-ssw-bonus"
)

func (c *char) a1Init() {
	if c.Base.Ascension < 1 {
		return
	}

	c.a1Stacks = NewRingQueue[int](6)
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...any) {
		// do nothing if previous char wasn't vesna
		prev := args[0].(int)
		if prev != c.Index() {
			return
		}
		c.a1Stacks.Clear()
	}, "vesna-a1-exit")
}

func (c *char) a1OnSkill() {
	if c.Base.Ascension < 1 {
		return
	}

	c.a1Stacks.Clear()
}

func (c *char) a1OnSpecialSkillOrBurst() {
	if c.Base.Ascension < 1 {
		return
	}

	c.a1AddStacks()
}

func (c *char) a1AddStacks() {
	if c.Base.Ascension < 1 {
		return
	}

	c.a1Stacks.PushOverwrite(c.TimePassed)
}

func (c *char) a1StackCount() int {
	if c.Base.Ascension < 1 {
		return 0
	}
	filter := func(src int) bool {
		return c.TimePassed <= src+20*60
	}

	return c.a1Stacks.Count(filter)
}

func (c *char) a1Mult() float64 {
	if c.Base.Ascension < 1 {
		return 1.0
	}

	return 1.0 + 0.1*float64(c.a1StackCount())
}

func (c *char) a4Init() {
	if c.Base.Ascension < 4 {
		return
	}

	m := make([]float64, attributes.EndStatType)
	cryoOrAnemo := 0
	other := 0
	for _, char := range c.Core.Player.Chars() {
		if char.Index() == c.Index() {
			continue
		}
		switch char.Base.Element {
		case attributes.Cryo, attributes.Anemo:
			cryoOrAnemo += 1
		default:
			other += 1
		}
	}

	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase(a4Key, -1),
		AffectedStat: attributes.NoStat,
		Amount: func() []float64 {
			if !c.isRadianceSSw() {
				return nil
			}

			m[attributes.ATKP] = 0.06 * float64(cryoOrAnemo) * c.c4a4Mult()
			m[attributes.EM] = 25 * float64(other) * c.c4a4Mult()
			return m
		},
	})
}

func (c *char) stellarInit() {
	c.Core.Flags.Custom[reactable.StellarSwirlEnableKey] = 1
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) {
		atk := args[1].(*info.AttackEvent)

		if atk.Info.AttackTag != attacks.AttackTagDirectStellarSwirl {
			return
		}

		bonus := min(c.TotalAtk()/100.0*0.007, 0.14)

		if c.Core.Flags.LogDebug {
			c.Core.Log.NewEvent("vesna adding stellar swirl base damage", glog.LogCharacterEvent, c.Index()).Write("bonus", bonus)
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
			c.Core.Log.NewEvent("vesna adding stellar swirl base damage", glog.LogCharacterEvent, c.Index()).Write("bonus", bonus)
		}

		atk.Info.BaseDmgBonus += bonus
	}, stellarBonusKey+"-reaction")
}
