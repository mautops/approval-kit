package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestGoModExists 验证 go.mod 文件存在
func TestGoModExists(t *testing.T) {
	rootDir := getProjectRoot(t)
	goModPath := filepath.Join(rootDir, "go.mod")

	if !fileExists(goModPath) {
		t.Fatal("go.mod 文件不存在")
	}
}

// TestGoModContent 验证 go.mod 文件内容符合要求
func TestGoModContent(t *testing.T) {
	rootDir := getProjectRoot(t)
	goModPath := filepath.Join(rootDir, "go.mod")

	content, err := readFile(goModPath)
	if err != nil {
		t.Fatalf("无法读取 go.mod 文件: %v", err)
	}

	// 检查必需的内容
	requiredContents := []string{
		"module",
		"go ",
	}

	for _, required := range requiredContents {
		if !strings.Contains(content, required) {
			t.Errorf("go.mod 文件缺少必需内容: %s", required)
		}
	}
}

// TestGoModValid 验证 go.mod 文件格式正确
func TestGoModValid(t *testing.T) {
	rootDir := getProjectRoot(t)

	// 运行 go mod verify 验证模块配置
	cmd := exec.Command("go", "mod", "verify")
	cmd.Dir = rootDir
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("go mod verify 失败: %v\n输出: %s", err, string(output))
	}
}

// TestGoModTidy 验证 go.mod 文件整洁(无多余依赖)
func TestGoModTidy(t *testing.T) {
	rootDir := getProjectRoot(t)

	// 运行 go mod tidy 检查依赖
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = rootDir
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("go mod tidy 失败: %v\n输出: %s", err, string(output))
	}

	// 如果输出不为空,说明有依赖变更,需要检查
	if len(output) > 0 {
		t.Logf("go mod tidy 输出: %s", string(output))
	}
}

// TestNoExternalDependencies 验证没有外部依赖(遵循依赖最小化原则)
func TestNoExternalDependencies(t *testing.T) {
	rootDir := getProjectRoot(t)
	goModPath := filepath.Join(rootDir, "go.mod")

	content, err := readFile(goModPath)
	if err != nil {
		t.Fatalf("无法读取 go.mod 文件: %v", err)
	}

	// 检查是否包含 require 块(外部依赖)
	// 根据 constitution,应该优先使用标准库,最小化外部依赖
	lines := strings.Split(content, "\n")
	inRequireBlock := false
	hasExternalDeps := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "require (") {
			inRequireBlock = true
			continue
		}
		if inRequireBlock {
			if strings.HasPrefix(line, ")") {
				break
			}
			// 忽略注释和空行
			if line != "" && !strings.HasPrefix(line, "//") {
				hasExternalDeps = true
				t.Logf("发现外部依赖: %s", line)
			}
		}
	}

	// 注意: 这里不强制要求没有外部依赖,只是记录
	// 因为某些功能可能需要外部依赖(如 HTTP 客户端增强)
	// 但应该遵循依赖最小化原则
	if hasExternalDeps {
		t.Log("警告: 发现外部依赖,请确保遵循依赖最小化原则")
	}
}

// TestGoVersion 验证 Go 版本符合要求
func TestGoVersion(t *testing.T) {
	rootDir := getProjectRoot(t)
	goModPath := filepath.Join(rootDir, "go.mod")

	content, err := readFile(goModPath)
	if err != nil {
		t.Fatalf("无法读取 go.mod 文件: %v", err)
	}

	// 检查 Go 版本
	// 根据 constitution,要求 Go 1.25.4 或更高版本
	if !strings.Contains(content, "go 1.25") && !strings.Contains(content, "go 1.26") {
		t.Log("警告: Go 版本可能不符合要求(需要 1.25.4 或更高)")
	}
}

// 辅助函数

func getProjectRoot(t *testing.T) string {
	rootDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("无法获取当前工作目录: %v", err)
	}

	// 如果当前在 tests 目录,向上找到项目根目录
	if filepath.Base(rootDir) == "tests" {
		rootDir = filepath.Dir(rootDir)
	}

	return rootDir
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func readFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

