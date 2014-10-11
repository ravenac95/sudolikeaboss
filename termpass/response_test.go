package termpass

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const SAMPLE_RESPONSE_0 = `
{
  "action": "fillItem",
  "payload": {
    "openInTabMode": "NewTab",
    "options": {
      "animate":true,
      "autosubmit":true
    }, 
    "item": {
      "uuid":"someuuid",
      "nakedDomains": ["somedomain.com"],
      "overview": {
        "title": "title",
        "url": "url"
      },
      "secureContents": {
        "htmlForm": {"htmlMethod":"post"},
        "fields": [
          {
            "value":"username",
            "id":"email",
            "name":"email",
            "type":"T",
            "designation":"username"
          },
          {
            "value":"password",
            "id":"password",
            "name":"password",
            "type":"P",
            "designation":"password"
          },
          {
            "value":"Login",
            "id":"",
            "name":"",
            "type":"I"
          }
        ]
      }
    },
    "action":"fillLogin"
  },
  "version":"01"
}
`

func TestParse1passwordResponse(t *testing.T) {
	actualPassword, err := getPasswordFromResponse(SAMPLE_RESPONSE_0)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, actualPassword, "password")
}
