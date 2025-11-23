package athameartis

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.AthameArtis, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	var w Weapon
	refine := p.Refine

	mCD := make([]float64, attributes.EndStatType)
	mAtkp := make([]float64, attributes.EndStatType)

	cd := 0.12 + 0.04*float64(refine)
	atkp := 0.15 + 0.05*float64(refine)
	atkpTeam := 0.12 + 0.04*float64(refine)

	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("athame-artis-burst-cdmg", -1),
		Amount: func(atk *info.AttackEvent, t info.Target) ([]float64, bool) {
			if atk.Info.AttackTag == attacks.AttackTagElementalBurst {
				mCD[attributes.CD] = cd * getBonus(c)
				return mCD, true
			}

			return nil, false
		},
	})

	c.Events.Subscribe(event.OnEnemyHit, func(args ...any) bool {
		// If attack does not belong to the equipped character then ignore
		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return false
		}

		// If this is not a burst then ignore
		if atk.Info.AttackTag != attacks.AttackTagElementalBurst {
			return false
		}

		for _, chars := range c.Player.Chars() {
			buff := atkpTeam
			if chars.Index() == char.Index() {
				buff = atkp
			}
			chars.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag("athame-artis-atkp", 3*60),
				AffectedStat: attributes.ATKP,
				Amount: func() ([]float64, bool) {
					mAtkp[attributes.ATKP] = buff * getBonus(c)
					return mAtkp, true
				},
			})
		}

		return false
	}, fmt.Sprintf("athame-artis-hook-%v", char.Base.Key.String()))
	return &w, nil
}

func getBonus(c *core.Core) float64 {
	if getMagicCount(c) < 2 {
		return 1.0
	}
	return 1.75
}

func getMagicCount(c *core.Core) int {
	count := 0
	for _, c := range c.Player.Chars() {
		if c.IsMagic {
			count += 1
		}
	}
	return count
}
