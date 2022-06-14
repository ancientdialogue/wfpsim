package qiqi

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

var burstFrames []int

const burstHitmark = 112

func init() {
	burstFrames = frames.InitAbilSlice(112)
}

// Only applies burst damage. Main Talisman functions are handled in qiqi.go
func (c *char) Burst(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Fortune-Preserving Talisman",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 50,
		Mult:       burstDmg[c.TalentLvlBurst()],
	}

	c.Core.QueueAttack(ai, combat.NewDefCircHit(5, false, combat.TargettableEnemy), burstHitmark, burstHitmark)

	c.ConsumeEnergy(10)
	c.SetCDWithDelay(action.ActionBurst, 20*60, 10)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstHitmark,
		Post:            burstHitmark,
		State:           action.BurstState,
	}
}

func (c *char) talismanHealHook() {
	c.Core.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		e, ok := args[0].(*enemy.Enemy)
		atk := args[1].(*combat.AttackEvent)
		if !ok {
			return false
		}

		//do nothing if talisman expired
		if e.GetTag(talismanKey) < c.Core.F {
			return false
		}
		//do nothing if talisman still on icd
		if e.GetTag(talismanICDKey) >= c.Core.F {
			return false
		}

		healAmt := c.healDynamic(burstHealPer, burstHealFlat, c.TalentLvlBurst())
		c.Core.Player.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  atk.Info.ActorIndex,
			Message: "Fortune-Preserving Talisman",
			Src:     healAmt,
			Bonus:   c.Stat(attributes.Heal),
		})
		e.SetTag(talismanICDKey, c.Core.F+60)

		return false
	}, "talisman-heal-hook")
}

// Handles C2, A4, and skill NA/CA on hit hooks
// Additionally handles burst Talisman hook - can't be done another way since Talisman is applied before the burst damage is dealt
func (c *char) onNACAHitHook() {
	c.Core.Events.Subscribe(event.OnAttackWillLand, func(args ...interface{}) bool {
		e, ok := args[0].(*enemy.Enemy)
		atk := args[1].(*combat.AttackEvent)
		if !ok {
			return false
		}
		if atk.Info.ActorIndex != c.Index {
			return false
		}

		// Talisman is applied before the damage is dealt
		if atk.Info.Abil == "Fortune-Preserving Talisman" {
			// c.talismanExpiry[t.Index()] = c.Core.F + 15*60
			e.SetTag(talismanKey, c.Core.F+15*60)
		}

		// All of the below only occur on Qiqi NA/CA hits
		if atk.Info.AttackTag != combat.AttackTagNormal || atk.Info.AttackTag != combat.AttackTagExtra {
			return false
		}

		// A4
		// When Qiqi hits opponents with her Normal and Charged Attacks, she has a 50% chance to apply a Fortune-Preserving Talisman to them for 6s. This effect can only occur once every 30s.
		if (c.c4ICDExpiry <= c.Core.F) && (c.Core.Rand.Float64() < 0.5) {
			// Don't want to overwrite a longer burst duration talisman with a shorter duration one
			// TODO: Unclear how the interaction works if there is already a talisman on enemy
			// TODO: Being generous for now and not putting it on CD if there is a conflict
			if e.GetTag(talismanKey) < c.Core.F+360 {
				e.SetTag(talismanKey, c.Core.F+360)
				c.c4ICDExpiry = c.Core.F + 30*60
				c.Core.Log.NewEvent(
					"Qiqi A4 Adding Talisman",
					glog.LogCharacterEvent,
					c.Index,
					"target", e.Index(),
					"talisman_expiry", e.GetTag(talismanKey),
					"c4_icd_expiry", c.c4ICDExpiry,
				)
			}
		}

		// Qiqi NA/CA healing proc in skill duration
		if c.Core.Status.Duration("qiqiskill") > 0 {
			c.Core.Player.Heal(player.HealInfo{
				Caller:  c.Index,
				Target:  -1,
				Message: "Herald of Frost (Attack)",
				Src:     c.healSnapshot(&c.skillHealSnapshot, skillHealOnHitPer, skillHealOnHitFlat, c.TalentLvlSkill()),
				Bonus:   c.skillHealSnapshot.Stats[attributes.Heal],
			})
		}

		return false
	}, "qiqi-onhit-naca-hook")
}
