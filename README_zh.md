 <div align="center">
   <img src="https://readme-typing-svg.demolab.com/?lines=Chat+2+Data&size=50&height=80&center=true&vCenter=true&&duration=1000&pause=5000">
</div>

[English](README.md) | 中文文档

>  🗣 📊 Chat2Data是一个与数据交互的工具,支持MySQL、PostgreSQL、SQLite3、CSV、文本、PDF和HTML页面。
 
[![Go Report Card](https://goreportcard.com/badge/github.com/byebyebruce/chat2data)](https://goreportcard.com/report/github.com/byebyebruce/chat2data)
![GitHub release (with filter)](https://img.shields.io/github/v/release/byebyebruce/chat2data)
[![Docker Pulls](https://img.shields.io/docker/pulls/bailu1901/chat2data)](https://hub.docker.com/r/bailu1901/chat2data/)
![](https://hits.sh/github.com/byebyebruce/chat2data/doc/hits.svg?label=visit)


## 特性
* 🗣 易于交互:Chat2Data允许你通过聊天的方式与数据交互,使得使用很直观。
* 🔗 支持多种数据库:它支持MySQL、PostgreSQL、SQLite3、CSV、文本、PDF和HTML页面。
* 🐳 Docker支持:提供Docker镜像方便部署。
* 💻 命令行和Web界面:同时提供命令行和Web界面。
* ⚙️ 简单安装:通过Go命令可以很容易安装。
* 🧠 AI集成:利用OpenAI API进行高级自然语言处理。

## 预览
![CLI](doc/cli.jpg)
![Web UI](doc/web-ui.png)

## 安装
#### 下载编译好的程序  
[Releases Page](https://github.com/byebyebruce/chat2data/releases)
  
#### Go安装 
`go install github.com/byebyebruce/chat2data/cmd/chat2data@latest`

## 快速运行
* 二进制程序
```bash
OPENAI_API_KEY=xxx chat2data db -c testdata/world_happiness_2015.db
```
输入: `Which is the highest happiness country?`

* Docker
```bash
docker run --rm -it -e OPENAI_API_KEY=xxx -p 8088:8088 bailu1901/chat2data html 'https://github.com/byebyebruce/chat2data'
```
在浏览器打开 http://localhost:8088,然后询问:chat2data的特性是什么?

## 配置
* 使用本地.env文件 cp .env.template .env 然后编辑它。
* 也可以使用 export OPENAI_API_KEY=xxx 来指定环境变量。
* 或者带着环境变量运行 OPENAI_API_KEY=xxx OPENAI_BASE_URL=https://api.openai.com/v1 chat2data db root:pwd@tcp(localhost:3306)/mydb    
 
## 用法
* 帮助信息 chat2data --help 
* 全局参数
```bash
      --web  -w  web ui port
      --cli  -c  cli mode
```
1. 运行命令行CLI(Command Line Interface)
   * mysql `chat2data db -c root:pwd@tcp(localhost:3306)/mydb` 
   * postgre `chat2data db -c postgres://db_user:mysecretpassword@localhost:5438/test?sslmode=disable`
   * sqlite3 `chat2data db -c sqlite.db`
   * csv `chat2data csv -c csvfile.csv` or `chat2data csv csvdir`
   * html `chat2data html -c https://github.com/byebyebruce/chat2data`
   * text `chat2data txt -c textfile.txt`
   * with env `OPENAI_API_KEY=xxx chat2data db -c root:pwd@tcp(localhost:3306)/mydb`
2. 运行Web界面
   * mysql `chat2data db root:example@tcp(10.12.21.101:3306)/mydb`
   * html `chat2data html https://github.com/byebyebruce/chat2data`
   * pdf `chat2data pdf testdata/sample.pdf`
   * sqlite3 `chat2data db -w=:0.0.0.0:8088 mytest.db`

## 构建
`git clone github.com/byebyebruce/chat2data`
* 构建二进制程序
```base
make build
```
* 构建Docker镜像
```bash 
docker build -t chat2data .
```

## TODO
- [x] 支持Docker
- [x] 支持PostgreSQL数据库
- [x] 支持加载CSV
- [x] 添加Web界面
- [x] 本地向量数据库
- [x] 支持加载HTML
- [x] 支持加载PDF
- [x] 文档问答
- [ ] 支持Word 文档
- [ ] 更优雅的命令行界面

## [更新日志](CHANGELOG.md)

## 特别感谢
* [🦜️🔗 LangChain Go](https://github.com/tmc/langchaingo)
 

