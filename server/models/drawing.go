package models

type Drawing struct {
	BackgroundColor *Color `json:"background_color"`
	Lines           []Line `json:"lines"`
}

type Color struct {
	Alpha float64 `json:"alpha"`
	Red   float64 `json:"red"`
	Green float64 `json:"green"`
	Blue  float64 `json:"blue"`
}

type Line struct {
	Color  *Color  `json:"color"`
	Size   float64 `json:"size"`
	Points []Point `json:"points"`
}

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
