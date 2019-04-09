package sortedlist

import "container/list"

//Insert inserts 'slotsNeeded' elements in ascending order from 'e'
func Insert(sortedList *list.List, e *list.Element, value int, slotsNeeded int) {
	if slotsNeeded == 0 {
		return
	}
	switch {
	case e == nil:
		e = sortedList.PushFront(value)
		Insert(sortedList, e, value+1, slotsNeeded-1)
	case e.Value.(int) < value:
		e = sortedList.InsertAfter(value, e)
		Insert(sortedList, e, value+1, slotsNeeded-1)
	default:
		Insert(sortedList, e.Prev(), value, slotsNeeded)
	}
}

//Remove removes 'slotsNeeded' elements in ascending order from 'e'
func Remove(sortedList *list.List, e *list.Element, slotsNeeded int) {
	if slotsNeeded > 1 {
		Remove(sortedList, e.Next(), slotsNeeded-1)
	}
	sortedList.Remove(e)
}

//FindSeq returns the first ascending sequence of 'slotsNeeded' values
func FindSeq(sortedList *list.List, slotsNeeded int) *list.Element {
	if sortedList.Len() < slotsNeeded {
		return nil
	}
	e1 := sortedList.Front()
	eN := e1
	for ii := 1; ii < slotsNeeded; ii++ {
		eN = eN.Next()
	}
	for ; eN != nil; e1, eN = e1.Next(), eN.Next() {
		if eN.Value.(int)-e1.Value.(int) == slotsNeeded-1 {
			return e1
		}
	}
	return nil
}
