package components

type Dropdown struct {
	SelectedOption DropdownOption
	AllOptions     []DropdownOption
	InputName      string
}

type DropdownOption struct {
	Name  string
	Link  string
	Value string
}
