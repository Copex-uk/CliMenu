package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"
)

// MenuItem represents a menu item with a label, command, and color.
type MenuItem struct {
	Label   string `json:"label"`
	Color   string `json:"color"`
	Command string `json:"command"` // New field for command
	//Fn      func()
}

// Menu represents the structure of the configuration file.
type Menu struct {
	Title struct {
		Label string `json:"label"`
		Color string `json:"color"`
	} `json:"title"`
	Items []MenuItem `json:"items"`
}

func main() {
	createExample := flag.Bool("create", false, "Create an example config.json")
	help := flag.Bool("help", false, "Display usage instructions")
	filePath := flag.String("file", "config.json", "Path to the config.json file")

	flag.Parse()

	if *help {
		fmt.Println("Usage:")
		fmt.Println("  menu -create : Create an example config.json")
		fmt.Println("  menu -file <path> : Specify a custom config.json file path")
		fmt.Println("  menu -help   : Display usage instructions")
		fmt.Println("\nColor Options:")
		fmt.Println("  Colors can be specified in the JSON file using the 'color' field for each menu item.")
		fmt.Println("  Available color options are:")
		fmt.Println("  - black")
		fmt.Println("  - red")
		fmt.Println("  - green")
		fmt.Println("  - yellow")
		fmt.Println("  - blue")
		fmt.Println("  - magenta")
		fmt.Println("  - cyan")
		fmt.Println("  - white")
		fmt.Println("\nMenu Item Limit:")
		fmt.Println("  The maximum limit for menu items is 10.")
		fmt.Println("  If the menu exceeds this limit, a message will be displayed.")
		fmt.Println("\nExample JSON usage:")
		fmt.Println("  {")
		fmt.Println("    \"title\": {")
		fmt.Println("      \"label\": \"My Menu\",")
		fmt.Println("      \"color\": \"blue\"")
		fmt.Println("    },")
		fmt.Println("    \"items\": [")
		fmt.Println("      {")
		fmt.Println("        \"label\": \"Option 1\",")
		fmt.Println("        \"command\": \"ls -lha\",")
		fmt.Println("        \"color\": \"green\"")
		fmt.Println("      }")
		fmt.Println("    ]")
		fmt.Println("  }")
		return
	}

	if *createExample {

		createExampleMenu(*filePath)
		fmt.Println("Example config.json created.")
		return
	}

	menu, err := loadMenu(*filePath)
	if err != nil {
		fmt.Printf("config.json not found at %s. Run 'menu -help' for usage instructions.\n", *filePath)
		return
	}

	if len(menu.Items) > 10 {
		fmt.Println("Maximum limit of 10 menu items exceeded.")
		fmt.Println("Please reduce the number of items in the menu.")
		return
	}

	quit := false
	for !quit {
		clearScreen()
		printMenu(menu.Title.Label, menu.Title.Color, menu.Items)

		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Select an option (1 - %d, or q to quit): ", len(menu.Items))
		input, _, err := reader.ReadLine()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if string(input) == "q" {
			quit = true
			continue
		}

		option, err := strconv.Atoi(string(input))
		if err != nil || option < 1 || option > len(menu.Items) {
			fmt.Println(">" + string(len(menu.Items)))
			fmt.Println("Invalid option. Please try again.")
			time.Sleep(2 * time.Second)
			continue
		}
		// use last option as Quit - uncomment
		//if option == len(menu.Items) {
		//    quit = true
		//    continue
		//}

		idx := option - 1
		clearScreen()
		printColoredText(menu.Items[idx].Color, menu.Items[idx].Label)
		fmt.Println()
		runCommand(menu.Items[idx].Command)
		fmt.Println("The command [" + menu.Items[idx].Command + "] has finshed executing")
		time.Sleep(5 * time.Second)
	}
	clearScreen()
}

func createExampleMenu(filePath string) {

	exampleMenu := Menu{
		Title: struct {
			Label string `json:"label"`
			Color string `json:"color"`
		}{
			Label: "My Menu",
			Color: "blue",
		},
		Items: []MenuItem{
			{
				Label:   "Option 1",
				Command: "ls -lha",
				Color:   "green",
			},
		},
	}

	saveMenu(filePath, exampleMenu)
}

func loadMenu(filePath string) (*Menu, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var menu Menu
	err = json.Unmarshal(file, &menu)
	if err != nil {
		return nil, err
	}

	return &menu, nil
}

func saveMenu(filePath string, menu Menu) error {

	fmt.Println(menu)

	menuJSON, err := json.MarshalIndent(menu, "", "  ")

	if err != nil {
		fmt.Println("Error marshaling menu to JSON:", err)
		return err
	}
	fmt.Println("writing example config to " + filePath)

	err = os.WriteFile(filePath, menuJSON, 0644)

	if err != nil {
		fmt.Println("Error writeing file")
		return err
	}

	return nil
}

func printColoredText(color, text string) {
	colorCode := getColorCode(color)
	fmt.Printf("\033[%sm%s\033[0m", colorCode, text)
}

func getColorCode(color string) string {
	switch color {
	case "black":
		return "30"
	case "red":
		return "31"
	case "green":
		return "32"
	case "yellow":
		return "33"
	case "blue":
		return "34"
	case "magenta":
		return "35"
	case "cyan":
		return "36"
	case "white":
		return "37"
	default:
		return "0"
	}
}

func printMenu(title, titleColor string, items []MenuItem) {

	// Calculate the padding on both sides to center the title
	padding := (51 - len(title)) / 2

	fmt.Printf("╔════════════════════════════════════════════════════╗\n")
	fmt.Printf("║%*s", padding, " ")
	printColoredText(titleColor, title)
	fmt.Printf("%*s║\n", 51-padding-len(title)+1, " ")
	fmt.Printf("╟────────────────────────────────────────────────────╢\n")
	for i, item := range items {
		colorCode := item.Color
		fmt.Printf("╟─── %d. ", i+1)
		printColoredText(colorCode, item.Label)
		fmt.Printf("%*s║\n", 45-len(item.Label), " ")
		if i < len(items)-1 {
			fmt.Printf("╟────────────────────────────────────────────────────╢\n")
		}
	}

	fmt.Printf("╚════════════════════════════════════════════════════╝\n")
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func runCommand(command string) {
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error running command: %v\n", err)
	}
}
