package reactable

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	// this stuff should be added to a polestar field construct/gadget instead of being global
	// but doing global here to be lazy
	PolestarFieldKey        = "polestar-field"
	StellarConductStacksKey = "stellar-conduct-stacks"
	StellarConductShredKey  = PolestarFieldKey + "-phys-shred"

	polestarFieldSrcKey      = "polestar-field-src"
	polestarFieldDur         = 6 * 60
	polestarFieldStacksKey   = "polestar-field-stacks"
	polestarFieldStackICDKey = "polestar-field-stacks-icd"
	polestarFieldStackICDDur = 0.1 * 60
	polestarFieldMaxStacks   = 12
)

var stellarConductBuff = []float64{
	0.2,
	0.29,
	0.3,
	0.31,
	0.32,
	0.33,
	0.34,
	0.35,
	0.36,
	0.37,
	0.38,
	0.39,
	0.4,
}

// this is global because it enables StellarConduct for all reactables
func EnableStellarConduct(core *core.Core) {
	core.Flags.Custom[StellarConductEnableKey] = 1

	core.Events.Subscribe(event.OnElementApplied, func(args ...any) {
		if core.Status.Duration(PolestarFieldKey) <= 0 {
			return
		}

		target := args[0].(info.Target)
		if target.Type() != info.TargettableEnemy {
			return
		}

		element := args[1].(attributes.Element)

		switch element {
		case attributes.Electro:
		case attributes.Cryo:
		// TODO: does adding frozen aura count?
		default:
			return
		}

		if core.Status.Duration(polestarFieldStackICDKey) > 0 {
			return
		}

		core.Status.Add(polestarFieldStackICDKey, polestarFieldStackICDDur)

		if core.Flags.Custom[polestarFieldStacksKey] > polestarFieldMaxStacks {
			return
		}

		core.Flags.Custom[polestarFieldStacksKey] += 1
		core.Log.NewEvent("Adding polestar field stored stacks", glog.LogElementEvent, -1).Write("new_stacks", core.Flags.Custom[polestarFieldStacksKey])
	}, "stellarconduct-hook")

	mCryo := make([]float64, attributes.EndStatType)
	mElectro := make([]float64, attributes.EndStatType)
	for _, char := range core.Player.Chars() {
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase(PolestarFieldKey+"-cryo", -1),
			AffectedStat: attributes.CryoP,
			Amount: func() []float64 {
				if !char.StatusIsActive(PolestarFieldKey) {
					return nil
				}
				buffLevel := int(core.Flags.Custom[StellarConductStacksKey])
				buff := stellarConductBuff[buffLevel]
				mCryo[attributes.CryoP] = buff
				return mCryo
			},
		})

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase(PolestarFieldKey+"-electro", -1),
			AffectedStat: attributes.ElectroP,
			Amount: func() []float64 {
				if !char.StatusIsActive(PolestarFieldKey) {
					return nil
				}
				buffLevel := int(core.Flags.Custom[StellarConductStacksKey])
				buff := stellarConductBuff[buffLevel]
				mElectro[attributes.ElectroP] = buff
				return mElectro
			},
		})
	}
}

func (r *Reactable) TryStellarConduct(a *info.AttackEvent) bool {
	if a.Info.Durability < info.ZeroDur {
		return false
	}
	// this is for non frozen one
	if r.GetAuraDurability(info.ReactionModKeyFrozen) >= info.ZeroDur {
		return false
	}
	var consumed info.Durability
	switch a.Info.Element {
	case attributes.Electro:
		if r.GetAuraDurability(info.ReactionModKeyCryo) < info.ZeroDur {
			return false
		}
		consumed = r.reduce(attributes.Cryo, a.Info.Durability, 1)
	case attributes.Cryo:
		// could be ec potentially
		if r.GetAuraDurability(info.ReactionModKeyElectro) < info.ZeroDur {
			return false
		}
		consumed = r.reduce(attributes.Electro, a.Info.Durability, 1)
	default:
		return false
	}

	a.Info.Durability -= consumed
	a.Info.Durability = max(a.Info.Durability, 0)
	a.Reacted = true
	r.queueStellarConduct(a)
	return true
}

func (r *Reactable) TryFrozenStellarConduct(a *info.AttackEvent) bool {
	if a.Info.Durability < info.ZeroDur {
		return false
	}
	// this is for frozen
	if r.GetAuraDurability(info.ReactionModKeyFrozen) < info.ZeroDur {
		return false
	}
	switch a.Info.Element {
	case attributes.Electro:
		// TODO: the assumption here is we first reduce cryo, and if there's any
		// src durability left, we reduce frozen. note that it's still only one
		// superconduct reaction
		a.Info.Durability -= r.reduce(attributes.Cryo, a.Info.Durability, 1)
		r.reduce(attributes.Frozen, a.Info.Durability, 1)
		a.Info.Durability = 0
		a.Reacted = true
	default:
		return false
	}

	r.queueStellarConduct(a)

	return false
}

func (r *Reactable) queueStellarConduct(a *info.AttackEvent) {
	r.core.Events.Emit(event.OnSuperconduct, r.self, a)
	// setup polestar field

	newPolestarField := r.core.Status.Duration(PolestarFieldKey) <= 0
	r.core.Status.Add(PolestarFieldKey, polestarFieldDur)
	if newPolestarField {
		r.startPolestarField()
	}
}

func (r *Reactable) startPolestarField() {
	r.core.Flags.Custom[polestarFieldSrcKey] = float64(r.core.F)
	r.polestarFieldTicker(r.core.F)
	r.core.Flags.Custom[polestarFieldStacksKey] = 0
}

func (r *Reactable) polestarFieldTicker(src int) {
	if r.core.Flags.Custom[polestarFieldSrcKey] != float64(src) {
		return
	}

	if r.core.Status.Duration(PolestarFieldKey) <= 0 {
		return
	}

	oldStacks := r.core.Flags.Custom[StellarConductStacksKey]
	newStacks := r.core.Flags.Custom[polestarFieldStacksKey]
	r.core.Flags.Custom[StellarConductStacksKey] = newStacks
	r.core.Flags.Custom[polestarFieldStacksKey] = 0

	for _, char := range r.core.Player.Chars() {
		char.AddStatus(PolestarFieldKey, 4*60, false)
	}

	r.core.Log.NewEvent("Updating polestar field buff stacks", glog.LogElementEvent, -1).Write("old_stacks", oldStacks).Write("new_stacks", newStacks)

	for _, e := range r.core.Combat.Enemies() {
		e, ok := e.(info.Enemy)
		if !ok {
			panic("core.Combat.Enemies() contains enemies that don't implement info.Enemy")
		}
		e.AddResistMod(info.ResistMod{
			Base:  modifier.NewBaseWithHitlag(StellarConductShredKey, 4*60),
			Ele:   attributes.Physical,
			Value: -0.40,
		})
	}

	r.core.Tasks.Add(func() {
		r.polestarFieldTicker(src)
	}, 4*60)
}
