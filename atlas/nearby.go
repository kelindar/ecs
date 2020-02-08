// Copyright (c) 2018 Josh Baker
// Copyright (c) 2020 Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.
// This is a fork of https://github.com/tidwall/rbang adapted for int32 coordinate system.

package atlas

// Nearby performs a kNN-type operation on the index.
// It's expected that the caller provides its own the `algo` function, which
// is used to calculate a distance to data. The `add` function should be
// called by the caller to "return" the data item along with a distance.
// The `iter` function will return all items from the smallest dist to the
// largest dist.
// Take a look at the SimpleBoxAlgo function for a usage example.
func (tr *Atlas) Nearby(
	algo DistanceFunc,
	iter func(min, max [2]int32, data interface{}, dist int32) bool,
) {
	var q queue
	var parent interface{}
	var children []Child

	for {
		// gather all children for parent
		children = tr.Children(parent, children[:0])
		for _, child := range children {
			q.push(qnode{
				dist:  algo(child.Min, child.Max, child.Data, child.Item),
				child: child,
			})
		}
		for {
			node, ok := q.pop()
			if !ok {
				// nothing left in queue
				return
			}
			if node.child.Item {
				if !iter(node.child.Min, node.child.Max,
					node.child.Data, node.dist) {
					return
				}
			} else {
				// gather more children
				parent = node.child.Data
				break
			}
		}
	}
}

// Priority Queue ordered by dist (smallest to largest)
type qnode struct {
	dist  int32
	child Child
}

type queue struct {
	nodes []qnode
	len   int
	size  int
}

func (q *queue) push(node qnode) {
	if q.nodes == nil {
		q.nodes = make([]qnode, 2)
	} else {
		q.nodes = append(q.nodes, qnode{})
	}
	i := q.len + 1
	j := i / 2
	for i > 1 && q.nodes[j].dist > node.dist {
		q.nodes[i] = q.nodes[j]
		i = j
		j = j / 2
	}
	q.nodes[i] = node
	q.len++
}

func (q *queue) pop() (qnode, bool) {
	if q.len == 0 {
		return qnode{}, false
	}
	n := q.nodes[1]
	q.nodes[1] = q.nodes[q.len]
	q.len--
	var j, k int
	i := 1
	for i != q.len+1 {
		k = q.len + 1
		j = 2 * i
		if j <= q.len && q.nodes[j].dist < q.nodes[k].dist {
			k = j
		}
		if j+1 <= q.len && q.nodes[j+1].dist < q.nodes[k].dist {
			k = j + 1
		}
		q.nodes[i] = q.nodes[k]
		i = k
	}
	return n, true
}
