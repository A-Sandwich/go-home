package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"strings"
	"time"
    "os"
    "strconv"
)

/*
* This function is my main.
* Function iterates over and county checks and sleeps
 */
func main() {
	emailData := parseFlags()

	if emailData.Sender == "" {
		fmt.Println("Missing sender email address!")
		log.Fatal("Missing sender email address!")
	}

	for {
		checkMonitoredCountiesWeather(emailData)
		time.Sleep(time.Duration(emailData.MinuteDelta) * time.Minute)
	}
}

/*
* Function checks county data and sends an email if it is triggered.
 */
func checkMonitoredCountiesWeather(emailData emailStruct) {
	counties := retrieveMonitoredCountiesData()
	for _, county := range counties.Counties {
		if areMonitoredCountiesDangerous(county, emailData.MonitoredCounties) {
			fmt.Println(county.Name + " triggered email with status of '" +
				county.Status)
			emailData.Message = "The county " + county.Name +
				" has the weather status of " + county.Status +
				" at " + county.Time
			emailData.Subject = county.Name + ": " + county.Status
			send(emailData)
		}
	}
	fmt.Println("sleep")
}

/*
* Function verifies if a county should trigger an email.
 */
func areMonitoredCountiesDangerous(county countyStruct, counties string) bool {
	triggerStatuses := []string{"watch", "warning"}
	triggerCounties := strings.Split(strings.ToLower(counties), ",")
	return arrayContains(triggerCounties,
		strings.ToLower(county.Name)) &&
		arrayContains(triggerStatuses, strings.ToLower(county.Status))
}

/*
* Function checks if an array contains a string.
 */
func arrayContains(array []string, value string) bool {
	for _, county := range array {
		if county == value {
			return true
		}
	}
	return false
}

/*
* Function gets county data from .gov endpoint.
 */
func retrieveMonitoredCountiesData() countiesStruct {
	resp, err := http.Get("https://www.in.gov/ai/dhs/dhs_travel_advisory.txt")
	if err != nil {
		fmt.Println("A HELP, get failed!")
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Borked http body.")
		log.Fatal(err)
	}

	var counties countiesStruct
	err = xml.Unmarshal(body, &counties)
	if err != nil {
		fmt.Println("Unable to unmrashal XML")
		log.Fatal(err)
	}
	return counties
}

/*
* This function sends an email with a counties current status.
* This function was initially taken from:
* https://gist.github.com/jpillora/cb46d183eca0710d909a
* and has been modified for use in this project. Thanks jpillora!
 */
func send(email emailStruct) {
	msg := "From: " + email.Sender + "\n" +
		"To: " + email.Recipient + "\n" +
		"Subject: " + email.Subject + "\n\n" +
		email.Message

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", email.Sender, email.Password, "smtp.gmail.com"),
		email.Sender, []string{email.Recipient}, []byte(msg))

	if err != nil {
		log.Fatal(err)
	}
}

/*
* This function parses flags and populates a struct containing the parsed data.
 */
func parseFlags() emailStruct {
    defaultData := parseEnvironmentVariables()

	countyptr := flag.String("MonitoredCounties",
		defaultData.MonitoredCounties,
		"MonitoredCounties you want to know the weather"+
			"for. E.G. 'Hamilton' or 'Cass,Whitley'")

	minuteptr := flag.Int("Minutes",
		defaultData.MinuteDelta,
		"How often you want to check for weather updates in minutes.")

	senderemailptr := flag.String("Sender",
		defaultData.Sender,
		"email to send notification emails."+
			" (Enable less secure apps.)")

	passwordptr := flag.String("Password",
		defaultData.Password,
		"Password for your sending e-mail")

	recipientemailptr := flag.String("Recipient",
		defaultData.Recipient,
		"email to send notification emails."+
			" (Enable less secure apps.)")

	flag.Parse()
    fmt.Println(*minuteptr)
	return emailStruct {
		Sender:            *senderemailptr,
		Password:          *passwordptr,
		Recipient:         *recipientemailptr,
		MonitoredCounties: *countyptr,
		MinuteDelta:       *minuteptr,
	}
}

/*
* Function parses environment variables and puts them into an emailStruct.
*/
func parseEnvironmentVariables() emailStruct {
    minutesStr := getEnvironmentDefault("GO_HOME_MINUTES", "15")
    minutes, _ := strconv.ParseInt(minutesStr, 10, 32)

    return emailStruct {
        Sender:            getEnvironmentDefault("GO_HOME_SENDER", ""),
        Password:          getEnvironmentDefault("GO_HOME_PASSWORD", ""),
        Recipient:         getEnvironmentDefault("GO_HOME_RECIPIENT", ""),
        MonitoredCounties: getEnvironmentDefault("GO_HOME_COUNTIES",
                                                 "Hamilton"),
        MinuteDelta:       int(minutes),
    }
}

/*
* Function returns environment variable value if it exists, otherwise this
* function returns the default value passed in.
*/
func getEnvironmentDefault(environmentKey string, defaultValue string) string {
    value, valid := os.LookupEnv(environmentKey)
    if valid {
        return value
    }
    return defaultValue
}

type countiesStruct struct {
	Counties []countyStruct `xml:"PLA_BurnBan.dbo.Status"`
}

type countyStruct struct {
	Name   string `xml:"county"`
	Status string `xml:"travel_status"`
	Time   string `xml:"posted_date"`
}

type emailStruct struct {
	Sender            string
	Password          string
	Recipient         string
	MonitoredCounties string
	MinuteDelta       int
	Message           string
	Subject           string
}
