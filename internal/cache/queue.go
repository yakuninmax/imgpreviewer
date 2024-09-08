package cache

type queue struct {
	size  int64 // current cache size in bytes
	front *item
	back  *item
}
type item struct {
	file file
	next *item
	prev *item
}

// Create new list.
func newQueue() *queue {
	return new(queue)
}

// Get first element.
func (q *queue) getFront() *item {
	return q.front
}

// Get last element.
func (q *queue) getBack() *item {
	return q.back
}

// Insert item to first position.
func (q *queue) pushFront(file file) *item {
	// New item.
	item := new(item)
	item.file = file

	// Check if list is empty.
	if q.size == 0 {
		q.back = item
	} else {
		item.next = q.front
		q.front.prev = item
	}

	// New item first.
	q.front = item

	q.size = q.size + file.size

	return item
}

// Remove item.
func (q *queue) remove(item *item) {
	// Get prev and next items of removing item.
	prev := item.prev
	next := item.next

	switch {
	// If item only.
	case prev == nil && next == nil:
		q.front = nil
		q.back = nil
		item = nil

	// If item is first.
	case prev == nil:
		next.prev = nil
		q.front = next

	// If item is last.
	case next == nil:
		prev.next = nil
		q.back = prev

	default:
		next.prev = prev
		prev.next = next
	}

	q.size = q.size - item.file.size
}

// Move item to front.
func (q *queue) moveToFront(item *item) {
	// If item is last.
	if item.next == nil {
		q.back = item.prev
	}

	q.pushFront(item.file)
	q.remove(item)
}
