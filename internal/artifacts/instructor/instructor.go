package instructor

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/artifact"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterSetFunc(keys.Instructor, NewSet)
}

type Set struct {
	Index int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) Init() error      { return nil }

// Implements Instructor artifact set:
// 2-Piece Bonus: Increases Elemental Mastery by 80.
// 4-Piece Bonus: Upon triggering an Elemental Reaction, increases all party members' Elemental Mastery by 120 for 8s.
func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (artifact.Set, error) {
	s := Set{}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.EM] = 80
		char.AddStatMod("instructor-2pc", -1, attributes.EM, func() ([]float64, bool) {
			return m, true
		})
	}
	if count >= 4 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.EM] = 120

		add := func(args ...interface{}) bool {
			atk := args[1].(*combat.AttackEvent)
			// Character must be on field to proc bonus
			if c.Player.Active() != char.Index {
				return false
			}
			// Source of elemental reaction must be the character with instructor
			if atk.Info.ActorIndex != char.Index {
				return false
			}

			// Add 120 EM to all characters except the one with instructor
			for i, this := range c.Player.Chars() {
				// Skip the one with instructor
				if i == char.Index {
					continue
				}

				this.AddStatMod("instructor-4pc", c.F+480, attributes.EM, func() ([]float64, bool) {
					return m, true
				})
			}
			return false
		}

		for i := event.ReactionEventStartDelim + 1; i < event.ReactionEventEndDelim; i++ {
			c.Events.Subscribe(i, add, fmt.Sprintf("instructor-4pc-%v", char.Base.Name))
		}
	}

	return &s, nil
}
