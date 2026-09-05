// pModeAdj performs two‑phase balancing of active‑color runs.
// Returns (adjustedString, excessCode):
//
//	excess == 0   → all internal groups and tail ended balanced
//	excess > 0    → balancing occurred, tail still has leftover
//	excess == -1  → no balancing possible (tail excess too small)
//
// Big‑group rules:
//
//	Active color = last rune.
//
//	ICS == 1:
//	    If active=='O':
//	        All O runs (len>=1) are big groups.
//	    If active=='L':
//	        L runs len>=2 are big groups.
//	        Single‑L runs are big only if in changeCandidates.
//
//	ICS > 1:
//	    Big groups = runs of active color with len >= ICS.
//	    changeCandidates ignored.
//
// Two‑pass balancing:
//
//	PASS 1: Equalize all big groups (small + large).
//	        Tail must remain >= largest internal group.
//	PASS 2: Equalize only the largest internal groups.
//	        Tail must remain >= largest internal group.
package main

func pModeAdj(inStr string, interCharSpaces int, changeCandidates []int) (string, int) {
	if len(inStr) == 0 {
		return inStr, 0
	}

	// Convert changeCandidates to a set.
	candSet := make(map[int]struct{}, len(changeCandidates))
	for _, idx := range changeCandidates {
		candSet[idx] = struct{}{}
	}

	// Parse runs.
	type run struct {
		char rune
		len  int
	}

	var runs []run
	var last rune
	count := 0

	for _, r := range inStr {
		if count == 0 {
			last = r
			count = 1
			continue
		}
		if r == last {
			count++
		} else {
			runs = append(runs, run{char: last, len: count})
			last = r
			count = 1
		}
	}
	if count > 0 {
		runs = append(runs, run{char: last, len: count})
	}

	// Active color.
	active := runs[len(runs)-1].char
	tail := len(runs) - 1

	// Identify big groups.
	big := make([]bool, len(runs))

	if interCharSpaces == 1 {
		if active == 'O' {
			for i := range runs {
				if runs[i].char == 'O' {
					big[i] = true
				}
			}
		} else {
			for i := range runs {
				if runs[i].char != 'L' {
					continue
				}
				if runs[i].len >= 2 {
					big[i] = true
				} else if _, ok := candSet[i]; ok {
					big[i] = true
				}
			}
		}
	} else {
		for i := range runs {
			if runs[i].char == active && runs[i].len >= interCharSpaces {
				big[i] = true
			}
		}
	}

	// Collect internal big groups.
	var bigIdx []int
	for i := 0; i < len(runs)-1; i++ {
		if big[i] && runs[i].char == active {
			bigIdx = append(bigIdx, i)
		}
	}

	if len(bigIdx) == 0 {
		return inStr, 0
	}

	didBalance := false

	// -------------------------
	// PASS 1 — Equalize all big groups
	// -------------------------
	for {
		if runs[tail].len <= len(bigIdx) {
			break
		}

		// largest internal group
		largest := runs[bigIdx[0]].len
		for _, idx := range bigIdx[1:] {
			if runs[idx].len > largest {
				largest = runs[idx].len
			}
		}

		nextTail := runs[tail].len - len(bigIdx)
		if nextTail < largest {
			break
		}

		// Perform one full round.
		for _, idx := range bigIdx {
			runs[idx].len++
		}
		runs[tail].len -= len(bigIdx)
		didBalance = true
	}

	// -------------------------
	// PASS 2 — Equalize only the largest internal groups
	// -------------------------
	// Identify largest internal size.
	largest := runs[bigIdx[0]].len
	for _, idx := range bigIdx[1:] {
		if runs[idx].len > largest {
			largest = runs[idx].len
		}
	}

	// Collect only the largest groups.
	var largeIdx []int
	for _, idx := range bigIdx {
		if runs[idx].len == largest {
			largeIdx = append(largeIdx, idx)
		}
	}

	if len(largeIdx) > 0 {
		for {
			if runs[tail].len <= len(largeIdx) {
				break
			}

			nextTail := runs[tail].len - len(largeIdx)
			if nextTail < largest {
				break
			}

			// Perform one full round on largest groups only.
			for _, idx := range largeIdx {
				runs[idx].len++
			}
			runs[tail].len -= len(largeIdx)
			didBalance = true

			// Update largest size.
			largest++
		}
	}

	// -------------------------
	// Determine excess code
	// -------------------------
	var excess int
	if !didBalance {
		excess = -1
	} else {
		// smallest internal group:
		smallest := runs[bigIdx[0]].len
		for _, idx := range bigIdx[1:] {
			if runs[idx].len < smallest {
				smallest = runs[idx].len
			}
		}
		if runs[tail].len == smallest {
			excess = 0
		} else {
			excess = runs[tail].len - smallest
		}
	}

	// Rebuild output.
	var out []rune
	for _, r := range runs {
		for i := 0; i < r.len; i++ {
			out = append(out, r.char)
		}
	}

	return string(out), excess
}
