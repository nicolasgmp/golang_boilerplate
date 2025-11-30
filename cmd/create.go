package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

func CreateBoilerplate() *cobra.Command {
	projectName, projectPath := "", ""

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create boilerplate for a new project",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := checkPathAndName(projectPath, projectName); err != nil {
				return err
			}

			globalPath := filepath.Join(projectPath, projectName)
			if _, err := os.Stat(globalPath); err == nil {
				fmt.Println("Project already exists at", globalPath)
				return err
			}

			if err := createDir(globalPath); err != nil {
				return fmt.Errorf("failed to create dir: %w", err)
			}

			if err := initGo(projectName, globalPath); err != nil {
				return fmt.Errorf("failed to init go module: %w", err)
			}

			if err := createInitialDirs(globalPath); err != nil {
				return fmt.Errorf("failed to create dirs: %w", err)
			}
			if err := createInitialFiles(globalPath); err != nil {
				return fmt.Errorf("failed to create files: %w", err)
			}
			if err := writeMainFile(globalPath); err != nil {
				return fmt.Errorf("failed to write main.go: %w", err)
			}

			if err := writeRoutesFile(globalPath); err != nil {
				return fmt.Errorf("failed to write routes.go: %w", err)
			}

			fmt.Println("Creating Project", projectName, "at", projectPath)

			return nil
		},
	}

	cmd.Flags().StringVarP(&projectName, "name", "n", "", "Name of the project")
	cmd.Flags().StringVarP(&projectPath, "path", "p", "", "Path where the project will be created")

	return cmd
}

func checkPathAndName(path, name string) error {
	if path == "" {
		return fmt.Errorf("you must supply a project path")
	}
	if name == "" {
		return fmt.Errorf("you must supply a project name")
	}
	return nil
}

func createDir(globalPath string) error {
	if err := os.Mkdir(globalPath, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func initGo(name, globalPath string) error {
	startGo := exec.Command("go", "mod", "init", name)
	startGo.Dir = globalPath
	startGo.Stdout = os.Stdout
	startGo.Stderr = os.Stderr

	err := startGo.Run()
	if err != nil {
		return err
	}

	return nil
}

func createInitialDirs(globalPath string) error {
	cmdPath := filepath.Join(globalPath, "cmd")
	if err := os.Mkdir(cmdPath, os.ModePerm); err != nil {
		return err
	}

	internalPath := filepath.Join(globalPath, "internal")
	if err := os.Mkdir(internalPath, os.ModePerm); err != nil {
		return err
	}

	handlerPath := filepath.Join(internalPath, "handler")
	if err := os.Mkdir(handlerPath, os.ModePerm); err != nil {
		return err
	}

	routesPath := filepath.Join(handlerPath, "routes")
	if err := os.Mkdir(routesPath, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func createInitialFiles(globalPath string) error {
	mainPath := filepath.Join(globalPath, "cmd", "main.go")
	mainFile, err := os.Create(mainPath)
	if err != nil {
		return err
	}
	defer mainFile.Close()

	routesFilePath := filepath.Join(globalPath, "internal", "handler", "routes", "routes.go")
	routesFile, err := os.Create(routesFilePath)
	if err != nil {
		return err
	}
	defer routesFile.Close()

	return nil
}

func writeMainFile(globalPath string) error {
	packageContent := []byte(`package main

	import "fmt"

	func main() {
		fmt.Println("Hello, World!")
	}
	`)

	mainPath := filepath.Join(globalPath, "cmd", "main.go")
	if err := os.WriteFile(mainPath, packageContent, 0o666); err != nil {
		return err
	}

	return nil
}

func writeRoutesFile(globalPath string) error {
	packageContent := []byte(`package routes

	// Type your code here
	`)

	routesFilePath := filepath.Join(globalPath, "internal", "handler", "routes", "routes.go")
	if err := os.WriteFile(routesFilePath, packageContent, 0o666); err != nil {
		return err
	}

	return nil
}
