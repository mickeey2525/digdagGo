package main

import (
	"context"
	digdaggo "digdagGo"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main() {
	tdToken := os.Getenv("TD_API_KEY")
	ctx := context.Background()

	client, err := digdaggo.New("https://api-workflow.treasuredata.com/api", tdToken, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	currentPath, err := filepath.Abs("./")
	if err != nil {
		fmt.Println(err)
	}
	project, err := client.PutProject(ctx, fmt.Sprintf("%s/testFiles/sample/sample.tar.gz", currentPath), "mikio-test")
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	fmt.Println(project.ID)

	projects, err := client.GetProjects(ctx, "mikio-test")
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	fmt.Println(projects.Projects[0].ID)

	schedules, err := client.GetProjectsSchedules(ctx, "119468", "", "")
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	fmt.Println(schedules)

	projectWithId, err := client.GetProjectsWithID(ctx, project.ID)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	fmt.Printf("%s\n", projectWithId.Revision)

	revisions, err := client.GetListRevisions(ctx, project.ID)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	fmt.Printf("%+v\n", revisions)

	secrets, err := client.GetSecrets(ctx, project.ID)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	fmt.Printf("%+v\n", secrets.Secrets...)

	sessions, err := client.GetProjectSessions(ctx, project.ID, "", "", "")
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	fmt.Printf("%+v\n", sessions)
	workflows, err := client.GetProjectWorkflows(ctx, project.ID, project.Revision, "")
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	fmt.Printf("the workflos are %+v\n", workflows)
	attempts, err := client.GetAttempts(ctx, project.Name, "", "", "", false)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	fmt.Printf("the attemts are %+v\n", attempts)
	version, err := client.GetServerVersion(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	fmt.Printf("the server version is %+v\n", version)
	//client.DownloadProjectFiles(ctx, project.ID, project.Revision, fmt.Sprintf("%s/test", currentPath), true)
	// ok, err := client.DeleteProjectsWithID(ctx, project.ID)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(ok.ID)
	sessionTime := time.Now()
	params := struct{}{}
	invokeAttempt, err := client.StartAttempt(ctx, params, 12367598, sessionTime)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(invokeAttempt)
}
