package absolution

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.Absolution, NewWeapon)
}

const buffKey = "absolution-buff"

type Weapon struct {
	Index int
	c     *core.Core
	char  *character.CharWrapper
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{
		c:    c,
		char: char,
	}
	buffs := make([]float64, attributes.EndStatType)
	buffs[attributes.CD] = 0.15 + 0.05*float64(p.Refine)
	dmgP := 0.16 + 0.04*float64(p.Refine-1)
	stacks := 0

	// Increasing the value of a Bond of Life increases the DMG the equipping character deals by 12% for 6s. Max 3 stacks.
	c.Events.Subscribe(event.OnHPDebt, func(args ...interface{}) bool {
		index := args[0].(int)
		if index != char.Index {
			return false
		}
		amt := args[1].(float64)
		if amt <= 0 {
			return false
		}
		// reset stacks if nothing active
		if !char.StatModIsActive(buffKey) {
			stacks = 0
		}
		stacks++
		if stacks > 3 {
			stacks = 3
		}
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase(buffKey, 360), // 6 sec?
			AffectedStat: attributes.DmgP,
			Amount: func() ([]float64, bool) {
				buffs[attributes.DmgP] = float64(stacks) * dmgP
				return buffs, true
			},
		})

		return false
	}, fmt.Sprintf("absolution-%v", char.Base.Key))

	return w, nil
}
