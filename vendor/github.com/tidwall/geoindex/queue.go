package geoindex

import "github.com/tidwall/geoindex/child"

// Priority Queue ordered by dist (smallest to largest)

type qnode struct {
	dist  float64
	child child.Child
}

type queue []qnode

func (q *queue) push(node qnode) {
	*q = append(*q, node)
	nodes := *q
	i := len(nodes) - 1
	parent := (i - 1) / 2
	for ; i != 0 && nodes[parent].dist > nodes[i].dist; parent = (i - 1) / 2 {
		nodes[parent], nodes[i] = nodes[i], nodes[parent]
		i = parent
	}
}

func (q *queue) pop() (qnode, bool) {
	nodes := *q
	if len(nodes) == 0 {
		return qnode{}, false
	}
	var n qnode
	n, nodes[0] = nodes[0], nodes[len(*q)-1]
	nodes = nodes[:len(nodes)-1]
	*q = nodes

	i := 0
	for {
		smallest := i
		left := i*2 + 1
		right := i*2 + 2
		if left < len(nodes) && nodes[left].dist <= nodes[smallest].dist {
			smallest = left
		}
		if right < len(nodes) && nodes[right].dist <= nodes[smallest].dist {
			smallest = right
		}
		if smallest == i {
			break
		}
		nodes[smallest], nodes[i] = nodes[i], nodes[smallest]
		i = smallest
	}
	return n, true
}
