package uptimekumaapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Tag struct {
	Name  string `json:"name"`
	Color string `json:"color"`
	ID    int    `json:"id,omitempty"`
}

type TagsResponse struct {
	Tags []Tag `json:"tags"`
}

type TagDeleteRequest struct {
	TagID int `json:"tag_id"`
}

func (api *UptimeKumaAPI) GetTags() (*[]Tag, error) {
	endpoint := fmt.Sprintf("%s/tags", api.Host)

	req, _ := http.NewRequest("GET", endpoint, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", api.Token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error while getting tags")
		return nil, err
	}

	var body TagsResponse
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		log.Println("Error while decoding tags")
		return nil, err
	}

	return &body.Tags, nil
}

func (api *UptimeKumaAPI) GetTag(name string) (*Tag, error) {
	tags, err := api.GetTags()
	if err != nil {
		return nil, err
	}

	for _, tag := range *tags {
		if tag.Name == name {
			return &tag, nil
		}
	}

	return nil, fmt.Errorf("Tag %s not found", name)
}

func (api *UptimeKumaAPI) CreateTag(name string, color string) (*Tag, error) {
	foundTag, _ := api.GetTag(name)
	if foundTag != nil {
		log.Printf("Tag %s already exists\n", name)
		return foundTag, nil
	}

	endpoint := fmt.Sprintf("%s/tags", api.Host)

	tag := Tag{
		Name:  name,
		Color: color,
	}

	reqMarshalled, _ := json.Marshal(tag)

	resp, err := api.postRequest(endpoint, reqMarshalled)
	if err != nil {
		log.Println("Error while creating tag")
		return nil, err
	}

	if resp.StatusCode != 200 {
		log.Println("Error while creating tag")
		log.Fatalln(resp.StatusCode)
	}

	return &tag, nil
}

func (api *UptimeKumaAPI) DeleteTag(name string) error {
	tags, _ := api.GetTags()
	for _, tag := range *tags {
		if tag.Name == name {
			log.Printf("Deleting tag %d %s\n", tag.ID, tag.Name)
			api.DeleteTagByID(tag.ID)
		}
	}
	return nil
}

func (api *UptimeKumaAPI) DeleteTagByID(id int) error {
	endpoint := fmt.Sprintf("%s/tags/%d", api.Host, id)

	body := TagDeleteRequest{
		TagID: id,
	}
	reqMarshalled, _ := json.Marshal(body)

	resp, _ := api.deleteRequest(endpoint, reqMarshalled)
	if resp.StatusCode != 200 {
		log.Println("Error while deleting tag")
		log.Fatalln(resp.StatusCode)
	}

	return nil
}
