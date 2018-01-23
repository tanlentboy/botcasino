package expiretimer

// 过期信息
type expire struct {
	ID        uint64 // 红包ID
	Timestamp int64  // 时间戳
}

type expireHeap []expire

func (h expireHeap) Len() int { return len(h) }
func (h expireHeap) Less(i, j int) bool {
	if h[i].Timestamp == h[j].Timestamp {
		return h[i].ID < h[j].ID
	}
	return h[i].Timestamp < h[j].Timestamp
}
func (h expireHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *expireHeap) Push(x interface{}) {
	*h = append(*h, x.(expire))
}

func (h *expireHeap) Pop() interface{} {
	if h.Len() == 0 {
		return nil
	}
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (h *expireHeap) Front() *expire {
	if h.Len() == 0 {
		return nil
	}
	return &(*h)[0]
}
