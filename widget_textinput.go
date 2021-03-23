package simple

var _ Widget = &TextInputWidget{}

type UpdateHandler func(a *App, t *TextInputWidget, value string) error

type TextInputWidget struct {
	ID       string
	Value    string
	OnUpdate UpdateHandler
	Position
}

func TextInput(id string, pos Position, value string, onUpdate UpdateHandler) *TextInputWidget {
	return &TextInputWidget{id, value, onUpdate, pos}
}

func (t *TextInputWidget) Render() (string, error) {
	return WidgetCommand{
		Name:     "textinput",
		ID:       t.ID,
		Position: t.Position,
		Extra:    t.Value,
	}.Render()
}

func (t *TextInputWidget) Update(e Event) ([]BoundEventHandler, error) {
	ie, ok := e.(InputEvent)
	if !ok || t.OnUpdate == nil || ie.ID != t.ID {
		return nil, nil
	}

	return []BoundEventHandler{
		func(a *App) error {
			return t.OnUpdate(a, t, ie.Value)
		},
	}, nil
}
