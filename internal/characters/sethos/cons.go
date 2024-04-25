package sethos

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const c2Key = "sethos-c2"
const c2Dur = 10 * 60

func (c *char) c1() {
	if c.Base.Cons < 1 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	m[attributes.CR] = 0.15
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("sethos-c1", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagExtra {
				return nil, false
			}
			if atk.Info.Abil != shadowPierceShotAil {
				return nil, false
			}
			return m, true
		},
	})
}

func (c *char) c2() {
	if c.Base.Cons < 2 {
		return
	}
	mElectro := make([]float64, attributes.EndStatType)
	c.AddStatMod(character.StatMod{
		Base: modifier.NewBase(c2Key, -1),
		Amount: func() ([]float64, bool) {
			stackCount := min(c.c2Stacks, 2.0)
			if stackCount == 0 {
				return nil, false
			}
			mElectro[attributes.ElectroP] = 0.15 * float64(stackCount)
			return mElectro, true
		},
	})
}

func (c *char) c2AddStack() {
	if c.Base.Cons < 2 {
		return
	}
	c.c2Stacks += 1
	c.SetTag(c2Key, min(c.c2Stacks, 2))
	c.QueueCharTask(func() {
		// tags currently aren't visible in the results UI
		// the user can still access it using .char.tags.xianyun-a1
		c.c2Stacks -= 1
		c.SetTag(c2Key, min(c.c2Stacks, 2))
	}, c2Dur)
}
