package model

import "github.com/LamkasDev/seal/cmd/engine/vulkan/vertex"

type VulkanModel struct {
	Vertices []vertex.VulkanVertex
	Indices  []uint16
}

func ConvertModel(rawModel Model) VulkanModel {
	model := VulkanModel{
		Vertices: []vertex.VulkanVertex{},
		Indices:  []uint16{},
	}
	for i := uint16(0); i < uint16(len(rawModel.VecIndices)); i++ {
		model.Vertices = append(model.Vertices, vertex.VulkanVertex{
			Position: rawModel.Vecs[uint16(rawModel.VecIndices[i])-1],
			TexCoord: rawModel.Uvs[uint16(rawModel.UvIndices[i])-1],
		})
		model.Indices = append(model.Indices, i)
	}

	return model
}
