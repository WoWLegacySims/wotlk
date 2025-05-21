package rogue

import (
	"strconv"
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/rogueinfo"
)

func (rogue *Rogue) applyPoisons() {
	rogue.applyDeadlyPoison()
	rogue.applyInstantPoison()
	rogue.applyWoundPoison()
}

func (rogue *Rogue) registerPoisonAuras() {
	if rogue.Talents.SavageCombat > 0 {
		rogue.savageCombatDebuffAuras = rogue.NewEnemyAuraArray(func(target *core.Unit, _ int32) *core.Aura {
			return core.SavageCombatAura(target, rogue.Talents.SavageCombat)
		})
	}
	if rogue.Talents.MasterPoisoner > 0 {
		rogue.masterPoisonerDebuffAuras = rogue.NewEnemyAuraArray(func(target *core.Unit, _ int32) *core.Aura {
			aura := core.MasterPoisonerDebuff(target, rogue.Talents.MasterPoisoner)
			aura.Duration = core.NeverExpires
			return aura
		})
	}
}

func (rogue *Rogue) registerDeadlyPoisonSpell() {
	dbc := rogueinfo.DeadlyPoison.GetMaxRank(rogue.Level)
	if dbc == nil {
		return
	}
	bp, _ := dbc.GetBPDie(0, rogue.Level)

	var energyMetrics *core.ResourceMetrics
	if rogue.HasSetBonus(Tier8, 2) {
		energyMetrics = rogue.NewEnergyMetrics(core.ActionID{SpellID: 64913})
	}

	rogue.DeadlyPoison = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbc.SpellID},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskWeaponProc,

		DamageMultiplier: []float64{1, 1.07, 1.14, 1.20}[rogue.Talents.VilePoisons],
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:     "DeadlyPoison",
				MaxStacks: 5,
				Duration:  time.Second * 12,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					if rogue.Talents.SavageCombat > 0 {
						rogue.savageCombatDebuffAuras.Get(aura.Unit).Activate(sim)
					}
					if rogue.Talents.MasterPoisoner > 0 {
						rogue.masterPoisonerDebuffAuras.Get(aura.Unit).Activate(sim)
					}
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					if rogue.Talents.SavageCombat > 0 {
						rogue.savageCombatDebuffAuras.Get(aura.Unit).Deactivate(sim)
					}
					if rogue.Talents.MasterPoisoner > 0 {
						rogue.masterPoisonerDebuffAuras.Get(aura.Unit).Deactivate(sim)
					}
				},
			},
			NumberOfTicks: 4,
			TickLength:    time.Second * 3,

			OnSnapshot: func(_ *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				if stacks := dot.GetStacks(); stacks > 0 {
					dot.SnapshotBaseDamage = (bp + 0.027*dot.Spell.MeleeAttackPower()) * float64(stacks)
					attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
					dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
				}
			},

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				result := dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
				if energyMetrics != nil && result.Landed() {
					rogue.AddEnergy(sim, 1, energyMetrics)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)
			if !result.Landed() {
				return
			}

			dot := spell.Dot(target)
			if !dot.IsActive() {
				dot.Apply(sim)
				dot.SetStacks(sim, 1)
				dot.TakeSnapshot(sim, false)
				return
			}

			if dot.GetStacks() < 5 {
				dot.Refresh(sim)
				dot.AddStack(sim)
				dot.TakeSnapshot(sim, false)
				return
			}

			if rogue.lastDeadlyPoisonProcMask.Matches(core.ProcMaskMeleeMH) {
				switch rogue.Options.OhImbue {
				case proto.Rogue_Options_InstantPoison:
					rogue.InstantPoison[DeadlyProc].Cast(sim, target)
				case proto.Rogue_Options_WoundPoison:
					rogue.WoundPoison[DeadlyProc].Cast(sim, target)
				}
			}
			if rogue.lastDeadlyPoisonProcMask.Matches(core.ProcMaskMeleeOH) {
				switch rogue.Options.MhImbue {
				case proto.Rogue_Options_InstantPoison:
					rogue.InstantPoison[DeadlyProc].Cast(sim, target)
				case proto.Rogue_Options_WoundPoison:
					rogue.WoundPoison[DeadlyProc].Cast(sim, target)
				}
			}
			dot.Refresh(sim)
			dot.TakeSnapshot(sim, false)
		},
	})
}

func (rogue *Rogue) procDeadlyPoison(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
	rogue.lastDeadlyPoisonProcMask = spell.ProcMask
	rogue.DeadlyPoison.Cast(sim, result.Target)
}

func (rogue *Rogue) getPoisonProcMask(imbue proto.Rogue_Options_PoisonImbue) core.ProcMask {
	var mask core.ProcMask
	if rogue.Options.MhImbue == imbue {
		mask |= core.ProcMaskMeleeMH
	}
	if rogue.Options.OhImbue == imbue {
		mask |= core.ProcMaskMeleeOH
	}
	return mask
}

func (rogue *Rogue) applyDeadlyPoison() {
	procMask := rogue.getPoisonProcMask(proto.Rogue_Options_DeadlyPoison)
	if procMask == core.ProcMaskUnknown || rogue.Level < 30 {
		return
	}

	rogue.RegisterAura(core.Aura{
		Label:    "Deadly Poison",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || !spell.ProcMask.Matches(procMask) {
				return
			}
			if sim.RandomFloat("Deadly Poison") < rogue.GetDeadlyPoisonProcChance() {
				rogue.procDeadlyPoison(sim, spell, result)
			}
		},
	})
}

func (rogue *Rogue) applyWoundPoison() {
	procMask := rogue.getPoisonProcMask(proto.Rogue_Options_WoundPoison)
	if procMask == core.ProcMaskUnknown || rogue.Level < 32 {
		return
	}

	const basePPM = 0.5 / (1.4 / 60) // ~21.43, the former 50% normalized to a 1.4 speed weapon
	rogue.woundPoisonPPMM = rogue.AutoAttacks.NewPPMManager(basePPM, procMask)

	rogue.RegisterAura(core.Aura{
		Label:    "Wound Poison",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}

			if rogue.woundPoisonPPMM.Proc(sim, spell.ProcMask, "Wound Poison") {
				rogue.WoundPoison[NormalProc].Cast(sim, result.Target)
			}
		},
	})
}

type PoisonProcSource int

const (
	NormalProc PoisonProcSource = iota
	DeadlyProc
	ShivProc
)

func (rogue *Rogue) makeInstantPoison(procSource PoisonProcSource) *core.Spell {
	dbc := rogueinfo.InstantPoison.GetMaxRank(rogue.Level)
	if dbc == nil {
		return nil
	}
	bp, die := dbc.GetBPDie(0, rogue.Level)

	isShivProc := procSource == ShivProc

	return rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbc.SpellID, Tag: int32(procSource)},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskWeaponProc,

		DamageMultiplier: []float64{1, 1.07, 1.14, 1.20}[rogue.Talents.VilePoisons],
		CritMultiplier:   rogue.SpellCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(bp, die) + 0.09*spell.MeleeAttackPower()
			if isShivProc {
				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHit)
			} else {
				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})
}

func (rogue *Rogue) makeWoundPoison(procSource PoisonProcSource, dmg float64, id int32) *core.Spell {
	isShivProc := procSource == ShivProc

	return rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: id, Tag: int32(procSource)},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskWeaponProc,

		DamageMultiplier: []float64{1, 1.07, 1.14, 1.20}[rogue.Talents.VilePoisons],
		CritMultiplier:   rogue.SpellCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dmg + 0.036*spell.MeleeAttackPower()

			var result *core.SpellResult
			if isShivProc {
				result = spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHit)
			} else {
				result = spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			}

			if result.Landed() {
				rogue.woundPoisonDebuffAuras.Get(target).Activate(sim)
			}
		},
	})
}

func (rogue *Rogue) registerWoundPoisonSpell() {
	dbc := rogueinfo.WoundPoison.GetMaxRank(rogue.Level)
	if dbc == nil {
		return
	}
	bp, _ := dbc.GetBPDie(1, rogue.Level)
	id := dbc.SpellID

	woundPoisonDebuffAura := core.Aura{
		Label:    "WoundPoison-" + strconv.Itoa(int(rogue.Index)),
		ActionID: core.ActionID{SpellID: id},
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if rogue.Talents.SavageCombat > 0 {
				rogue.savageCombatDebuffAuras.Get(aura.Unit).Activate(sim)
			}
			if rogue.Talents.MasterPoisoner > 0 {
				rogue.masterPoisonerDebuffAuras.Get(aura.Unit).Activate(sim)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if rogue.Talents.SavageCombat > 0 {
				rogue.savageCombatDebuffAuras.Get(aura.Unit).Deactivate(sim)
			}
			if rogue.Talents.MasterPoisoner > 0 {
				rogue.masterPoisonerDebuffAuras.Get(aura.Unit).Deactivate(sim)
			}
		},
	}

	rogue.woundPoisonDebuffAuras = rogue.NewEnemyAuraArray(func(target *core.Unit, _ int32) *core.Aura {
		return target.RegisterAura(woundPoisonDebuffAura)
	})
	rogue.WoundPoison = [3]*core.Spell{
		rogue.makeWoundPoison(NormalProc, bp, id),
		rogue.makeWoundPoison(DeadlyProc, bp, id),
		rogue.makeWoundPoison(ShivProc, bp, id),
	}
}

func (rogue *Rogue) registerInstantPoisonSpell() {
	rogue.InstantPoison = [3]*core.Spell{
		rogue.makeInstantPoison(NormalProc),
		rogue.makeInstantPoison(DeadlyProc),
		rogue.makeInstantPoison(ShivProc),
	}
}

func (rogue *Rogue) GetDeadlyPoisonProcChance() float64 {
	return 0.3 + 0.04*float64(rogue.Talents.ImprovedPoisons) + rogue.deadlyPoisonProcChanceBonus
}

func (rogue *Rogue) UpdateInstantPoisonPPM(bonusChance float64) {
	procMask := rogue.getPoisonProcMask(proto.Rogue_Options_InstantPoison)
	if procMask == core.ProcMaskUnknown {
		return
	}

	const basePPM = 0.2 / (1.4 / 60) // ~8.57, the former 20% normalized to a 1.4 speed weapon

	ppm := basePPM * (1 + float64(rogue.Talents.ImprovedPoisons)*0.1 + bonusChance)
	rogue.instantPoisonPPMM = rogue.AutoAttacks.NewPPMManager(ppm, procMask)
}

func (rogue *Rogue) applyInstantPoison() {
	procMask := rogue.getPoisonProcMask(proto.Rogue_Options_InstantPoison)
	if procMask == core.ProcMaskUnknown || rogue.Level < 20 {
		return
	}

	rogue.UpdateInstantPoisonPPM(0)

	rogue.RegisterAura(core.Aura{
		Label:    "Instant Poison",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}

			if rogue.instantPoisonPPMM.Proc(sim, spell.ProcMask, "Instant Poison") {
				rogue.InstantPoison[NormalProc].Cast(sim, result.Target)
			}
		},
	})
}
