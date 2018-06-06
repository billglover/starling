package starling

import (
	"bytes"
	"net/http"
	"testing"
)

var validateTestCases = []struct {
	body      string
	secret    string
	signature string
	valid     bool
}{
	{
		body:      "this is the request body",
		secret:    "1234567890",
		signature: "05pnHTd02EsPBgaFi7EFB7lUHQo1RTKVFUBrcXLbrHNNft6G34v4qMAak6rjO0hqwoW9a4DpQ5X8Hc65/JHuUw==",
		valid:     true,
	},
	{
		body:      "this is the request body",
		secret:    "1234567890",
		signature: "[invalid]05pnHTd02EsPBgaFi7EFB7lUHQo1RTKVFUBrcXLbrHNNft6G34v4qMAak6rjO0hqwoW9a4DpQ5X8Hc65/JHuUw==",
		valid:     false,
	},
	{
		body:      "[invalid]this is the request body",
		secret:    "1234567890",
		signature: "05pnHTd02EsPBgaFi7EFB7lUHQo1RTKVFUBrcXLbrHNNft6G34v4qMAak6rjO0hqwoW9a4DpQ5X8Hc65/JHuUw==",
		valid:     false,
	},
	{
		body:      "this is the request body",
		secret:    "[invalid]1234567890",
		signature: "05pnHTd02EsPBgaFi7EFB7lUHQo1RTKVFUBrcXLbrHNNft6G34v4qMAak6rjO0hqwoW9a4DpQ5X8Hc65/JHuUw==",
		valid:     false,
	},
}

func TestValidate_Valid(t *testing.T) {
	for _, tc := range validateTestCases {
		body := []byte(tc.body)
		req, err := http.NewRequest("POST", "http://localhost/callback", bytes.NewBuffer(body))
		if err != nil {
			t.Error("should create a request without error:", err)
		}

		req.Header.Set("X-Hook-Signature", tc.signature)
		valid, err := Validate(req, tc.secret)
		if err != nil {
			t.Error("should be able to perform validation without error:", err)
		}

		if valid != tc.valid {
			t.Error("should return:", tc.valid)
		}
	}
}
