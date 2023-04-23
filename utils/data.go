package utils

func RemoveDup(this []*float64) []*float64 {
	tgt := this[:1]
	for _, v := range this {
		tgt = append(tgt, v)
	}
	return tgt
}
