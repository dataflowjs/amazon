package main

import "google.golang.org/api/sheets/v4"

var preUngatedColor = &sheets.Color{
	Blue:  142,
	Green: 124,
	Red:   195,
}

var automaticUngatingColor = &sheets.Color{
	Blue:  147,
	Green: 196,
	Red:   125,
}

var additionalTestColor = &sheets.Color{
	Blue:  56,
	Green: 118,
	Red:   29,
}

var documentRequiredColor = &sheets.Color{
	Blue:  255,
	Green: 229,
	Red:   153,
}

var notQualifiedColor = &sheets.Color{
	Blue:  234,
	Green: 153,
	Red:   153,
}
