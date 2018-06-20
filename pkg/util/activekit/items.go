package activekit

type MenuItems []*MenuItem

func (items MenuItems) Labels() []string {
	labels := make([]string, 0, len(items))
	for _, item := range items {
		labels = append(labels, item.Label)
	}
	return labels
}

func (items MenuItems) Copy() MenuItems {
	return append(MenuItems{}, items...)
}

func (items MenuItems) Delete(i int) MenuItems {
	cp := items.Copy()
	return append(cp[:i], cp[i+1:]...)
}

func (items MenuItems) Append(newItems ...*MenuItem) MenuItems {
	return append(items.Copy(), newItems...)
}

func (items MenuItems) Len() int {
	return len(items)
}

func (items MenuItems) NotNil() MenuItems {
	cp := items.Copy()
	for i, item := range cp {
		if item == nil {
			cp = append(cp[:i], cp[i+1:]...)
		}
	}
	return cp
}

func StringSelector(items []string, action func(string) error) MenuItems {
	var menuItems = make(MenuItems, 0, len(items))
	for _, item := range items {
		menuItems = menuItems.Append(&MenuItem{
			Label: item,
			Action: func(item string) func() error {
				return func() error { return action(item) }
			}(item),
		})
	}
	return menuItems
}

func HandleString(promt string, handler func(string) bool) func() error {
	return func() error {
		for {
			var str = Promt(promt)
			if !handler(str) {
				break
			}
		}
		return nil
	}
}
