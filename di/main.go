package di

func SetupDependencies() error {
	if err := SetupFirestore(); err != nil {
		return err
	}

	return nil
}
