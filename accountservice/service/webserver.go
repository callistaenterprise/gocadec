package service

import (
        "net/http"
        "log"
        ct "github.com/eriklupander/cloudtoolkit"
)

func StartWebServer(port string) {

        r := NewRouter()
        http.Handle("/", r)

        err := http.ListenAndServe(":" + port, nil)

        ct.Log.Println("Starting HTTP service at " + port)

        if err != nil {
                log.Println("An error occured starting HTTP listener at port " + port)
                log.Println("Error: " + err.Error())
        }
}
