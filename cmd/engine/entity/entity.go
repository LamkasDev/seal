package entity

import "github.com/EngoEngine/glm"

type Entity struct {
	Position   glm.Vec3
	Rotation   float32
	Components []EntityComponent
}

func NewEntity(position glm.Vec3) (Entity, error) {
	entity := Entity{
		Position:   position,
		Rotation:   0,
		Components: []EntityComponent{},
	}

	return entity, nil
}

func RenderEntity(entity *Entity) error {
	for _, component := range entity.Components {
		if err := RenderEntityComponent(&component); err != nil {
			return err
		}
	}

	return nil
}

func FreeEntity(entity *Entity) error {
	for _, component := range entity.Components {
		if err := FreeEntityComponent(&component); err != nil {
			return err
		}
	}

	return nil
}
