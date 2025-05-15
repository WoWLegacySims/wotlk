package spellinfo

import "time"

type SpellEffect struct {
	BasePoints   float64
	Die          float64
	LevelScaling float64
	Coefficient  float64
}

type SpellInfo struct {
	SpellID  int32
	MinLevel int32
	MaxLevel int32
	Duration int32
	Period   int32
	CastTime time.Duration
	BaseCost float64
	Effects  [3]SpellEffect
}

type Spell struct {
	SpellInfos []*SpellInfo
}

func (spell *Spell) FindMaxRank(level int32) *SpellInfo {
	var info *SpellInfo
	for _, s := range spell.SpellInfos {
		if level >= s.MinLevel {
			if info == nil || info.MinLevel < s.MinLevel {
				info = s
			}
		}
	}
	return info
}

func (spell *Spell) FindDownRank(level int32) *SpellInfo {
	var info *SpellInfo
	var downRankIndex int
	for i, s := range spell.SpellInfos {
		if level >= s.MinLevel {
			if info == nil || info.MinLevel < s.MinLevel {
				info = s
				downRankIndex = i
			}
		}
	}
	if downRankIndex > 0 {
		return spell.SpellInfos[downRankIndex-1]
	}
	return nil
}
