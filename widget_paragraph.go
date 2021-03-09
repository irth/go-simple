package simple

var _ Widget = &ParagraphWidget{}

type ParagraphWidget struct {
	Text string
	Position
}

func Paragraph(pos Position, text string) *ParagraphWidget {
	return &ParagraphWidget{text, pos}
}

func (b *ParagraphWidget) Render() (string, error) {
	return WidgetCommand{
		Multiline: true,
		Name:      "paragraph",
		Position:  b.Position,
		Extra:     b.Text,
	}.Render()
}

func (b *ParagraphWidget) Update(out Output) ([]BoundEventHandler, error) {
	return nil, nil
}
