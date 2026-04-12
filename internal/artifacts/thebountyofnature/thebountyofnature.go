package thebountyofnature

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
	core.RegisterSetFunc(keys.TheBountyOfNature, NewSet)
}

type Set struct {
	char  *character.CharWrapper
	core  *core.Core
	Index int
	Count int
	buff  []float64
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) GetCount() int    { return s.Count }
func (s *Set) Init() error {
	return nil
}

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{
		char:  char,
		core:  c,
		Count: count,
		buff:  make([]float64, attributes.EndStatType),
	}

	if count < 2 {
		return &s, nil
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.ER] = 0.2
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("thebountyofnature-2pc", -1),
		AffectedStat: attributes.ER,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})

	if count < 4 || !char.IsHexerei {
		return &s, nil
	}

	s.buff = make([]float64, attributes.EndStatType)
	c.Events.Subscribe(event.OnSkill, func(args ...any) bool {
		if c.Player.Active() != char.Index() {
			return false
		}
		duration := 60 * 20
		for _, x := range s.core.Player.Chars() {
			x.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag("thebountyofnature-4pc", duration),
				AffectedStat: attributes.NoStat,
				Amount: func() ([]float64, bool) {
					for _, elementP := range [...]attributes.Stat{
						attributes.PyroP,
						attributes.HydroP,
						attributes.AnemoP,
						attributes.ElectroP,
						attributes.DendroP,
						attributes.CryoP,
						attributes.GeoP,
					} {
						s.buff[elementP] = 0
					}
					buffStrength := 0.2
					if s.core.Player.GetHexereiCount() >= 2 {
						buffStrength = 0.4
						s.buff[attributes.EleToDmgP(s.core.Player.ActiveChar().Base.Element)] = buffStrength
					}
					s.buff[attributes.EleToDmgP(char.Base.Element)] = buffStrength
					return s.buff, true
				},
			})
		}
		return false
	}, fmt.Sprintf("thebountyofnature-4pc-%v", char.Base.Key.String()))

	return &s, nil
}
