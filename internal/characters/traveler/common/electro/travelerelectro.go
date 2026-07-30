package electro

import (
	"github.com/genshinsim/gcsim/internal/characters/traveler/common"
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

type Traveler struct {
	*tmpl.Character
	abundanceAmulets      int
	burstC6Hits           int
	burstC6WillGiveEnergy bool
	burstSnap             info.Snapshot
	burstAtk              *info.AttackEvent
	burstSrc              int
	gender                int
	trueMoonBuff          bool
	trueMoonStacks        int
}

func NewTraveler(s *core.Core, w *character.CharWrapper, p info.CharacterProfile, gender int) (*Traveler, error) {
	c := Traveler{
		gender: gender,
	}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.Base.Element = attributes.Electro
	c.EnergyMax = 80
	c.BurstCon = 3
	c.SkillCon = 5
	c.NormalHitNum = normalHitNum

	common.TravelerStoryBuffs(w, p)
	trueMoonBuff, okTrueMoonBuff := p.Params["true_moon_story_buff"]

	if !okTrueMoonBuff {
		trueMoonBuff = 1
	}
	c.trueMoonBuff = trueMoonBuff > 0
	return &c, nil
}

func (c *Traveler) Init() error {
	c.trueMoonInit()
	c.burstProc()
	return nil
}

func (c *Traveler) AnimationStartDelay(k info.AnimationDelayKey) int {
	switch k {
	case info.AnimationXingqiuN0StartDelay:
		if c.gender == 0 {
			return 8
		}
		return 7
	default:
		return c.Character.AnimationStartDelay(k)
	}
}
