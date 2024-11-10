package chasca

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/internal/template/nightsoul"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/model"
)

func init() {
	core.RegisterCharFunc(keys.Chasca, NewChar)
}

type char struct {
	*tmpl.Character
	nightsoulState           *nightsoul.State
	nightsoulSrc             int
	shadowhuntShells         []attributes.Element
	maxshadowhuntShell       int
	conversionElements       []attributes.Element
	uniqueConversionElements map[attributes.Element]bool
	a4Dmg                    float64
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.SkillCon = 3
	c.BurstCon = 5
	c.HasArkhe = false

	c.maxshadowhuntShell = 6

	c.conversionElements = make([]attributes.Element, 0)
	c.uniqueConversionElements = make(map[attributes.Element]bool)
	c.shadowhuntShells = make([]attributes.Element, 6)
	w.Character = &c

	c.nightsoulState = nightsoul.New(s, w)

	return nil
}

func (c *char) Init() error {
	c.loadTeamElement()
	c.a1Amount()
	c.a4()
	return nil
}

func (c *char) loadTeamElement() {
	for _, char := range c.Core.Player.Chars() {
		if char.Base.Key != c.Base.Key {
			switch char.Base.Element {
			case attributes.Pyro, attributes.Cryo, attributes.Hydro, attributes.Electro:
				c.conversionElements = append(c.conversionElements, char.Base.Element)
				c.uniqueConversionElements[char.Base.Element] = true
			}
		}
	}
}

func (c *char) AnimationStartDelay(k model.AnimationDelayKey) int {
	if c.nightsoulState.HasBlessing() {
		switch k {
		case model.AnimationXingqiuN0StartDelay:
			return 33
		default:
			return 30
		}
	}
	return c.Character.AnimationStartDelay(k)
}
