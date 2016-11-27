package onepass

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	log "github.com/Sirupsen/logrus"

	"github.com/satori/go.uuid"

	"github.com/ravenac95/sudolikeaboss/websocketclient"
)

type Command struct {
	Action   string  `json:"action"`
	Number   int     `json:"number,omitempty"`
	Version  string  `json:"version,omitempty"`
	BundleID string  `json:"bundleId,omitempty"`
	Payload  Payload `json:"payload"`
}

type Payload struct {
	Version      string            `json:"version,omitempty"`
	ExtID        string            `json:"extId,omitempty"`
	Method       string            `json:"method,omitempty"`
	Secret       string            `json:"secret,omitempty"`
	Capabilities []string          `json:"capabilities,omitempty"`
	URL          string            `json:"url,omitempty"`
	Options      map[string]string `json:"options,omitempty"`
	CC           string            `json:"cc,omitempty"`
	CS           string            `json:"cs,omitempty"`
	M4           string            `json:"M4,omitempty"`
	M3           string            `json:"M3,omitempty"`
	Algorithm    string            `json:"alg,omitempty"`
	Iv           string            `json:"iv,omitempty"`
	Data         string            `json:"data,omitempty"`
	Hmac         string            `json:"hmac,omitempty"`
}

type WebsocketClient interface {
	Connect() error
	Receive(v interface{}) error
	Send(v interface{}) error
}

// Configuration struct
type Configuration struct {
	WebsocketURI      string `json:"websocketUri"`
	WebsocketProtocol string `json:"websocketProtocol"`
	WebsocketOrigin   string `json:"websocketOrigin"`
	DefaultHost       string `json:"defaultHost"`
	StateDirectory    string `json:"stateDirectory"`
}

type OnePasswordClient struct {
	DefaultHost             string
	websocketClient         WebsocketClient
	StateDirectory          string
	number                  int
	extID                   string
	secret                  []byte
	cc                      []byte
	cs                      []byte
	sessionHmacK            []byte
	sessionEncK             []byte
	base64urlWithoutPadding *b64.Encoding
}

type StateFileConfig struct {
	Secret string `json:"secret"`
	ExtID  string `json:"extID"`
}

func NewClientWithConfig(configuration *Configuration) (*OnePasswordClient, error) {
	return NewClient(configuration.WebsocketURI, configuration.WebsocketProtocol, configuration.WebsocketOrigin, configuration.DefaultHost, configuration.StateDirectory)
}

func NewClient(websocketUri string, websocketProtocol string, websocketOrigin string, defaultHost string, stateDirectory string) (*OnePasswordClient, error) {
	websocketClient := websocketclient.NewClient(websocketUri, websocketProtocol, websocketOrigin)

	return NewCustomClient(websocketClient, defaultHost, stateDirectory)
}

func NewCustomClient(websocketClient WebsocketClient, defaultHost string, stateDirectory string) (*OnePasswordClient, error) {
	client := OnePasswordClient{
		websocketClient: websocketClient,
		DefaultHost:     defaultHost,
		StateDirectory:  stateDirectory,
	}

	base64urlWithoutPadding := b64.URLEncoding.WithPadding(b64.NoPadding)
	client.base64urlWithoutPadding = base64urlWithoutPadding

	// Load the state directory if stuff is in there
	err := client.LoadOrSetupState()

	if err != nil {
		return nil, err
	}

	err = client.Connect()

	if err != nil {
		return nil, err
	}

	return &client, nil
}

func (client *OnePasswordClient) LoadOrSetupState() error {
	stateFilePath := path.Join(client.StateDirectory, "state.json")
	stateFileExists, err := Exists(stateFilePath)

	if err != nil {
		return err
	}

	if stateFileExists {
		stateFileStr, err := ioutil.ReadFile(stateFilePath)
		if err != nil {
			return err
		}

		var stateFileConfig StateFileConfig

		err = json.Unmarshal(stateFileStr, &stateFileConfig)

		if err != nil {
			return err
		}

		secret, err := client.base64urlWithoutPadding.DecodeString(stateFileConfig.Secret)

		if err != nil {
			return err
		}

		client.extID = stateFileConfig.ExtID
		client.secret = secret
	} else {
		err := EnsureDir(client.StateDirectory)

		if err != nil {
			return err
		}

		extIDBytes := uuid.NewV4()
		extID := extIDBytes.String()
		client.extID = extID

		secret, err := GenerateRandomBytes(32)

		if err != nil {
			return err
		}
		client.secret = secret

		stateFileConfig := StateFileConfig{
			ExtID:  extID,
			Secret: client.base64urlWithoutPadding.EncodeToString(secret),
		}

		stateFileStr, err := json.Marshal(&stateFileConfig)
		if err != nil {
			return err
		}

		return ioutil.WriteFile(stateFilePath, []byte(stateFileStr), 0700)
	}
	return nil
}

func (client *OnePasswordClient) Connect() error {
	err := client.websocketClient.Connect()

	return err
}

func (client *OnePasswordClient) SendShowPopupCommand() (*Response, error) {
	payload := Payload{
		URL:     client.DefaultHost,
		Options: map[string]string{"source": "toolbar-button"},
	}

	command := client.createCommand("showPopup", payload)

	response, err := client.SendEncryptedCommand(command)

	if err != nil {
		return nil, err
	}

	decryptedPayloadRaw, err := client.decryptResponse(response)

	if err != nil {
		return nil, err
	}

	var decryptedPayload ResponsePayload

	json.Unmarshal(decryptedPayloadRaw, &decryptedPayload)

	response.Payload = decryptedPayload

	return response, nil
}

func (client *OnePasswordClient) createCommand(action string, payload Payload) *Command {
	command := Command{
		Action: action,
		//Number:   client.number,
		Version: "4.6.2.90",
		//BundleID: "com.sudolikeaboss.sudolikeaboss",
		Payload: payload,
	}

	// Increment the number (it's a 1password thing that I saw whilst listening
	// to their commands
	client.number += 1
	return &command
}

func (client *OnePasswordClient) SendHelloCommand() (*Response, error) {
	capabilities := make([]string, 2)
	capabilities[0] = "auth-sma-hmac256"
	capabilities[1] = "aead-cbchmac-256"

	payload := Payload{
		Version:      "4.6.2.90",
		ExtID:        client.extID,
		Capabilities: capabilities,
	}

	command := client.createCommand("hello", payload)

	response, err := client.SendCommand(command)
	if err != nil {
		return nil, err
	}

	if response.Action != "authNew" && response.Action != "authBegin" {
		errorMsg := fmt.Sprintf("Unexpected response: %s", response.Action)
		err = errors.New(errorMsg)
		return nil, err
	}

	return response, nil
}

func (client *OnePasswordClient) hmacSignWithSecret(dataToSign ...[]byte) []byte {
	return HmacSha256(client.secret, dataToSign...)
}

func (client *OnePasswordClient) hmacSignWithSession(dataToSign ...[]byte) []byte {
	return HmacSha256(client.sessionHmacK, dataToSign...)
}

func (client *OnePasswordClient) authRegister() (*Response, error) {
	secretB64 := b64.URLEncoding.EncodeToString(client.secret)

	authRegisterPayload := Payload{
		ExtID:  client.extID,
		Method: "auth-sma-hmac256",
		Secret: secretB64,
	}

	authRegisterCommand := client.createCommand("authRegister", authRegisterPayload)

	registerResponse, err := client.SendCommand(authRegisterCommand)

	if err != nil {
		return nil, err
	}

	if registerResponse.Action != "authRegistered" {
		errorMsg := fmt.Sprintf("Unexpected response: %s", registerResponse.Action)
		err = errors.New(errorMsg)
		return nil, err
	}

	return registerResponse, nil
}

func (client *OnePasswordClient) authBegin(cc []byte) (*Response, error) {
	ccB64 := client.base64urlWithoutPadding.EncodeToString(cc)

	authBeginPayload := Payload{
		Method: "auth-sma-hmac256",
		ExtID:  client.extID,
		CC:     ccB64,
	}

	authBeginCommand := client.createCommand("authBegin", authBeginPayload)

	authBeginResponse, err := client.SendCommand(authBeginCommand)

	if err != nil {
		return nil, err
	}

	if authBeginResponse.Action != "authContinue" {
		errorMsg := fmt.Sprintf("Unexpected response: %s", authBeginResponse.Action)
		err = errors.New(errorMsg)
		return nil, err
	}

	return authBeginResponse, nil
}

func (client *OnePasswordClient) generateM4(m3 []byte) []byte {
	return client.hmacSignWithSecret(m3)
}

func (client *OnePasswordClient) generateM3(cs []byte, cc []byte) []byte {
	csAndCc := append(cs[:], cc[:]...)

	csAndCcSha := sha256.New()
	csAndCcSha.Write(csAndCc)

	h := hmac.New(sha256.New, client.secret)
	h.Write(csAndCcSha.Sum(nil))
	return h.Sum(nil)
}

func (client *OnePasswordClient) generateEncK(m3 []byte, m4 []byte) []byte {
	return client.hmacSignWithSecret(m3, m4, []byte("encryption"))
}

func (client *OnePasswordClient) generateHmacK(m3 []byte, m4 []byte) []byte {
	return client.hmacSignWithSecret(m4, m3, []byte("hmac"))
}

func (client *OnePasswordClient) signMessageHmac(iv []byte, data []byte, adata []byte) []byte {
	return client.hmacSignWithSession(iv, data)
}

func (client *OnePasswordClient) Debug(secretB64 string, csB64 string, ccB64 string, m3B64 string, m4B64 string, encKB64 string, hmacKB64 string, ivB64 string, plaintext string, adata string, ciphertextB64 string, hmacB64 string) {
	secret, _ := client.base64urlWithoutPadding.DecodeString(secretB64)

	client.secret = secret

	cs, _ := client.base64urlWithoutPadding.DecodeString(csB64)
	cc, _ := client.base64urlWithoutPadding.DecodeString(ccB64)

	m3, _ := client.base64urlWithoutPadding.DecodeString(m3B64)
	m4, _ := client.base64urlWithoutPadding.DecodeString(m4B64)

	encK, _ := client.base64urlWithoutPadding.DecodeString(encKB64)
	client.sessionEncK = encK

	hmacK, _ := client.base64urlWithoutPadding.DecodeString(hmacKB64)
	client.sessionHmacK = hmacK

	//iv, _ := client.base64urlWithoutPadding.DecodeString(ivB64)

	//ciphertext, _ := client.base64urlWithoutPadding.DecodeString(ciphertextB64)

	hmac, _ := client.base64urlWithoutPadding.DecodeString(hmacB64)

	// Generate M3
	if bytes.Compare(client.generateM3(cs, cc), m3) != 0 {
		log.Printf("Bad M3 logic")
	}

	// Generate M4
	if bytes.Compare(client.generateM4(m3), m4) != 0 {
		log.Printf("Bad M4 logic")
	}

	// Generate EncK
	//if bytes.Compare(client.generateEncK(m3, m4), encK) != 0 {
	//log.Printf("Bad EncK Logic")
	//}

	client.sessionEncK = client.generateEncK(m3, m4)

	//if bytes.Compare(client.generateHmacK(m3, m4), hmacK) != 0 {
	//log.Printf("Bad HmacK Logic")
	//}
	client.sessionHmacK = client.generateHmacK(m3, m4)

	generatedHmac := client.signMessageHmac([]byte(ivB64), []byte(ciphertextB64), []byte(adata))

	if bytes.Compare(generatedHmac, hmac) != 0 {
		log.Printf("Bad Hmac Session Logic")
		log.Printf("%s != %s", client.base64urlWithoutPadding.EncodeToString(generatedHmac), client.base64urlWithoutPadding.EncodeToString(hmac))
	}

	log.Printf("Done")
}

func (client *OnePasswordClient) Register(code string) (*Response, error) {
	fmt.Printf("The 1password helper will request registration of code: %s\n", code)
	fmt.Println("To complete registration. You must accept that code from the helper.")
	_, err := client.authRegister()

	if err != nil {
		fmt.Printf("Registration failed with %s\n", err)
		return nil, err
	}

	return nil, nil
}

func (client *OnePasswordClient) Authenticate(register bool) (*Response, error) {
	helloResponse, err := client.SendHelloCommand()

	if err != nil {
		return nil, err
	}

	if register {
		if helloResponse.Action != "authNew" {
			fmt.Println("sudolikeaboss is already registered.")
			os.Exit(0)
		}

		_, err := client.Register(helloResponse.Payload.Code)

		if err != nil {
			return nil, err
		}
	}

	cc, err := GenerateRandomBytes(16)

	if err != nil {
		return nil, err
	}

	authBeginResponse, err := client.authBegin(cc)

	if err != nil {
		return nil, err
	}

	m3, err := client.base64urlWithoutPadding.DecodeString(authBeginResponse.Payload.M3)

	// Verify M3
	cs, _ := client.base64urlWithoutPadding.DecodeString(authBeginResponse.Payload.CS)

	expectedM3Bytes := client.generateM3(cs, cc)

	if bytes.Compare(expectedM3Bytes, m3) != 0 {
		errorMsg := fmt.Sprintf("M3 not expected value")
		err = errors.New(errorMsg)
		return nil, err
	}

	m4 := client.generateM4(m3)
	m4B64 := client.base64urlWithoutPadding.EncodeToString(m4)

	authVerifyPayload := Payload{
		Method: "auth-sma-hmac256",
		M4:     m4B64,
		ExtID:  client.extID,
	}

	authVerifyCommand := client.createCommand("authVerify", authVerifyPayload)

	authVerifyResponse, err := client.SendCommand(authVerifyCommand)

	if err != nil {
		return nil, err
	}

	if authVerifyResponse.Action != "welcome" {
		errorMsg := fmt.Sprintf("Unexpected response: %s", authVerifyResponse.Action)
		err = errors.New(errorMsg)
		return nil, err
	}

	// Generate the keys
	//
	// encK = HMAC-SHA256(secret, M3|M4|"encryption")
	client.sessionEncK = client.generateEncK(m3, m4)

	// hmacK = HMAC-SHA256(secret, M4|M3|"hmac")
	client.sessionHmacK = client.generateHmacK(m3, m4)

	log.Printf("hmacK = %s", client.base64urlWithoutPadding.EncodeToString(client.sessionHmacK))

	decryptedPayload, err := client.decryptResponse(authVerifyResponse)
	if err != nil {
		return nil, err
	}
	log.Printf("%s", b64.StdEncoding.EncodeToString(decryptedPayload))

	return authVerifyResponse, nil
}

func (client *OnePasswordClient) decryptResponse(response *Response) ([]byte, error) {
	iv, err := client.base64urlWithoutPadding.DecodeString(response.Payload.Iv)

	if err != nil {
		return nil, err
	}

	data, err := client.base64urlWithoutPadding.DecodeString(response.Payload.Data)

	if err != nil {
		return nil, err
	}

	hmac, err := client.base64urlWithoutPadding.DecodeString(response.Payload.Hmac)

	if err != nil {
		return nil, err
	}

	// Verify hmac
	expectedHmac := client.hmacSignWithSession([]byte(response.Payload.Iv), []byte(response.Payload.Data))

	log.Printf(
		"%s == %s",
		client.base64urlWithoutPadding.EncodeToString(hmac),
		client.base64urlWithoutPadding.EncodeToString(expectedHmac),
	)

	if bytes.Compare(expectedHmac, hmac) != 0 {
		errorMsg := fmt.Sprintf("Hmac unexpected")
		err = errors.New(errorMsg)
		return nil, err
	}

	// Decrypt
	payload, err := Decrypt(client.sessionEncK, iv, data)

	return payload, err
}

func (client *OnePasswordClient) encryptPayload(payload *Payload) (*Payload, error) {
	iv, err := GenerateRandomBytes(16)

	if err != nil {
		return nil, err
	}

	// Encrypt the payload
	payloadJsonStr, err := json.Marshal(payload)

	if err != nil {
		return nil, err
	}

	// Encrypt the payload
	encryptedPayload, err := Encrypt(client.sessionEncK, iv, payloadJsonStr)

	encryptedPayloadB64 := client.base64urlWithoutPadding.EncodeToString(encryptedPayload)

	// Generate HMAC for the message
	ivB64 := client.base64urlWithoutPadding.EncodeToString(iv)

	payloadHmac := client.hmacSignWithSession([]byte(ivB64), []byte(encryptedPayloadB64))

	payloadHmacB64 := client.base64urlWithoutPadding.EncodeToString(payloadHmac)

	newPayload := Payload{
		Iv:        ivB64,
		Data:      encryptedPayloadB64,
		Algorithm: "aead-cbchmac-256",
		Hmac:      payloadHmacB64,
	}

	return &newPayload, nil
}

func (client *OnePasswordClient) SendCommand(command *Command) (*Response, error) {
	jsonStr, err := json.Marshal(command)

	if err != nil {
		return nil, err
	}

	err = client.SendJSON(jsonStr)

	if err != nil {
		return nil, err
	}

	response, err := client.ReceiveJSON()

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (client *OnePasswordClient) SendEncryptedCommand(command *Command) (*Response, error) {
	// Create the encrypted payload
	plaintextPayload := command.Payload

	encryptedPayload, err := client.encryptPayload(&plaintextPayload)

	if err != nil {
		return nil, err
	}

	command.Payload = *encryptedPayload

	jsonStr, err := json.Marshal(command)

	if err != nil {
		return nil, err
	}

	err = client.SendJSON(jsonStr)

	if err != nil {
		return nil, err
	}

	response, err := client.ReceiveJSON()

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (client *OnePasswordClient) SendJSON(jsonStr []byte) error {
	log.Printf("Sending: %s", jsonStr)

	err := client.websocketClient.Send(jsonStr)

	if err != nil {
		return err
	}
	return nil
}

func (client *OnePasswordClient) ReceiveJSON() (*Response, error) {
	var rawResponseStr string

	err := client.websocketClient.Receive(&rawResponseStr)

	if err != nil {
		return nil, err
	}

	log.Printf("Received: %s", rawResponseStr)

	response, err := LoadResponse(rawResponseStr)

	if err != nil {
		return nil, err
	}

	return response, nil
}
