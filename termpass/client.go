package termpass

type Payload interface{}

type Command struct {
	Action   string  `json:"action"`
	Number   int     `json:"number"`
	Version  string  `json:"version"`
	BundleId string  `json:"bundleId"`
	Payload  Payload `json:"payload"`
}

type HelloPayload struct {
	Version      string   `json:"version"`
	Capabilities []string `json:"capabilities"`
}

type ShowPopupPayload struct {
	Url     string `json:"url"`
	Options string `json:"options"`
}

type CommandFactory struct {
}

func (commandFactory *CommandFactory) New(commandType string) Command {
}

func SendCommand(commandType string) {
	// Create the command with a command factory
	// JSON encode the command object
}
