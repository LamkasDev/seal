package uniform

import (
	"unsafe"

	"github.com/EngoEngine/glm"
	"github.com/vulkan-go/vulkan"
)

type VulkanUniform struct {
	Model      glm.Mat4
	View       glm.Mat4
	Projection glm.Mat4
}

const VulkanUniformSize = unsafe.Sizeof(VulkanUniform{})
const VulkanUniformModelOffset = unsafe.Offsetof(VulkanUniform{}.Model)
const VulkanUniformViewOffset = unsafe.Offsetof(VulkanUniform{}.View)
const VulkanUniformProjectionOffset = unsafe.Offsetof(VulkanUniform{}.Projection)

func NewVulkanUniform3D(extent vulkan.Extent2D, camera glm.Vec3, position glm.Vec3, rotation glm.Vec3) VulkanUniform {
	uniform := VulkanUniform{
		View:       glm.Translate3D(-camera.X(), -camera.Y(), -camera.Z()),
		Projection: glm.Perspective(glm.DegToRad(60), float32(extent.Width)/float32(extent.Height), 0.1, 100),
	}
	uniform.Projection[5] *= -1

	translationMat := glm.Translate3D(position.X(), position.Y(), position.Z())
	rotationMat := glm.Ident4()
	currentRotMat := glm.HomogRotate3DX(glm.DegToRad(rotation.X()))
	rotationMat = currentRotMat.Mul4(&rotationMat)
	currentRotMat = glm.HomogRotate3DY(glm.DegToRad(rotation.Y()))
	rotationMat = currentRotMat.Mul4(&rotationMat)
	currentRotMat = glm.HomogRotate3DZ(glm.DegToRad(rotation.Z()))
	rotationMat = currentRotMat.Mul4(&rotationMat)
	uniform.Model = translationMat.Mul4(&rotationMat)

	return uniform
}
