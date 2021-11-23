package testdata

// Array struct to help test arrays
type Array struct {
	String            []string
	StringPtr         []*string
	StringOptionalPtr []*string `fmgen:"optional"`
}
