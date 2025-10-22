package mona

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

const (
	bubbleKey = "mona-bubble"
	omenKey   = "omen-debuff"
)

func init() {
	core.RegisterCharFunc(keys.Mona, NewChar)
}

type char struct {
	*tmpl.Character
	a4Stats                []float64
	phantasmalBubbleStacks int
	omenStartingBonusDur   int
	magicOmenExtension     int
	c2icd                  int
	c2AfterBurst           bool
	c2Buff                 []float64
	c6Src                  int
	c6Stacks               int
}

func NewChar(s *core.Core, w *character.CharWrapper, p info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.BurstCon = 3
	c.SkillCon = 5

	c.c2icd = -1
	c.c6Src = -1

	w.Character = &c

	magic, ok := p.Params["magic"]
	if !ok {
		magic = 1
	}
	c.IsMagic = magic > 0

	return nil
}

func (c *char) Init() error {
	c.burstHook()
	c.burstDamageBonus()
	c.a4()

	c.magicInit()

	if c.Base.Cons >= 1 {
		c.c1()
	}
	c.c2Init()
	if c.Base.Cons >= 4 {
		c.c4()
	}
	c.c6Init()
	return nil
}
