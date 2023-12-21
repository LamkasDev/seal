package scene

import (
	"github.com/LamkasDev/seal/cmd/engine/entity"
	sealEntity "github.com/LamkasDev/seal/cmd/engine/entity"
	sealComponent "github.com/LamkasDev/seal/cmd/engine/entity/component"
	"github.com/LamkasDev/seal/cmd/engine/renderer"
)

type Scene struct {
	Entities []sealEntity.Entity
}

func NewScene() (Scene, error) {
	scene := Scene{
		Entities: []sealEntity.Entity{},
	}

	return scene, nil
}

func SpawnSceneModel(scene *Scene) error {
	centity, err := entity.NewEntity()
	if err != nil {
		return err
	}
	component, err := sealComponent.NewEntityComponent(renderer.RendererInstance.Pipeline.Mesh)
	if err != nil {
		return err
	}
	centity.Components = append(centity.Components, component)
	scene.Entities = append(scene.Entities, centity)

	return nil
}

func RenderScene(scene *Scene) error {
	if err := renderer.BeginRendererFrame(&renderer.RendererInstance); err != nil {
		return err
	}
	for _, entity := range scene.Entities {
		if err := sealEntity.RenderEntity(&entity); err != nil {
			return err
		}
	}
	if err := renderer.EndRendererFrame(&renderer.RendererInstance); err != nil {
		return err
	}

	return nil
}

func FreeScene(scene *Scene) error {
	for _, entity := range scene.Entities {
		if err := sealEntity.FreeEntity(&entity); err != nil {
			return err
		}
	}

	return nil
}
