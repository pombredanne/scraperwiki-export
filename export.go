package main

import (
	"flag"
	"fmt"
	"os"
	"path"
)

const USAGE = `
scraperwiki-export is a simple program that will allow you to export your 
data and code from ScraperWiki Classic. You may export just a single scraper
or with the correct arguments download all of your work from ScraperWiki.

To download a single scraper, you could try the following line in your terminal:

    scraperwiki-export -output mydata -single <short name of scraper>

You can get the scraper's short name from the url, in the example below the 
short name of the scraper is 'hospital_smr_data'

    https://scraperwiki.com/scrapers/hospital_smr_data/

To download all of your data, you need to know your ScraperWiki Classic username
and replace it in the following line:

    scraperwiki-export --output mydata -username <username>

Bulk uploads may take a while if you have a lot of data in your scrapers, but
there should be a simple progress meter to show you how far along the process 
is for each scraper.

WARNING: This program does not currently allow you to retrieve scrapers and
data that are stored in a 'vault'. 

`

var output_folder string
var single string
var username string

func init() {
	flag.StringVar(&single, "single", "", "Optionally specify the name of a single scraper")
	flag.StringVar(&username, "user", "", "Specify the ScraperWiki username whose scrapers/data to fetch")
	flag.StringVar(&output_folder, "output", "", "Specify the output folder for your downloads")

}

func usage() {
	fmt.Fprintf(os.Stderr, "\nusage: scraperwiki-export [-single <scraper_name>] -username <username> -output <folder>\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, USAGE)
	os.Exit(0)
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

func process_scraper(name string) {
	fmt.Printf("+ %s\n", name)

	folder := path.Join(output_folder, name)
	check_folder(folder)

	if err := getCode(name, folder); err != nil {
		fmt.Println("  - Failed to get code:", err)
		return
	}

	if err := getDB(name, folder); err != nil {
		fmt.Println("  - Failed to get database:", err)
		return
	}
}

func main() {
	var missing_args bool = false

	flag.Usage = usage
	flag.Parse()

	// Ensure we have enough arguments to run properly.
	if username == "" && single == "" {
		fmt.Println("Error: You must specify a username to bulk download")
		missing_args = true
	}

	if output_folder == "" {
		fmt.Println("Error: You must specify an output folder")
		missing_args = true
	}

	if missing_args {
		fmt.Println(USAGE)
		return
	}

	// Create the output folder if it doesn't exist
	fmt.Println("Checking output folder:", output_folder)
	check_folder(output_folder)

	// We  don't care about the user if we are fetching a single
	// scraper, so let's save ourselves a http request by not
	// verifying the username.
	if single != "" {
		process_scraper(single)
		return
	}

	fmt.Println("Getting user details for:", username)
	info, err := getInfo(username)
	if err != nil {
		fmt.Println(err)
		return
	}

	for role := range info.CodeRoles {
		fmt.Printf("Processing %s scrapers\n", role)
		for p := range info.CodeRoles[role] {
			scraper_name := info.CodeRoles[role][p]
			process_scraper(scraper_name)
		}
	}
}
