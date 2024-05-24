package components

type Navbar struct {
	ShouldShowActions bool
	ProjectsDropdown  Dropdown
	DateRange         DateRange
	SettingsButton    Button
}

func NewNavbar(ShouldShowActions bool) Navbar {
	return Navbar{
		ShouldShowActions: ShouldShowActions,
		SettingsButton:    Button{Text: "Settings", Icon: "gear", Link: "/settings", IsSubmit: false},
	}
}
