package service

type MapRepository interface {
	GetObjects()
}

type Map struct {
	repo MapRepository
}

func New(repo MapRepository) *Map {
	return &Map{repo: repo}
}
