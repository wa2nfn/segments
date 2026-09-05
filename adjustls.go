package main

import "strings"

// adjustLs redistributes trailing Ls into internal L-runs of length >= 3,
// preserving at least a distance of 2 between small and large groups.
//
// Signature:
//
//	str       → input string (contains L and D, tail is L-only)
//	wordspace → minimum tail length to retain
//
// Returns:
//
//	adjusted string, excess code:
//	  excess == 0  → tail ended exactly at wordspace
//	  excess > 0   → balancing occurred, tail still has leftover beyond wordspace
//	  excess == -1 → no balancing possible (excess too small)
func adjustLs(str string, wordspace int) (string, int) {
	if len(str) == 0 {
		return str, 0
	}

	// Count trailing Ls.
	endCount := 0
	for i := len(str) - 1; i >= 0 && str[i] == 'L'; i-- {
		endCount++
	}
	if endCount == 0 {
		return str, 0
	}

	// If tail is already at or below wordspace, nothing to do.
	if endCount <= wordspace {
		return str, -1
	}

	body := str[:len(str)-endCount]

	// Parse runs of L and D in the body.
	type run struct {
		char byte
		len  int
	}

	var runs []run
	if len(body) > 0 {
		last := body[0]
		count := 1
		for i := 1; i < len(body); i++ {
			if body[i] == last {
				count++
			} else {
				runs = append(runs, run{char: last, len: count})
				last = body[i]
				count = 1
			}
		}
		runs = append(runs, run{char: last, len: count})
	}

	// Collect indexes of L-runs with len >= 3.
	var lIdx []int
	for i := range runs {
		if runs[i].char == 'L' && runs[i].len >= 3 {
			lIdx = append(lIdx, i)
		}
	}
	if len(lIdx) == 0 {
		return str, -1
	}

	tailLen := endCount
	didBalance := false

	// Helper to recompute small and large sizes among L-groups.
	getSmallLarge := func() (int, int) {
		if len(lIdx) == 0 {
			return 0, 0
		}
		small := runs[lIdx[0]].len
		large := runs[lIdx[0]].len

		for _, idx := range lIdx[1:] {
			n := runs[idx].len
			if n < small {
				small = n
			}
			if n > large {
				large = n
			}
		}
		return small, large
	}

	// -------------------------
	// PASS 0 — Ensure distance >= 2 (grow only largest groups if needed)
	// -------------------------
	{
		small, large := getSmallLarge()
		if large-small < 2 && large > 0 && small > 0 {
			// Collect largest groups.
			var largeIdx []int
			for _, idx := range lIdx {
				if runs[idx].len == large {
					largeIdx = append(largeIdx, idx)
				}
			}
			for large-small < 2 && len(largeIdx) > 0 {
				if tailLen <= len(largeIdx) {
					break
				}
				nextTail := tailLen - len(largeIdx)
				// Tail must remain >= new largest and >= wordspace.
				if nextTail < large+1 || nextTail < wordspace {
					break
				}
				for _, idx := range largeIdx {
					runs[idx].len++
				}
				tailLen = nextTail
				didBalance = true
				large++
				// small unchanged; distance increases.
				small, _ = getSmallLarge()
			}
		}
	}

	// -------------------------
	// PASS 1 — Grow all eligible L-groups equally
	// -------------------------
	for {
		if tailLen <= len(lIdx) {
			break
		}

		_, large := getSmallLarge()
		nextTail := tailLen - len(lIdx)

		// After increment, all groups grow by 1, so largest becomes large+1.
		// Tail must remain >= new largest and >= wordspace.
		if nextTail < large+1 || nextTail < wordspace {
			break
		}

		for _, idx := range lIdx {
			runs[idx].len++
		}
		tailLen = nextTail
		didBalance = true
	}

	// -------------------------
	// PASS 2 — Grow only the largest L-groups
	// -------------------------
	for {
		small, large := getSmallLarge()
		_ = small // small is not used directly here, but kept for clarity.

		// Collect current largest groups.
		var largeIdx []int
		for _, idx := range lIdx {
			if runs[idx].len == large {
				largeIdx = append(largeIdx, idx)
			}
		}
		if len(largeIdx) == 0 {
			break
		}
		if tailLen <= len(largeIdx) {
			break
		}

		nextTail := tailLen - len(largeIdx)
		// Tail must remain >= new largest and >= wordspace.
		if nextTail < large+1 || nextTail < wordspace {
			break
		}

		for _, idx := range largeIdx {
			runs[idx].len++
		}
		tailLen = nextTail
		didBalance = true
	}

	// -------------------------
	// Determine excess code
	// -------------------------
	var excess int
	if !didBalance {
		excess = -1
	} else {
		if tailLen <= wordspace {
			excess = 0
		} else {
			excess = tailLen - wordspace
		}
	}

	// Rebuild body from runs.
	var b strings.Builder
	for _, r := range runs {
		b.WriteString(strings.Repeat(string(r.char), r.len))
	}
	// Append tail.
	b.WriteString(strings.Repeat("L", tailLen))

	return b.String(), excess
}
