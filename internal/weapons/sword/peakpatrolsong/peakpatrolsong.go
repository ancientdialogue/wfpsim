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

	def := 0.06 + float64(r)*0.02
	dmgP := 0.075 + float64(r)*0.025
	eleDmgBonus := 0.06 + float64(r)*0.02
	maxEleDmgBonus := 0.192 + float64(r)*0.064

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
		}

		d := make([]float64, attributes.EndStatType)
		d[attributes.DEFP] = def * float64(w.stacks)
		d[attributes.DmgP] = dmgP * float64(w.stacks)

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(stackKey, 6*60),
			AffectedStat: attributes.NoStat,
			Amount: func() ([]float64, bool) {
				return d, true
			},
		})

		m := make([]float64, attributes.EndStatType)
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
