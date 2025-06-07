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
	"math"
	"slices"
	"strings"
)

const maxPriority = 100

var prioritiesCache = map[string][]uint{}

type Buff struct {
	name        string
	category    string
	damage      string
	description string
	skillGroups map[string]struct{}
	priority    uint
}

func zeroAsMax(i uint) uint {
	if i == 0 {
		return math.MaxUint
	}
	return i
}

func compareBuffByPriorityName(a *Buff, b *Buff) int {
	return cmp.Or(
		cmp.Compare(zeroAsMax(a.priority), zeroAsMax(b.priority)),
		cmp.Compare(a.name, b.name),
	)
}

type SkillGroup struct {
	name  string
	buffs map[string]*Buff
}

func (sg SkillGroup) ToCSV(priorityFlag bool) []string {
	res := make([]string, 3, len(sg.buffs)+3)
	copy(res, strings.Split(sg.name, ","))
	if priorityFlag {
		for _, buff := range slices.SortedFunc(maps.Values(sg.buffs), compareBuffByPriorityName) {
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

func compareSkillGroupByPriority(a SkillGroup, b SkillGroup) int {
	aCount := countByPriority(a)
	bCount := countByPriority(b)

	for i := 1; i < maxPriority; i++ {
		if c := cmp.Compare(bCount[i], aCount[i]); c != 0 {
			return c
		}
	}

	return cmp.Compare(bCount[0], aCount[0])
}

func countByPriority(group SkillGroup) []uint {
	groupName := group.name
	priorities, ok := prioritiesCache[groupName]
	if ok {
		return priorities
	}

	priorities = make([]uint, maxPriority)
	for _, buff := range group.buffs {
		priorities[buff.priority]++
	}

	prioritiesCache[groupName] = priorities

	return priorities
}
