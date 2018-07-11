package activekit

func ItemsFromIter(maxIndex uint, next func(index uint) *MenuItem) MenuItems {
	var items = make(MenuItems, 0, maxIndex)
	for i := uint(0); i < maxIndex; i++ {
		var item = next(i)
		if item != nil {
			items = append(items, item)
		}
	}
	return items
}
