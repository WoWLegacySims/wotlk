package helpers

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
)

type WeaponProc struct {
	Name         string
	PPM          float64
	Chance       float64
	ICD          time.Duration
	OnlyWhiteHit bool
}

func NewWeaponProc(itemID int32, conf WeaponProc, procspellConfig core.SpellConfig) {
	core.NewItemEffect(itemID, func(a core.Agent) {
		character := a.GetCharacter()
		procMask := character.GetProcMaskForItem(itemID)

		procspellConfig.CritMultiplier = character.DefaultSpellCritMultiplier()
		procspellConfig.DamageMultiplier = 1
		procspellConfig.ThreatMultiplier = 1

		procSpell := character.RegisterSpell(procspellConfig)
		if conf.PPM == 0 && conf.Chance == 0 {
			conf.OnlyWhiteHit = true
			switch procMask {
			case core.ProcMaskMeleeMH:
				conf.Chance = character.MainHand().SwingSpeed * 1.8
				procMask = core.ProcMaskMeleeMHAuto
			case core.ProcMaskMeleeOH:
				conf.Chance = character.OffHand().SwingSpeed * 1.6
				procMask = core.ProcMaskMeleeOHAuto
			case core.ProcMaskRanged:
				conf.Chance = character.Ranged().SwingSpeed * 1.8
				procMask = core.ProcMaskRangedAuto
			}
		}

		if conf.OnlyWhiteHit {
			procMask = procMask &^ (core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskMeleeOrRangedSpecial)
		}

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       conf.Name,
			Callback:   core.CallbackOnSpellHitDealt,
			ProcMask:   procMask,
			Outcome:    core.OutcomeLanded,
			PPM:        conf.PPM,
			ICD:        conf.ICD,
			ProcChance: conf.Chance,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procSpell.Cast(sim, result.Target)
			},
		})
	})
}

type WeaponDamageProc struct {
	WeaponProc
	SpellID     int32
	SpellSchool core.SpellSchool
	BasePoints  float64
	Die         float64
}

func NewWeaponDamageProc(itemID int32, conf WeaponDamageProc) {
	var actionid core.ActionID
	if conf.SpellID > 0 {
		actionid = core.ActionID{SpellID: conf.SpellID}
	} else {
		actionid = core.ActionID{ItemID: itemID}
	}

	NewWeaponProc(itemID, conf.WeaponProc, core.SpellConfig{
		ActionID:    actionid,
		SpellSchool: conf.SpellSchool,
		ProcMask:    core.ProcMaskEmpty,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dmg := sim.Roll(conf.BasePoints, conf.Die)
			spell.CalcAndDealDamage(sim, target, dmg, spell.OutcomeMagicHitAndCrit)
		},
	})
}

type WeaponDotProc struct {
	WeaponProc
	SpellID     int32
	SpellSchool core.SpellSchool
	Ticks       int32
	Interval    time.Duration
	BasePoints  float64
	Aura        core.Aura
}

func NewWeaponDotProc(itemID int32, conf WeaponDotProc) {
	var actionid core.ActionID
	if conf.SpellID > 0 {
		actionid = core.ActionID{SpellID: conf.SpellID}
	} else {
		actionid = core.ActionID{ItemID: itemID}
	}

	NewWeaponProc(itemID, conf.WeaponProc, core.SpellConfig{
		ActionID:    actionid,
		SpellSchool: conf.SpellSchool,
		ProcMask:    core.ProcMaskEmpty,
		Dot: core.DotConfig{
			NumberOfTicks: conf.Ticks,
			TickLength:    conf.Interval,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
			Aura: conf.Aura,
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				dot := spell.Dot(target)
				dot.SnapshotBaseDamage = conf.BasePoints
				dot.Apply(sim)
			}
		},
	})
}

type WeaponDamageWithDotProc struct {
	WeaponProc
	SpellID     int32
	SpellSchool core.SpellSchool
	Ticks       int32
	Interval    time.Duration
	BasePoints  float64
	Die         float64
	DotBP       float64
}

func NewWeaponDamageWithDotProc(itemID int32, conf WeaponDamageWithDotProc) {
	var actionid core.ActionID
	if conf.SpellID > 0 {
		actionid = core.ActionID{SpellID: conf.SpellID}
	} else {
		actionid = core.ActionID{ItemID: itemID}
	}

	NewWeaponProc(itemID, conf.WeaponProc, core.SpellConfig{
		ActionID:    actionid,
		SpellSchool: conf.SpellSchool,
		ProcMask:    core.ProcMaskEmpty,
		Dot: core.DotConfig{
			NumberOfTicks: conf.Ticks,
			TickLength:    conf.Interval,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dmg := sim.Roll(conf.BasePoints, conf.Die)
			result := spell.CalcAndDealDamage(sim, target, dmg, spell.OutcomeMagicHitAndCrit)
			if result.Landed() {
				dot := spell.Dot(target)
				dot.SnapshotBaseDamage = conf.DotBP
				dot.Apply(sim)
			}
		},
	})
}

type WeaponExtraAttack struct {
	WeaponProc
}

func NewWeaponExtraAttackProc(itemID int32, conf WeaponExtraAttack) {
	core.NewItemEffect(itemID, func(agent core.Agent) {
		character := agent.GetCharacter()
		procMask := character.GetProcMaskForItem(itemID)

		var config core.SpellConfig
		switch procMask {
		case core.ProcMaskMeleeMH:
			config = *character.AutoAttacks.MHConfig()
		case core.ProcMaskMeleeOH:
			config = *character.AutoAttacks.OHConfig()
		}
		config.ActionID = core.ActionID{ItemID: itemID}
		extraAttack := character.GetOrRegisterSpell(config)

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       conf.Name,
			Callback:   core.CallbackOnSpellHitDealt,
			ProcMask:   procMask,
			Outcome:    core.OutcomeLanded,
			ICD:        conf.ICD,
			ProcChance: conf.Chance,
			PPM:        conf.PPM,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				extraAttack.Cast(sim, result.Target)
			},
		})
	})
}

type WeaponDotProcWithExtraDamage struct {
	WeaponDotProc
	ExtraBP          float64
	ExtraDie         float64
	ExtraSpellSchool core.SpellSchool
}

func NewWeaponDotProcWithExtraDamage(itemID int32, conf WeaponDotProcWithExtraDamage) {
	core.NewItemEffect(itemID, func(a core.Agent) {
		character := a.GetCharacter()
		procMask := character.GetProcMaskForItem(itemID)

		extradamage := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:         core.ActionID{ItemID: itemID},
			SpellSchool:      conf.ExtraSpellSchool,
			ProcMask:         core.ProcMaskMelee,
			CritMultiplier:   1,
			DamageMultiplier: 1,
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				dmg := sim.Roll(conf.ExtraBP, conf.ExtraDie)
				spell.CalcAndDealDamage(sim, target, dmg, spell.OutcomeMagicHitAndCrit)
			},
		})

		bolt := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:         core.ActionID{SpellID: conf.SpellID},
			SpellSchool:      conf.SpellSchool,
			ProcMask:         core.ProcMaskEmpty,
			CritMultiplier:   1,
			DamageMultiplier: 1,
			Dot: core.DotConfig{
				NumberOfTicks: conf.Ticks,
				TickLength:    conf.Interval,
				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
				},
			},
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
				if result.Landed() {
					dot := spell.Dot(target)
					dot.SnapshotBaseDamage = conf.BasePoints
					dot.Apply(sim)
				}
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     conf.Name,
			ActionID: core.ActionID{ItemID: itemID},
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: procMask,
			Outcome:  core.OutcomeLanded,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				extradamage.Cast(sim, result.Target)
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       conf.Name + " Proc",
			ActionID:   core.ActionID{SpellID: conf.SpellID},
			Callback:   core.CallbackOnSpellHitDealt,
			ProcMask:   procMask,
			Outcome:    core.OutcomeLanded,
			PPM:        conf.PPM,
			ProcChance: conf.Chance,
			ICD:        conf.ICD,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				bolt.Cast(sim, result.Target)
			},
		})
	})
}

type WeaponDamageProcWithExtraDamage struct {
	WeaponDamageProc
	ExtraBP          float64
	ExtraDie         float64
	ExtraSpellSchool core.SpellSchool
}

func NewWeaponDamageProcWithExtraDamage(itemID int32, conf WeaponDamageProcWithExtraDamage) {
	core.NewItemEffect(itemID, func(a core.Agent) {
		character := a.GetCharacter()
		procMask := character.GetProcMaskForItem(itemID)

		extradamage := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:         core.ActionID{ItemID: itemID},
			SpellSchool:      conf.ExtraSpellSchool,
			ProcMask:         core.ProcMaskMelee,
			CritMultiplier:   1,
			DamageMultiplier: 1,
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				dmg := sim.Roll(conf.ExtraBP, conf.ExtraDie)
				spell.CalcAndDealDamage(sim, target, dmg, spell.OutcomeMagicHitAndCrit)
			},
		})

		bolt := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:         core.ActionID{SpellID: conf.SpellID},
			SpellSchool:      conf.SpellSchool,
			ProcMask:         core.ProcMaskEmpty,
			CritMultiplier:   1,
			DamageMultiplier: 1,
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				dmg := sim.Roll(conf.BasePoints, conf.Die)
				spell.CalcAndDealDamage(sim, target, dmg, spell.OutcomeMagicHitAndCrit)
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     conf.Name,
			ActionID: core.ActionID{ItemID: itemID},
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: procMask,
			Outcome:  core.OutcomeLanded,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				extradamage.Cast(sim, result.Target)
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       conf.Name + " Proc",
			ActionID:   core.ActionID{SpellID: conf.SpellID},
			Callback:   core.CallbackOnSpellHitDealt,
			ProcMask:   procMask,
			Outcome:    core.OutcomeLanded,
			PPM:        conf.PPM,
			ProcChance: conf.Chance,
			ICD:        conf.ICD,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				bolt.Cast(sim, result.Target)
			},
		})
	})
}
