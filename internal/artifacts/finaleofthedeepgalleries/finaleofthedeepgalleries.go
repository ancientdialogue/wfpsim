package finaleofthedeepgalleries

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

const galleries4pcNormalDmgICDKey = "galleries-4pc-normal-icd"
const galleries4pcBurstDmgICDKey = "galleries-4pc-burst-icd"

func init() {
	core.RegisterSetFunc(keys.FinaleOfTheDeepGalleries, NewSet)
}

type Set struct {
	core  *core.Core
	char  *character.CharWrapper
	buff  []float64
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
		core:  c,
		char:  char,
		Count: count,
	}

	// Cryo DMG Bonus +15%
	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.CryoP] = 0.15
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("galleries-2pc", -1),
			AffectedStat: attributes.CryoP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}

	// When the equipping Character has 0 Elemental Energy, Normal Attack DMG
	// is increased by 60% and Elemental Burst DMG is increased by 60%. After
	// the equipping character deals Normal Attack DMG, the aforementioned
	// Elemental Burst effect will stop applying for 6s. After the equipping
	// character deals Elemental Burst DMG, the aforementioned Normal Attack
	// effect will stop applying for 6s. This effect can trigger even if the
	// equipping character is off the field.
	if count >= 4 {
		s.buff = make([]float64, attributes.EndStatType)
		s.buff[attributes.DmgP] = 0.60

		c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
			atk := args[1].(*combat.AttackEvent)
			if atk.Info.ActorIndex != char.Index {
				return false
			}
			switch atk.Info.AttackTag {
			case attacks.AttackTagNormal:
				char.AddStatus(galleries4pcBurstDmgICDKey, 6*60, true)

			case attacks.AttackTagElementalBurst:
				char.AddStatus(galleries4pcNormalDmgICDKey, 6*60, true)
			}
			return false
		}, fmt.Sprintf("troupe-4pc-%v", char.Base.Key.String()))

		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase("galleries-4pc", -1),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				if char.Energy > 0 {
					return nil, false
				}
				switch atk.Info.AttackTag {
				case attacks.AttackTagNormal:
					if !char.StatusIsActive(galleries4pcNormalDmgICDKey) {
						return s.buff, true
					}
				case attacks.AttackTagElementalBurst:
					if !char.StatusIsActive(galleries4pcBurstDmgICDKey) {
						return s.buff, true
					}
				}
				return nil, false
			},
		})
	}

	return &s, nil
}
