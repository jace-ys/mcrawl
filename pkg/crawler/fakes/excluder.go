package fakes

type FakeExcluder struct {
	exclude bool
}

func NewFakeExcluder(exclude bool) *FakeExcluder {
	return &FakeExcluder{
		exclude: exclude,
	}
}

func (f *FakeExcluder) Exclude(path string) bool {
	return f.exclude
}
