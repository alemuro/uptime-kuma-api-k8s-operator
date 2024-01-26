package uptimekumaapi

import (
	"encoding/json"
	"fmt"
	"log"
)

func (api *UptimeKumaAPI) areDifferentMonitors(monitor1 *Monitor, monitor2 *Monitor) bool {
	if monitor1.Name != monitor2.Name {
		return true
	}
	if monitor1.URL != monitor2.URL {
		return true
	}
	if monitor1.Interval != monitor2.Interval {
		return true
	}
	// Compare tag arrays
	if len(monitor1.Tags) != len(monitor2.Tags) {
		return true
	}

	return false
}

func (api *UptimeKumaAPI) CreateMonitor(name string, url string, interval int, tags []string) (*Monitor, error) {
	endpoint := fmt.Sprintf("%s/monitors", api.Host)
	req := Monitor{
		Type:     "http",
		Name:     name,
		URL:      url,
		Interval: interval,
	}

	// Review if the monitor already exists
	foundMonitor, err := api.GetMonitor(name)
	if err == nil {
		if api.areDifferentMonitors(foundMonitor, &req) {
			api.PatchMonitor(foundMonitor.ID, name, url, interval, tags)
		}

		log.Printf("Monitor %s already exists\n", name)
		return foundMonitor, nil
	}

	body, _ := json.Marshal(req)
	resp, err := api.postRequest(endpoint, body)
	if err != nil || resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error while creating monitor")
	}

	return &req, nil
}

func (api *UptimeKumaAPI) PatchMonitor(ID int64, name string, url string, interval int, tags []string) (*Monitor, error) {
	log.Printf("Updating monitor %s\n", name)
	endpoint := fmt.Sprintf("%s/monitors/%d", api.Host, ID)

	req := Monitor{
		Type:     "http",
		Name:     name,
		URL:      url,
		Interval: interval,
	}
	body, _ := json.Marshal(req)
	resp, err := api.patchRequest(endpoint, body)
	if err != nil || resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error while updating monitor")
	}
	return &req, nil
}

func (api *UptimeKumaAPI) GetMonitors() (*[]Monitor, error) {
	endpoint := fmt.Sprintf("%s/monitors", api.Host)

	resp, err := api.getRequest(endpoint)
	if err != nil || resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error while creating monitor")
	}

	var body MonitorsResponse
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return nil, fmt.Errorf("Error while decoding monitors: %w", err)
	}

	return &body.Monitors, nil
}

func (api *UptimeKumaAPI) GetMonitor(name string) (*Monitor, error) {
	monitors, err := api.GetMonitors()
	if err != nil {
		return nil, fmt.Errorf("Error while getting monitor: %w", err)
	}

	for _, monitor := range *monitors {
		if monitor.Name == name {
			return &monitor, nil
		}
	}

	return nil, fmt.Errorf("Monitor %s not found", name)
}

func (api *UptimeKumaAPI) DeleteMonitor(name string) error {
	monitors, err := api.GetMonitors()
	if err != nil {
		return fmt.Errorf("Error while deleting monitor: %w", err)
	}

	for _, monitor := range *monitors {
		if monitor.Name == name {
			log.Printf("Deleting monitor %s\n", name)
			err := api.DeleteMonitorByID(monitor.ID)
			if err != nil {
				return fmt.Errorf("Error while deleting monitor: %w", err)
			}
		}
	}

	return nil
}

func (api *UptimeKumaAPI) DeleteMonitorByID(id int64) error {
	endpoint := fmt.Sprintf("%s/monitors/%d", api.Host, id)

	resp, err := api.deleteRequest(endpoint, nil)
	if err != nil || resp.StatusCode != 200 {
		return fmt.Errorf("Error while deleting monitor: %w", err)
	}

	return nil
}

type MonitorsResponse struct {
	Monitors []Monitor `json:"monitors"`
}

type Monitor struct {
	ID       int64    `json:"id,omitempty"`
	Tags     []string `json:"tags,omitempty"`
	Type     string   `json:"type"`
	Name     string   `json:"name"`
	URL      string   `json:"url"`
	Interval int      `json:"interval"`
}
