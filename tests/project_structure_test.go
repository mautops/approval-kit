package tests

import (
	"os"
	"path/filepath"
	"testing"
)

// TestProjectStructure 验证项目目录结构是否符合规范
func TestProjectStructure(t *testing.T) {
	requiredDirs := []string{
		"internal/task",
		"internal/template",
		"internal/statemachine",
		"internal/node",
		"internal/event",
		"tests/task",
		"tests/template",
		"tests/statemachine",
		"tests/node",
		"tests/event",
		"examples",
	}

	requiredFiles := []string{
		"go.mod",
		"README.md",
	}

	// 获取项目根目录
	// 如果从 tests 目录运行,需要向上找到项目根目录
	rootDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("无法获取当前工作目录: %v", err)
	}
	
	// 如果当前在 tests 目录,向上找到项目根目录
	if filepath.Base(rootDir) == "tests" {
		rootDir = filepath.Dir(rootDir)
	}

	// 检查必需的目录
	for _, dir := range requiredDirs {
		dirPath := filepath.Join(rootDir, dir)
		info, err := os.Stat(dirPath)
		if err != nil {
			t.Errorf("必需的目录不存在: %s", dir)
			continue
		}
		if !info.IsDir() {
			t.Errorf("路径不是目录: %s", dir)
		}
	}

	// 检查必需的文件
	for _, file := range requiredFiles {
		filePath := filepath.Join(rootDir, file)
		info, err := os.Stat(filePath)
		if err != nil {
			t.Errorf("必需的文件不存在: %s", file)
			continue
		}
		if info.IsDir() {
			t.Errorf("路径不是文件: %s", file)
		}
	}
}

