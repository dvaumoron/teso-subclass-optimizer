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
	"encoding/csv"
	"maps"
	"os"
	"slices"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		panic("Please provide a path to the sub-class data file.")
	}

	filter := noFilter
	if len(os.Args) > 2 {
		filter = parseFilter(os.Args[2])
	}

	skillGroups := readSkillGroupData(os.Args[1], filter)

	mergedBy3 := mergeBy3(skillGroups)

	slices.SortFunc(mergedBy3, compareSkillGroup)

	csvWriter := csv.NewWriter(os.Stdout)

	csvWriter.Write([]string{
		"Group Name 1",
		"Group Name 2",
		"Group Name 3",
		"Buffs",
	})

	for _, group := range mergedBy3 {
		csvWriter.Write(group.ToCSV())
	}
	csvWriter.Flush()
}

func readSkillGroupData(path string, buffFilter BuffFilter) []SkillGroup {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	data, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}

	groupNames := data[0][4:]
	skillGroupCount := len(groupNames)
	buffs := make([]*Buff, len(data)-1)
	for buffIndex, buffLine := range data[1:] {
		catStr := buffLine[1]
		category, ok := ParseCategory(catStr)
		if !ok {
			panic("Invalid category: " + catStr)
		}

		buff := Buff{
			name:        buffLine[0],
			category:    category,
			damage:      ParseDamageType(buffLine[2]),
			description: buffLine[3],
			skillGroups: make(map[string]struct{}, skillGroupCount),
		}

		for groupIndex, okStr := range buffLine[4:] {
			if okStr != "" {
				buff.skillGroups[groupNames[groupIndex]] = struct{}{}
			}
		}

		buff.valid = buffFilter(&buff)

		buffs[buffIndex] = &buff
	}

	skillGroups := make([]SkillGroup, skillGroupCount)
	for groupIndex, groupName := range groupNames {
		groupBuffs := map[string]*Buff{}
		for _, buff := range buffs {
			if _, ok := buff.skillGroups[groupName]; ok && buff.valid {
				groupBuffs[buff.name] = buff
			}
		}

		skillGroups[groupIndex] = SkillGroup{
			name:  groupName,
			buffs: groupBuffs,
		}
	}

	return skillGroups
}

func mergeSkillGroups(groups []SkillGroup) SkillGroup {
	names := make([]string, len(groups))
	for i, group := range groups {
		names[i] = group.name
	}

	merged := SkillGroup{
		name:  strings.Join(names, ","),
		buffs: map[string]*Buff{},
	}

	for _, group := range groups {
		maps.Copy(merged.buffs, group.buffs)
	}

	return merged
}

func mergeBy3(groups []SkillGroup) []SkillGroup {
	end := len(groups)
	if end < 3 {
		panic("Not enough skill groups to merge. At least 3 are required.")
	}

	iEnd := end - 2
	jEnd := end - 1

	var mergedGroups []SkillGroup
	for i := 0; i < iEnd; i++ {
		for j := i + 1; j < jEnd; j++ {
			for k := j + 1; k < end; k++ {
				groupedBy3 := []SkillGroup{
					groups[i],
					groups[j],
					groups[k],
				}

				mergedGroups = append(mergedGroups, mergeSkillGroups(groupedBy3))
			}
		}
	}

	return mergedGroups
}
