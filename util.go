package main

func dedupInt(v []int) []int {
	ex := map[int]bool{}
	r := v[:0]
	for _, a := range v {
		if ex[a] {
			continue
		}
		r = append(r, a)
		ex[a] = true
	}
	return r
}
