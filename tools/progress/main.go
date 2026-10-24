// Command progress reports checked roadmap items per module and per tier.
//
// Usage: go run ./tools/progress
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

var (
	tiers  = []string{"E", "S", "D", "T"}
	itemRe = regexp.MustCompile("^\\s*(?:- |\\| )\\[([ xX])\\] `\\[([ESDT])\\]`")
)

type tally struct {
	done  map[string]int
	total map[string]int
}

func newTally() tally {
	return tally{done: map[string]int{}, total: map[string]int{}}
}

func (t tally) add(tier string, checked bool) {
	t.total[tier]++
	if checked {
		t.done[tier]++
	}
}

func main() {
	roadmaps, err := filepath.Glob("[0-9][0-9]-*/roadmap.md")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	slices.Sort(roadmaps)

	overall := newTally()
	fmt.Printf("%-30s %8s %8s %8s %8s\n", "module", "E", "S", "D", "T")
	for _, path := range roadmaps {
		module, err := countFile(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("%-30s", filepath.Dir(path))
		for _, tier := range tiers {
			fmt.Printf(" %8s", cell(module, tier))
			overall.done[tier] += module.done[tier]
			overall.total[tier] += module.total[tier]
		}
		fmt.Println()
	}

	fmt.Printf("%-30s", "TOTAL")
	for _, tier := range tiers {
		fmt.Printf(" %8s", cell(overall, tier))
	}
	fmt.Println()
}

func countFile(path string) (tally, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return tally{}, err
	}

	t := newTally()
	for line := range strings.Lines(string(data)) {
		match := itemRe.FindStringSubmatch(line)
		if match == nil {
			continue
		}
		checked := strings.EqualFold(match[1], "x")
		t.add(match[2], checked)
	}
	return t, nil
}

func cell(t tally, tier string) string {
	if t.total[tier] == 0 {
		return "-"
	}
	return fmt.Sprintf("%d/%d", t.done[tier], t.total[tier])
}
