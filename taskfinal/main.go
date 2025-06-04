package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

type RepoInfo struct {
	Name          string `json:"name"`
	FullName      string `json:"full_name"`
	Stars         int    `json:"stargazers_count"`
	Forks         int    `json:"forks_count"`
	OpenIssues    int    `json:"open_issues_count"`
	Language      string `json:"language"`
	HTMLURL       string `json:"html_url"`
	DefaultBranch string `json:"default_branch"`
}

type PageData struct {
	Error    string
	Repo     *RepoInfo
	Username string
	RepoName string
}

var tmpl = template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<html>
<head>
	<title>GitHub Repo Analyzer</title>
	<style>
		body { font-family: Arial, sans-serif; max-width: 600px; margin: auto; }
		.error { color: red; }
	</style>
</head>
<body>
	<h1>GitHub Repo Analyzer</h1>
	<form method="GET" action="/">
		<label>GitHub Username:<br><input name="username" value="{{.Username}}" required></label><br><br>
		<label>Repo Name:<br><input name="repo" value="{{.RepoName}}" required></label><br><br>
		<button type="submit">Analyze</button>
	</form>

	{{if .Error}}<p class="error">{{.Error}}</p>{{end}}

	{{if .Repo}}
		<h2>Repository: <a href="{{.Repo.HTMLURL}}" target="_blank">{{.Repo.FullName}}</a></h2>
		<ul>
			<li>⭐ Stars: {{.Repo.Stars}}</li>
			<li>🍴 Forks: {{.Repo.Forks}}</li>
			<li>🐞 Open Issues: {{.Repo.OpenIssues}}</li>
			<li>📝 Language: {{.Repo.Language}}</li>
			<li>🌿 Default Branch: {{.Repo.DefaultBranch}}</li>
		</ul>
	{{end}}
</body>
</html>
`))

func handler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	repoName := r.URL.Query().Get("repo")

	data := PageData{
		Username: username,
		RepoName: repoName,
	}

	if username == "" || repoName == "" {
		tmpl.Execute(w, data)
		return
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", username, repoName)
	resp, err := http.Get(url)
	if err != nil {
		data.Error = "Failed to fetch repo data."
		tmpl.Execute(w, data)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		data.Error = "Repo not found or GitHub API error."
		tmpl.Execute(w, data)
		return
	}

	var repo RepoInfo
	err = json.NewDecoder(resp.Body).Decode(&repo)
	if err != nil {
		data.Error = "Failed to parse GitHub API response."
		tmpl.Execute(w, data)
		return
	}

	data.Repo = &repo
	tmpl.Execute(w, data)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
