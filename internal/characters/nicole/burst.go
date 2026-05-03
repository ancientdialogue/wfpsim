package nicole

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

const (
	burstHitmark      = 118
	projectionHitmark = 20
	burstKey          = "silent-contemplation"
	burstICDKey       = "nicole-burst-projection-icd"
)

var burstFrames []int

func init() {
	burstFrames = frames.InitAbilSlice(113) // Q -> Jump
	burstFrames[action.ActionAttack] = 112
	burstFrames[action.ActionCharge] = 112
	burstFrames[action.ActionSkill] = 111
	burstFrames[action.ActionDash] = 112
	burstFrames[action.ActionWalk] = 112
	burstFrames[action.ActionSwap] = 110
}

func (c *char) Burst(_ map[string]int) (action.Info, error) {
	c.burstHits = 0

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Revelation: Ladder of Divine Ascent",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}

	// with delay
	c.ConsumeEnergy(8)
	c.SetCD(action.ActionBurst, 15*60)

	// initial hit
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 6.5), burstHitmark, burstHitmark)
	c.QueueCharTask(func() { c.AddStatus(burstKey, 20*60, true) }, burstHitmark+1)
	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap],
		State:           action.BurstState,
	}, nil
}

func (c *char) burstInit() {
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) bool {
		c.Core.Log.NewEvent("Nicole Burst test", glog.LogCharacterEvent, c.Index())
		t, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}

		ae := args[1].(*info.AttackEvent)

		if ae.Info.ActorIndex != c.Core.Player.Active() {
			return false
		}

		if c.burstHits > 4 {
			return false
		}

		if !c.StatusIsActive(burstKey) {
			return false
		}

		if c.StatusIsActive(burstICDKey) {
			return false
		}

		c.AddStatus(burstICDKey, 3*60, true)
		c.burstHits += 1

		char := c.Core.Player.Chars()[ae.Info.ActorIndex]

		ai := info.AttackInfo{
			ActorIndex: char.Index(),
			Abil:       "Arcane Projection",
			AttackTag:  attacks.AttackTagNone,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    char.Base.Element,
			Mult:       burstProj[c.TalentLvlBurst()],
			FlatDmg:    c.hexereiOnProjection(char),
		}

		ap := combat.NewCircleHitOnTarget(t, nil, 3)

		c.Core.QueueAttack(ai, ap, projectionHitmark, projectionHitmark)

		return false
	}, "nicole-burst-hook")
}
