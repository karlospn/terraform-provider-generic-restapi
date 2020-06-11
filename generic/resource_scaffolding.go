package generic

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"github.com/tidwall/gjson"
)

func resourceScaffolding() *schema.Resource {
	return &schema.Resource{
		Create: resourceScaffoldingCreate,
		Read:   resourceScaffoldingRead,
		Update: resourceScaffoldingUpdate,
		Delete: resourceScaffoldingDelete,

		Schema: map[string]*schema.Schema{
			"id_attribute": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Allows per-resource override of `id_attribute` ",
				Required:    true,
			},
			"data": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceScaffoldingCreate(d *schema.ResourceData, meta interface{}) error {

	id := d.Get("id_attribute").(string)
	data := d.Get("data").(string)

	client := meta.(*api_client)

	requestID := gjson.Get(data, id)

	if requestID.String() == "" {
		return fmt.Errorf(" id not found on response")
	}

	b, _ := json.Marshal(json.RawMessage(data))
	_, err := send(client, "POST", strings.Replace(client.create_method, "{id}", requestID.String(), -1), string(b))

	if err != nil {
		return fmt.Errorf("Failed to create record: %s", err)
	}

	log.Printf("resource create called. Object build:\n%s\n", data)

	d.SetId(requestID.String())
	return resourceScaffoldingRead(d, meta)
}

func resourceScaffoldingRead(d *schema.ResourceData, meta interface{}) error {

	id := d.Id()
	client := meta.(*api_client)
	result, err := send(client, "GET", strings.Replace(client.read_method, "{id}", id, -1), "")

	if err != nil {

		if strings.Contains(err.Error(), "not found") {
			log.Printf("resource read called. No id found:\n%s\n", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf(
			"There was a problem when trying to find object with ID: %s",
			d.Id())
	}

	log.Printf("resource read called. Object built:\n%s\n", result)

	norm, _ := structure.NormalizeJsonString(result)

	d.Set("data", norm)

	return nil
}

func resourceScaffoldingUpdate(d *schema.ResourceData, meta interface{}) error {

	id := d.Get("id_attribute").(string)
	data := d.Get("data").(string)
	client := meta.(*api_client)
	requestID := gjson.Get(data, id)

	if requestID.String() == "" {
		return fmt.Errorf(" id not found on response")
	}

	log.Printf("resource updated called. Object built:\n%s\n", data)

	b, _ := json.Marshal(json.RawMessage(data))

	_, err := send(client, "PUT", strings.Replace(client.update_method, "{id}", requestID.String(), -1), string(b))

	if err != nil {
		return fmt.Errorf("Failed to update record: %s", err)
	}

	return resourceScaffoldingRead(d, meta)
}

func resourceScaffoldingDelete(d *schema.ResourceData, meta interface{}) error {

	id := d.Id()
	client := meta.(*api_client)
	_, err := send(client, "DELETE", strings.Replace(client.destroy_method, "{id}", id, -1), "")

	if err != nil {
		return err
	}

	return nil
}

func send(client *api_client, method string, path string, data string) (string, error) {

	fulluri := client.uri + path
	var req *http.Request
	var err error

	buffer := bytes.NewBuffer([]byte(data))

	if data == "" {
		req, err = http.NewRequest(method, fulluri, nil)
	} else {
		req, err = http.NewRequest(method, fulluri, buffer)

		if err == nil {
			req.Header.Set("Content-Type", "application/json")
		}
	}

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	if client.username != "" && client.password != "" {
		req.SetBasicAuth(client.username, client.password)
	}

	resp, err := client.http_client.Do(req)

	if err != nil {
		return "", err
	}

	bodyBytes, err2 := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err2 != nil {
		return "", err2
	}

	body := string(bodyBytes)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return body, errors.New(fmt.Sprintf("Unexpected response code '%d': %s", resp.StatusCode, body))
	}

	return body, nil

}
