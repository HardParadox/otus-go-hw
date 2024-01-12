package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	head   *ListItem
	tail   *ListItem
	length int
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	currenItem := &ListItem{
		Value: v,
	}

	if l.length > 0 {
		currenItem.Next = l.head
		l.head.Prev = currenItem
	} else {
		l.tail = currenItem
	}

	l.head = currenItem

	l.length++

	return l.head
}

func (l *list) PushBack(v interface{}) *ListItem {
	currenItem := &ListItem{
		Value: v,
	}

	if l.length > 0 {
		l.tail.Next = currenItem
		currenItem.Prev = l.tail
	} else {
		l.head = currenItem
	}

	l.tail = currenItem

	l.length++

	return l.tail
}

func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.head = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.tail = i.Prev
	}

	i.Next = nil
	i.Prev = nil
	i.Value = nil
	l.length--
}

func (l *list) MoveToFront(i *ListItem) {
	if i.Prev == nil || l.head == i {
		return
	}

	i.Prev.Next = i.Next

	if i.Next != nil {
		i.Next.Prev = i.Prev
	}

	i.Prev = nil
	i.Next = l.head
	l.head.Prev = i
	l.head = i
}
