package vertex

import (
	"unsafe"

	"github.com/EngoEngine/glm"
)

type VulkanVertex struct {
	Position glm.Vec3
	Color    glm.Vec3
	TexCoord glm.Vec2
}

const VulkanVertexSize = unsafe.Sizeof(VulkanVertex{})
const VulkanVertexPositionOffset = unsafe.Offsetof(VulkanVertex{}.Position)
const VulkanVertexColorOffset = unsafe.Offsetof(VulkanVertex{}.Color)
const VulkanVertexTexCoordOffset = unsafe.Offsetof(VulkanVertex{}.TexCoord)

var DefaultVertices = []VulkanVertex{
	{Position: glm.Vec3{-0.5, -0.5, 0.5}, Color: glm.Vec3{1, 0, 0}, TexCoord: glm.Vec2{1, 0}},
	{Position: glm.Vec3{0.5, -0.5, 0.5}, Color: glm.Vec3{0, 1, 0}, TexCoord: glm.Vec2{0, 0}},
	{Position: glm.Vec3{0.5, 0.5, 0.5}, Color: glm.Vec3{0, 0, 1}, TexCoord: glm.Vec2{0, 1}},
	{Position: glm.Vec3{-0.5, 0.5, 0.5}, Color: glm.Vec3{1, 1, 1}, TexCoord: glm.Vec2{1, 1}},
	{Position: glm.Vec3{-0.5, -0.5, -0.5}, Color: glm.Vec3{1, 0, 0}, TexCoord: glm.Vec2{1, 0}},
	{Position: glm.Vec3{0.5, -0.5, -0.5}, Color: glm.Vec3{0, 1, 0}, TexCoord: glm.Vec2{0, 0}},
	{Position: glm.Vec3{0.5, 0.5, -0.5}, Color: glm.Vec3{0, 0, 1}, TexCoord: glm.Vec2{0, 1}},
	{Position: glm.Vec3{-0.5, 0.5, -0.5}, Color: glm.Vec3{1, 1, 1}, TexCoord: glm.Vec2{1, 1}},
}
var DefaultVerticesIndex = []uint16{
	0, 1, 2, 2, 3, 0, // Front face
	4, 7, 6, 6, 5, 4, // Back face
	7, 4, 0, 0, 3, 7, // Left face
	6, 2, 1, 1, 5, 6, // Right face
	4, 5, 1, 1, 0, 4, // Bottom face
	7, 3, 2, 2, 6, 7, // Top face
}
