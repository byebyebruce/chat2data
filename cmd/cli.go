package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/byebyebruce/chat2data/datachain"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

func CLI(chain *datachain.DataChain) error {
	qa := func(str string) {
		ctx1, cancel1 := context.WithTimeout(context.Background(), time.Minute*2)
		defer cancel1()

		answer, err := chain.Run(ctx1, str)
		if err != nil {
			fmt.Println(color.RedString("error:%s", err))
		} else {
			fmt.Println(color.GreenString("Answer:\n"), color.GreenString(answer))
			//fmt.Println(color.GreenString("RefTables:\n"), color.GreenString("%v", refTables))
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	tbs, err := chain.SQLChain.Database.Engine.TableNames(ctx)
	if err != nil {
		return err
	}
	if len(tbs) == 0 {
		return fmt.Errorf("no tables")
	}
	color.Green("There are %d tables", len(tbs))

	defaultQuestion := fmt.Sprintf(" How may records are there in the table %s?", tbs[0])
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
		qa(str)
		fmt.Println()
	}
}
