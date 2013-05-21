package main

import (
	"flag"
	"fmt"
	"os"
	"path"
)

var output_folder string
var single string
var username string

func init() {
	flag.StringVar(&single, "single", "", "Optionally specify the name of a single scraper")
	flag.StringVar(&username, "user", "", "Specify the ScraperWiki username whose scrapers/data to fetch")
	flag.StringVar(&output_folder, "output", "", "Specify the output folder for your downloads")

}

// Checks if the foldername provided exists, and if not then attempts to create it
func check_folder(folder string) {
	if _, err := os.Stat(folder); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(folder, os.ModeDir|os.ModePerm); err != nil {
				fmt.Println("  Could not create folder ", folder)
				panic(err)
			}
		}
	}
}

func main() {
	var missing_args bool = false

	flag.Parse()

	if username == "" {
		fmt.Println("You must specify a username")
		missing_args = true
	}

	if output_folder == "" {
		fmt.Println("You must specify an output folder")
		missing_args = true
	}

	if missing_args {
		return
	}

	fmt.Println("Checking output folder:", output_folder)
	check_folder(output_folder)

	fmt.Println("Checking the username:", username)
	info, err := getInfo(username)
	if err != nil {
		fmt.Println(err)
		return
	}

	for role := range info.CodeRoles {
		fmt.Printf("Processing scrapers where the user is the %s\n", role)
		for p := range info.CodeRoles[role] {
			scraper_name := info.CodeRoles[role][p]
			fmt.Printf(" + %s\n", scraper_name)

			folder := path.Join(output_folder, scraper_name)
			check_folder(folder)

			if err := getCode(scraper_name, folder); err != nil {
				fmt.Println("    - Failed to get code:", err)
			}

			if err := getDB(scraper_name, folder); err != nil {
				fmt.Println("    - Failed to get database:", err)
			}

		}
	}
}
