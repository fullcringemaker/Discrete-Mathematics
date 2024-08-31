package main

func add(a, b []int32, p int) []int32 {
	var res []int32
	var q int32 = 0
	var i int

	maxLen := len(a)
	if len(b) > maxLen {
		maxLen = len(b)
	}
	res = make([]int32, 0, maxLen)

	for i = 0; i < maxLen; i++ {
		var aVal, bVal int32
		if i < len(a) {
			aVal = a[i]
		}
		if i < len(b) {
			bVal = b[i]
		}
		sum := aVal + bVal + q
		current := sum % int32(p)
		q = sum / int32(p)
		res = append(res, current)
	}

	if q > 0 {
		res = append(res, q)
	}

	return res
}

func main() {
	
}
