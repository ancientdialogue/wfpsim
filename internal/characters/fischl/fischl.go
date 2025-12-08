package fischl

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/internal/template/minazuki"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Fischl, NewChar)
}

type char struct {
	*tmpl.Character
	// field use for calculating oz damage
	ozPos           info.Point
	ozSnapshot      info.AttackEvent
	ozSource        int  // keep tracks of source of oz aka resets
	ozActive        bool // purely used for gscl conditional purposes
	ozTickSrc       int  // used for oz recast attacks
	ozTravel        int
	burstOzSpawnSrc int // prevent double oz spawn from burst
	c6Watcher       *minazuki.Watcher
}

func NewChar(s *core.Core, w *character.CharWrapper, p info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.SkillCon = 3
	c.BurstCon = 5

	c.ozSource = -1
	c.ozActive = false
	c.ozTickSrc = -1

	c.ozTravel = 10
	travel, ok := p.Params["oz_travel"]
	if ok {
		c.ozTravel = travel
	}

	hexerei, ok := p.Params["hexerei"]
	if !ok {
		hexerei = 1
	}
	c.IsHexerei = hexerei > 0

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a4()
	c.hexereiInit()
	return c.c6Init()
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "oz":
		return c.ozActive, nil
	case "oz-source":
		return c.ozSource, nil
	case "oz-duration":
		return c.StatusDuration(ozActiveKey), nil
	default:
		return c.Character.Condition(fields)
	}
}

func (c *char) ActionReady(a action.Action, p map[string]int) (bool, action.Failure) {
	// check if it is possible to recast oz
	if a == action.ActionSkill && p["recast"] != 0 && c.ozActive {
		return !c.StatusIsActive(skillRecastCDKey), action.SkillCD
	}
	// check if cast skill with oz on-field
	if a == action.ActionSkill && c.ozActive {
		return false, action.NoFailure
	}
	return c.Character.ActionReady(a, p)
}

func (c *char) AnimationStartDelay(k info.AnimationDelayKey) int {
	if k == info.AnimationXingqiuN0StartDelay {
		return 9
	}
	return c.Character.AnimationStartDelay(k)
}
