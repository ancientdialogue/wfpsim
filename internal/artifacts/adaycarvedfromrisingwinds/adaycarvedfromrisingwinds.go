package adaycarvedfromrisingwinds

import (
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
	core.RegisterSetFunc(keys.ADayCarvedFromRisingWinds, NewSet)
}

type Set struct {
	char  *character.CharWrapper
	Index int
	Count int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) GetCount() int    { return s.Count }
func (s *Set) Init() error {
	return nil
}

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{
		char:  char,
		Count: count,
	}

	if count < 2 {
		return &s, nil
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = 0.18
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("risingwinds-2pc", -1),
		AffectedStat: attributes.ATKP,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})

	if count < 4 {
		return &s, nil
	}

	m2 := make([]float64, attributes.EndStatType)
	m2[attributes.ATKP] = 0.25
	if char.IsHexerei {
		m2[attributes.CR] = 0.2
	}
	c.Events.Subscribe(event.OnEnemyDamage, func(args ...any) bool {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return false
		}
		switch atk.Info.AttackTag {
		case attacks.AttackTagExtra:
		case attacks.AttackTagNormal:
		case attacks.AttackTagElementalArt:
		case attacks.AttackTagElementalArtHold:
		case attacks.AttackTagElementalBurst:
		default:
			return false
		}

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("risingwinds-4pc-buff", 6*60),
			AffectedStat: attributes.NoStat,
			Amount: func() ([]float64, bool) {
				return m2, true
			},
		})

		return false
	}, "breeze-4pc-"+char.Base.Key.String())

	return &s, nil
}
