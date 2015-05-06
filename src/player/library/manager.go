package library

import (
	"errors"
)

type MusicManager struct {
	musics []Music
}

func NewMusicManager() *MusicManager {
	return &MusicManager{make([]Music, 0)}
}

func (m *MusicManager) Len() int {
	return len(m.musics)
}

func (m *MusicManager) Find(name string) *Music {
	if len(m.musics) == 0 {
		return nil
	}
	for _, m := range m.musics {
		if m.Name == name {
			return &m
		}
	}

	return nil
}

func (m *MusicManager) Add(music *Music) {
	m.musics = append(m.musics, music)
}

func (m *MusicManager) Get(index int) (*Music, error) {
	if index < 0 || index > len(m.musics) {
		return nil
	}

	return &m.musics[index]
}

func (m *MusicManager) Remove(index int) (*Music, error) {
	if index < 0 || index > len(m.musics) {
		return nil
	}

	removedMusic := &m.musics[index]

	//从数组切片中删除元素
	if index < len(m.musics)-1 { //中间元素
		m.musics = append(m.musics[:index-1], m.musics[index+1:]...)
	} else if index == 0 { //删除仅有的一个元素
		m.musics = make([]Music, 0)
	} else { //删除的是最后一个元素
		m.musics = m.musics[:index-1]
	}

	return removedMusic
}
