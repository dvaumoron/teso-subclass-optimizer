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
	"strconv"
	"strings"
)

var scoreCache = map[string]uint{}

type Buff struct {
	name        string
	category    string
	damage      string
	description string
	skillGroups map[string]struct{}
	score       uint
}

func compareBuffByScoreName(a *Buff, b *Buff) int {
	return cmp.Or(
		cmp.Compare(b.score, a.score),
		cmp.Compare(a.name, b.name),
	)
}

type SkillGroup struct {
	name  string
	buffs map[string]*Buff
}

func (sg SkillGroup) ToCSV(scoreFlag bool) []string {
	res := make([]string, 3, len(sg.buffs)+3)
	copy(res, strings.Split(sg.name, ","))
	if scoreFlag {
		res = append(res, strconv.FormatUint(uint64(sumScore(sg)), 10))
		for _, buff := range slices.SortedFunc(maps.Values(sg.buffs), compareBuffByScoreName) {
			res = append(res, buff.name)
		}
	} else {
		res = append(res, slices.Sorted(maps.Keys(sg.buffs))...)
	}
	return res
}

func compareSkillGroupByNumber(a SkillGroup, b SkillGroup) int {
	return cmp.Compare(len(b.buffs), len(a.buffs))
}

func compareSkillGroupByScore(a SkillGroup, b SkillGroup) int {
	return cmp.Compare(sumScore(b), sumScore(a))
}

func sumScore(group SkillGroup) uint {
	groupName := group.name
	score, ok := scoreCache[groupName]
	if ok {
		return score
	}

	for _, buff := range group.buffs {
		score += buff.score
	}

	scoreCache[groupName] = score

	return score
}
