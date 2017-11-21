package test

// didPanic returns true if the function passed to it panics. Otherwise, it returns false.
func didPanic(panicTestFunc func()) (bool, string) {

	didPanic := false
	var message interface{}
	func() {

		defer func() {
			if message = recover(); message != nil {
				didPanic = true
			}
		}()

		// call the target function
		panicTestFunc()

	}()
	if message == nil {
		message = ""
	}
	return didPanic, message.(string)
}