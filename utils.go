package main

import (
	"image/color"

	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

func Map[T any,K any] (arr []T, f func(T) K) []K {
	out := make([]K, len(arr))
	for i, elem := range arr {
		out[i] = f(elem)
	}
	return out
}

func Rainbow(base lipgloss.Style, s string, colors []color.Color) string {
	var str string
	for i, ss := range s {
		color, _ := colorful.MakeColor(colors[i%len(colors)])
		str = str + base.Foreground(lipgloss.Color(color.Hex())).Render(string(ss))
	}
	return str
}
func Clamp(low, input, high int) int {
	if input < low {
		return low
	}
	if input > high {
		return high
	}
	return input
}
