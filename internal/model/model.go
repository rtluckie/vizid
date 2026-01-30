package model

type Components struct {
	Year   bool
	Month  bool
	Day    bool
	Hour   bool
	Minute bool
	Second bool
	Ms     bool
	UUID   bool
}

type Options struct {
	Timezone   string
	Warn       bool
	Custom     bool
	Components Components
}
