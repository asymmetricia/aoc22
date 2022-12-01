package aoc

import "fmt"

type PQueue[Node any] struct {
	Head *PQueueNode[Node]
}

type PQueueNode[Node any] struct {
	Node     Node
	Priority int
	Next     *PQueueNode[Node]
}

func (pq *PQueue[Node]) Print() {
	cursor := pq.Head
	for cursor != nil {
		fmt.Println(cursor.Node, " ", cursor.Priority)
		cursor = cursor.Next
	}
}

func (pq *PQueue[Node]) Pop() Node {
	if pq.Head == nil {
		panic("pop on empty pqueue")
	}
	ret := pq.Head.Node
	pq.Head = pq.Head.Next
	return ret
}

func (pq *PQueue[Node]) AddWithPriority(node Node, prio int) {
	newnode := &PQueueNode[Node]{
		Node:     node,
		Priority: prio,
	}

	if pq.Head == nil {
		pq.Head = newnode
		return
	}

	if pq.Head.Priority > prio {
		newnode.Next = pq.Head
		pq.Head = newnode
		return
	}

	cursor := pq.Head
	for {
		if cursor.Next == nil || cursor.Next.Priority > prio {
			break
		}
		cursor = cursor.Next
	}

	newnode.Next, cursor.Next = cursor.Next, newnode
	return
}
