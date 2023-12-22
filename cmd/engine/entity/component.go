package entity

type EntityComponent struct {
	Entity *Entity
	Data   EntityComponentMeshData
	Render func(component *EntityComponent) error
}

func RenderEntityComponent(component *EntityComponent) error {
	return component.Render(component)
}

func FreeEntityComponent(component *EntityComponent) error {
	return nil
}
