package durin

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var a1ReactToElements = map[event.Event][]attributes.Element{
	event.OnOverload:        {attributes.Electro, attributes.Pyro},
	event.OnSwirlPyro:       {attributes.Anemo, attributes.Pyro},
	event.OnCrystallizePyro: {attributes.Geo, attributes.Pyro},
	event.OnBurning:         {attributes.Dendro, attributes.Pyro},
}

const a1BlackKey = "durin-a1-black"

func (c *char) a1Init() {
	if c.Base.Ascension < 1 {
		return
	}

	for event, elements := range a1ReactToElements {
		c.Core.Events.Subscribe(event, c.a1MakeResShred(elements), fmt.Sprintf("durin-a1-hook-%v", event))
	}

	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...any) bool {
		t, ok := args[0].(*enemy.Enemy)
		atk := args[1].(*info.AttackEvent)
		if !ok {
			return false
		}
		if !t.IsBurning() {
			return false
		}
		switch atk.Info.Element {
		case attributes.Dendro:
		case attributes.Pyro:
		default:
			return false
		}

		if !c.StatusIsActive(burstKeyWhite) {
			return false
		}

		t.AddResistMod(info.ResistMod{
			Base:  modifier.NewBaseWithHitlag("durin-a1-"+attributes.Dendro.String(), 6*60),
			Ele:   attributes.Dendro,
			Value: -0.20 * c.magicA1Bonus(),
		})

		t.AddResistMod(info.ResistMod{
			Base:  modifier.NewBaseWithHitlag("durin-a1-"+attributes.Pyro.String(), 6*60),
			Ele:   attributes.Pyro,
			Value: -0.20 * c.magicA1Bonus(),
		})

		return false
	}, "durin-a1-hook-on-dmg")
}

func (c *char) a1OnBurst(isWhite bool) {
	if c.Base.Ascension < 1 {
		return
	}

	if isWhite {
		c.DeleteReactBonusMod(a1BlackKey)
		return
	}

	c.AddReactBonusMod(character.ReactBonusMod{
		Base: modifier.NewBaseWithHitlag(a1BlackKey, 20*60),
		Amount: func(ai info.AttackInfo) (float64, bool) {
			if ai.Amped {
				return 0.40 * c.magicA1Bonus(), false
			}
			return 0, false
		},
	})
}

func (c *char) a1MakeResShred(elements []attributes.Element) func(args ...any) bool {
	return func(args ...any) bool {
		t, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}

		if !c.StatusIsActive(burstKeyWhite) {
			return false
		}

		for _, ele := range elements {
			t.AddResistMod(info.ResistMod{
				Base:  modifier.NewBaseWithHitlag("durin-a1-"+ele.String(), 6*60),
				Ele:   ele,
				Value: -0.20 * c.magicA1Bonus(),
			})
		}
		return false
	}
}

func (c *char) a4Dmg() float64 {
	if c.Base.Ascension < 4 {
		return 1.0
	}

	if !c.StatusIsActive(blackKey) && !c.StatusIsActive(whiteKey) {
		return 1.0
	}

	bonus := min(c.TotalAtk()/100*0.01, 0.25)

	return 1 + bonus
}

func (c *char) magicA1Bonus() float64 {
	if !c.IsMagic {
		return 1
	}

	if c.getMagicCount() < 2 {
		return 1
	}

	return 1.75
}

func (c *char) getMagicCount() int {
	count := 0
	for _, c := range c.Core.Player.Chars() {
		if c.IsMagic {
			count += 1
		}
	}
	return count
}
