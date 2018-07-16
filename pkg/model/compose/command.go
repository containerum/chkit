package compose

import "gopkg.in/yaml.v2"

var (
	_ yaml.Unmarshaler = new(Command)
)

type Command []string

func (command *Command) UnmarshalYAML(unmarshal Unmarshaler) error {
	var value interface{}
	if err := unmarshal(&value); err != nil {
		return err
	}
	switch value := value.(type) {
	case []string:
		*command = Command(value)
	case string:
		*command = Command{value}
	}
	return nil
}
