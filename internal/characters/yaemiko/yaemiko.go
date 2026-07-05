package yaemiko

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

const (
	yaeTotemCount  = "totems"
	yaeTotemStatus = "yae_oldest_totem_expiry"
)

func init() {
	core.RegisterCharFunc(keys.YaeMiko, NewChar)
}

type char struct {
	*tmpl.Character
	skillDur               int
	kitsuneDetectionRadius float64
	kitsunes               []*kitsune
	c4buff                 []float64
	revelation             bool
	revelationEnhanceTick  bool
}

func NewChar(s *core.Core, w *character.CharWrapper, p info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 90
	c.NormalHitNum = normalHitNum
	c.BurstCon = 5
	c.SkillCon = 3

	c.skillDur = 900

	c.SetNumCharges(action.ActionSkill, 3)

	revelation, ok := p.Params["revelation"]
	if !ok {
		revelation = 1
	}
	c.revelation = revelation > 0

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a4()
	if c.Base.Cons >= 2 {
		c.kitsuneDetectionRadius = 20
	} else {
		c.kitsuneDetectionRadius = 12.5
	}
	if c.Base.Cons >= 4 {
		c.c4buff = make([]float64, attributes.EndStatType)
		c.c4buff[attributes.ElectroP] = .20
	}

	c.revelationInit()
	c.c1Init()
	c.c2Init()
	c.c6Init()
	return nil
}

func (c *char) AnimationStartDelay(k info.AnimationDelayKey) int {
	if k == info.AnimationXingqiuN0StartDelay {
		return 7
	}
	return c.Character.AnimationStartDelay(k)
}

func (c *char) isRadianceSSC() bool {
	if !c.revelation {
		return false
	}
	return c.StatusIsActive(reactable.PolestarFieldKey)
}
