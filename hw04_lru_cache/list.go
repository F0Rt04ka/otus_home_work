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
	Length    int
	FirstItem *ListItem
	LastItem  *ListItem
}

func (l *list) addItem(v interface{}) *ListItem {
	newItem := new(ListItem)
	newItem.Value = v

	l.Length++
	if l.FirstItem == nil {
		l.FirstItem = newItem
		l.LastItem = newItem
	}

	return newItem
}

func (l *list) Len() int {
	return l.Length
}

func (l *list) Front() *ListItem {
	return l.FirstItem
}

func (l *list) Back() *ListItem {
	return l.LastItem
}

func (l *list) PushFront(v interface{}) *ListItem {
	newItem := l.addItem(v)

	if l.FirstItem == newItem {
		return newItem
	}

	newItem.Prev = nil
	newItem.Next = l.FirstItem
	newItem.Next.Prev = newItem
	l.FirstItem = newItem

	return newItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newItem := l.addItem(v)

	if l.LastItem == newItem {
		return newItem
	}

	newItem.Next = nil
	newItem.Prev = l.LastItem
	newItem.Prev.Next = newItem
	l.LastItem = newItem

	return newItem
}

func (l *list) Remove(i *ListItem) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}

	switch i {
	case l.FirstItem:
		l.FirstItem = i.Next
	case l.LastItem:
		l.LastItem = i.Prev
	}

	l.Length--
}

func (l *list) MoveToFront(i *ListItem) {
	if i.Prev == nil {
		return
	}

	i.Prev.Next = i.Next

	if i.Next == nil {
		l.LastItem = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}

	i.Next = l.FirstItem
	i.Next.Prev = i
	l.FirstItem = i
}

func NewList() List {
	return new(list)
}
