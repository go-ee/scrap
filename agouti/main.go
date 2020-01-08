package main

import (
	"github.com/sclevine/agouti"
	"log"
)

func main() {
	//driver := agouti.GeckoDriver()

	command := []string{"/home/z000ru5y/dev/d/ee-grab/ee-grab/drivers/geckodriver", "--port={{.Port}}"}
	driver := agouti.NewWebDriver("http://{{.Address}}", command,
		"browser.download.folderList=2")

	if err := driver.Start(); err != nil {
		log.Fatal("Failed to start Selenium:", err)
	}
	page, err := driver.NewPage(agouti.Browser("firefox"))
	if err != nil {
		log.Fatal("Failed to open page:", err)
	}

	folderName := "SSWA12"
	url := "https://wse02.siemens.com/content/P0002864/SSWA/Shared Documents/Forms/AllItems.aspx?" +
		"RootFolder=%2Fcontent%2FP0002864%2FSSWA%2FShared Documents" +
		"%2FSSWA Participant Folders%2F" + folderName

	if err := page.Navigate(url); err != nil {
		//log.Fatal("Failed to navigate:", err)
	}

	loginURL, err := page.URL()
	if err != nil {
		log.Fatal("Failed to get page URL:", err)
	}

	expectedLoginURL := "http://localhost:3000/login"
	if loginURL != expectedLoginURL {
		log.Fatal("Expected URL to be", expectedLoginURL, "but got", loginURL)
	}

	loginPrompt, err := page.Find("#prompt").Text()
	if err != nil {
		log.Fatal("Failed to get login prompt text:", err)
	}

	expectedPrompt := "Please login."
	if loginPrompt != expectedPrompt {
		log.Fatal("Expected login prompt to be", expectedPrompt, "but got", loginPrompt)
	}

	if err := driver.Stop(); err != nil {
		log.Fatal("Failed to close pages and stop WebDriver:", err)
	}
}
