package core

import (
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
)

type BaseStatsKey struct {
	Race  proto.Race
	Class proto.Class
	Level int32
}

var BaseStats = map[BaseStatsKey]stats.Stats{}

// To calculate base stats, get a naked level 70 of the race/class you want, ideally without any talents to mess up base stats.
//  Basic stats are as-shown (str/agi/stm/int/spirit)

// Base Spell Crit is calculated by
//   1. Take as-shown value (troll shaman have 3.5%)
//   2. Calculate the bonus from int (for troll shaman that would be 104/78.1=1.331% crit)
//   3. Subtract as-shown from int bouns (3.5-1.331=2.169)
//   4. 2.169*22.08 (rating per crit percent) = 47.89 crit rating.

// Base mana can be looked up here: https://wowwiki-archive.fandom.com/wiki/Base_mana

// These are also scattered in various dbc/casc files,
// `octbasempbyclass.txt`, `combatratings.txt`, `chancetospellcritbase.txt`, etc.

var RaceOffsets = map[proto.Race]stats.Stats{
	proto.Race_RaceUnknown: stats.Stats{},
	proto.Race_RaceHuman:   stats.Stats{},
	proto.Race_RaceOrc: {
		stats.Agility:   -3,
		stats.Strength:  3,
		stats.Intellect: -3,
		stats.Spirit:    2,
		stats.Stamina:   1,
	},
	proto.Race_RaceDwarf: {
		stats.Agility:   -4,
		stats.Strength:  5,
		stats.Intellect: -1,
		stats.Spirit:    -1,
		stats.Stamina:   1,
	},
	proto.Race_RaceNightElf: {
		stats.Agility:   4,
		stats.Strength:  -4,
		stats.Intellect: 0,
		stats.Spirit:    0,
		stats.Stamina:   0,
	},
	proto.Race_RaceUndead: {
		stats.Agility:   -2,
		stats.Strength:  -1,
		stats.Intellect: -2,
		stats.Spirit:    5,
		stats.Stamina:   0,
	},
	proto.Race_RaceTauren: {
		stats.Agility:   -4,
		stats.Strength:  5,
		stats.Intellect: -4,
		stats.Spirit:    2,
		stats.Stamina:   1,
	},
	proto.Race_RaceGnome: {
		stats.Agility:   2,
		stats.Strength:  -5,
		stats.Intellect: 3,
		stats.Spirit:    0,
		stats.Stamina:   0,
	},
	proto.Race_RaceTroll: {
		stats.Agility:   2,
		stats.Strength:  1,
		stats.Intellect: -4,
		stats.Spirit:    1,
		stats.Stamina:   0,
	},
	proto.Race_RaceBloodElf: {
		stats.Agility:   2,
		stats.Strength:  -3,
		stats.Intellect: 3,
		stats.Spirit:    -2,
		stats.Stamina:   0,
	},
	proto.Race_RaceDraenei: {
		stats.Agility:   -3,
		stats.Strength:  1,
		stats.Intellect: 0,
		stats.Spirit:    2,
		stats.Stamina:   0,
	},
}

var ApBonus = map[proto.Class]stats.Stats{
	proto.Class_ClassUnknown:     {},
	proto.Class_ClassWarrior:     {stats.AttackPower: -20},
	proto.Class_ClassPaladin:     {stats.AttackPower: -20},
	proto.Class_ClassHunter:      {stats.AttackPower: -20, stats.RangedAttackPower: -10},
	proto.Class_ClassRogue:       {stats.AttackPower: -20},
	proto.Class_ClassPriest:      {},
	proto.Class_ClassDeathknight: {stats.AttackPower: -20},
	proto.Class_ClassShaman:      {stats.AttackPower: -20},
	proto.Class_ClassMage:        {},
	proto.Class_ClassWarlock:     {stats.AttackPower: -10},
	proto.Class_ClassDruid:       {stats.AttackPower: -20},
}

var ApScaling = map[proto.Class]stats.Stats{
	proto.Class_ClassUnknown:     {},
	proto.Class_ClassWarrior:     {stats.AttackPower: 3.0},
	proto.Class_ClassPaladin:     {stats.AttackPower: 3.0},
	proto.Class_ClassHunter:      {stats.AttackPower: 2.0, stats.RangedAttackPower: 2.0},
	proto.Class_ClassRogue:       {stats.AttackPower: 2.0},
	proto.Class_ClassPriest:      {},
	proto.Class_ClassDeathknight: {stats.AttackPower: 3.0},
	proto.Class_ClassShaman:      {stats.AttackPower: 2.0},
	proto.Class_ClassMage:        {},
	proto.Class_ClassWarlock:     {},
	proto.Class_ClassDruid:       {},
}

func MakeBaseStats(r proto.Race, c proto.Class, l int32) stats.Stats {
	return ClassBaseStats[c][l].Add(RaceOffsets[r]).Add(BaseCrit[c][l].Multiply(CritRatingPerCritChance[l])).Add(ApBonus[c]).Add(ApScaling[c].Multiply(float64(l)))
}
