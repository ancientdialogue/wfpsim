package klee

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Klee, NewChar)
}

type char struct {
	*tmpl.Character
	a1CurrentStack     int
	a1MaxStack         int
	c1Chance           float64
	savedNormalCounter int
}

func NewChar(s *core.Core, w *character.CharWrapper, p info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.SkillCon = 3
	c.BurstCon = 5

	c.SetNumCharges(action.ActionSkill, 2)

	w.Character = &c

	magic, ok := p.Params["magic"]
	if !ok {
		magic = 1
	}
	c.IsMagic = magic > 0

	c.a1MaxStack = 1
	if c.IsMagic {
		c.a1MaxStack = 3
	}

	return nil
}

func (c *char) Init() error {
	c.onExitField()
	c.magicInit()
	return nil
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "spark-stacks":
		return c.a1CurrentStack, nil
	default:
		return c.Character.Condition(fields)
	}
}

// Magic: Secret Rite
// During Klee's Elemental Burst, her Normal Attack sequence does not reset.
// TODO: This causes N1Q to not reset counter, is this accurate behavior?
func (c *char) ResetNormalCounter() {
	if c.IsMagic && c.getMagicCount() > 1 && c.Core.Player.Active() == c.Index() && c.Core.Status.Duration(burstKey) > 0 {
		c.NormalCounter = c.savedNormalCounter
		return
	}
	c.Character.ResetNormalCounter()
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	if a == action.ActionCharge {
		if c.StatusIsActive(a1SparkKey) {
			return 0
		}
		return 50
	}
	return c.Character.ActionStam(a, p)
}
