package main

import ("encoding/xml"
    "fmt"
    "io/ioutil"
    "net/http"
    "flag"
    "time"
    "net/smtp"
    "strings"
    "log"
)

type countiesStruct struct {
    Counties []countyStruct `xml:"PLA_BurnBan.dbo.Status"`
}

type countyStruct struct {
    County string `xml:"county"`
    Status string `xml:"travel_status"`
    Time string `xml:"posted_date"`
}

type emailStruct struct {
    Sender string
    Password string
    Recipient string
    County string
    MinuteDelta int
    Message string
    Subject string
}


func main() {
    triggerStatuses := []string{"watch", "warning"}
    emailData := parseFlags()
    if emailData.Sender == "" {
        fmt.Println("Missing sender email address!")
        log.Fatal("Missing sender email address!")
    }
    triggerCounties := strings.Split(strings.ToLower(emailData.County), ",")
    for {
        counties := retrieveCountyData()
        for _, element := range counties.Counties {
            if arrayContains(triggerCounties, strings.ToLower(element.County)) {
                if arrayContains(triggerStatuses, strings.ToLower(element.Status)) {
                    emailData.Message = "The county " + element.County +
                               " has the weather status of " + element.Status +
                               " at " + element.Time
                    emailData.Subject = element.County + ": " + element.Status
                    send(emailData)
                }
            }
        }

        time.Sleep(time.Duration(emailData.MinuteDelta) * time.Minute)
    }
}

func arrayContains(array []string, value string) bool {
    for _, element := range array {
        if element == value {
            return true
        }
    }
    return false
}

func retrieveCountyData() countiesStruct {
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

func send(email emailStruct) {
    msg := "From: " + email.Sender + "\n" +
        "To: " + email.Recipient + "\n" +
        email.Subject +
        email.Message

    err := smtp.SendMail("smtp.gmail.com:587",
        smtp.PlainAuth("", email.Sender, email.Password, "smtp.gmail.com"),
        email.Sender, []string{email.Recipient}, []byte(msg))

    if err != nil {
        log.Fatal(err)
    }
}

func parseFlags() emailStruct {
    countyptr := flag.String("County", "Hamilton", "County you want to know the weather for.")
    minuteptr := flag.Int("Minutes", 15, "How often you want to check for weather updates.")
    senderemailptr := flag.String("Sender", "", "email to send notification emails. (Enable less secure apps.)")
    passwordptr := flag.String("Password", "", "Password for your sending e-mail")
    recipientemailptr := flag.String("Recipient", "", "email to send notification emails. (Enable less secure apps.)")

    flag.Parse()
    return emailStruct{
        Sender: *senderemailptr,
        Password: *passwordptr,
        Recipient: *recipientemailptr,
        County: *countyptr,
        MinuteDelta: *minuteptr,
    }
}
