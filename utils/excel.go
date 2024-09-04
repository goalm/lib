package utils

// 1 -> A, 2 -> B, 3 -> C, ..., 26 -> Z, 27 -> AA, 28 -> AB, ...
func IntToClmTitle(n int) string {
	res := ""

	for n > 0 {
		n--
		res = string(byte(n%26)+'A') + res
		n /= 26
	}

	return res
}
