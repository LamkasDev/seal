package shader

const SHADER_BASIC = "basic"

var DefaultShaders = []string{
	SHADER_BASIC,
}

type VulkanShaderContainer struct {
	Shaders map[string]VulkanShader
}

func NewVulkanShaderContainer() (VulkanShaderContainer, error) {
	container := VulkanShaderContainer{
		Shaders: map[string]VulkanShader{},
	}
	for _, shader := range DefaultShaders {
		if _, err := CreateShaderWithContainer(&container, shader); err != nil {
			return container, err
		}
	}

	return container, nil
}

func CreateShaderWithContainer(container *VulkanShaderContainer, id string) (VulkanShader, error) {
	shader, err := NewShader(id)
	if err != nil {
		return shader, err
	}
	container.Shaders[id] = shader

	return shader, nil
}

func FreeVulkanShaderContainer(container *VulkanShaderContainer) error {
	return nil
}
