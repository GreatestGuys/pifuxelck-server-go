package models

type Drawing struct {
	BackgroundColor *Color `json:"background_color"`
	Lines           []Line `json:"lines"`
}

type Color struct {
	Alpha string `json:"alpha"`
	Red   string `json:"red"`
	Green string `json:"green"`
	Blue  string `json:"blue"`
}

type Line struct {
	Color  *Color  `json:"color"`
	Size   string  `json:"size"`
	Points []Point `json:"points"`
}

type Point struct {
	X string `json:"x"`
	Y string `json:"y"`
}
