package main

import (
	"fmt"
	"go-PlumAIO/src/bot/log"
	"go-PlumAIO/src/bot/stores/ambush"
	"go-PlumAIO/src/bot/stores/ldlc"
	"go-PlumAIO/src/bot/stores/stylefile"
	"go-PlumAIO/src/bot/stores/swatch"
	"go-PlumAIO/src/bot/tasks"
	"os"
	"sync"
)

var (
	format = fmt.Sprintf
)

func main() {

	for {

		log.Infoln(format("PlumAIO BOT"), "-")
		log.Infoln(format("Welcome user.."), "-")

		config, err := tasks.ReadConfig()

		if err != nil {
			log.Error("Error reading config file", "-")
			os.Exit(0)
		}

		proxies, err := tasks.ReadProxies()

		if err != nil {
			log.Error("Error reading proxies file", "-")
			os.Exit(0)
		}

		rows, err := tasks.ReadFile("tasks/tasks.csv")

		if err != nil {
			fmt.Println(err)
			log.Error("Error reading tasks file", "-")
			os.Exit(0)
		}

		if len(rows) > 500 {
			log.Error("Max 500 tasks", "-")
			os.Exit(0)
		}

		wg := sync.WaitGroup{}

		var choice int
		fmt.Printf("\n\n1 - Ambush\n2 - LDLC\n3 - Stylefile\n4 - Swatch\n5 - Exit\n\nChoose : ")

		fmt.Scanf("%d", &choice)

		switch choice {
		case 1:
			for id, row := range rows {
				wg.Add(1)
				ambush.Start(row, config, proxies, id, &wg)
			}
			return
		case 2:
			fmt.Println("Ldlc currently unavailable ...")
			return
		case 3:
			for id, row := range rows {
				wg.Add(1)
				stylefile.Start(row, config, proxies, id, &wg)
			}
			return
		case 4:
			fmt.Println("Swatch currently unavailable ...")
		case 5:
			os.Exit(0)
			return
		default:
			continue
		}
	}
}
