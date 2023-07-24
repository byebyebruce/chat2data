 <div align="center">
   <img src="https://readme-typing-svg.demolab.com/?lines=Chat+2+Data&size=50&height=80&center=true&vCenter=true&&duration=1000&pause=5000">
</div>

> 🗣 📊Chat2Data is tool for interacting with databases, supporting MySQL, PostgreSQL, SQLite3, and CSV files
## Feature
* 🗣 Easy Interaction: Chat2Data lets you chat with your databases, making it intuitive to use.
* 🔗 Multiple Databases: It supports MySQL, PostgreSQL, SQLite3, and CSV files.
* 🐳 Docker Support: It provides a Docker image for easy deployment.
* 💻 CLI and Web UI: It offers both a command line and a web interface.
* ⚙️ Simple Installation: It's easy to install with Go command.
* 🧠 AI Integration: It leverages OpenAI API for advanced natural language processing.
 
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
OPENAI_API_KEY=xxx chat2data --mysql=root:pwd@tcp(localhost:3306)/mydb
```
* Docker
```bash
docker run --rm -it -e OPENAI_API_KEY=xxx bailu1901/chat2data --mysql=root:pwd@tcp(localhost:3306)/mydb
```

## Config
   * Use local `.env` file `cp .env.template .env` then edit it.  
   * You can also use `export OPENAI_API_KEY=xxx` to specify the environment variables.
   * Or run with env `OPENAI_API_KEY=xxx OPENAI_BASE_URL=https://api.openai.com chat2data --mysql=root:pwd@tcp(localhost:3306)/mydb`
    
## Command
1. Run CLI(command line interface)
   * help `chat2data --help`
   * mysql `chat2data --mysql=root:pwd@tcp(localhost:3306)/mydb` 
   * postgre `chat2data --postgre=postgres://db_user:mysecretpassword@localhost:5438/test?sslmode=disable`
   * sqlite3 `chat2data --sqlite3=mytest.db`
   * csv `chat2data --csv=csvfile.csv` or `chat2data --csv=csvdir`
   * with env `OPENAI_API_KEY=xxx chat2data --mysql=root:pwd@tcp(localhost:3306)/mydb`
   * 
2. Run Web UI
   * mysql `chat2data --mysql="root:example@tcp(10.12.21.101:3306)/mydb" --web=8088`
   * sqlite3 `chat2data --sqlite3=./mytest.db  --web=8088`

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
- [ ] Local vector database
- [ ] Doc QA

## [Change Log](CHANGELOG.md)

## Special Thanks
* [🦜️🔗 LangChain Go](https://github.com/tmc/langchaingo)
 
![](https://hits.sh/github.com/byebyebruce/chat2data/doc/hits.svg?label=%F0%9F%91%80)

