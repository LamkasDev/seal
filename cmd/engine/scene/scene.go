package scene

import (
	"github.com/LamkasDev/seal/cmd/engine/entity"
	sealEntity "github.com/LamkasDev/seal/cmd/engine/entity"
	sealMesh "github.com/LamkasDev/seal/cmd/engine/vulkan/mesh"
	"github.com/LamkasDev/seal/cmd/engine/vulkan/renderer"
	sealTransform "github.com/LamkasDev/seal/cmd/engine/vulkan/transform"
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

func SpawnSceneModel(scene *Scene, mesh *sealMesh.VulkanMesh, transform sealTransform.VulkanTransform) (entity.Entity, error) {
	centity, err := entity.NewEntity(transform)
	if err != nil {
		return centity, err
	}
	component, err := entity.NewEntityComponentMesh(&centity, mesh)
	if err != nil {
		return centity, err
	}
	centity.Components = append(centity.Components, component)
	scene.Entities = append(scene.Entities, centity)

	return centity, nil
}

func UpdateScene(scene *Scene) error {
	for _, entity := range scene.Entities {
		if err := sealEntity.UpdateEntity(&entity); err != nil {
			return err
		}
	}

	return nil
}

func RenderScene(scene *Scene) error {
	if err := renderer.BeginRendererFrame(renderer.RendererInstance); err != nil {
		return err
	}
	for _, entity := range scene.Entities {
		if err := sealEntity.RenderEntity(&entity); err != nil {
			return err
		}
	}
	if err := renderer.EndRendererFrame(renderer.RendererInstance); err != nil {
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
