package service

// Services all service object injected here
type Services struct {
	User  UserService
	Log   LogService
	Movie MovieService
	Genre GenreService
}
