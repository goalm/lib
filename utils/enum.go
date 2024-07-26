package utils

// bi-directional map
type Enum struct {
	idxName string
	VarName string
	fw      map[int]string
	bw      map[string]int
}

func NewEnum() *Enum {
	return &Enum{
		fw: make(map[int]string),
		bw: make(map[string]int),
	}
}

func (m *Enum) SetNames(idxName, varName string) {
	m.idxName = idxName
	m.VarName = varName
}

func (m *Enum) Add(a int, b string) {
	// optionally verify uniqueness constraint
	m.fw[a] = b
	m.bw[b] = a
}

func (m *Enum) IntToStr(a int) (string, bool) {
	b, ok := m.fw[a]
	return b, ok
}

func (m *Enum) StrToInt(b string) (int, bool) {
	a, ok := m.bw[b]
	return a, ok
}
