package main

// Subject - это структура, определяющая элементы списка учебных предметов.
type Subject struct {
	SubjectTitle string
	Teacher      string
	Hours        float32
}

var SubjectMap = make(map[string]Subject)

func (m SubjectMap) Del(s string) {
	if _, ok := m[s]; ok {
		delete(m, s)
	}
}
