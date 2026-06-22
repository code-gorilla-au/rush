package games

func NewService(store Store) Service {
	return Service{
		Store: store,
	}
}
