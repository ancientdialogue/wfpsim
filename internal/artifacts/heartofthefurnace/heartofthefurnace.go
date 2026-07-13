package heartofthefurnace

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
	core.RegisterSetFunc(keys.HeartOfTheFurnace, NewSet)
}

const stellar4pcKey = "heart-of-the-furnace-4pc"

type Set struct {
	Index int
	Count int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) GetCount() int    { return s.Count }
func (s *Set) Init() error      { return nil }

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{Count: count}

	if count < 2 {
		return &s, nil
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = 0.18
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("heart-of-the-furnace-2pc", -1),
		AffectedStat: attributes.ATKP,
		Amount: func() []float64 {
			return m
		},
	})

	if count < 4 {
		return &s, nil
	}

	m2 := make([]float64, attributes.EndStatType)
	m2[attributes.ATKP] = 0.12

	gainBuffs := func() {
		for _, otherChars := range c.Player.Chars() {
			otherChars.AddReactBonusMod(character.ReactBonusMod{
				Base: modifier.NewBaseWithHitlag(stellar4pcKey, 12*60),
				Amount: func(ai info.AttackInfo) float64 {
					switch ai.AttackTag {
					case attacks.AttackTagDirectStellarConduct,
						attacks.AttackTagDirectStellarSwirl,
						attacks.AttackTagReactionStellarSwirl:
						return 0.5
					default:
						return 0
					}
				},
			})
		}
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase(stellar4pcKey+"-atk", 12*60),
			AffectedStat: attributes.ATKP,
			Amount: func() []float64 {
				return m2
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
		gainBuffs()
	}

	hookDmg := func(args ...any) {
		atk := args[1].(*info.AttackEvent)

		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}

		if atk.Info.ActorIndex != char.Index() {
			return
		}

		switch atk.Info.AttackTag {
		case attacks.AttackTagDirectStellarConduct:
		case attacks.AttackTagDirectStellarSwirl:
		case attacks.AttackTagReactionStellarSwirl:
		default:
			return
		}

		gainBuffs()
	}

	c.Events.Subscribe(event.OnStellarConduct, hook, stellar4pcKey+"-"+char.Base.Key.String())
	c.Events.Subscribe(event.OnStellarSwirl, hook, stellar4pcKey+"-"+char.Base.Key.String())
	c.Events.Subscribe(event.OnEnemyDamage, hookDmg, stellar4pcKey+"-"+char.Base.Key.String())

	return &s, nil
}
