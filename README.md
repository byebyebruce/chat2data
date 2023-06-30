 <div align="center">
            <img src="https://readme-typing-svg.demolab.com/?lines=Chat+2+Data&size=50&height=80&center=true&vCenter=true&&duration=1000&pause=5000">
        </div>

> Use AI to chat to your mysql database or sqlite3 database.

## PREVIEW
![](doc/preview.jpg)

## BUILD
* go >= 1.19
```bash
make build
```

## RUN
1. Config environment variables 
   * Use local `.env` file `cp .env.template .env` then edit it.  
   * You can also use `export OPENAI_API_KEY=xxx` to specify the environment variables.
2. Run CLI(command line interface)
   * mysql `chat2data --mysql=root:pwd@tcp(localhost:3306)/mydb` 
   * sqlite3 `chat2data --sqlite3=mytest.db`

## TODO
- [ ] Support csv
- [ ] Web ui

## SPECIAL THANKS
* [🦜️🔗 LangChain Go](https://github.com/tmc/langchaingo)

![](https://hits.sh/github.com/byebyebruce/chat2data/doc/hits.svg?label=%F0%9F%91%80)
