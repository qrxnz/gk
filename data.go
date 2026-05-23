package main

func (b *Board) initLists() {
	b.cols = []column{
		newColumn(todo),
		newColumn(inProgress),
		newColumn(done),
	}
	b.cols[todo].list.Title = "To Do"
	b.cols[inProgress].list.Title = "In Progress"
	b.cols[done].list.Title = "Done"

	items, err := b.storage.Load()
	if err != nil {
		b.err = err
		return
	}

	b.cols[todo].list.SetItems(items[todo])
	b.cols[inProgress].list.SetItems(items[inProgress])
	b.cols[done].list.SetItems(items[done])

	b.habits, err = b.storage.LoadHabits()
	if err != nil {
		b.err = err
		return
	}
}
