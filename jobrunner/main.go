package main

import (
	"fmt"
	"html/template"
	"time"

	"github.com/bamzi/jobrunner"
	"github.com/gin-gonic/gin"
)

var tpl, _ = template.New("test").Parse(`<html>
	<head>
		<style>
body {
	margin: 30px 0 0 0;
  font-size: 11px;
  font-family: sans-serif;
  color: #345;

}
h1 {
	font-size: 24px;
	text-align: center;
	padding: 10 0 30px;
}
table {
	/*max-width: 80%;*/
	margin: 0 auto;
  border-collapse: collapse;
  border: none;
}
table td, table th {
	min-width: 25px;
	width: auto;
  padding: 15px 20px;
  border: none;
}


table tr:nth-child(odd) {
  background-color: #f0f0f0;
}
table tr:nth-child(1) {
  background-color: #345;
  color: white;
}
th {
  text-align: left;
}

		</style>
	</head>
	<body>

<h1>JobRunner Status Report</h1>

<table>
	<tr><th>ID</th><th>Name</th><th>Status</th><th>Last run</th><th>Next run</th><th>Latency</th></tr>
{{range .}}

	<tr>
		<td>{{.Id}}</td>
		<td>{{.JobRunner.Name}}</td>
		<td>{{.JobRunner.Status}}</td>
		<td>{{if not .Prev.IsZero}}{{.Prev.Format "2006-01-02 15:04:05"}}{{end}}</td>
		<td>{{if not .Next.IsZero}}{{.Next.Format "2006-01-02 15:04:05"}}{{end}}</td>
		<td>{{.JobRunner.Latency}}</td>
	</tr>
{{end}}
</table>
`)

type GreetingJob struct {
	Name string
}

func (g GreetingJob) Run() {
	fmt.Println("Hello,", g.Name)
}

type EmailJob struct {
	Email string
}

func (e EmailJob) Run() {
	fmt.Println("Send,", e.Email)
}

func main() {
	jobrunner.Start()
	jobrunner.Every(5*time.Second, GreetingJob{Name: "dj"})
	jobrunner.Every(10*time.Second, EmailJob{Email: "935653229@qq.com"})

	r := gin.Default()
	r.GET("/jobrunner/json", JobJson)
	r.GET("/jobrunner/html", JobHtml)
	r.Run(":8888")

}

func JobJson(c *gin.Context) {
	c.JSON(200, jobrunner.StatusJson())
}
func JobHtml(c *gin.Context) {
	tpl.Execute(c.Writer, jobrunner.StatusPage())
}
