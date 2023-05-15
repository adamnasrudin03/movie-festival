package repository

// Repositories all repo object injected here
type Repositories struct {
	User  UserRepository
	Log   LogRepository
	Movie MovieRepository
}
