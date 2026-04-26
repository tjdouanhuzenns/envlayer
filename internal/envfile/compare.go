package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// CompareResult holds the result of comparing two named environments.
type CompareResult struct {
	LeftName  string
	RightName string
	OnlyLeft  []string
	OnlyRight []string
	Differ    map[string][2]string // key -> [leftVal, rightVal]
	Same      []string
}

// Summary returns a human-readable summary of the comparison.
func (r *CompareResult) Summary() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Comparing [%s] vs [%s]\n", r.LeftName, r.RightName)
	fmt.Fprintf(&sb, "  Only in %s (%d): %s\n", r.LeftName, len(r.OnlyLeft), strings.Join(r.OnlyLeft, ", "))
	fmt.Fprintf(&sb, "  Only in %s (%d): %s\n", r.RightName, len(r.OnlyRight), strings.Join(r.OnlyRight, ", "))
	fmt.Fprintf(&sb, "  Different (%d):\n", len(r.Differ))
	keys := make([]string, 0, len(r.Differ))
	for k := range r.Differ {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		pair := r.Differ[k]
		fmt.Fprintf(&sb, "    %s: %q -> %q\n", k, pair[0], pair[1])
	}
	fmt.Fprintf(&sb, "  Same (%d)\n", len(r.Same))
	return sb.String()
}

// Compare performs a detailed comparison between two EnvMaps.
func Compare(leftName string, left EnvMap, rightName string, right EnvMap) *CompareResult {
	result := &CompareResult{
		LeftName:  leftName,
		RightName: rightName,
		Differ:    make(map[string][2]string),
	}

	allKeys := make(map[string]struct{})
	for k := range left {
		allKeys[k] = struct{}{}
	}
	for k := range right {
		allKeys[k] = struct{}{}
	}

	for k := range allKeys {
		lv, inLeft := left[k]
		rv, inRight := right[k]
		switch {
		case inLeft && !inRight:
			result.OnlyLeft = append(result.OnlyLeft, k)
		case !inLeft && inRight:
			result.OnlyRight = append(result.OnlyRight, k)
		case lv == rv:
			result.Same = append(result.Same, k)
		default:
			result.Differ[k] = [2]string{lv, rv}
		}
	}

	sort.Strings(result.OnlyLeft)
	sort.Strings(result.OnlyRight)
	sort.Strings(result.Same)
	return result
}
