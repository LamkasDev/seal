package entity

import (
	sealComponent "github.com/LamkasDev/seal/cmd/engine/entity/component"
)

type Entity struct {
	Components []sealComponent.EntityComponent
}

func NewEntity() (Entity, error) {
	entity := Entity{
		Components: []sealComponent.EntityComponent{},
	}

	return entity, nil
}

func RenderEntity(entity *Entity) error {
	for _, component := range entity.Components {
		if err := sealComponent.RenderEntityComponent(&component); err != nil {
			return err
		}
	}

	return nil
}

func FreeEntity(entity *Entity) error {
	for _, component := range entity.Components {
		if err := sealComponent.FreeEntityComponent(&component); err != nil {
			return err
		}
	}

	return nil
}
