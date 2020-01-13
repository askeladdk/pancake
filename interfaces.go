package pancake

type Window interface {
	ShouldClose() bool
	Size() (int, int)
	Update()
}

type WindowOptions struct {
	Title  string
	Width  int
	Height int
}
