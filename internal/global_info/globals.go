package global_info

import (
	"os"
	"path/filepath"
)

var GO_BAAS_EXECUTABLE_PATH string
var GO_BAAS_EXECUTABLE_DIR string
var GO_BAAS_CONFIG_DIR string
var GO_BAAS_OUTPUT_DIR string
var GO_BAAS_DEFAULT_CONFIG_PATH string

func InitGlobalInfo() {
	GO_BAAS_EXECUTABLE_PATH, _ = os.Executable()
	GO_BAAS_EXECUTABLE_DIR = filepath.Dir(GO_BAAS_EXECUTABLE_PATH)
	GO_BAAS_OUTPUT_DIR = filepath.Join(GO_BAAS_EXECUTABLE_DIR, "output")
	GO_BAAS_CONFIG_DIR = filepath.Join(GO_BAAS_EXECUTABLE_DIR, "config")
	GO_BAAS_DEFAULT_CONFIG_PATH = filepath.Join(GO_BAAS_CONFIG_DIR, "global_config.yaml")
}
