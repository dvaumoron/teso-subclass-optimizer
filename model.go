/*
 *
 * Copyright 2025 dvaumoron.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package main

import (
	"cmp"
	"maps"
	"slices"
	"strings"
)

const (
	att Category = iota
	def
	move
)

type Category int

// parse from french ^^;
func ParseCategory(s string) (Category, bool) {
	switch strings.ToLower(s) {
	case "offensif":
		return att, true
	case "défensif":
		return def, true
	case "déplacement":
		return move, true
	default:
		return 0, false
	}
}

func (c Category) Filter(b *Buff) bool {
	return c == b.category
}

const (
	none DamageType = iota
	physical
	magical
	all
)

type DamageType int

// parse from french ^^;
func ParseDamageType(s string) DamageType {
	switch strings.ToLower(s) {
	case "physique":
		return physical
	case "magique":
		return magical
	case "tout":
		return all
	default:
		return none
	}
}

func (dt DamageType) Filter(b *Buff) bool {
	switch dt2 := b.damage; dt {
	case physical:
		return dt2 == physical || dt2 == all
	case magical:
		return dt2 == magical || dt2 == all
	case all:
		return dt2 != none
	default:
		return true
	}
}

type Buff struct {
	name        string
	category    Category
	damage      DamageType
	description string
	skillGroups map[string]struct{}
	valid       bool
}

type SkillGroup struct {
	name  string
	buffs map[string]*Buff
}

func (sg SkillGroup) ToCSV() []string {
	res := make([]string, len(sg.buffs)+3)
	copy(res, strings.Split(sg.name, ","))
	copy(res[3:], slices.Sorted(maps.Keys(sg.buffs)))
	return res
}

func compareSkillGroup(a, b SkillGroup) int {
	return cmp.Compare(len(b.buffs), len(a.buffs))
}
