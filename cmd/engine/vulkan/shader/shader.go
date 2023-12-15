package shader

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/LamkasDev/seal/cmd/common/arch"
	"github.com/LamkasDev/seal/cmd/logger"
)

type VulkanShader struct {
	Vertex   []byte
	Fragment []byte
}

func CompileShader(source string, target string) error {
	sourceStat, sourceErr := os.Stat(source)
	targetStat, targetErr := os.Stat(target)
	if sourceErr == nil && targetErr == nil && targetStat.ModTime().After(sourceStat.ModTime()) {
		return nil
	}
	logger.DefaultLogger.Debugf("compiling shader '%s'...", source)
	cmd := exec.Command("../../resources/glslc", source, "-o", target)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func NewShader(id string) (VulkanShader, error) {
	var err error
	shader := VulkanShader{}

	vertexPath := fmt.Sprintf("../../resources/shaders/%s.vert.spv", id)
	if arch.SealDebug {
		if err = CompileShader(fmt.Sprintf("../../resources/shaders/%s.vert", id), vertexPath); err != nil {
			return shader, err
		}
	}
	if _, err = os.Stat(vertexPath); err != nil {
		return shader, err
	}
	if shader.Vertex, err = os.ReadFile(vertexPath); err != nil {
		return shader, err
	}

	fragmentPath := fmt.Sprintf("../../resources/shaders/%s.frag.spv", id)
	if arch.SealDebug {
		if err = CompileShader(fmt.Sprintf("../../resources/shaders/%s.frag", id), fragmentPath); err != nil {
			return shader, err
		}
	}
	if _, err = os.Stat(fragmentPath); err != nil {
		return shader, err
	}
	if shader.Fragment, err = os.ReadFile(fragmentPath); err != nil {
		return shader, err
	}

	return shader, nil
}
