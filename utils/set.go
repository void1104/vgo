package utils

type Empty struct{}

var empty Empty

type Set struct {
	m map[string]Empty
}

// Add 添加元素
func (s *Set) Add(key string) {
	s.m[key] = empty
}

// Remove 删除元素
func (s *Set) Remove(key string) {
	delete(s.m, key)
}

// Exist 检查元素是否存在
func (s *Set) Exist(key string) bool {
	if _, ok := s.m[key]; ok {
		return true
	}
	return false
}

// Travel 遍历set
func (s *Set) Travel() (keys []string) {
	for key, _ := range s.m {
		keys = append(keys, key)
	}
	return
}
