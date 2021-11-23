package testdata

type iface interface {
	Run()
}

type impl struct {
	i iface
}
