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
	"strconv"
	"strings"
)

var (
	headerScored = []string{"Group Name 1", "Group Name 2", "Group Name 3", "Score", "Buffs"}
	headerBasic  = []string{"Group Name 1", "Group Name 2", "Group Name 3", "Buffs"}
)

func main() {
	skillGroupFilePath, scoreFilePath := "", ""
	switch lenArgs := len(os.Args); {
	case lenArgs < 2:
		panic("Please provide a path to the sub-class data file.")
	default:
		scoreFilePath = os.Args[2]
		fallthrough
	case lenArgs == 2:
		skillGroupFilePath = os.Args[1]
	}

	skillGroups, buffs := readSkillGroupData(skillGroupFilePath)

	mergedBy3 := mergeBy3(skillGroups)

	compareSkillGroup := compareSkillGroupByNumber
	scoreFlag := scoreFilePath != ""
	if scoreFlag {
		readScore(scoreFilePath, buffs)
		compareSkillGroup = compareSkillGroupByScore
	}

	slices.SortFunc(mergedBy3, compareSkillGroup)

	printSkillGroupsCSV(mergedBy3, scoreFlag)
}

func readSkillGroupData(path string) ([]SkillGroup, []*Buff) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	data, err := csv.NewReader(file).ReadAll()
	if err != nil {
		panic(err)
	}

	groupNames := data[0][4:]
	skillGroupCount := len(groupNames)
	buffs := make([]*Buff, len(data)-1)
	for buffIndex, buffLine := range data[1:] {
		buff := Buff{
			name:        buffLine[0],
			category:    buffLine[1],
			damage:      buffLine[2],
			description: buffLine[3],
			skillGroups: make(map[string]struct{}, skillGroupCount),
		}

		for groupIndex, okStr := range buffLine[4:] {
			if okStr != "" {
				buff.skillGroups[groupNames[groupIndex]] = struct{}{}
			}
		}

		buffs[buffIndex] = &buff
	}

	skillGroups := make([]SkillGroup, skillGroupCount)
	for groupIndex, groupName := range groupNames {
		groupBuffs := map[string]*Buff{}
		for _, buff := range buffs {
			if _, ok := buff.skillGroups[groupName]; ok {
				groupBuffs[buff.name] = buff
			}
		}

		skillGroups[groupIndex] = SkillGroup{
			name:  groupName,
			buffs: groupBuffs,
		}
	}

	return skillGroups, buffs
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

func readScore(path string, buffs []*Buff) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	data, err := csv.NewReader(file).ReadAll()
	if err != nil {
		panic(err)
	}

	buffNameToScore := make(map[string]uint, len(data))
	for buffIndex, buffLine := range data {
		buffScoreStr := buffLine[1]
		if buffScoreStr == "" {
			continue
		}

		buffScore, err := strconv.ParseUint(buffScoreStr, 10, 32)
		if err != nil {
			if buffIndex == 0 { // could have a header line
				continue
			}
			panic(err)
		}

		buffNameToScore[buffLine[0]] = uint(buffScore)
	}

	for _, buff := range buffs {
		if score, ok := buffNameToScore[buff.name]; ok {
			buff.score = score
		}
	}
}

func printSkillGroupsCSV(mergedBy3 []SkillGroup, scoreFlag bool) {
	csvWriter := csv.NewWriter(os.Stdout)

	header := headerBasic
	if scoreFlag {
		header = headerScored
	}
	csvWriter.Write(header)

	for _, group := range mergedBy3 {
		csvWriter.Write(group.ToCSV(scoreFlag))
	}
	csvWriter.Flush()
}
