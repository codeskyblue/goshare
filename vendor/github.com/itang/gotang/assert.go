package gotang

func AssertNoError(err error, message string) {
	if err != nil {
		Assert(false, message+" ("+err.Error()+")")
	}
}

func Assert(assertion bool, message string) {
	if !assertion {
		err := "assertion failed"
		if message != "" {
			err += ": " + message
		}
		panic(err)
	}
}
