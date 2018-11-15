package main

import ("encoding/xml"
    "fmt"
    "io/ioutil"
    "net/http"
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
            if element.County == "Hamilton" {
                fmt.Println(element.County)
                fmt.Println(element.Status)
                fmt.Println(element.Time)
            }

        }
    }

}
