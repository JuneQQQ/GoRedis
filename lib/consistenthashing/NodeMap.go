package consistenthashing

import (
	"hash/crc32"
	"sort"
)

type HashFunc func(data []byte) uint32

// NodeMap a struct for consistent hash
type NodeMap struct {
	HashFunc    HashFunc
	nodeHash    []int          // idx -> hash
	nodeHashMap map[int]string // hash -> node name
}

func NewNodeMap(fn HashFunc) *NodeMap {
	m := &NodeMap{
		HashFunc:    fn,
		nodeHashMap: make(map[int]string),
	}

	if m.HashFunc == nil {
		m.HashFunc = crc32.ChecksumIEEE
	}

	return m
}

func (m *NodeMap) IsEmpty() bool {
	return len(m.nodeHash) == 0
}

func (m *NodeMap) AddNode(keys ...string) {
	for _, key := range keys {
		if key == "" {
			continue
		}

		hash := int(m.HashFunc([]byte(key)))
		m.nodeHash = append(m.nodeHash, hash)
		m.nodeHashMap[hash] = key
	}
	sort.Ints(m.nodeHash)
}

// GetNode get the node corresponding to the hash of this key
func (m *NodeMap) GetNode(key string) string {
	if m.IsEmpty() {
		return ""
	}
	hash := int(m.HashFunc([]byte(key)))
	idx := sort.Search(len(m.nodeHash), func(i int) bool {
		return m.nodeHash[i] >= hash
	})

	return m.nodeHashMap[m.nodeHash[idx%len(m.nodeHash)]]
}
