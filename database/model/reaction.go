package model

import "encoding/json"

func (r *Reactions) GetClients() ([]string, error) {
	if len(r.Clients) == 0 {
		return []string{}, nil
	}
	var ret []string
	err := json.Unmarshal([]byte(r.Clients), &ret)
	return ret, err
}

func (r *Reactions) AppendClient(client string) error {
	clients, err := r.GetClients()
	if err != nil {
		return err
	}
	clients = append(clients, client)
	clientsJson, err := json.Marshal(clients)
	if err != nil {
		return err
	}
	r.Clients = string(clientsJson)
	return nil
}

func (r *Reactions) RemoveClient(client string) error {
	clients, err := r.GetClients()
	if err != nil {
		return err
	}
	clients = remove(clients, client)
	clientsJson, err := json.Marshal(clients)
	if err != nil {
		return err
	}
	r.Clients = string(clientsJson)
	return nil
}

func remove(slice []string, element string) []string {
	for i, v := range slice {
		if v == element {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
