package error

// check err
func Check(e error) {
	if e != nil {
		panic(e)
	}
}
