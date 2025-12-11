package archaic

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterSetFunc(keys.ArchaicPetra, NewSet)
}

type Set struct {
	element attributes.Element
	Index   int
	Count   int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) GetCount() int    { return s.Count }
func (s *Set) Init() error      { return nil }

func NewSet(core *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{Count: count}

	if count < 2 {
		return &s, nil
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.GeoP] = 0.15
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("archaic-2pc", -1),
		AffectedStat: attributes.GeoP,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})

	if count < 4 {
		return &s, nil
	}

	m2 := make([]float64, attributes.EndStatType)
	core.Events.Subscribe(event.OnLunarCrystallize, func(args ...any) bool {
		if core.Player.Active() != char.Index() {
			return false
		}

		if _, ok := args[0].(*enemy.Enemy); !ok {
			return false
		}

		ae, ok := args[1].(*info.AttackEvent)
		if !ok {
			return false
		}
		if ae.Info.ActorIndex != char.Index() {
			return false
		}

		m2[attributes.HydroP] = 0.35
		// Apply mod to all characters
		for _, c := range core.Player.Chars() {
			c.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag("archaic-4pc", 10*60),
				AffectedStat: attributes.NoStat,
				Amount: func() ([]float64, bool) {
					return m2, true
				},
			})
		}
		return false
	}, fmt.Sprintf("archaic-4pc-lcr-%v", char.Base.Key.String()))

	core.Events.Subscribe(event.OnShielded, func(args ...any) bool {
		// Character that picks it up must be the petra set holder
		if core.Player.Active() != char.Index() {
			return false
		}

		// Check shield
		shd := args[0].(shield.Shield)
		if shd.Type() != shield.Crystallize {
			return false
		}
		s.element = shd.Element()

		// Activate
		// TODO: cd for proc?
		core.Log.NewEvent("archaic petra proc'd", glog.LogArtifactEvent, char.Index()).
			Write("ele", s.element)

		m2[attributes.PyroP] = 0
		m2[attributes.HydroP] = 0
		m2[attributes.CryoP] = 0
		m2[attributes.ElectroP] = 0
		m2[attributes.AnemoP] = 0
		m2[attributes.GeoP] = 0
		m2[attributes.DendroP] = 0
		m2[attributes.EleToDmgP(s.element)] = 0.35 // 35%

		// Apply mod to all characters
		for _, c := range core.Player.Chars() {
			c.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag("archaic-4pc", 10*60),
				AffectedStat: attributes.NoStat,
				Amount: func() ([]float64, bool) {
					return m2, true
				},
			})
		}

		return false
	}, fmt.Sprintf("archaic-4pc-%v", char.Base.Key.String()))

	return &s, nil
}
