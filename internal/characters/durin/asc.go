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

var reactToElements = map[event.Event][]attributes.Element{
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

	for event, elements := range reactToElements {
		c.Core.Events.Subscribe(event, c.a1MakeBuff(elements), fmt.Sprintf("durin-a1-hook-%v", event))
	}
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

func (c *char) a1MakeBuff(elements []attributes.Element) func(args ...any) bool {
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
		return 0
	}

	if !c.StatusIsActive(blackKey) && !c.StatusIsActive(whiteKey) {
		return 0
	}

	return min(c.TotalAtk()*0.5, 1000)
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
