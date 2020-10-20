package main

var pluginName string = "conversion_four"

// Converter is an implementation of plugin's interface
type Converter struct{}

func (Converter) Convert(in float64) (out float64, err error) {
	// Celsius to Fahrenheit
	out = in*9/5 + 32
	return
}

func main() {
	panic("blah blah blah") // OK, die straight away
}
