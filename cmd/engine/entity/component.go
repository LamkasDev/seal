package entity

type EntityComponent struct {
	Entity *Entity
	Data   interface{}
	Update func(component *EntityComponent) error
	Render func(component *EntityComponent) error
}

func UpdateEntityComponent(component *EntityComponent) error {
	return component.Update(component)
}

func RenderEntityComponent(component *EntityComponent) error {
	return component.Render(component)
}

func FreeEntityComponent(component *EntityComponent) error {
	return nil
}
