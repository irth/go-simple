package simple

import "fmt"

var _ Widget = &ButtonWidget{}

type ClickHandler func(a *App, b *ButtonWidget) error

type ButtonWidget struct {
	ID      string
	Name    string
	OnClick ClickHandler
	Position
}

func Button(id string, pos Position, name string, onClick ClickHandler) *ButtonWidget {
	return &ButtonWidget{id, name, onClick, pos}
}

func (b *ButtonWidget) Render() (string, error) {
	return WidgetCommand{
		Name:     "button",
		ID:       b.ID,
		Position: b.Position,
		Extra:    b.Name,
	}.Render()
}

func (b *ButtonWidget) Update(e Event) ([]BoundEventHandler, error) {
	fmt.Println(e)
	ie, ok := e.(SelectedEvent)
	if !ok || b.OnClick == nil || ie.ID != b.ID {
		return nil, nil
	}

	return []BoundEventHandler{
		func(a *App) error {
			return b.OnClick(a, b)
		},
	}, nil
}
