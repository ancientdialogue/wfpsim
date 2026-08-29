package vesna

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c1Key = "vesna-c1"
	c2Key = "vesna-c2"
	c4Key = "vesna-c4"
	c6Key = "vesna-c6"

	c6TransposeHitmark   = 30
	c6SpiritBladeHitmark = 45
	c6PinionHitmark      = 40
)

var c6Frames []int

func init() {
	c6Frames = frames.InitAbilSlice(60)
}

func (c *char) c1Init() {
	if c.Base.Cons < 1 {
		return
	}

	c.AddReactBonusMod(character.ReactBonusMod{
		Base: modifier.NewBase(c1Key, -1),
		Amount: func(ai info.AttackInfo) float64 {
			switch ai.AttackTag {
			case attacks.AttackTagDirectStellarSwirl, attacks.AttackTagReactionStellarSwirl:
			default:
				return 0
			}

			if !c.StatusIsActive(skillKey) {
				return 0
			}
			return 0.2
		},
	})
}

func (c *char) c1UseSkillStack() bool {
	if c.Base.Cons < 1 {
		return true
	}

	if c.skillsMaxLvlUsed != 0 {
		return true
	}

	if c.skillLvl != 3 {
		return true
	}

	return false
}

func (c *char) c1ExtraMaxLvlSkills() int {
	if c.Base.Cons < 1 {
		return 0
	}
	return 1
}

func (c *char) c2Init() {
	if c.Base.Cons < 2 {
		return
	}

	if c.Base.Ascension < 1 {
		return
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = 0.6

	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase(c2Key, -1),
		AffectedStat: attributes.ATKP,
		Amount: func() []float64 {
			if c.a1StackCount() < c.a1Stacks.Len() {
				return nil
			}
			return m
		},
	})
}

// must be called after c.a1OnSkill()
func (c *char) c2OnSkill() {
	if c.Base.Cons < 2 {
		return
	}

	if c.Base.Ascension < 1 {
		return
	}

	for range c.a1Stacks.Len() {
		c.a1AddStacks()
	}
}

func (c *char) c4a4Mult() float64 {
	if c.Base.Cons < 4 {
		return 1.0
	}

	return 3.0
}

func (c *char) c6Init() {
	if c.Base.Cons < 6 {
		return
	}

	hook := func(aeInd int, tag attacks.AttackTag) func(args ...any) {
		return func(args ...any) {
			atk := args[aeInd].(*info.AttackEvent)
			if atk.Info.AttackTag != tag {
				return
			}

			if atk.Info.ActorIndex != c.Index() {
				return
			}

			if c.Core.Flags.LogDebug {
				c.Core.Log.NewEvent("Adding vesna c6 stellar swirl elevation", glog.LogCharacterEvent, c.Index()).Write("amt", 0.2)
			}

			atk.Info.Elevation += 0.2
		}
	}

	c.Core.Events.Subscribe(event.OnApplyAttack, hook(0, attacks.AttackTagDirectStellarSwirl), c6Key+"-direct")
	c.Core.Events.Subscribe(event.OnLunarReactionAttack, hook(1, attacks.AttackTagReactionStellarSwirl), c6Key+"-reaction")
}

func (c *char) c6OnMaxLvlSkill() {
	if c.Base.Cons < 6 {
		return
	}

	c.AddStatus(c6Key, 5*60, true)
}

func (c *char) c6Attack() (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Windborne Blade: Transpose (C6)",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagElementalArt,
		ICDGroup:   attacks.ICDGroupVesnaSkill,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       1.5,
	}
	ap := combat.NewCircleHitOnTargetFanAngle(c.Core.Combat.Player(), nil, 5, 60)
	c.Core.QueueAttack(ai, ap, c6TransposeHitmark, c6TransposeHitmark, c.particleCB)

	ai.Abil = "Windborne Blade: Transpose Spirit Blade (C6)"
	ai.Mult = 2

	if c.isRadianceSSw() {
		ai.Abil += stellarSwirlText
		ai.AttackTag = attacks.AttackTagDirectStellarSwirl
		ai.IgnoreDefPercent = 1
		ai.Durability = 0
	}
	c.Core.QueueAttack(ai, ap, c6SpiritBladeHitmark, c6SpiritBladeHitmark, c.particleCB)

	if c.StatusIsActive(skillKey) {
		c.pinionAttack(c6PinionHitmark)
	}

	return action.Info{
		Frames:          func(next action.Action) int { return c6Frames[next] },
		AnimationLength: c6Frames[action.InvalidAction],
		CanQueueAfter:   c6Frames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}, nil
}
