package athameartis

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	weaponKey = "exaiphanesblade"
	icdKey    = weaponKey + "-icd"
)

func init() {
	core.RegisterWeaponFunc(keys.ExaiphanesBlade, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}

	refine := p.Refine

	switch char.Base.Key {
	case keys.Lumine:
	case keys.LumineAnemo:
	case keys.LumineCryo:
	case keys.LumineElectro:
	case keys.LumineDendro:
	case keys.LumineGeo:
	case keys.LumineHydro:
	case keys.LuminePyro:
	case keys.Aether:
	case keys.AetherAnemo:
	case keys.AetherCryo:
	case keys.AetherElectro:
	case keys.AetherDendro:
	case keys.AetherGeo:
	case keys.AetherHydro:
	case keys.AetherPyro:
	default:
		return w, nil
	}

	critDmg := 0.0
	if refine > 1 {
		critDmg = 0.06 * 7 // we assume traveler has all resonances done
	}

	energy := 0
	switch {
	case 1 <= refine && refine < 3:
		energy = 3
	case 3 <= refine:
		energy = 5
	}

	atkp := 0.12 + 0.04*float64(refine)
	if 4 <= refine {
		atkp = 0.08 * float64(refine)
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = atkp
	statMod := character.StatMod{
		Base:         modifier.NewBaseWithHitlag(weaponKey, 8*60),
		AffectedStat: attributes.ATKP,
		Amount: func() []float64 {
			return m
		},
	}

	c.Events.Subscribe(event.OnEnemyHit, func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}

		ae := args[1].(*info.AttackEvent)

		if ae.Info.ActorIndex != char.Index() {
			return
		}

		if char.StatusIsActive(icdKey) {
			return
		}

		char.AddStatus(icdKey, 5*60, true)

		char.AddStatMod(statMod)
		if energy > 0 {
			char.AddEnergy(weaponKey, float64(energy))
		}
	}, weaponKey+"-"+char.Base.Key.String())

	if critDmg == 0 {
		return w, nil
	}

	mCD := make([]float64, attributes.EndStatType)
	mCD[attributes.CD] = critDmg
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase(weaponKey+"-cd", -1),
		AffectedStat: attributes.CD,
		Amount: func() []float64 {
			return mCD
		},
	})

	return w, nil
}
