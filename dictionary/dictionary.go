package dictionary

type Dictionary interface {
	Search(world string) (*World, error)
}
