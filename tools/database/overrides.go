package database

import (
	"regexp"

	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
)

var OtherItemIdsToFetch = []string{
	// Hallow's End Ilvl bumped rings
	"211817",
	"211844",
	"211847",
	"211850",
	"211851",
}

var ItemOverrides = []*proto.UIItem{
	{ /** Destruction Holo-gogs */ Id: 32494, ClassAllowlist: []proto.Class{proto.Class_ClassMage, proto.Class_ClassPriest, proto.Class_ClassWarlock}},
	{ /** Gadgetstorm Goggles */ Id: 32476, ClassAllowlist: []proto.Class{proto.Class_ClassShaman}},
	{ /** Magnified Moon Specs */ Id: 32480, ClassAllowlist: []proto.Class{proto.Class_ClassDruid}},
	{ /** Quad Deathblow X44 Goggles */ Id: 34353, ClassAllowlist: []proto.Class{proto.Class_ClassDruid, proto.Class_ClassRogue}},
	{ /** Hyper-Magnified Moon Specs */ Id: 35182, ClassAllowlist: []proto.Class{proto.Class_ClassDruid}},
	{ /** Lightning Etched Specs */ Id: 34355, ClassAllowlist: []proto.Class{proto.Class_ClassShaman}},
	{ /** Annihilator Holo-Gogs */ Id: 34847, ClassAllowlist: []proto.Class{proto.Class_ClassMage, proto.Class_ClassPriest, proto.Class_ClassWarlock}},

	// Balance T9 "of Conquest" Alliance set
	{Id: 48158, SetName: "Malfurion's Regalia"},
	{Id: 48159, SetName: "Malfurion's Regalia"},
	{Id: 48160, SetName: "Malfurion's Regalia"},
	{Id: 48161, SetName: "Malfurion's Regalia"},
	{Id: 48162, SetName: "Malfurion's Regalia"},

	// Deathknight T9 "of Conquest" Horde set
	{Id: 48501, SetName: "Koltira's Battlegear"},
	{Id: 48502, SetName: "Koltira's Battlegear"},
	{Id: 48503, SetName: "Koltira's Battlegear"},
	{Id: 48504, SetName: "Koltira's Battlegear"},
	{Id: 48505, SetName: "Koltira's Battlegear"},

	// Deathknight T9 "of Conquest" Tank Horde set
	{Id: 48558, SetName: "Koltira's Plate"},
	{Id: 48559, SetName: "Koltira's Plate"},
	{Id: 48560, SetName: "Koltira's Plate"},
	{Id: 48561, SetName: "Koltira's Plate"},
	{Id: 48562, SetName: "Koltira's Plate"},

	// love is in the air loot
	{Id: 51804, Expansion: proto.Expansion_ExpansionWotlk},
	{Id: 51805, Expansion: proto.Expansion_ExpansionWotlk},
	{Id: 51806, Expansion: proto.Expansion_ExpansionWotlk},
	{Id: 51807, Expansion: proto.Expansion_ExpansionWotlk},

	//headless horseman
	{Id: 49126, Expansion: proto.Expansion_ExpansionWotlk},
	{Id: 49128, Expansion: proto.Expansion_ExpansionWotlk},
	{Id: 49121, Expansion: proto.Expansion_ExpansionWotlk},
	{Id: 49123, Expansion: proto.Expansion_ExpansionWotlk},
	{Id: 49124, Expansion: proto.Expansion_ExpansionWotlk},

	//coren direbrew
	{Id: 48663, Expansion: proto.Expansion_ExpansionWotlk},
	{Id: 49074, Expansion: proto.Expansion_ExpansionWotlk},
	{Id: 49076, Expansion: proto.Expansion_ExpansionWotlk},
	{Id: 49078, Expansion: proto.Expansion_ExpansionWotlk},
	{Id: 49080, Expansion: proto.Expansion_ExpansionWotlk},
	{Id: 49116, Expansion: proto.Expansion_ExpansionWotlk},
	{Id: 49118, Expansion: proto.Expansion_ExpansionWotlk},

	// Heirloom Dwarven Handcannon, Wowhead partially glitchs out and shows us some other lvl calc for this
	{Id: 44093, Stats: stats.Stats{stats.MeleeCrit: 30, stats.SpellCrit: 30, stats.Resilience: 13, stats.AttackPower: 34}.ToFloatArray()},

	// T7+T8 10 Normal stuff
	{Id: 39139, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15956, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39140, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15956, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39141, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15956, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39146, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15956, Difficulty: 3, ZoneId: 3456}}}, {Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15932, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39188, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15956, Difficulty: 3, ZoneId: 3456}}}, {Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15932, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39189, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15956, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39190, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15956, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39191, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15956, Difficulty: 3, ZoneId: 3456}}}, {Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15932, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39192, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15956, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39193, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15956, Difficulty: 3, ZoneId: 3456}}}, {Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15932, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39194, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15932, Difficulty: 3, ZoneId: 3456}}}, {Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15953, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39195, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15953, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39196, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15953, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39197, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15953, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39198, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15953, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39199, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15953, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39200, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15953, Difficulty: 3, ZoneId: 3456}}}, {Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15932, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39215, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15953, Difficulty: 3, ZoneId: 3456}}}, {Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15932, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39216, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15953, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39217, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15953, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39221, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15952, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39224, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15952, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39225, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15952, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39226, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15952, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39228, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15952, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39229, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15952, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39230, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15952, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39231, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15952, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39232, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15952, Difficulty: 3, ZoneId: 3456}}}, {Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15932, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39233, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15952, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39234, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15954, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39235, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15954, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39236, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15954, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39237, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15954, Difficulty: 3, ZoneId: 3456}}}, {Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15932, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39239, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15954, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39240, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15954, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39241, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15954, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39242, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15954, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39243, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15954, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39244, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15954, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39245, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15936, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39246, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15936, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39247, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15936, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39249, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15936, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39250, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15936, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39251, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15936, Difficulty: 3, ZoneId: 3456}}}, {Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15932, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39252, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15936, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39255, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15936, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39256, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 16011, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39257, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 16011, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39258, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 16011, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39259, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 16011, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39260, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 16011, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39261, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 16028, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39267, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 16028, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39270, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 16028, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39271, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 16028, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39272, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 16028, Difficulty: 3, ZoneId: 3456}}}, {Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15932, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39273, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 16028, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39274, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 16028, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39275, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 16028, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39276, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15931, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39277, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15931, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39278, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15931, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39279, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15931, Difficulty: 3, ZoneId: 3456}}}, {Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15932, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39280, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15931, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39281, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15931, Difficulty: 3, ZoneId: 3456}}}, {Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15932, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39282, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15931, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39283, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15931, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39284, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15931, Difficulty: 3, ZoneId: 3456}}}, {Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15932, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39285, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15931, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39291, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15928, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39292, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15928, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39293, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15928, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39294, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15928, Difficulty: 3, ZoneId: 3456}}}, {Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15932, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39295, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15928, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39296, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 16061, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39297, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 16061, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39298, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 16061, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39299, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 16061, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39306, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 16061, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39307, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 16061, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39308, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 16061, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39309, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 16061, Difficulty: 3, ZoneId: 3456}}}, {Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15932, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39310, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 16061, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39311, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 16061, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39344, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 16060, Difficulty: 3, ZoneId: 3456}}}, {Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15932, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39345, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 16060, Difficulty: 3, ZoneId: 3456}}}, {Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15932, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39369, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 16060, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39379, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 16060, Difficulty: 3, ZoneId: 3456}}}, {Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15932, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39386, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 16060, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39388, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 16060, Difficulty: 3, ZoneId: 3456}}}, {Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15932, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39389, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 16060, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39390, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 16060, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39391, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 16060, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39392, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 16060, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39393, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 16064, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39394, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 16064, Difficulty: 3, ZoneId: 3456}}}, {Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15932, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39395, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 16064, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39396, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 16064, Difficulty: 3, ZoneId: 3456}}}, {Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15932, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39397, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 16064, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39398, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15989, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39399, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15989, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39401, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15989, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39403, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15989, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39404, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15989, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39405, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15989, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39407, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15989, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39408, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15989, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39409, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15989, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39415, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15989, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39416, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15990, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39417, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15990, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39419, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15990, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39420, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15990, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39421, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15990, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39422, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15990, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39423, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15990, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39424, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15990, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39425, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15990, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39426, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 15990, Difficulty: 3, ZoneId: 3456}}}}},
	{Id: 39427, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{Difficulty: 3, ZoneId: 3456, OtherName: "Trash"}}}}},
	{Id: 39467, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{Difficulty: 3, ZoneId: 3456, OtherName: "Trash"}}}}},
	{Id: 39468, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{Difficulty: 3, ZoneId: 3456, OtherName: "Trash"}}}}},
	{Id: 39470, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{Difficulty: 3, ZoneId: 3456, OtherName: "Trash"}}}}},
	{Id: 39472, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{Difficulty: 3, ZoneId: 3456, OtherName: "Trash"}}}}},
	{Id: 39473, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{Difficulty: 3, ZoneId: 3456, OtherName: "Trash"}}}}},
	{Id: 40426, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 28860, Difficulty: 3, ZoneId: 4493}}}}},
	{Id: 40427, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 28860, Difficulty: 3, ZoneId: 4493}}}}},
	{Id: 40428, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 28860, Difficulty: 3, ZoneId: 4493}}}}},
	{Id: 40429, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 28860, Difficulty: 3, ZoneId: 4493}}}}},
	{Id: 40430, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 28860, Difficulty: 3, ZoneId: 4493}}}}},
	{Id: 40474, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 28859, Difficulty: 3, ZoneId: 4500}}}}},
	{Id: 40475, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 28859, Difficulty: 3, ZoneId: 4500}}}}},
	{Id: 40486, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 28859, Difficulty: 3, ZoneId: 4500}}}}},
	{Id: 40488, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 28859, Difficulty: 3, ZoneId: 4500}}}}},
	{Id: 40489, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 28859, Difficulty: 3, ZoneId: 4500}}}}},
	{Id: 40491, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 28859, Difficulty: 3, ZoneId: 4500}}}}},
	{Id: 40497, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 28859, Difficulty: 3, ZoneId: 4500}}}}},
	{Id: 40511, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 28859, Difficulty: 3, ZoneId: 4500}}}}},
	{Id: 40519, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 28859, Difficulty: 3, ZoneId: 4500}}}}},
	{Id: 40526, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 28859, Difficulty: 3, ZoneId: 4500}}}}},
	{Id: 43988, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 28860, Difficulty: 3, ZoneId: 4493, Category: "One Drake Left"}}}}},
	{Id: 43989, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 28860, Difficulty: 3, ZoneId: 4493, Category: "One Drake Left"}}}}},
	{Id: 43990, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 28860, Difficulty: 3, ZoneId: 4493, Category: "One Drake Left"}}}}},
	{Id: 43991, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 28860, Difficulty: 3, ZoneId: 4493, Category: "One Drake Left"}}}}},
	{Id: 43992, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 28860, Difficulty: 3, ZoneId: 4493, Category: "One Drake Left"}}}}},
	{Id: 43993, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 28860, Difficulty: 3, ZoneId: 4493, Category: "Two Drakes Left"}}}}},
	{Id: 43994, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 28860, Difficulty: 3, ZoneId: 4493, Category: "Two Drakes Left"}}}}},
	{Id: 43995, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 28860, Difficulty: 3, ZoneId: 4493, Category: "Two Drakes Left"}}}}},
	{Id: 43996, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 28860, Difficulty: 3, ZoneId: 4493, Category: "Two Drakes Left"}}}}},
	{Id: 43998, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 28860, Difficulty: 3, ZoneId: 4493, Category: "Two Drakes Left"}}}}},
	{Id: 44657, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Quest{Quest: &proto.QuestSource{Id: 13384, Name: "Judgment at the Eye of Eternity"}}}}},
	{Id: 44658, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Quest{Quest: &proto.QuestSource{Id: 13384, Name: "Judgment at the Eye of Eternity"}}}}},
	{Id: 44659, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Quest{Quest: &proto.QuestSource{Id: 13384, Name: "Judgment at the Eye of Eternity"}}}}},
	{Id: 44660, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Quest{Quest: &proto.QuestSource{Id: 13384, Name: "Judgment at the Eye of Eternity"}}}}},
	{Id: 45282, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33113, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45283, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33113, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45284, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33113, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45285, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33113, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45286, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33113, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45287, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33113, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45288, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33113, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45289, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33113, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45291, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33113, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45292, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33113, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45298, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33186, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45299, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33186, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45301, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33186, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45302, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33186, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45304, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33186, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45305, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33186, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45306, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33186, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45307, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33186, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45308, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33186, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45309, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33118, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45310, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33118, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45311, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33118, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45312, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33118, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45313, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33118, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45314, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33118, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45316, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33118, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45317, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33118, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45318, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33118, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45321, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33118, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45322, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 32857, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45324, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 32857, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45329, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 32857, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45330, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 32857, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45331, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 32857, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45332, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 32857, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45333, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 32857, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45378, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 32857, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45418, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 32857, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45423, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 32857, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45458, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 32845, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45464, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 32845, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45675, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33293, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45676, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33293, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45677, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33293, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45679, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33293, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45680, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33293, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45682, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33293, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45685, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33293, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45686, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33293, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45687, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33293, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45694, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33293, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45695, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 32930, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45696, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 32930, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45697, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 32930, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45698, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 32930, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45699, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 32930, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45701, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 32930, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45702, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 32930, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45703, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 32930, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45704, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 32930, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45707, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33515, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45708, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33515, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45709, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33515, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45711, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33515, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45712, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33515, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45713, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33515, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45832, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33515, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45864, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33515, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45865, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33515, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45866, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33515, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45872, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 32845, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45873, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 32845, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45874, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 32845, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45892, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 32865, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45893, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 32865, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45894, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 32865, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45927, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 32865, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45973, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 32865, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45975, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33350, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45976, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33350, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45996, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33271, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 45997, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33271, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 46008, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33271, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 46009, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33271, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 46010, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33271, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 46011, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33271, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 46012, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33271, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 46013, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33271, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 46014, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33271, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 46015, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33271, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 46016, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33288, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 46018, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33288, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 46019, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33288, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 46021, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33288, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 46022, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33288, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 46024, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33288, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 46025, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33288, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 46028, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33288, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 46030, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33288, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 46031, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{NpcId: 33288, Difficulty: 3, ZoneId: 4273}}}}},
	{Id: 46339, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{Difficulty: 3, ZoneId: 4273, OtherName: "Trash"}}}}},
	{Id: 46340, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{Difficulty: 3, ZoneId: 4273, OtherName: "Trash"}}}}},
	{Id: 46341, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{Difficulty: 3, ZoneId: 4273, OtherName: "Trash"}}}}},
	{Id: 46342, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{Difficulty: 3, ZoneId: 4273, OtherName: "Trash"}}}}},
	{Id: 46343, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{Difficulty: 3, ZoneId: 4273, OtherName: "Trash"}}}}},
	{Id: 46344, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{Difficulty: 3, ZoneId: 4273, OtherName: "Trash"}}}}},
	{Id: 46345, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{Difficulty: 3, ZoneId: 4273, OtherName: "Trash"}}}}},
	{Id: 46346, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{Difficulty: 3, ZoneId: 4273, OtherName: "Trash"}}}}},
	{Id: 46347, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{Difficulty: 3, ZoneId: 4273, OtherName: "Trash"}}}}},
	{Id: 46350, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{Difficulty: 3, ZoneId: 4273, OtherName: "Trash"}}}}},
	{Id: 46351, Sources: []*proto.UIItemSource{{Source: &proto.UIItemSource_Drop{Drop: &proto.DropSource{Difficulty: 3, ZoneId: 4273, OtherName: "Trash"}}}}},
}

// Keep these sorted by item ID.
var ItemAllowList = map[int32]struct{}{
	9380:  {}, //Jang'traze
	11815: {}, // Hand of Justice
	12631: {}, // Fiery Plate Gauntlets
	12590: {}, // Felstriker
	15808: {}, // Fine Light Crossbow (for hunter testing).
	18843: {},
	18844: {},
	18847: {},
	18848: {},
	19019: {}, // Thunderfury
	19808: {}, // Rockhide Strongfish
	20837: {}, // Sunstrider Axe
	20966: {}, // Jade Pendant of Blasting
	21625: {}, // Scarab Brooch
	21685: {}, // Petrified Scarab
	24114: {}, // Braided Eternium Chain
	28572: {}, // Blade of the Unrequited
	28830: {}, // Dragonspine Trophy
	29383: {}, // Bloodlust Brooch
	29387: {}, // Gnomeregan Auto-Blocker 600
	29994: {}, // Thalassian Wildercloak
	29996: {}, // Rod of the Sun King
	30032: {}, // Red Belt of Battle
	30627: {}, // Tsunami Talisman
	30720: {}, // Serpent-Coil Braid
	31193: {}, // Blade of Unquenched Thirst
	32387: {}, // Idol of the Raven Goddess
	32658: {}, // Badge of Tenacity
	33135: {}, // Falling Star
	33140: {}, // Blood of Amber
	33143: {}, // Stone of Blades
	33144: {}, // Facet of Eternity
	33504: {}, // Libram of Divine Purpose
	33506: {}, // Skycall Totem
	33507: {}, // Stonebreaker's Totem
	33508: {}, // Idol of Budding Life
	33510: {}, // Unseen moon idol
	33829: {}, // Hex Shrunken Head
	33831: {}, // Berserkers Call
	34472: {}, // Shard of Contempt
	34473: {}, // Commendation of Kael'thas
	37032: {}, // Edge of the Tuskarr
	37574: {}, // Libram of Furious Blows
	38072: {}, // Thunder Capacitor
	38212: {}, // Death Knight's Anguish
	38287: {}, // Empty Mug of Direbrew
	38289: {}, // Coren's Lucky Coin
	39208: {}, // Sigil of the Dark Rider
	41752: {}, // Brunnhildar Axe
	6360:  {}, // Steelscale Crushfish
	8345:  {}, // Wolfshead Helm
	9449:  {}, // Manual Crowd Pummeler

	// Sets
	27510: {}, // Tidefury Gauntlets
	27802: {}, // Tidefury Shoulderguards
	27909: {}, // Tidefury Kilt
	28231: {}, // Tidefury Chestpiece
	28349: {}, // Tidefury Helm

	15056: {}, // Stormshroud Armor
	15057: {}, // Stormshroud Pants
	15058: {}, // Stormshroud Shoulders
	21278: {}, // Stormshroud Gloves

	// Undead Slaying Sets
	// Plate
	43068: {},
	43069: {},
	43070: {},
	43071: {},
	// Cloth
	43072: {},
	43073: {},
	43074: {},
	43075: {},
	// Mail
	43076: {},
	43077: {},
	43078: {},
	43079: {},
	//Leather
	43080: {},
	43081: {},
	43082: {},
	43083: {},
}

// Keep these sorted by item ID.
var ItemDenyList = map[int32]struct{}{
	17782: {}, // talisman of the binding shard
	17783: {}, // talisman of the binding fragment
	17802: {}, // Deprecated version of Thunderfury
	18582: {},
	18583: {},
	18584: {},
	22736: {},
	24265: {},
	33046: {}, //PvP Pwn
	32384: {},
	32421: {},
	32466: {},
	32422: {},
	33482: {},
	33350: {},
	34576: {}, // Battlemaster's Cruelty
	34577: {}, // Battlemaster's Depreavity
	34578: {}, // Battlemaster's Determination
	34579: {}, // Battlemaster's Audacity
	34580: {}, // Battlemaster's Perseverence
	38282: {},

	38694: {}, // "Family" Shoulderpads heirloom
	39263: {}, // "Dissevered Leggings"

	41587: {},
	41588: {},
	41589: {},
	41590: {},
	45084: {}, // 'Book of Crafting Secrets' heirloom

	// '10 man' onyxia head rewards
	49312: {},
	49313: {},
	49314: {},

	50251: {}, // 'one hand shadows edge'
	53500: {}, // Tectonic Plate

	48880: {}, // DK's Tier 9 Duplicates
	48881: {}, // DK's Tier 9 Duplicates
	48882: {}, // DK's Tier 9 Duplicates
	48883: {}, // DK's Tier 9 Duplicates
	48884: {}, // DK's Tier 9 Duplicates
	48885: {}, // DK's Tier 9 Duplicates
	48886: {}, // DK's Tier 9 Duplicates
	48887: {}, // DK's Tier 9 Duplicates
	48888: {}, // DK's Tier 9 Duplicates
	48889: {}, // DK's Tier 9 Duplicates
	48890: {}, // DK's Tier 9 Duplicates
	48891: {}, // DK's Tier 9 Duplicates
	48892: {}, // DK's Tier 9 Duplicates
	48893: {}, // DK's Tier 9 Duplicates
	48894: {}, // DK's Tier 9 Duplicates
	48895: {}, // DK's Tier 9 Duplicates
	48896: {}, // DK's Tier 9 Duplicates
	48897: {}, // DK's Tier 9 Duplicates
	48898: {}, // DK's Tier 9 Duplicates
	48899: {}, // DK's Tier 9 Duplicates

	50741: {}, //rp item
}

// Item icons to include in the DB, so they don't need to be separately loaded in the UI.
var ExtraItemIcons = []int32{
	// Pet foods
	33874,
	43005,

	// Spellstones
	41174,
	41196,

	// Demonic Rune
	12662,

	// Food IDs
	27655,
	27657,
	27658,
	27664,
	33052,
	33825,
	33872,
	34753,
	34754,
	34756,
	34758,
	34767,
	34769,
	42994,
	42995,
	42996,
	42998,
	42999,
	43000,
	43015,

	// Flask IDs
	13512,
	22851,
	22853,
	22854,
	22861,
	22866,
	33208,
	40079,
	44939,
	46376,
	46377,
	46378,
	46379,

	// Elixer IDs
	40072,
	40078,
	40097,
	40109,
	44328,
	44332,

	// Elixer IDs
	13452,
	13454,
	22824,
	22827,
	22831,
	22833,
	22834,
	22835,
	22840,
	28103,
	28104,
	31679,
	32062,
	32067,
	32068,
	39666,
	40068,
	40070,
	40073,
	40076,
	44325,
	44327,
	44329,
	44330,
	44331,
	9088,
	9224,

	// Potions / In Battle Consumes
	13442,
	20520,
	22105,
	22788,
	22828,
	22832,
	22837,
	22838,
	22839,
	22849,
	31677,
	33447,
	33448,
	36892,
	40093,
	40211,
	40212,
	40536,
	40771,
	41119,
	41166,
	42545,
	42641,

	// Poisons
	43231,
	43233,
	43235,

	// Thistle Tea
	7676,

	// Scrolls
	37094,
	43466,
	43464,
	37092,
	37098,
	43468,

	// Drums
	49633,
	49634,
}

// Raid buffs / debuffs
var SharedSpellsIcons = []int32{
	// Revitalize, Rejuv, WG
	48545,
	26982,
	53251,

	// Registered CD's
	49016,
	57933,
	64382,
	10060,
	16190,
	29166,
	53530,
	33206,
	2825,
	54758,

	// Raid Buffs
	43002,
	57567,
	54038,

	48470,
	17051,

	25898,
	25899,

	48942,
	20140,
	58753,
	16293,

	48161,
	14767,

	58643,
	52456,
	57623,

	48073,

	48934,
	20045,
	47436,

	53138,
	30809,
	19506,

	31869,
	31583,
	34460,

	57472,
	50720,

	53648,

	47440,
	12861,
	47982,
	18696,

	48938,
	20245,
	58774,
	16206,

	17007,
	34300,
	29801,

	55610,
	65990,
	29193,

	48160,
	31878,
	53292,
	54118,
	44561,

	24907,
	48396,
	51470,

	3738,
	47240,
	57722,
	58656,

	54043,
	48170,
	31025,
	31035,
	6562,
	31033,
	53307,
	16840,
	54648,

	// Raid Debuffs
	8647,
	7386,
	55754,

	770,
	33602,
	50511,
	18180,
	56631,
	53598,

	26016,
	47437,
	12879,
	48560,
	16862,
	55487,

	48566,
	46855,
	57393,

	30706,
	20337,
	58410,

	47502,
	12666,
	55095,
	51456,
	53696,
	48485,

	3043,
	29859,
	58413,
	65855,

	17800,
	17803,
	12873,
	28593,

	33198,
	51161,
	48511,
	47865,

	20271,
	53408,

	11374,
	15235,

	27013,

	58749,
	49071,

	30708,
}

// If any of these match the item name, don't include it.
var DenyListNameRegexes = []*regexp.Regexp{
	regexp.MustCompile(`30 Epic`),
	regexp.MustCompile(`63 Blue`),
	regexp.MustCompile(`63 Green`),
	regexp.MustCompile(`66 Epic`),
	regexp.MustCompile(`90 Epic`),
	regexp.MustCompile(`90 Green`),
	regexp.MustCompile(`Boots 1`),
	regexp.MustCompile(`Boots 2`),
	regexp.MustCompile(`Boots 3`),
	regexp.MustCompile(`Bracer 1`),
	regexp.MustCompile(`Bracer 2`),
	regexp.MustCompile(`Bracer 3`),
	regexp.MustCompile(`DB\d`),
	regexp.MustCompile(`DEPRECATED`),
	regexp.MustCompile(`Deprecated: Keanna`),
	regexp.MustCompile(`Indalamar`),
	regexp.MustCompile(`Monster -`),
	regexp.MustCompile(`NEW`),
	regexp.MustCompile(`PH`),
	regexp.MustCompile(`QR XXXX`),
	regexp.MustCompile(`TEST`),
	regexp.MustCompile(`Test`),
	regexp.MustCompile(`zOLD`),
}

// Allows manual overriding for Gem fields in case WowHead is wrong.
var GemOverrides = []*proto.UIGem{
	{Id: 33131, RequiredProfession: proto.Profession_Jewelcrafting, Unique: true, Expansion: proto.Expansion_ExpansionTbc},
	{Id: 33133, RequiredProfession: proto.Profession_Jewelcrafting, Unique: true, Expansion: proto.Expansion_ExpansionTbc},
	{Id: 33134, RequiredProfession: proto.Profession_Jewelcrafting, Unique: true, Expansion: proto.Expansion_ExpansionTbc},
	{Id: 33135, RequiredProfession: proto.Profession_Jewelcrafting, Unique: true, Expansion: proto.Expansion_ExpansionTbc},
	{Id: 33140, RequiredProfession: proto.Profession_Jewelcrafting, Unique: true, Expansion: proto.Expansion_ExpansionTbc},
	{Id: 33143, RequiredProfession: proto.Profession_Jewelcrafting, Unique: true, Expansion: proto.Expansion_ExpansionTbc},
	{Id: 34256, Unique: true, Expansion: proto.Expansion_ExpansionTbc},
	{Id: 34831, Unique: true, Expansion: proto.Expansion_ExpansionTbc},
}
var GemAllowList = map[int32]struct{}{
	22459: {}, // Void Sphere
	36766: {}, // Bright Dragon's Eye
	36767: {}, // Solid Dragon's Eye
}
var GemDenyList = map[int32]struct{}{
	// pvp non-unique gems not in game currently.
	32735: {},
	34142: {}, // Infinite Sphere
	34143: {}, // Chromatic Sphere
	35489: {},
	37430: {}, // Solid Sky Sapphire (Unused)
	38545: {},
	38546: {},
	38547: {},
	38548: {},
	38549: {},
	38550: {},
}
