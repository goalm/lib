package esg

// refer to yaml ?
func init() {
	var CurveData = make([][]float64, 10)
	for i := range CurveData {
		CurveData[i] = make([]float64, 1200)
	}
}
