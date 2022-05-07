package superapp

func (s *SuperApp) Run() error {
	return s.router.Run()
}
