package main

// check if its a std wedge seg count
func mkStd(segs int) string {
	_, exists := stdWedges[segs]
	if exists {
		return "*,"
	} else {
		return ", "
	}
}
