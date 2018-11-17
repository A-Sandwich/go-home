package main

import ("encoding/xml"
    "fmt"
    "io/ioutil"
    "net/http"
    "flag"
    "time"
    "net/smtp"
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
    emailData := parseFlags()

    for {
        resp, err := http.Get("https://www.in.gov/ai/dhs/dhs_travel_advisory.txt")
        if err != nil {
            fmt.Println("A HELP, get failed!")
        }

        defer resp.Body.Close()
        body, err := ioutil.ReadAll(resp.Body)

        if err != nil {
            fmt.Println("Borked http body.")
        }

        var counties countiesStruct
        err = xml.Unmarshal(body, &counties)
        if err != nil {
            fmt.Println("Unable to unmrashal XML")
        } else {
            for _, element := range counties.Counties {
                if element.County == emailData.County {
                    fmt.Println(element.County)
                    fmt.Println(element.Status)
                    fmt.Println(element.Time)
                    if emailData.Sender != "" {
                        emailData.Message = "The county " + element.County +
                                   " has the weather status of " + element.Status +
                                   " at " + element.Time
                        emailData.Subject = element.County + ": " + element.Status
                        send(emailData)
                    }
                }
            }
        }

        time.Sleep(time.Duration(emailData.MinuteDelta) * time.Minute)
    }
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
        return
    }
}

func parseFlags() emailStruct {
    countyptr := flag.String("County", "Hamilton", "County you want to know the weather for.")
    minuteptr := flag.Int("Minutes", 15, "How often you want to check for weather updates.")
    senderemailptr := flag.String("e-mail", "", "e-mail to send notification emails. (Enable less secure apps.)")
    passwordptr := flag.String("e-mail password", "", "Password for your sending e-mail")
    recipientemailptr := flag.String("e-mail", "", "e-mail to send notification emails. (Enable less secure apps.)")
    flag.Parse()

    return emailStruct{
        Sender: *senderemailptr,
        Password: *passwordptr,
        Recipient: *recipientemailptr,
        County: *countyptr,
        MinuteDelta: *minuteptr,
    }
}
