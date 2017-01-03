package model

type UserAccount struct {
        Id              string `json:"id"`
        Name            string  `json:"name"`
        ImageData       []byte `json:"imageData"`
        ImageUrl        string  `json:"imageUrl"`
        AccountServedBy string `json:"accountServedBy"`
        ServedBy        string  `json:"servedBy"`
}

type Account struct {
        Id       string `json:"id"`
        Name     string  `json:"name"`
        ServedBy string `json:"servedBy"`
}
