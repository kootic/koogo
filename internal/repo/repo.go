package repo

// Repositories holds all repository implementations.
type Repositories struct {
	DB DatabaseConnection

	User   KooUserRepository
	Pet    KooPetRepository
	Health HealthRepository
}

type DatabaseConnection interface {
	Close() error
}

func (r *Repositories) Close() error {
	return r.DB.Close()
}
