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
                }

            }
        }

        time.Sleep(time.Duration(*minuteptr) * time.Minute)
    }
}

func send(body string) {
    from := "@gmail.com"
    pass := ""
    to := "kbburkholder@gmail.com"

    msg := "From: " + from + "\n" +
        "To: " + to + "\n" +
        "Subject: Hello there\n\n" +
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
