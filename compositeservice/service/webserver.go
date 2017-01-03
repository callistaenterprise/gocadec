package service

import (
        "net/http"
        ct "github.com/eriklupander/cloudtoolkit"
)

func StartWebServer(port string) {
        r := NewRouter()
        http.Handle("/", r)

        ct.Log.Println("Starting HTTP service at " + port)
        err := http.ListenAndServe(":" + port, nil)

        if err != nil {
                ct.Log.Errorln("An error occured starting HTTP listener at port " + port)
                ct.Log.Errorln("Error: " + err.Error())
        }
}
