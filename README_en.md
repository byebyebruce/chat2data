 <div align="center">
   <img src="https://readme-typing-svg.demolab.com/?lines=Chat+2+Data&size=50&height=80&center=true&vCenter=true&&duration=1000&pause=5000">
</div>

[中文](README.md) | English

> 🗣 📊 Chat2Data is a tool for interacting with your data, including MySQL, PostgreSQL, SQLite3, CSV, Text, PDF, and HTML pages.

[![Go Report Card](https://goreportcard.com/badge/github.com/byebyebruce/chat2data)](https://goreportcard.com/report/github.com/byebyebruce/chat2data)
![GitHub release (with filter)](https://img.shields.io/github/v/release/byebyebruce/chat2data)
[![Docker Pulls](https://img.shields.io/docker/pulls/bailu1901/chat2data)](https://hub.docker.com/r/bailu1901/chat2data/)
![](https://hits.sh/github.com/byebyebruce/chat2data/doc/hits.svg?label=visit)


## Feature
* 🗣 Easy Interaction: Chat2Data allows you to chat with your data, making it intuitive to use.
* 🔗 Multiple Databases: It supports MySQL, PostgreSQL, SQLite3, CSV, Text, PDF, and HTML pages.
* 🐳 Docker Support: It provides a Docker image for easy deployment.
* 💻 CLI and Web UI: It offers both a command line and a web interface.
* ⚙️ Simple Installation: It's easy to install with a Go command.
* 🧠 AI Integration: It leverages the OpenAI API for advanced natural language processing. 
 
## Preview
![CLI](doc/cli.jpg)
![Web UI](doc/web-ui.png)

## Install
#### Download  
[Releases Page](https://github.com/byebyebruce/chat2data/releases)
  
#### Go install  
`go install github.com/byebyebruce/chat2data/cmd/chat2data@latest`

## Quick Run
* Binary
```bash
OPENAI_API_KEY=xxx chat2data db -c testdata/world_happiness_2015.db
```
Ask: `Which is the highest happiness country?`

* Docker
```bash
docker run --rm -it -e OPENAI_API_KEY=xxx -p 8088:8088 bailu1901/chat2data html 'https://github.com/byebyebruce/chat2data'
```
Open `http://localhost:8088` in browser, then ask: `What is the feature of chat2data?`

## Config
   * Use local `.env` file `cp .env.template .env` then edit it.  
   * You can also use `export OPENAI_API_KEY=xxx` to specify the environment variables.
   * Or run with env `OPENAI_API_KEY=xxx OPENAI_BASE_URL=https://api.openai.com/v1 chat2data db root:pwd@tcp(localhost:3306)/mydb`
    
## Usage
* help `chat2data --help`  
global flags
```bash
      --web  -w  web ui port
      --cli  -c  cli mode
```
1. Run CLI(command line interface)
   * mysql `chat2data db -c root:pwd@tcp(localhost:3306)/mydb` 
   * postgre `chat2data db -c postgres://db_user:mysecretpassword@localhost:5438/test?sslmode=disable`
   * sqlite3 `chat2data db -c sqlite.db`
   * csv `chat2data csv -c csvfile.csv` or `chat2data csv csvdir`
   * html `chat2data html -c https://github.com/byebyebruce/chat2data`
   * text `chat2data txt -c textfile.txt`
   * with env `OPENAI_API_KEY=xxx chat2data db -c root:pwd@tcp(localhost:3306)/mydb`
2. Run Web UI
   * mysql `chat2data db root:example@tcp(10.12.21.101:3306)/mydb`
   * html `chat2data html https://github.com/byebyebruce/chat2data`
   * pdf `chat2data pdf testdata/sample.pdf`
   * sqlite3 `chat2data db -w=:0.0.0.0:8088 mytest.db`

## Build 
`git clone github.com/byebyebruce/chat2data`
* build binary
```base
make build
```
* build docker image
```bash 
docker build -t chat2data .
```

## TODO
- [x] Support Docker
- [x] Support Postgre Database
- [x] Support load csv
- [x] Add Web ui
- [x] Local vector database
- [x] Support load html
- [x] Support load pdf
- [x] Doc QA
- [ ] Support word
- [ ] Beautiful CLI 

## [Change Log](CHANGELOG.md)

## Special Thanks
* [🦜️🔗 LangChain Go](https://github.com/tmc/langchaingo)
 

