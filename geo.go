package mmdbgeo

func InitDB() {
	errs := make(chan error)

	go func() {
		defer close(errs)
		checkDBs(errs)
	}()

	for err := range errs {
		if err != nil {
			panic(err)
		}
	}
}
