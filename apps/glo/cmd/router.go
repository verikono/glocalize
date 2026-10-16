func GetRoutedCommandMethod(route []string) func {

	source := route[0]
	path := route[1:]

	
	switch source {
	case "parse":
		return func() {
}