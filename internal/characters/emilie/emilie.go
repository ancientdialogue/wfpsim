package emilie

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
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
	scentCount         int
	lumidouceSrc       int
	lumidoucePos       geometry.Point
	lumidouceLvl       int
	lumidouceScents    int
	lumidouceScentLast int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 50
	c.NormalHitNum = normalHitNum
	c.SkillCon = 5
	c.BurstCon = 3
	c.HasArkhe = true

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a4()
	return nil
}
func (c *char) AnimationStartDelay(k model.AnimationDelayKey) int {
	if k == model.AnimationXingqiuN0StartDelay {
		return 13
	}
	return c.Character.AnimationStartDelay(k)
}
