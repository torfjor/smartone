package smartone

import "time"

// Cache represents a response that should be cached by Tidomat smartOne
type Cache struct {
	Id         string    `json:"id"`
	ValidUntil time.Time `json:"validUntil"`
	Label1     string    `json:"label1,omitempty"`
	Label2     string    `json:"label2,omitempty"`
	Label3     string    `json:"label3,omitempty"`
	Label4     string    `json:"label4,omitempty"`
	Label5     string    `json:"label5,omitempty"`
	Unique     string    `json:"unique,omitempty"`
	PIN        string    `json:"pin,omitempty"`
}

type OnCardResponse struct {
	Result        bool   `json:"result"`
	ResultMessage string `json:"resultMessage,omitempty"`
	PIN           string `json:"pin,omitempty"`
	Cache         *Cache `json:"cache,omitempty"`
}
