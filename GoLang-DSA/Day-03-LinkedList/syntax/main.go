package main

type Node struct {
	Data int
	next  *Node
}

type LinkedList struct {
	head *Node
	size int
}

func (ll *LinkedList) InsertAtEnd(data int) {
	newNode := &Node{Data: data}
	if ll.head == nil {
		ll.head = newNode
	} else {
		current := ll.head
		for current.next != nil {
			current = current.next
		}
		current.next = newNode
	}
	ll.size++
}

func (ll *LinkedList) InsertAtBeginning(data int) {
	newNode := &Node{Data: data, next: ll.head}
	ll.head = newNode
	ll.size++
}

func (ll *LinkedList) InsertAtPosition(data int, position int) {
	if position < 0 || position > ll.size {
		println("Invalid position")
		return
	}
	newNode := &Node{Data: data}
	if position == 0 {
		newNode.next = ll.head
		ll.head = newNode
	} else {
		current := ll.head
		for i := 0; i < position-1; i++ {
			current = current.next
		}
		newNode.next = current.next
		current.next = newNode
	}
	ll.size++
}

func (ll *LinkedList) DeleteAtPosition(position int) {
	if position < 0 || position >= ll.size {
		println("Invalid position")
		return
	}
	if position == 0 {
		ll.head = ll.head.next
	} else {
		current := ll.head
		for i := 0; i < position-1; i++ {
			current = current.next
		}
		current.next = current.next.next
	}
	ll.size--
}

func (ll *LinkedList) DeleteAtEnd() {
	if ll.head == nil {
		return
	}
	if ll.head.next == nil {
		ll.head = nil
	} else {
		current := ll.head
		for current.next.next != nil {
			current = current.next
		}
		current.next = nil
	}
	ll.size--
}

func (ll *LinkedList) Size() int {
	return ll.size
}

func (ll *LinkedList) Search(data int) int {
	current := ll.head
	position := 0
	for current != nil {
		if current.Data == data {
			return position
		}
		current = current.next
		position++
	}
	return -1 // not found
}

func (ll *LinkedList) Display() {
	current := ll.head
	for current != nil {
		print(current.Data, " -> ")
		current = current.next
	}
	println("nil")
}

func main() {
	ll := &LinkedList{}
	ll.InsertAtEnd(10) // Time Complexity: O(1) if we maintain a tail pointer, otherwise O(n)
	ll.InsertAtEnd(20)
	ll.InsertAtEnd(30)
	ll.Display() // Time Complexity: O(n)
	ll.InsertAtBeginning(5) // Time Complexity: O(1)
	ll.Display()
	ll.InsertAtPosition(15, 2) // Time Complexity: O(n) in the worst case (inserting at the end), O(1) if inserting at the beginning
	ll.Display()
	ll.DeleteAtPosition(1) // Time Complexity: O(n) in the worst case (deleting the last node), O(1) if deleting the first node
	ll.Display()
	ll.DeleteAtEnd() // Time Complexity: O(n) if we don't maintain a tail pointer, otherwise O(1)
	ll.Display()
	println("Size of linked list:", ll.Size()) // Time Complexity: O(1) if we maintain a size variable, otherwise O(n) to count the nodes
	println("Position of 15:", ll.Search(15)) // Time Complexity: O(n) in the worst case (searching for the last node), O(1) if the element is at the head
	println("Position of 100:", ll.Search(100))
}