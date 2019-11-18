package gober

// ComplexData ..
type ComplexData struct {
	N int
	S string
	M map[string]int
	P []byte
	C *ComplexData
}

const (
	// Port ..
	Port = ":61000"
)
