package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/byebyebruce/chat2data/assets"
	"github.com/byebyebruce/chat2data/datachain"
	"github.com/fatih/color"
)

func Web(addr string, chain *datachain.DataChain) {
	f, err := fs.Sub(assets.Web, "web")
	if err != nil {
		printError(err)
		return
	}
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.FS(f))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFS(assets.Web, "web/index.html")
		if err != nil {
			printError(err)
			return
		}
		if err := tmpl.Execute(w, nil); err != nil {
			printError(err)
			return
		}
	})
	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		data, err := io.ReadAll(r.Body)
		if err != nil {
			printError(err)
			return
		}

		req := &struct{ SearchText string }{}
		if err := json.Unmarshal(data, req); err != nil {
			printError(err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
		defer cancel()

		log.Println("SearchText:", req.SearchText)
		answer, err := chain.Run(ctx, req.SearchText)
		if err != nil {
			printError(err)
			return
		}
		fmt.Println(color.GreenString("answer:\n"), color.GreenString(answer))
		if _, err := w.Write([]byte(answer)); err != nil {
			printError(err)
			return
		}
	})

	host := fmt.Sprintf("http://localhost:%s", addr)
	fmt.Println(color.GreenString("web"), color.GreenString(host))

	if err := http.ListenAndServe(":"+addr, nil); err != nil {
		log.Fatal(err)
	}
}

func printError(err error) {
	fmt.Println(color.RedString("error:%s", err))
}
