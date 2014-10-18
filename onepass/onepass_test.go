package onepass_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/ravenac95/sudolikeaboss/onepass"
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

type MockWebsocketClient struct {
	responseString string
	sendHandler    func(v interface{}) error
}

func (mock *MockWebsocketClient) Connect() error {
	return nil
}

func (mock *MockWebsocketClient) Receive(v interface{}) error {
	switch data := v.(type) {
	case *string:
		*data = mock.responseString
		return nil
	case *[]byte:
		*data = []byte(mock.responseString)
		return nil
	}
	return nil
}

func (mock *MockWebsocketClient) Send(v interface{}) error {
	return nil
}

var _ = Describe("Termpass", func() {
	Describe("Response", func() {
		var (
			response *Response
			err      error
		)

		BeforeEach(func() {
			response, err = LoadResponse(SAMPLE_RESPONSE_0)
			if err != nil {
				panic(err)
			}
		})

		It("should do something", func() {
			password, err := response.GetPassword()
			Expect(password).To(Equal("password"))
			Expect(err).To(BeNil())
		})
	})

	Describe("Client", func() {
		var (
			client              *OnePasswordClient
			mockWebsocketClient *MockWebsocketClient
			err                 error
		)

		BeforeEach(func() {
			mockWebsocketClient = &MockWebsocketClient{}
			client, err = NewCustomClient(mockWebsocketClient, "fakehost")
		})

		It("should connect", func() {
			err := client.Connect()
			Expect(err).To(BeNil())
		})

		It("should send hello command to 1password", func() {
			err := client.Connect()

			mockWebsocketClient.responseString = "{}"

			response, err := client.SendHelloCommand()

			Expect(err).To(BeNil())
			Expect(response).ToNot(BeNil())
		})

		It("should send showPopup command to 1password", func() {
			err := client.Connect()

			mockWebsocketClient.responseString = SAMPLE_RESPONSE_0

			response, err := client.SendShowPopupCommand()

			Expect(err).To(BeNil())
			Expect(response.GetPassword()).To(Equal("password"))
		})
	})
})
