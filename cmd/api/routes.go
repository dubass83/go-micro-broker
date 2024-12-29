package api

func (s *Server) MountHandlers() {

	// Mount all handlers here
	s.Router.Post("/", Broker)
	s.Router.Post("/log-grpc", s.LogItemViaGRPC)
	s.Router.Post("/handle", s.HandleSubmission)

}
