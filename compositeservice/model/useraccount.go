package model

type UserAccount struct {
        Id              string `json:"id"`
        Name            string `json:"name"`
        ImageData       []byte `json:"imageData"`
        ImageUrl        string `json:"imageUrl"`
        QuoteOfTheDay   string `json:"quoteOfTheDay"`
        AccountServedBy string `json:"accountServedBy"`
        ImageServedBy   string `json:"imageServedBy"`
        QuoteServedBy   string `json:"quoteServedBy"`
}

type Account struct {
        Id       string `json:"id"`
        Name     string `json:"name"`
        ServedBy string `json:"servedBy"`
}

type Quote struct {
        IpAddress string `json:"ipAddress"`
        Quote     string `json:"quote"`
        Language  string `json:"language"`
}
