package main

import (
"encoding/json"
"fmt"
"io/ioutil"
"net/http"
"net/url"
"strings"
"time"
)

const (
indeedAPIEndpoint = "http://api.indeed.com/ads/apisearch"
)

type jobListing struct {
JobTitle  string `json:"jobtitle"`
Company   string `json:"company"`
Location  string `json:"formattedLocation"`
Salary    string `json:"salary"`
URL       string `json:"url"`
Age       int    `json:"age"`
Score     int
}

type jobListings struct {
Listings []jobListing `json:"results"`
}

type resume struct {
Name        string
Email       string
PhoneNumber string
Experience  []string
Education   []string
Skills      []string
}

func main() {
// Prompt user for resume details
fmt.Println("Enter resume details:")
resume := promptResume()

// Fetch job listings from the Indeed API
listings, err := fetchJobListings(resume)
if err != nil {
fmt.Println("Error fetching job listings:", err)
return
}

// Score job listings based on user's qualifications
for i := range listings.Listings {
score := calculateJobScore(&listings.Listings[i], resume)
listings.Listings[i].Score = score
}

// Sort job listings by score
sortJobListings(&listings.Listings)

// Display top job listings
fmt.Println("Top job listings:")
for i, listing := range listings.Listings[:10] {
fmt.Printf("%d. %s (%s) - %s - %s - %d days old\n", i+1, listing.JobTitle, listing.Company, listing.Location, listing.Salary, listing.Age)
}
}

func promptResume() resume {
var experience, education, skills string
fmt.Print("Name: ")
name := scanLine()
fmt.Print("Email: ")
email := scanLine()
fmt.Print("Phone number: ")
phone := scanLine()
fmt.Print("Experience (comma-separated list): ")
experience = scanLine()
fmt.Print("Education (comma-separated list): ")
education = scanLine()
fmt.Print("Skills (comma-separated list): ")
skills = scanLine()
return resume{
Name:        name,
Email:       email,
PhoneNumber: phone,
Experience:  strings.Split(experience, ","),
Education:   strings.Split(education, ","),
Skills:      strings.Split(skills, ","),
}
}

func scanLine() string {
var line string
fmt.Scanln(&line)
return line
}

func fetchJobListings(r resume) (*jobListings, error) {
// Build API request URL
query := url.Values{}
query.Set("q", "")
query.Set("l", "")
query.Set("sort", "date")
query.Set("format", "json")
query.Set("publisher", "YOUR_PUBLISHER_ID") // replace with your Indeed API publisher ID
url := fmt.Sprintf("%s?%s", indeedAPIEndpoint, query.Encode())

// Send API request
resp, err := http.Get(url)
if err != nil {
return nil, err
}
defer resp.Body.Close()

// Parse API response
body, err := ioutil.ReadAll(resp.Body)
if err != nil {
return nil, err
}
var listings jobListings
err = json.Unmarshal(body, &listings)
if err != nil {
return nil, err
}

return

&listings, nil
}

func calculateJobScore(listing *jobListing, r resume) int {
score := 0
score += matchKeywords(listing.JobTitle, r.Experience, 5)
score += matchKeywords(listing.JobTitle, r.Education, 5)
score += matchKeywords(listing.JobTitle, r.Skills, 3)
score += matchKeywords(listing.Company, r.Experience, 2)
score += matchKeywords(listing.Company, r.Education, 2)
score += matchKeywords(listing.Company, r.Skills, 1)
score += matchKeywords(listing.Location, r.Experience, 1)
score += matchKeywords(listing.Location, r.Education, 1)
score += matchKeywords(listing.Location, r.Skills, 1)
score += matchSalary(listing.Salary)
score += matchAge(listing.Age)
return score
}

func matchKeywords(str string, keywords []string, weight int) int {
score := 0
for _, keyword := range keywords {
if containsWord(str, keyword) {
score += weight
}
}
return score
}

func matchSalary(salary string) int {
if salary == "" {
return 0
}
return 10
}

func matchAge(age int) int {
ageDays := time.Now().Sub(time.Unix(int64(age), 0)).Hours() / 24
if ageDays < 7 {
return 10
} else if ageDays < 30 {
return 5
} else if ageDays < 90 {
return 2
}
return 0
}

func splitKeywords(keywords string) []string {
return strings.Split(strings.ToLower(keywords), " ")
}

func containsWord(str, word string) bool {
return strings.Contains(strings.ToLower(str), word)
}

func sortJobListings(listings *[]jobListing) {
sort.SliceStable(*listings, func(i, j int) bool {
return listings[i].Score > listings[j].Score
})
}
