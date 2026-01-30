package albedo

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Albedo, NewChar)
}

type char struct {
	*tmpl.Character
	lastConstruct int
	bloomSnapshot info.Snapshot
	// tracking skill information
	skillActive     bool
	skillArea       info.AttackPattern
	skillAttackInfo info.AttackInfo
	skillSnapshot   info.Snapshot
	// hexerei
	aureliths      [2]int
	oldestAurelith int
	hexereiBuff    []float64
	// c2 tracking
	c1Buff   []float64
	c2stacks int
}

func NewChar(s *core.Core, w *character.CharWrapper, p info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 40
	c.NormalHitNum = normalHitNum
	c.SkillCon = 3
	c.BurstCon = 5

	w.Character = &c

	hexerei, ok := p.Params["hexerei"]
	if !ok {
		hexerei = 1
	}
	c.IsHexerei = hexerei > 0

	return nil
}

func (c *char) Init() error {
	c.skillHook()
	c.a1()
	c.hexereiInit()
	c.c1Init()
	c.c6Init()
	return nil
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "elevator":
		return c.skillActive, nil
	case "c2stacks":
		return c.c2stacks, nil
	default:
		return c.Character.Condition(fields)
	}
}

func (c *char) AnimationStartDelay(k info.AnimationDelayKey) int {
	if k == info.AnimationXingqiuN0StartDelay {
		return 9
	}
	return c.Character.AnimationStartDelay(k)
}
