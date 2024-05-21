package components

type Dropdown struct {
	SelectedOption DropdownOption
	AllOptions     []DropdownOption
}

type DropdownOption struct {
	Name string
	Link string
}
