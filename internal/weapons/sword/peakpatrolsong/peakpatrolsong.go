package peakpatrolsong

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.PeakPatrolSong, NewWeapon)
}

type Weapon struct {
	Index  int
	stacks int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	def := 0.045 + float64(r)*0.015
	eleDmgBonus := 0.09 + float64(r)*0.03
	maxEleDmgBonus := 0.27 + float64(r)*0.09

	m := make([]float64, attributes.EndStatType)

	const stackKey = "peakpatrolsong-stack"
	const buffKey = "peakpatrolsong-buff"
	const icdKey = "peakpatrolsong-icd"

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if atk.Info.AttackTag != attacks.AttackTagNormal && atk.Info.AttackTag != attacks.AttackTagPlunge {
			return false
		}
		if char.StatusIsActive(icdKey) {
			return false
		}

		char.AddStatus(icdKey, 0.1*60, true)

		if w.stacks < 2 {
			w.stacks++
			if w.stacks > 2 {
				w.stacks = 2
			}
		}
		m[attributes.DEFP] = def * float64(w.stacks)

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(stackKey, 6*60),
			AffectedStat: attributes.DEFP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})

		if w.stacks == 2 {
			for _, c := range c.Player.Chars() {
				c.AddStatMod(character.StatMod{
					Base:         modifier.NewBaseWithHitlag(buffKey, 15*60),
					AffectedStat: attributes.DmgP,
					Amount: func() ([]float64, bool) {
						bonus := eleDmgBonus * (char.TotalDef() / 1000)
						if bonus > maxEleDmgBonus {
							bonus = maxEleDmgBonus
						}
						m[attributes.DmgP] = bonus
						return m, true
					},
				})
			}
		}

		return false
	}, fmt.Sprintf("peakpatrolsong-%v", char.Base.Key.String()))

	return w, nil
}
