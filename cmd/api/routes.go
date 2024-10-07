package api

func (s *Server) MountHandlers() {

	// Mount all handlers here
	s.Router.Post("/", Broker)
	s.Router.Post("/handle", s.HandleSubmission)

}
