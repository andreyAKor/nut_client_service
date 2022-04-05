package server

import (
	"net/http"

	"github.com/pkg/errors"
	nut "github.com/robbiet480/go.nut"
)

// UPS contains information about a specific UPS provided by the NUT instance.
type UPS struct {
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	Master         bool       `json:"master"`
	NumberOfLogins int        `json:"numberOfLogins"`
	Clients        []string   `json:"clients"`
	Variables      []Variable `json:"variables"`
	Commands       []Command  `json:"commands"`
}

// Variable describes a single variable related to a UPS.
type Variable struct {
	Name          string      `json:"name"`
	Value         interface{} `json:"value"`
	Type          string      `json:"type"`
	Description   string      `json:"description"`
	Writeable     bool        `json:"writeable"`
	MaximumLength int         `json:"maximumLength"`
	OriginalType  string      `json:"originalType"`
}

// Command describes an available command for a UPS.
type Command struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// get
func (s *Server) get(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	list, err := s.nutClient.GetUPSList()
	if err != nil {
		return nil, errors.Wrap(err, "get UPS list fail")
	}

	return convertListToList(list), nil
}

func convertListToList(l []nut.UPS) []UPS {
	var res []UPS
	for _, v := range l {
		res = append(res, UPS{
			Name:           v.Name,
			Description:    v.Description,
			Master:         v.Master,
			NumberOfLogins: v.NumberOfLogins,
			Clients:        v.Clients,
			Variables:      convertVariableToVariable(v.Variables),
			Commands:       convertCommandsToCommands(v.Commands),
		})
	}
	return res
}

func convertVariableToVariable(l []nut.Variable) []Variable {
	var res []Variable
	for _, v := range l {
		res = append(res, Variable{
			Name:          v.Name,
			Value:         v.Value,
			Type:          v.Type,
			Description:   v.Description,
			Writeable:     v.Writeable,
			MaximumLength: v.MaximumLength,
			OriginalType:  v.OriginalType,
		})
	}
	return res
}

func convertCommandsToCommands(l []nut.Command) []Command {
	var res []Command
	for _, v := range l {
		res = append(res, Command{
			Name:        v.Name,
			Description: v.Description,
		})
	}
	return res
}
