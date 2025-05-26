package core

import (
	"slices"
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
)

func ApplyAlchemyBonus(stats stats.Stats, effect int32) stats.Stats {
	mod, ok := Mixology[effect]
	if !ok {
		mod = 1.3
	}
	return stats.Multiply(mod)

}

// Registers all consume-related effects to the Agent.
func applyConsumeEffects(agent Agent) {
	character := agent.GetCharacter()
	consumes := character.Consumes
	if consumes == nil {
		return
	}
	if consumes.Flask != proto.Flask_FlaskUnknown {
		flask := Flasks[consumes.Flask]
		if character.Level >= flask.Level {
			bonus := flask.Stats
			if character.HasProfession(proto.Profession_Alchemy) {
				bonus = ApplyAlchemyBonus(bonus, int32(consumes.Flask))
			}

			switch consumes.Flask {
			case proto.Flask_FlaskofBlindingLight:
				character.PseudoStats.ArcaneSpellPower += bonus[stats.SpellPower]
				character.PseudoStats.HolySpellPower += bonus[stats.SpellPower]
				character.PseudoStats.NatureSpellPower += bonus[stats.SpellPower]
			case proto.Flask_FlaskofPureDeath:
				character.PseudoStats.ShadowSpellPower += bonus[stats.SpellPower]
				character.PseudoStats.FireSpellPower += bonus[stats.SpellPower]
				character.PseudoStats.FrostSpellPower += bonus[stats.SpellPower]
			default:
				character.AddStats(bonus)
			}
		}
	} else {
		if consumes.BattleElixir != proto.BattleElixir_BattleElixirUnknown {
			elixir := BattleElixirs[consumes.BattleElixir]
			if character.Level >= elixir.Level {
				bonus := elixir.Stats
				if character.HasProfession(proto.Profession_Alchemy) {
					bonus = ApplyAlchemyBonus(bonus, int32(consumes.BattleElixir))
				}

				switch consumes.BattleElixir {
				case proto.BattleElixir_ElixirofFirepower:
					fallthrough
				case proto.BattleElixir_ElixirofMajorFirepower:
					fallthrough
				case proto.BattleElixir_ElixirofGreaterFirepower:
					character.PseudoStats.FireSpellPower += bonus[stats.SpellPower]
				case proto.BattleElixir_ElixirofFrostPower:
					fallthrough
				case proto.BattleElixir_ElixirofMajorFrostPower:
					character.PseudoStats.FrostSpellPower += bonus[stats.SpellPower]
				case proto.BattleElixir_ElixirofShadowPower:
					fallthrough
				case proto.BattleElixir_ElixirofMajorShadowPower:
					character.PseudoStats.ShadowSpellPower += bonus[stats.SpellPower]
				case proto.BattleElixir_ElixirofDemonslaying:
					if character.CurrentTarget.MobType == proto.MobType_MobTypeDemon {
						character.PseudoStats.MobTypeAttackPower = bonus[stats.AttackPower]
					}
				default:
					character.AddStats(bonus)
				}
			}

		}
		if consumes.GuardianElixir != proto.GuardianElixir_GuardianElixirUnknown {
			elixir := GuardianElixirs[consumes.GuardianElixir]
			if character.Level >= elixir.Level {
				bonus := elixir.Stats
				if character.HasProfession(proto.Profession_Alchemy) {
					bonus = ApplyAlchemyBonus(bonus, int32(consumes.GuardianElixir))
				}
				switch consumes.GuardianElixir {
				case proto.GuardianElixir_EarthenElixir:
					character.PseudoStats.BonusPhysicalDamageTaken -= 20
					character.PseudoStats.BonusSpellDamageTaken -= 20
				case proto.GuardianElixir_GiftofArthas:
					character.AddStats(bonus)
					debuffAuras := (&character.Unit).NewEnemyAuraArray(GiftOfArthasAura)

					actionID := ActionID{SpellID: 11374}
					goaProc := character.RegisterSpell(SpellConfig{
						ActionID:    actionID,
						SpellSchool: SpellSchoolNature,
						ProcMask:    ProcMaskEmpty,

						ThreatMultiplier: 1,
						FlatThreatBonus:  90,

						ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
							debuffAuras.Get(target).Activate(sim)
							spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
						},
					})

					character.RegisterAura(Aura{
						Label:    "Gift of Arthas",
						Duration: NeverExpires,
						OnReset: func(aura *Aura, sim *Simulation) {
							aura.Activate(sim)
						},
						OnSpellHitTaken: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
							if result.Landed() &&
								spell.SpellSchool == SpellSchoolPhysical &&
								sim.RandomFloat("Gift of Arthas") < 0.3 {
								goaProc.Cast(sim, spell.Unit)
							}
						},
					})
				default:
					character.AddStats(bonus)
				}
			}
		}
	}

	if consumes.Food != proto.Food_FoodUnknown {
		food := Foods[consumes.Food]
		if character.Level >= food.Level {
			character.AddStats(food.Stats)
		}
	}
	if character.Class != proto.Class_ClassRogue && character.Class != proto.Class_ClassShaman && character.Class != proto.Class_ClassWarlock {
		ApplyImbue(character, consumes.MhImbue, true)
		ApplyImbue(character, consumes.OhImbue, false)
	}

	registerPotionCD(agent, consumes)
	registerConjuredCD(agent, consumes)
	registerExplosivesCD(agent, consumes)
}

func ApplyImbue(character *Character, imbue proto.WeaponImbue, isMH bool) {
	if imbue == proto.WeaponImbue_ImbueUnknown && ((isMH && character.GetMHWeapon() == nil) || (!isMH && character.GetOHWeapon() == nil)) {
		return
	}
	bonusDamage := 0.0

	switch imbue {
	case proto.WeaponImbue_BrilliantWizardOil:
		character.AddStats(stats.Stats{
			stats.SpellCrit:  14,
			stats.MeleeCrit:  14,
			stats.SpellPower: 36,
		})
	case proto.WeaponImbue_BrilliantManaOil:
		character.AddStats(stats.Stats{
			stats.MP5:        12,
			stats.SpellPower: 13,
		})
	case proto.WeaponImbue_SuperiorWizardOil:
		character.AddStat(stats.SpellPower, 42)
	case proto.WeaponImbue_SuperiorManaOil:
		character.AddStat(stats.MP5, 14)
	case proto.WeaponImbue_WizardOil:
		character.AddStat(stats.SpellPower, 24)
	case proto.WeaponImbue_LesserManaOil:
		character.AddStat(stats.MP5, 8)
	case proto.WeaponImbue_LesserWizardOil:
		character.AddStat(stats.SpellPower, 16)
	case proto.WeaponImbue_MinorManaOil:
		character.AddStat(stats.MP5, 4)
	case proto.WeaponImbue_MinorWizardOil:
		character.AddStat(stats.SpellPower, 8)
	case proto.WeaponImbue_AdamantiteSharpeningStone:
		character.AddStats(stats.Stats{stats.MeleeCrit: 14, stats.SpellCrit: 14})
		bonusDamage = 12
	case proto.WeaponImbue_AdamantiteWeightStone:
		character.AddStats(stats.Stats{stats.MeleeCrit: 14, stats.SpellCrit: 14})
		bonusDamage = 12
	case proto.WeaponImbue_FelSharpeningStone:
		bonusDamage = 12
	case proto.WeaponImbue_FelWeightstone:
		bonusDamage = 12
	case proto.WeaponImbue_ElementalSharpeningStone:
		if character.Class != proto.Class_ClassHunter {
			character.AddStat(stats.MeleeCrit, 28)
		}
	case proto.WeaponImbue_DenseSharpeningStone:
		bonusDamage = 8
	case proto.WeaponImbue_DenseWeightstone:
		bonusDamage = 8
	case proto.WeaponImbue_SolidSharpeningStone:
		bonusDamage = 6
	case proto.WeaponImbue_SolidWeightStone:
		bonusDamage = 6
	case proto.WeaponImbue_HeavySharpeningStone:
		bonusDamage = 4
	case proto.WeaponImbue_HeavyWeightStone:
		bonusDamage = 4
	case proto.WeaponImbue_CoarseSharpeningStone:
		bonusDamage = 3
	case proto.WeaponImbue_CoarseWeightStone:
		bonusDamage = 3
	case proto.WeaponImbue_RoughSharpeningStone:
		bonusDamage = 2
	case proto.WeaponImbue_RoughWeightStone:
		bonusDamage = 2
	case proto.WeaponImbue_ConsecratedWeapon:
		if character.CurrentTarget.MobType == proto.MobType_MobTypeDemon || character.CurrentTarget.MobType == proto.MobType_MobTypeUndead {
			character.PseudoStats.MobTypeAttackPower += 150
		}
	case proto.WeaponImbue_BlessedWeaponCoating:
		blessedWeaponCoating(character)
	case proto.WeaponImbue_RighteousWeaponCoating:
		righteousWeaponCoating(character)
	}
	weapon := character.AutoAttacks.MH()
	if !isMH {
		weapon = character.AutoAttacks.OH()
	}
	weapon.BaseDamageMin += bonusDamage
	weapon.BaseDamageMax += bonusDamage
}

func ApplyPetConsumeEffects(pet *Character, ownerConsumes *proto.Consumes) {
	switch ownerConsumes.PetFood {
	case proto.PetFood_PetFoodSpicedMammothTreats:
		pet.AddStats(stats.Stats{
			stats.Strength: 30,
			stats.Stamina:  30,
		})
	case proto.PetFood_PetFoodKiblersBits:
		pet.AddStats(stats.Stats{
			stats.Strength: 20,
			stats.Stamina:  20,
		})
	}

	pet.AddStat(stats.Agility, []float64{0, 5, 9, 13, 17, 20}[ownerConsumes.PetScrollOfAgility])
	pet.AddStat(stats.Strength, []float64{0, 5, 9, 13, 17, 20}[ownerConsumes.PetScrollOfStrength])
}

var PotionAuraTag = "Potion"

func righteousWeaponCoating(character *Character) {
	procAura := character.NewTemporaryStatsAura("Righteous Weapon Coating Proc", ActionID{ItemID: 34539}, stats.Stats{stats.AttackPower: 300, stats.RangedAttackPower: 300}, time.Second*10)

	MakeProcTriggerAura(&character.Unit, ProcTrigger{
		Name:       "Righteous Weapon Coating",
		ActionID:   ActionID{ItemID: 34539},
		Callback:   CallbackOnSpellHitDealt,
		ProcMask:   ProcMaskMeleeOrRanged,
		ICD:        time.Second * 45,
		Outcome:    OutcomeLanded,
		ProcChance: 1,
		Handler: func(sim *Simulation, spell *Spell, result *SpellResult) {
			procAura.Activate(sim)
		},
	})
}

func blessedWeaponCoating(character *Character) {
	metrics := character.NewManaMetrics(ActionID{ItemID: 34538})
	MakeProcTriggerAura(&character.Unit, ProcTrigger{
		Name:     "Blessed Weapon Coating",
		ActionID: ActionID{ItemID: 34538},
		Callback: CallbackOnCastComplete,
		PPM:      1,
		ICD:      time.Second * 45,
		Handler: func(sim *Simulation, spell *Spell, result *SpellResult) {
			character.AddMana(sim, 165, metrics)
		},
	})
}

func registerPotionCD(agent Agent, consumes *proto.Consumes) {
	character := agent.GetCharacter()
	defaultPotion := consumes.DefaultPotion
	startingPotion := consumes.PrepopPotion

	potionCD := character.NewTimer()
	if character.Spec == proto.Spec_SpecBalanceDruid {
		// Create both pots spells so they will be selectable in APL UI regardless of settings.
		speedMCD := makePotionActivation(proto.Potions_PotionOfSpeed, character, potionCD)
		wildMagicMCD := makePotionActivation(proto.Potions_PotionOfWildMagic, character, potionCD)
		speedMCD.Spell.Flags |= SpellFlagAPL | SpellFlagMCD
		wildMagicMCD.Spell.Flags |= SpellFlagAPL | SpellFlagMCD
	}

	if defaultPotion == proto.Potions_UnknownPotion && startingPotion == proto.Potions_UnknownPotion {
		return
	}

	startingMCD := makePotionActivation(startingPotion, character, potionCD)
	if startingMCD.Spell != nil {
		startingMCD.Spell.Flags |= SpellFlagPrepullPotion
	}

	var defaultMCD MajorCooldown
	if defaultPotion == startingPotion {
		defaultMCD = startingMCD
	} else {
		defaultMCD = makePotionActivation(defaultPotion, character, potionCD)
	}
	if defaultMCD.Spell != nil {
		defaultMCD.Spell.Flags |= SpellFlagCombatPotion
		character.AddMajorCooldown(defaultMCD)
	}
}

var AlchStoneItemIDs = []int32{44322, 44323, 44324}

func (character *Character) HasAlchStone() bool {
	alchStoneEquipped := false
	for _, itemID := range AlchStoneItemIDs {
		alchStoneEquipped = alchStoneEquipped || character.HasTrinketEquipped(itemID)
	}
	return character.HasProfession(proto.Profession_Alchemy) && alchStoneEquipped
}

func makePotionActivation(potionType proto.Potions, character *Character, potionCD *Timer) MajorCooldown {
	mcd := makePotionActivationInternal(potionType, character, potionCD)
	if mcd.Spell != nil {
		// Mark as 'Encounter Only' so that users are forced to select the generic Potion
		// placeholder action instead of specific potion spells, in APL prepull. This
		// prevents a mismatch between Consumes and Rotation settings.
		mcd.Spell.Flags |= SpellFlagEncounterOnly | SpellFlagPotion
		oldApplyEffects := mcd.Spell.ApplyEffects
		mcd.Spell.ApplyEffects = func(sim *Simulation, target *Unit, spell *Spell) {
			oldApplyEffects(sim, target, spell)
			if sim.CurrentTime < 0 {
				if potionType == proto.Potions_IndestructiblePotion {
					potionCD.Set(sim.CurrentTime + 2*time.Minute)
				} else {
					potionCD.Set(sim.CurrentTime + time.Minute)
				}
				character.UpdateMajorCooldowns()
			}
		}
	}
	return mcd
}

func makePotionActivationInternal(potionType proto.Potions, character *Character, potionCD *Timer) MajorCooldown {
	alchStoneEquipped := character.HasAlchStone()
	hasEngi := character.HasProfession(proto.Profession_Engineering)

	potionCast := CastConfig{
		CD: Cooldown{
			Timer:    potionCD,
			Duration: time.Minute * 60, // Infinite CD
		},
	}

	if potionType == proto.Potions_RunicHealingPotion || potionType == proto.Potions_RunicHealingInjector {
		itemId := map[proto.Potions]int32{
			proto.Potions_RunicHealingPotion:   33447,
			proto.Potions_RunicHealingInjector: 41166,
		}[potionType]
		actionID := ActionID{ItemID: itemId}
		healthMetrics := character.NewHealthMetrics(actionID)
		return MajorCooldown{
			Type: CooldownTypeSurvival,
			Spell: character.GetOrRegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast:     potionCast,
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					healthGain := sim.RollWithLabel(2700, 4500, "RunicHealingPotion")

					if alchStoneEquipped && potionType == proto.Potions_RunicHealingPotion {
						healthGain *= 1.40
					} else if hasEngi && potionType == proto.Potions_RunicHealingInjector {
						healthGain *= 1.25
					}
					character.GainHealth(sim, healthGain*character.PseudoStats.HealingTakenMultiplier, healthMetrics)
				},
			}),
		}
	} else if potionType == proto.Potions_RunicManaPotion || potionType == proto.Potions_RunicManaInjector {
		itemId := map[proto.Potions]int32{
			proto.Potions_RunicManaPotion:   33448,
			proto.Potions_RunicManaInjector: 42545,
		}[potionType]
		actionID := ActionID{ItemID: itemId}
		manaMetrics := character.NewManaMetrics(actionID)
		return MajorCooldown{
			Type: CooldownTypeMana,
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				// Only pop if we have less than the max mana provided by the potion minus 1mp5 tick.
				totalRegen := character.ManaRegenPerSecondWhileCasting() * 5
				manaGain := 4400.0
				if alchStoneEquipped && potionType == proto.Potions_RunicManaPotion {
					manaGain *= 1.4
				} else if hasEngi && potionType == proto.Potions_RunicManaInjector {
					manaGain *= 1.25
				}
				return character.MaxMana()-(character.CurrentMana()+totalRegen) >= manaGain
			},
			Spell: character.GetOrRegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast:     potionCast,
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					manaGain := sim.RollWithLabel(4200, 4400, "RunicManaPotion")
					if alchStoneEquipped && potionType == proto.Potions_RunicManaPotion {
						manaGain *= 1.4
					} else if hasEngi && potionType == proto.Potions_RunicManaInjector {
						manaGain *= 1.25
					}
					character.AddMana(sim, manaGain, manaMetrics)
				},
			}),
		}
	} else if potionType == proto.Potions_IndestructiblePotion {
		actionID := ActionID{ItemID: 40093}
		aura := character.NewTemporaryStatsAura("Indestructible Potion", actionID, stats.Stats{stats.Armor: 3500}, time.Minute*2)
		return MajorCooldown{
			Type: CooldownTypeDPS,
			Spell: character.GetOrRegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast:     potionCast,
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					aura.Activate(sim)
				},
			}),
		}
	} else if potionType == proto.Potions_PotionOfSpeed {
		actionID := ActionID{ItemID: 40211}
		aura := character.NewTemporaryStatsAura("Potion of Speed", actionID, stats.Stats{stats.MeleeHaste: 500, stats.SpellHaste: 500}, time.Second*15)
		return MajorCooldown{
			Type: CooldownTypeDPS,
			Spell: character.GetOrRegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast:     potionCast,
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					aura.Activate(sim)
				},
			}),
		}
	} else if potionType == proto.Potions_PotionOfWildMagic {
		actionID := ActionID{ItemID: 40212}
		aura := character.NewTemporaryStatsAura("Potion of Wild Magic", actionID, stats.Stats{stats.SpellPower: 200, stats.SpellCrit: 200, stats.MeleeCrit: 200}, time.Second*15)
		return MajorCooldown{
			Type: CooldownTypeDPS,
			Spell: character.GetOrRegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast:     potionCast,
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					aura.Activate(sim)
				},
			}),
		}
	} else if potionType == proto.Potions_DestructionPotion {
		actionID := ActionID{ItemID: 22839}
		aura := character.NewTemporaryStatsAura("Destruction Potion", actionID, stats.Stats{stats.SpellPower: 120, stats.SpellCrit: 2 * character.CritRatingPerCritChance}, time.Second*15)
		return MajorCooldown{
			Type: CooldownTypeDPS,
			Spell: character.GetOrRegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast:     potionCast,
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					aura.Activate(sim)
				},
			}),
		}
	} else if potionType == proto.Potions_SuperManaPotion {
		alchStoneEquipped := character.HasAlchStone()
		actionID := ActionID{ItemID: 22832}
		manaMetrics := character.NewManaMetrics(actionID)
		return MajorCooldown{
			Type: CooldownTypeMana,
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				// Only pop if we have less than the max mana provided by the potion minus 1mp5 tick.
				totalRegen := character.ManaRegenPerSecondWhileCasting() * 5
				return character.MaxMana()-(character.CurrentMana()+totalRegen) >= 3000
			},
			Spell: character.GetOrRegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast:     potionCast,
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					// Restores 1800 to 3000 mana. (2 Min Cooldown)
					manaGain := sim.RollWithLabel(1800, 3000, "super mana")
					if alchStoneEquipped {
						manaGain *= 1.4
					}
					character.AddMana(sim, manaGain, manaMetrics)
				},
			}),
		}
	} else if potionType == proto.Potions_HastePotion {
		actionID := ActionID{ItemID: 22838}
		aura := character.NewTemporaryStatsAura("Haste Potion", actionID, stats.Stats{stats.MeleeHaste: 400}, time.Second*15)
		return MajorCooldown{
			Type: CooldownTypeDPS,
			Spell: character.GetOrRegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast:     potionCast,
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					aura.Activate(sim)
				},
			}),
		}
	} else if potionType == proto.Potions_MightyRagePotion {
		actionID := ActionID{ItemID: 13442}
		aura := character.NewTemporaryStatsAura("Mighty Rage Potion", actionID, stats.Stats{stats.Strength: 60}, time.Second*15)
		rageMetrics := character.NewRageMetrics(actionID)
		return MajorCooldown{
			Type: CooldownTypeDPS,
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				if character.Class == proto.Class_ClassWarrior {
					return character.CurrentRage() < 25
				}
				return true
			},
			Spell: character.GetOrRegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast:     potionCast,
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					aura.Activate(sim)
					if character.Class == proto.Class_ClassWarrior {
						bonusRage := sim.RollWithLabel(45, 75, "Mighty Rage Potion")
						character.AddRage(sim, bonusRage, rageMetrics)
					}
				},
			}),
		}
	} else if potionType == proto.Potions_FelManaPotion {
		actionID := ActionID{ItemID: 31677}

		// Restores 3200 mana over 24 seconds.
		manaGain := 3200.0
		alchStoneEquipped := character.HasAlchStone()
		if alchStoneEquipped {
			manaGain *= 1.4
		}
		mp5 := manaGain / 24 * 5

		buffAura := character.NewTemporaryStatsAura("Fel Mana Potion", actionID, stats.Stats{stats.MP5: mp5}, time.Second*24)
		debuffAura := character.NewTemporaryStatsAura("Fel Mana Potion Debuff", ActionID{SpellID: 38927}, stats.Stats{stats.SpellPower: -25}, time.Minute*15)

		return MajorCooldown{
			Type: CooldownTypeMana,
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				// Only pop if we have low enough mana. The potion takes effect over 24
				// seconds so we can pop it a little earlier than the full value.
				return character.MaxMana()-character.CurrentMana() >= 2000
			},
			Spell: character.GetOrRegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast:     potionCast,
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					buffAura.Activate(sim)
					debuffAura.Activate(sim)
					debuffAura.Refresh(sim)
				},
			}),
		}
	} else if potionType == proto.Potions_InsaneStrengthPotion {
		actionID := ActionID{ItemID: 22828}
		aura := character.NewTemporaryStatsAura("Insane Strength Potion", actionID, stats.Stats{stats.Strength: 120, stats.Defense: -75}, time.Second*15)
		return MajorCooldown{
			Type: CooldownTypeDPS,
			Spell: character.GetOrRegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast:     potionCast,
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					aura.Activate(sim)
				},
			}),
		}
	} else if potionType == proto.Potions_IronshieldPotion {
		actionID := ActionID{ItemID: 22849}
		aura := character.NewTemporaryStatsAura("Ironshield Potion", actionID, stats.Stats{stats.Armor: 2500}, time.Minute*2)
		return MajorCooldown{
			Type: CooldownTypeDPS,
			Spell: character.GetOrRegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast:     potionCast,
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					aura.Activate(sim)
				},
			}),
		}
	} else if potionType == proto.Potions_HeroicPotion {
		actionID := ActionID{ItemID: 22837}
		aura := character.NewTemporaryStatsAura("Heroic Potion", actionID, stats.Stats{stats.Strength: 70, stats.Health: 700}, time.Second*15)
		return MajorCooldown{
			Type: CooldownTypeDPS,
			Spell: character.GetOrRegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast:     potionCast,
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					aura.Activate(sim)
				},
			}),
		}
	} else {
		return MajorCooldown{}
	}
}

var ConjuredAuraTag = "Conjured"

func registerConjuredCD(agent Agent, consumes *proto.Consumes) {
	character := agent.GetCharacter()
	conjuredType := consumes.DefaultConjured

	if conjuredType == proto.Conjured_ConjuredDarkRune {
		actionID := ActionID{ItemID: 20520}
		manaMetrics := character.NewManaMetrics(actionID)
		// damageTakenManaMetrics := character.NewManaMetrics(ActionID{SpellID: 33776})
		spell := character.RegisterSpell(SpellConfig{
			ActionID: actionID,
			Flags:    SpellFlagNoOnCastComplete,
			Cast: CastConfig{
				CD: Cooldown{
					Timer:    character.GetConjuredCD(),
					Duration: time.Minute * 15,
				},
			},
			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				// Restores 900 to 1500 mana. (2 Min Cooldown)
				manaGain := sim.RollWithLabel(900, 1500, "dark rune")
				character.AddMana(sim, manaGain, manaMetrics)

				// if character.Class == proto.Class_ClassPaladin {
				// 	// Paladins gain extra mana from self-inflicted damage
				// 	// TO-DO: It is possible for damage to be resisted or to crit
				// 	// This would affect mana returns for Paladins
				// 	manaFromDamage := manaGain * 2.0 / 3.0 * 0.1
				// 	character.AddMana(sim, manaFromDamage, damageTakenManaMetrics, false)
				// }
			},
		})
		character.AddMajorCooldown(MajorCooldown{
			Spell: spell,
			Type:  CooldownTypeMana,
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				// Only pop if we have less than the max mana provided by the potion minus 1mp5 tick.
				totalRegen := character.ManaRegenPerSecondWhileCasting() * 5
				return character.MaxMana()-(character.CurrentMana()+totalRegen) >= 1500
			},
		})
	} else if conjuredType == proto.Conjured_ConjuredFlameCap {
		actionID := ActionID{ItemID: 22788}

		flameCapProc := character.RegisterSpell(SpellConfig{
			ActionID:    actionID,
			ProcMask:    ProcMaskEmpty,
			SpellSchool: SpellSchoolFire,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
				spell.CalcAndDealDamage(sim, target, 40, spell.OutcomeMagicHitAndCrit)
			},
		})

		const procChance = 0.185
		flameCapAura := character.RegisterAura(Aura{
			Label:    "Flame Cap",
			ActionID: actionID,
			Duration: time.Minute,
			OnGain: func(aura *Aura, sim *Simulation) {
				character.PseudoStats.FireSpellPower += 80
			},
			OnExpire: func(aura *Aura, sim *Simulation) {
				character.PseudoStats.FireSpellPower -= 80
			},
			OnSpellHitDealt: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
				if !result.Landed() || !spell.ProcMask.Matches(ProcMaskMeleeOrRanged) {
					return
				}
				if sim.RandomFloat("Flame Cap Melee") > procChance {
					return
				}

				flameCapProc.Cast(sim, result.Target)
			},
		})

		spell := character.RegisterSpell(SpellConfig{
			ActionID: actionID,
			Flags:    SpellFlagNoOnCastComplete,
			Cast: CastConfig{
				CD: Cooldown{
					Timer:    character.GetConjuredCD(),
					Duration: time.Minute * 3,
				},
			},
			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				flameCapAura.Activate(sim)
			},
		})
		character.AddMajorCooldown(MajorCooldown{
			Spell: spell,
			Type:  CooldownTypeDPS,
		})
	} else if conjuredType == proto.Conjured_ConjuredHealthstone {
		actionID := ActionID{ItemID: 36892}
		healthMetrics := character.NewHealthMetrics(actionID)

		spell := character.RegisterSpell(SpellConfig{
			ActionID: actionID,
			Flags:    SpellFlagNoOnCastComplete,
			Cast: CastConfig{
				CD: Cooldown{
					Timer:    character.GetConjuredCD(),
					Duration: time.Minute * 2,
				},
			},
			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				character.GainHealth(sim, 4280*character.PseudoStats.HealingTakenMultiplier, healthMetrics)
			},
		})
		character.AddMajorCooldown(MajorCooldown{
			Spell: spell,
			Type:  CooldownTypeSurvival,
		})
	}
}

var ThermalSapperActionID = ActionID{ItemID: 42641}
var ExplosiveDecoyActionID = ActionID{ItemID: 40536}
var SaroniteBombActionID = ActionID{ItemID: 41119}
var CobaltFragBombActionID = ActionID{ItemID: 40771}

func registerExplosivesCD(agent Agent, consumes *proto.Consumes) {
	if consumes.ExplosiveBig == proto.Explosive_Big_ExplosiveBigUnknown && consumes.ExplosiveMedium == proto.Explosive_Medium_ExplosiveMediumUnknown && consumes.ExplosiveSmall == proto.Explosive_Small_ExplosiveSmallUnknown {
		return
	}
	character := agent.GetCharacter()
	sharedTimer := character.NewTimer()

	var bigExplosive *Spell
	var mediumExplosive *Spell
	var smallExplosive *Spell

	switch consumes.ExplosiveBig {
	case proto.Explosive_Big_ThermalSapper:
		bigExplosive = character.newBigExplosiveSpell(sharedTimer, 42641, 2187, 625, 2187, 625)
	case proto.Explosive_Big_SuperSapperCharge:
		bigExplosive = character.newBigExplosiveSpell(sharedTimer, 23827, 899, 501, 674, 451)
	case proto.Explosive_Big_GoblinSapperCharge:
		bigExplosive = character.newBigExplosiveSpell(sharedTimer, 10646, 449, 301, 374, 251)
	}

	switch consumes.ExplosiveMedium {
	case proto.Explosive_Medium_ExplosiveDecoy:
		mediumExplosive = character.newMediumExplosiveSpell(sharedTimer, 40536, 1439, 721)
	}

	switch consumes.ExplosiveSmall {
	case proto.Explosive_Small_ExplosiveSaroniteBomb:
		smallExplosive = character.newSmallExplosiveSpell(sharedTimer, 41119, 1149, 251)
	case proto.Explosive_Small_ExplosiveCobaltFragBomb:
		smallExplosive = character.newSmallExplosiveSpell(sharedTimer, 40771, 749, 251)
	case proto.Explosive_Small_TheBiggerOne:
		smallExplosive = character.newSmallExplosiveSpell(sharedTimer, 23826, 599, 401)
	case proto.Explosive_Small_DenseDynamite:
		smallExplosive = character.newSmallExplosiveSpell(sharedTimer, 18641, 339, 121)
	case proto.Explosive_Small_HeavyDynamite:
		smallExplosive = character.newSmallExplosiveSpell(sharedTimer, 4378, 127, 45)
	case proto.Explosive_Small_EzThroDynamiteII:
		smallExplosive = character.newSmallExplosiveSpell(sharedTimer, 18588, 212, 75)
	case proto.Explosive_Small_EzThroDynamite:
		smallExplosive = character.newSmallExplosiveSpell(sharedTimer, 6714, 50, 19)
	}

	if bigExplosive != nil {
		character.AddMajorCooldown(MajorCooldown{
			Spell:    bigExplosive,
			Type:     CooldownTypeDPS | CooldownTypeExplosive,
			Priority: CooldownPriorityLow + 30,
		})
	}

	if mediumExplosive != nil {
		character.AddMajorCooldown(MajorCooldown{
			Spell:    mediumExplosive,
			Type:     CooldownTypeDPS | CooldownTypeExplosive,
			Priority: CooldownPriorityLow + 20,
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				return sim.GetRemainingDuration() < time.Minute || smallExplosive == nil
			},
		})
	}

	if smallExplosive != nil {
		character.AddMajorCooldown(MajorCooldown{
			Spell:    smallExplosive,
			Type:     CooldownTypeDPS | CooldownTypeExplosive,
			Priority: CooldownPriorityLow + 10,
		})
	}
}

// Creates a spell object for the common explosive case.
func (character *Character) newBasicExplosiveSpellConfig(sharedTimer *Timer, actionID ActionID, school SpellSchool, basepoints float64, die float64, cooldown Cooldown, selfBasepoints float64, selfDie float64) SpellConfig {
	dealSelfDamage := selfBasepoints > 0 && selfDie > 0
	return SpellConfig{
		ActionID:    actionID,
		SpellSchool: school,
		ProcMask:    ProcMaskEmpty,

		Cast: CastConfig{
			CD: cooldown,
			SharedCD: Cooldown{
				Timer:    sharedTimer,
				Duration: TernaryDuration(slices.Contains([]int32{40536}, actionID.ItemID), time.Minute*2, time.Minute),
			},
		},

		// Explosives always have 1% resist chance, so just give them hit cap.
		BonusHit:         100,
		DamageMultiplier: 1,
		CritMultiplier:   2,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				baseDamage := sim.Roll(basepoints, die) * sim.Encounter.AOECapMultiplier()
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}

			if dealSelfDamage {
				baseDamage := sim.Roll(selfBasepoints, selfDie)
				spell.CalcAndDealDamage(sim, &character.Unit, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
		},
	}
}
func (character *Character) newBigExplosiveSpell(sharedTimer *Timer, itemID int32, basepoints float64, die float64, selfBasepoints float64, selfDie float64) *Spell {
	return character.GetOrRegisterSpell(character.newBasicExplosiveSpellConfig(sharedTimer, ActionID{ItemID: itemID}, SpellSchoolFire, basepoints, die, Cooldown{Timer: character.NewTimer(), Duration: time.Minute * 5}, selfBasepoints, selfDie))
}

func (character *Character) newMediumExplosiveSpell(sharedTimer *Timer, itemID int32, basepoints float64, die float64) *Spell {
	return character.GetOrRegisterSpell(character.newBasicExplosiveSpellConfig(sharedTimer, ActionID{ItemID: itemID}, SpellSchoolPhysical, basepoints, die, Cooldown{Timer: character.NewTimer(), Duration: time.Minute * 2}, 0, 0))
}

func (character *Character) newSmallExplosiveSpell(sharedTimer *Timer, itemID int32, basepoints float64, die float64) *Spell {
	return character.GetOrRegisterSpell(character.newBasicExplosiveSpellConfig(sharedTimer, ActionID{ItemID: itemID}, SpellSchoolFire, basepoints, die, Cooldown{Timer: character.NewTimer(), Duration: time.Minute}, 0, 0))
}
