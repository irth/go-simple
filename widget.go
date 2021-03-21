package simple

type Widget interface {
	Render() (string, error)
	Update(stdout Event) ([]BoundEventHandler, error)
}
