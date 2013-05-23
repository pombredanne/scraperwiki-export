package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
)

type Scrapers struct {
	Owner  []string
	Editor []string
}

type InfoDict struct {
	Username    string
	ProfileName string
	CodeRoles   map[string][]string
	DateJoined  string
}

type UserInfoList struct {
	Information []InfoDict
}

// Retrieves information about the user, most notably a list of scraper names that
// they are either the owner or editors of.  These will be used to fetch the code+data
func getInfo(username string) (InfoDict, error) {
	address := fmt.Sprintf("https://api.scraperwiki.com/api/1.0/scraper/getuserinfo?format=jsondict&username=%s", username)

	cresp, err := http.Get(address)
	if err != nil {
		return InfoDict{}, err
	}
	defer cresp.Body.Close()

	dec := json.NewDecoder(cresp.Body)

	var items []InfoDict
	if err := dec.Decode(&items); err != nil {
		return InfoDict{}, err
	}

	if len(items) == 0 {
		return InfoDict{}, errors.New("User not found")
	}

	return items[0], nil
}

func copy_db(reader io.Reader, writer io.Writer, length int64) (int64, error) {
	var read int64
	var p float32
	for {
		buffer := make([]byte, 2097152)
		cBytes, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		read = read + int64(cBytes)
		p = float32(read) / float32(length) * 100
		fmt.Printf("\r  progress: %3.2f                ", p)
		writer.Write(buffer[0:cBytes])
	}
	fmt.Printf("\r\n")
	return read, nil
}

// Attempts to retrieve the database (if one exists) for the provided scraper
// downloading to a folder named after that scraper.  It won't fetch anything
// when SW says the file is empty, we already have a file with the same name
// and exact size or it fails.
func getDB(name string, output_folder string) error {

	address := fmt.Sprintf("https://scraperwiki.com/scrapers/export_sqlite/%s/", name)

	resp, err := http.Head(address)
	defer resp.Body.Close()

	length, _ := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 0)
	if length == 0 {
		fmt.Println("  Skipping download, no data")
		return nil
	}

	output_file := path.Join(output_folder, name+".sqlite")

	// Check if the file already exists and how large it is
	st, err := os.Stat(output_file)
	if err == nil {
		if st.Size() == length {
			fmt.Println("  Skipping download, already have data")
			return nil
		}
	}

	f, err := os.Create(output_file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	cresp, err := http.Get(address)
	defer cresp.Body.Close()

	_, err = copy_db(cresp.Body, f, length)
	fmt.Println("  Wrote DB to " + output_file)
	return err
}

// When provided with a scraper name, and an output folder this function makes
// a call to the ScraperWiki API to find information about the scraper.  It only
// uses the language and code sections of the response, but will then write out
// the code with the right extension to scraper_name/scraper_name.EXT where .EXT
// is determined by the language
func getCode(name string, output_folder string) error {
	address := "https://api.scraperwiki.com/api/1.0/scraper/getinfo?format=jsondict&version=-1&quietfields=attachable_here%7Cattachables%7Ctags%7Clast_run%7chistory%7Cdatasummary%7Cuserroles%7Crunevents%7Clast_run&name=" + name

	resp, err := http.Get(address)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var items []map[string]interface{}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&items); err != nil {
		return errors.New("Failed to retrieve info for that scraper, does it exist?")
	}

	languages := map[string]string{
		"python": ".py",
		"ruby":   ".rb",
		"php":    ".php",
		"html":   ".html",
	}

	code := fmt.Sprintf("%v", items[0]["code"])
	if len(code) == 0 {
		fmt.Println("  Skipping writing code as there is none")
		return nil
	}

	language := fmt.Sprintf("%v", items[0]["language"])
	output_file := path.Join(output_folder, name+languages[language])
	f, err := os.Create(output_file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	l, err := f.WriteString(code)
	if err != nil {
		panic("Failed to write code to file")
	}
	fmt.Printf("  Code is %d bytes in size : %s\n", l, output_file)
	return nil
}
