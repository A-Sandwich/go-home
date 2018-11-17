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


func main() {
    countyptr := flag.String("County", "Hamilton", "County you want to know the weather for.")
    minuteptr := flag.Int("Minutes", 15, "How often you want to check for weather updates.")
    senderemailptr := flag.String("e-mail", "", "e-mail to send notification emails. (Enable less secure apps.)")
    passwordptr := flag.String("e-mail password", "", "Password for your sending e-mail")
    recipientemailptr := flag.String("e-mail", "", "e-mail to send notification emails. (Enable less secure apps.)")
    flag.Parse()

    fmt.Println(*countyptr)
    fmt.Println(*minuteptr)

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
                if element.County == *countyptr {
                    fmt.Println(element.County)
                    fmt.Println(element.Status)
                    fmt.Println(element.Time)
                    if *senderemailptr != "" {
                        message := "The county " + element.County +
                                   " has the weather status of " + element.Status +
                                   " at " + element.Time
                        subject := element.County + ": " + element.Status
                        send(subject, message, *senderemailptr, *passwordptr, *recipientemailptr)
                    }
                }

            }
        }

        time.Sleep(time.Duration(*minuteptr) * time.Minute)
    }
}

func send(subject string, body string, sender string, password string, recipient string) {
    from := sender
    pass := password
    to := recipient

    msg := "From: " + from + "\n" +
        "To: " + to + "\n" +
        subject +
        body

    err := smtp.SendMail("smtp.gmail.com:587",
        smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
        from, []string{to}, []byte(msg))

    if err != nil {
//        log.Printf("smtp error: %s", err)
        return
    }
 //   log.Print("sent, visit http://foobarbazz.mailinator.com")
}
