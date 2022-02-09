package main

import (
	"context"
	"digdagGo"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	tdToken := os.Getenv("TD_API_KEY")
	ctx := context.Background()

	client, err := digdagGo.New("https://api-workflow.treasuredata.com/api", tdToken, nil)
	if err != nil {
		fmt.Println(err)
	}
	projects, err := client.GetProjects(ctx, "mikio-test")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", projects.Projects[0].DeletedAt)

	projectWithId, err := client.GetProjectsWithID(ctx, "590006")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", projectWithId.ID)

	currentPath, err := filepath.Abs(".")
	if err != nil {
		fmt.Println(err)
	}

	project, err := client.PutProject(ctx, fmt.Sprintf("%s/testFiles/sample/sample.tar.gz", currentPath), "mikio-test")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(project.ID)

	ok, err := client.DeleteProjectsWithID(ctx, "590006")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ok.ID)

}
