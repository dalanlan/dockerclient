package docker

import (
	"encoding/json"
	//"fmt"
)

type Image struct {
	RepoTags    []string    //`json:"RepoTags,omitempty"`
	Id          string      //`json:"Id"`
	Created     int64       //`json:"Created,omitempty"`
	Size        int64       //`json:"Size,omitempty"`
	VirtualSize int64       //`json:"VirtualSize,omitempty"`
	ParentId    string      //`json:"ParentId,omitempty"`
	RepoDigests []string    //`json:"RepoDigests,omitempty"`
	Labels      interface{} //`json:"Labels,omitempty"`
}

func (c *Client) ListImage() ([]Image, error) {
	method := "GET"
	url := c.endpoint + "/images/json"
	body, err := c.do(method, url, DoOption{})
	if err != nil {
		return nil, err
	}
	//fmt.Println(string(body))
	var img []Image
	err = json.Unmarshal(body, &img)
	if err != nil {
		return nil, err
	}
	return img, nil

}
