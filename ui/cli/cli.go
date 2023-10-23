package cli

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/byebyebruce/chat2data/qa"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

func CLI(qa qa.QA, info any) error {
	fmt.Println(info)
	_qa := func(str string) {
		ctx1, cancel1 := context.WithTimeout(context.Background(), time.Minute*2)
		defer cancel1()

		answer, err := qa.Answer(ctx1, str)
		if err != nil {
			fmt.Println(color.RedString("error:%s", err))
		} else {
			fmt.Println(color.GreenString("Answer:\n"), color.GreenString(answer))
		}
	}

	defaultQuestion := "" //fmt.Sprintf(" How may records are there in the table %s?", tbs[0])
	for {
		pt := promptui.Prompt{
			Label:   color.CyanString("Input your question"),
			Default: defaultQuestion,
		}

		str, err := pt.Run()
		if err != nil {
			fmt.Println(err)
			return err
		}
		str = strings.TrimSpace(str)
		if len(str) == 0 {
			continue
		}
		_qa(str)
		fmt.Println()
	}
}
