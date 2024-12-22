package athousandblazingsuns

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/template/nightsoul"
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

func init() {
	core.RegisterWeaponFunc(keys.AThousandBlazingSuns, NewWeapon)
}

type Weapon struct {
	Index          int
	totalExtension int
	savedDuration  int
}

const (
	scorchingBrilliance = "scorching-brilliance"
	buffIcdKey          = "thousand-blazing-suns-icd"
	extensionIcdKey     = "scorching-brilliance-extension-icd"

	cooldown      = 600 // 10s
	naCaExtension = 120 // 2s
	extensionIcd  = 60  // 1s
	maxExtension  = 360 // 6s
)

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	// Gain the "Scorching Brilliance" effect when using an Elemental Skill or Burst:
	// CRIT DMG increased by 20% and ATK increased by 28% for 6s. This effect can trigger once every 10s.
	// While a "Scorching Brilliance" instance is active, its duration is increased by 2s after Normal or Charged attacks deal Elemental DMG.
	// This effect can trigger once every second, and the max duration increase is 6s.

	// (this part is not implemented)
	// Additionally, when the equipping character is in the Nightsoul's Blessing state, "Scorching Brilliance" effects are increased by 75%,
	// and its duration will only count down when the equipping character is on the field.
	w := &Weapon{}
	r := p.Refine

	critDmg := .15 + float64(r)*.05
	atk := .21 + float64(r)*.07
	initialDur := 360 // 6s

	buff := make([]float64, attributes.EndStatType)
	buff[attributes.CD] = critDmg
	buff[attributes.ATKP] = atk

	buffNightsoul := make([]float64, attributes.EndStatType)
	buffNightsoul[attributes.CD] = critDmg * 1.75
	buffNightsoul[attributes.ATKP] = atk * 1.75
	char.AddStatMod(character.StatMod{
		Base: modifier.NewBaseWithHitlag("thousand-blazing-suns", -1),
		Amount: func() ([]float64, bool) {
			if !char.StatusIsActive(scorchingBrilliance) {
				return nil, false
			}
			if char.StatusIsActive(nightsoul.NightsoulBlessingStatus) {
				return buffNightsoul, true
			}
			return buff, true
		},
	})

	gainBuff := func(args ...interface{}) bool {
		if c.Player.Active() != char.Index {
			return false
		}
		if char.StatusIsActive(buffIcdKey) {
			return false
		}
		char.AddStatus(buffIcdKey, cooldown, true)
		char.AddStatus(scorchingBrilliance, initialDur, true)
		w.totalExtension = 0
		return false
	}
	c.Events.Subscribe(event.OnBurst, gainBuff, fmt.Sprintf("thousand-blazing-suns-%v-burst", char.Base.Key.String()))
	c.Events.Subscribe(event.OnSkill, gainBuff, fmt.Sprintf("thousand-blazing-suns-%v-skill", char.Base.Key.String()))

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if c.Player.Active() != char.Index {
			return false
		}
		if atk.Info.AttackTag != attacks.AttackTagNormal && atk.Info.AttackTag != attacks.AttackTagExtra {
			return false
		}
		if atk.Info.Element == attributes.Physical {
			return false
		}
		if char.StatusIsActive(extensionIcdKey) {
			return false
		}
		if w.totalExtension >= maxExtension {
			return false
		}

		char.AddStatus(extensionIcdKey, extensionIcd, true)
		char.ExtendStatus(scorchingBrilliance, naCaExtension)
		w.totalExtension += naCaExtension

		return false
	}, fmt.Sprintf("thousand-blazing-suns-%v-extend", char.Base.Key.String()))

	c.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		// swapping out
		prev := args[0].(int)
		if prev != char.Index {
			return false
		}
		if !char.StatusIsActive(scorchingBrilliance) {
			w.savedDuration = 0
			return false
		}
		w.savedDuration = char.StatusDuration(scorchingBrilliance)
		char.AddStatus(scorchingBrilliance, -1, false)

		return false
	}, fmt.Sprintf("thousand-blazing-suns-%v-swap-exit", char.Base.Key.String()))

	c.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		// swapping in
		next := args[1].(int)
		if next != char.Index {
			return false
		}
		if !char.StatusIsActive(scorchingBrilliance) {
			return false
		}
		char.DeleteStatus(scorchingBrilliance)
		char.AddStatus(scorchingBrilliance, w.savedDuration, true)

		return false
	}, fmt.Sprintf("thousand-blazing-suns-%v-swap-enter", char.Base.Key.String()))
	return w, nil
}
