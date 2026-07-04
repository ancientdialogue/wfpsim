package sandrone

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

const (
	maxDecode          = 100
	stellarConductText = " (Stellar-Conduct)"
	stellarSwirlText   = " (Stellar-Swirl)"
)

func init() {
	core.RegisterCharFunc(keys.Sandrone, NewChar)
}

type fagioState int

const (
	stateIdle fagioState = iota
	stateDecoding
	stateOverdriveCharge
	stateOverdriveIdle
)

type radianceState int

const (
	radianceNone radianceState = iota
	radianceStellarConduct
	radianceStellarSwirl
)

type char struct {
	*tmpl.Character

	decode               float64
	currSweepingFireTick int
	caLength             int
	caToMaxDecode        bool
	caSrc                int
	decodeReduceSrc      int
	currFagio            fagioState
	a1DecreasedDecode    float64
	a1RefinedTactics     int
	c2stacks             int
	c6stacks             int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.BurstCon = 5
	c.NormalCon = 3
	c.NormalHitNum = normalHitNum

	w.Character = &c

	c.decode = 0

	return nil
}

func (c *char) Init() error {
	c.a4Init()
	c.stellarInit()
	c.c1Init()
	c.c2Init()
	c.c4Init()
	c.c6Init()
	return nil
}

func (c *char) AnimationStartDelay(k info.AnimationDelayKey) int {
	if k == info.AnimationXingqiuN0StartDelay {
		return 10
	}
	return c.Character.AnimationStartDelay(k)
}
