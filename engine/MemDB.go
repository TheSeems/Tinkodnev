package engine

type MemDB interface {
	Get(id uint64) (Member, error)
	Add(member Member) (uint64, error)
	Init(data string)
	Search(data string, limit int) ([]Member, error)
}
