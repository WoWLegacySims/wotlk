#!/usr/bin/python

import mysql.connector
from typing import List
import csv

# Generates go/ts files for assorted values read from the acore db

def GenIntIndexedDb(file : str):
    db = {}
    with open(file) as tsv:
        first = True
        for line in csv.reader(tsv, delimiter="\t"):
            if first:
                first = False
                continue
            db[int(line[0])] = line[1:]
    return db

BASE_DIR = ""

DIR_PATH = "assets/db_inputs/spells/"
SPELL = "Spell.csv"

SPELLDBC = GenIntIndexedDb(DIR_PATH + SPELL)

Off = {
    "Effect": 68,
    "EffectDie": 71,
    "EffectScaling": 74,
    "EffectBasePoints": 77,
    "EffectAura": 92,
    "EffectMiscValueA": 107,
    "EffectTriggerSpell": 113,
}

AuraEffect = {
    10: "threat",
    14: "damage taken",
    35: "increase energy",
    39: "school immune",
    40: "damage immune",
    42: "proc trigger",
    59: "damage done creature",
    69: "school absorb",
    71: "crit chance school",
    72: "cost school pct",
    73: "cost school",
    79: "damage percent done",
    80: "percent stat",
    87: "damage percent taken",
    102: "attack power vs",
    107: "add flat mod",
    108: "add pct mod",
    110: "power regen pct",
    113: "ranged dmg taken",
    114: "ranged dmg taken pct",
    118: "healing pct",
    131: "rap vs",
    132: "increase energy pct",
    136: "healing done pct",
    137: "total stat pct",
}

SpellEffect = {
    "ApplyAura": 6,
}

ItemSpellTrigger= {
    "Use": 0,
    "Equip": 1,
    "Chance on hit": 2,
}

class ItemSpell:
    entry:int
    trigger:int

class Item:
    entry: int
    name: str
    spell: List[ItemSpell]


def FilterItems(item: Item):
    for spell in item.spell:
        if (spell.entry > 0):
            if(spell.trigger == ItemSpellTrigger["Use"]):
                return True
            elif(spell.trigger == ItemSpellTrigger["Chance on hit"]):
                return True
            elif(spell.trigger == ItemSpellTrigger["Equip"]):
                if (spell.entry in SPELLDBC.keys()):
                    spellinfo = SPELLDBC[spell.entry]
                    for i in range(0,3):
                        effect = int(spellinfo[Off["Effect"]+i])
                        if(effect == SpellEffect["ApplyAura"]):
                            auraeffect = int(spellinfo[Off["EffectAura"]])
                            if (auraeffect in AuraEffect.keys()):
                                print(f"{item.name}: {item.entry} was kept because {AuraEffect[auraeffect]} {spell.entry}")
                                return True

    return False



if __name__ == "__main__":
    db = mysql.connector.connect(host="localhost", user="root", password="root", database="acore_world")

    cursor = db.cursor()
    cursor.execute("SELECT entry,name,spellid_1,spelltrigger_1,spellid_2,spelltrigger_2,spellid_3,spelltrigger_3 from item_template WHERE InventoryType > 0  AND (spellid_1 > 0 OR  spellid_2 > 0 OR spellid_3 > 0)")
    results = cursor.fetchall()
    items = []
    for x in results:
        item = Item()
        item.entry = x[0]
        item.name = x[1]
        item.spell = []
        spell = ItemSpell()
        spell.entry = x[2]
        spell.trigger = x[3]
        item.spell.append(spell)
        spell = ItemSpell()
        spell.entry = x[4]
        spell.trigger = x[5]
        item.spell.append(spell)
        spell = ItemSpell()
        spell.entry = x[6]
        spell.trigger = x[7]
        item.spell.append(spell)
        items.append(item)

    items = filter(FilterItems,items)


    output = ""
    for item in items:
        output += f"{item.name}\t{item.entry}\n"

    fname = BASE_DIR + "sim/items.csv"

    f = open(fname, "w")
    f.write(output)
    f.close()

