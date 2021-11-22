package testdata

// LineNum struct to assist in testing lineNum
type LineNum struct {
	ID   int64 `fmgen:"-"`
	Name string
	Age  int64 `fmgen:"optional"`
}
