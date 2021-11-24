# Go Factory Method Generator
[![Go](https://github.com/ryan-holcombe/fmgen/actions/workflows/go.yml/badge.svg)](https://github.com/ryan-holcombe/fmgen/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/ryan-holcombe/fmgen/branch/main/graph/badge.svg?token=083O6ONW1P)](https://codecov.io/gh/ryan-holcombe/fmgen)
[![Go Report Card](https://goreportcard.com/badge/github.com/ryan-holcombe/fmgen)](https://goreportcard.com/report/github.com/ryan-holcombe/fmgen)

A Go factory method generator. Parses all packages to find all `struct` signatures. Then generates a `fm_gen.go` file for each package.

### Example
```
// Sample demo struct
type Sample struct {
    ID          int64 `fmgen:"-"`
    Name        string
    Age         int64 `fmgen:"optional"`
    LastUpdated time.Time
}
```

### Output
```
// NewSample generated factory method for Sample
func NewSample(Name string, LastUpdated time.Time, Age *int64) *Sample {
    result := &Sample{
        Name:        Name,
        LastUpdated: LastUpdated,
    }
    if Age != nil {
        result.Age = *Age
    }
    return result
}
```

### Flags
`-d` will specify the directory to parse (defaults to ./)

`-r` recursively process directories (defaults to true)

`-f` overrides `-d`, will process a single file instead

`-s` comma separated list of structs to include. all others will be ignored

`-v` verbose mode, will include additional logging

### Example
This will search the directory recursively and only process `Struct1`
```
fmgen -d ./pkg -s Struct1
```

### Options
Adding `fmgen:-` to a struct comment will exclude that struct from generation
```
// Sample demo struct, fmgen:-
type Sample struct {
    ...
}
```

Adding `fmgen:-` to a struct field will omit that field
```
type Sample struct {
    ID int64 `fmgen:"-"`
    ...
}
```

Adding `fmgen:optional` to a struct field will allow a `nil` to be passed in. A `nil` will be skipped when creating the struct
```
type Sample struct {
    ID int64 `fmgen:"optional"`
    ...
}
```