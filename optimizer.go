package main

import (
	"cmp"
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

	filter := ""
	if len(os.Args) > 2 {
		filter = os.Args[2]
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

type Buff struct {
	name        string
	category    string
	description string
	valid       bool
}

func readSkillGroupData(path string, filter string) []SkillGroup {
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

	buffs := make([]Buff, len(data)-1)
	for buffIndex, buffLine := range data[1:] {
		buff := Buff{
			name:        buffLine[0],
			category:    buffLine[1],
			description: buffLine[2],
		}

		buff.valid = filter == "" || buff.category == filter

		buffs[buffIndex] = buff
	}

	groupNames := data[0][3:]
	skillGroups := make([]SkillGroup, len(groupNames))
	for groupIndex, groupName := range groupNames {
		columnIndex := groupIndex + 3
		groupBuffs := map[string]*Buff{}
		for buffIndex, buffLine := range data[1:] {
			if buffLine[columnIndex] != "" && buffs[buffIndex].valid {
				buffName := buffLine[0]

				groupBuffs[buffName] = &buffs[buffIndex]
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

func compareSkillGroup(a, b SkillGroup) int {
	return cmp.Compare(len(b.buffs), len(a.buffs))
}
