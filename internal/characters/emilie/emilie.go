package emilie

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/model"
)

func init() {
	core.RegisterCharFunc(keys.Emilie, NewChar)
}

type char struct {
	*tmpl.Character
	lumidouceSrc int
	lumidoucePos geometry.Point
	lumidouceLvl int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.SkillCon = 5
	c.BurstCon = 3
	c.HasArkhe = true

	c.lumidouceSrc = -1

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a1()
	c.a4()

	c.c1()
	c.c2()
	return nil
}
func (c *char) AnimationStartDelay(k model.AnimationDelayKey) int {
	if k == model.AnimationXingqiuN0StartDelay {
		return 13
	}
	return c.Character.AnimationStartDelay(k)
}

func (c *char) getTotalAtk() float64 {
	stats, _ := c.Stats()
	return c.Base.Atk*(1+stats[attributes.ATKP]) + stats[attributes.ATK]
}
