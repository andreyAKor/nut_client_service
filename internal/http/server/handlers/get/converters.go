package get

import "github.com/andreyAKor/nut_client"

func convertListToList(l []*nut_client.UPS) []UPS {
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

func convertVariableToVariable(l []nut_client.Variable) []Variable {
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

func convertCommandsToCommands(l []nut_client.Command) []Command {
	var res []Command
	for _, v := range l {
		res = append(res, Command{
			Name:        v.Name,
			Description: v.Description,
		})
	}
	return res
}
