package components

type Navbar struct {
	ShouldShowMenu   bool
	ProjectsDropdown Dropdown
	DateRange        DateRange
	SettingsButton   Button
}

func NewNavbar(shouldShowMenu bool) Navbar {
	return Navbar{
		ShouldShowMenu: shouldShowMenu,
		SettingsButton: Button{Text: "Settings", Icon: "gear", Link: "/settings", IsSubmit: false},
	}
}
