package util

type Keyer interface {
	GetKey() int64
}

type LinkedHashMap struct {
	keys   []int64
	values map[int64]Keyer
}

func NewLinkedHashMap() *LinkedHashMap {
	return &LinkedHashMap{
		keys:   make([]int64, 0),
		values: make(map[int64]Keyer),
	}
}

func (lhm *LinkedHashMap) Put(value Keyer) {
	key := value.GetKey()
	if _, ok := lhm.values[key]; !ok {
		lhm.keys = append(lhm.keys, key)
	}
	lhm.values[key] = value
}

func (lhm *LinkedHashMap) Get(key int64) (Keyer, bool) {
	value, ok := lhm.values[key]
	return value, ok
}

func (lhm *LinkedHashMap) Keys() []int64 {
	return lhm.keys
}
