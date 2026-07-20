package scarletproof

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterSetFunc(keys.ScarletProof, NewSet)
}

type Set struct {
	stacks int
	Index  int
	Count  int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) GetCount() int    { return s.Count }
func (s *Set) Init() error      { return nil }

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{Count: count}
	s.stacks = 0

	if count < 2 {
		return &s, nil
	}
	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = 0.18
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("scarletproof-2pc", -1),
		AffectedStat: attributes.ATKP,
		Amount: func() []float64 {
			return m
		},
	})

	if count < 4 {
		return &s, nil
	}

	m2 := make([]float64, attributes.EndStatType)
	m2[attributes.CR] = 0.16

	gainBuff := func(char *character.CharWrapper) {
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("scarletproof-4pc-cr", 10*60),
			AffectedStat: attributes.CR,
			Amount: func() []float64 {
				return m2
			},
		})

		char.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBaseWithHitlag("scarletproof-4pc-react", 10*60),
			Amount: func(ai info.AttackInfo) float64 {
				switch ai.AttackTag {
				case attacks.AttackTagDirectStellarSwirl,
					attacks.AttackTagReactionStellarSwirl:
					return 0.4
				}
				return 0
			},
		})
	}
	hook := func(args ...any) {
		atk := args[1].(*info.AttackEvent)

		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}

		if atk.Info.ActorIndex != char.Index() {
			return
		}

		// ignore if character not on field
		if c.Player.Active() != char.Index() {
			return
		}

		gainBuff(char)
	}

	c.Events.Subscribe(event.OnStellarSwirl, hook, "scarletproof-4pc-"+char.Base.Key.String())

	// hookDmg := func(args ...any) {
	// 	atk := args[1].(*info.AttackEvent)

	// 	if _, ok := args[0].(*enemy.Enemy); !ok {
	// 		return
	// 	}

	// 	if atk.Info.ActorIndex != char.Index() {
	// 		return
	// 	}

	// 	// ignore if character not on field
	// 	if c.Player.Active() != char.Index() {
	// 		return
	// 	}

	// 	switch atk.Info.AttackTag {
	// 	case attacks.AttackTagDirectStellarSwirl:
	// 	case attacks.AttackTagReactionStellarSwirl:
	// 	default:
	// 		return
	// 	}

	// 	gainBuff(char)
	// }
	// c.Events.Subscribe(event.OnEnemyDamage, hookDmg, "scarletproof-4pc-"+char.Base.Key.String())
	return &s, nil
}
