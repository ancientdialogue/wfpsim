package aubadeofmorningstarandmoon

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

const onFieldGraceKey = "morningstar-4pc-onfield"

func init() {
	core.RegisterSetFunc(keys.AubadeOfMorningstarAndMoon, NewSet)
}

type Set struct {
	lastSwap int
	core     *core.Core
	char     *character.CharWrapper
	buff     []float64
	Index    int
	Count    int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) GetCount() int    { return s.Count }
func (s *Set) Init() error {
	if s.buff == nil { // no 4pc
		return nil
	}

	return nil
}

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{
		core:     c,
		char:     char,
		lastSwap: -1,
		Count:    count,
	}

	if count < 2 {
		return &s, nil
	}
	// Increases Elemental Mastery by 80
	m := make([]float64, attributes.EndStatType)
	m[attributes.EM] = 80
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("morningstar-2pc", -1),
		AffectedStat: attributes.EM,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})

	// When the equipping character is off-field, Lunar Reaction DMG is increased by 20%.
	// When the party's Moonsign Level is at least Ascendant Gleam, Lunar Reaction DMG will
	// be further increased by 40%. This effect will disappear after the equipping
	// character is active for 3s.
	if count >= 4 {
		s.buff = make([]float64, attributes.EndStatType)
		s.buff[attributes.DmgP] = 0.25

		c.Events.Subscribe(event.OnCharacterSwap, func(args ...any) bool {
			next := args[1].(int)
			if next == char.Index() {
				char.AddStatus(onFieldGraceKey, 3*60, true)
			}
			return false
		}, fmt.Sprintf("morningstar-4pc-%v", char.Base.Key.String()))

		char.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBase("morningstar-4pc", -1),
			Amount: func(ai info.AttackInfo) (float64, bool) {
				if !attacks.AttackTagIsLunar(ai.AttackTag) {
					return 0, false
				}
				if c.Player.Active() == char.Index() && !char.StatusIsActive(onFieldGraceKey) {
					return 0, false
				}

				val := 0.2
				if c.Player.GetMoonsignCount() >= 2 {
					val += 0.4
				}

				return val, false
			},
		})
	}

	return &s, nil
}
