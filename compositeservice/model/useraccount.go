package model

type UserAccount struct {
        Id              string `json:"id"`
        Name            string  `json:"name"`
        ImageData       []byte `json:"imageData"`
        ImageUrl        string  `json:"imageUrl"`
        AccountServedBy string `json:"accountServedBy"`
        ImageServedBy   string  `json:"imageServedBy"`
}

type Account struct {
        Id       string `json:"id"`
        Name     string  `json:"name"`
        ServedBy string `json:"servedBy"`
}
