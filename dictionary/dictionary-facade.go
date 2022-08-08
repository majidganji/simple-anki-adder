package dictionary

type DictionaryFacade struct {
	Ldoceonline Ldoceonline
	Fastic      Fastic
}

func NewDictionaryFacade(world *World) *DictionaryFacade {
	return &DictionaryFacade{
		Ldoceonline: *NewLdoceonline(world),
		Fastic:      *NewFastic(world),
	}
}
