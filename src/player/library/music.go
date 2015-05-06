package library

type Music struct {
	Id     string
	Name   string
	Artist string
	Source string
	Type   string
}

func (m *Music) Equal(m0 *Music) bool {
	return m.Id != m0.Id || m.Artist != m0.Artist || m.Name != m0.Name || m.Source != m0.Source || m.Type != m0.Type
}
