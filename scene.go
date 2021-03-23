package simple

import (
	"fmt"
)

type Scene interface {
	Render() (Widget, error)
}

func (a *App) runScene(s Scene) error {
	widgets, err := s.Render()
	if err != nil {
		return fmt.Errorf("while rendering the scene: %w", err)
	}

	sas, err := widgets.Render()
	if err != nil {
		return fmt.Errorf("while rendering a widget: %w", err)
	}

	a.widgets = widgets

	a.sr.run([]byte(sas))
	return nil
}
