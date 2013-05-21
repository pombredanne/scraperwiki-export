# scraperwiki-export

This tool provides an easy way for you to retrieve your data and code from ScraperWiki Classic.

The quickest way to run it, assuming you have the [go](http://golang.org/doc/install) compiler installed is to do:

  * git clone git@github.com:rossjones/scraperwiki-export.git
  * cd scraperwiki-export
  * export GOPATH=\`pwd\`
  * go run *.go --user YOUR_USERNAME --output ./data 
  * Go make a cup of tea.
  
  
If you don't have the go compiler installed, you'll need to wait for a binary.


## Still to-do

 * Progress info
 * Fetch a couple at a time.
 * See if SW supports range for re-starts
 * Implement fetching of a single item
 * Implement fetching a substring match of the name
 * Better error checking
 * See if we can get private scrapers