package entity

import sealTransform "github.com/LamkasDev/seal/cmd/engine/vulkan/transform"

type Entity struct {
	Layer      EntityLayer
	Transform  sealTransform.VulkanTransform
	Components []EntityComponent
}

func NewEntity(transform sealTransform.VulkanTransform) (Entity, error) {
	entity := Entity{
		Layer:      LAYER_DEFAULT,
		Transform:  transform,
		Components: []EntityComponent{},
	}

	return entity, nil
}

func UpdateEntity(entity *Entity) error {
	for _, component := range entity.Components {
		if err := UpdateEntityComponent(&component); err != nil {
			return err
		}
	}

	return nil
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
