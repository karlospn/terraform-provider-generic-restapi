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
)

func resourceScaffolding() *schema.Resource {
	return &schema.Resource{
		Create: resourceScaffoldingCreate,
		Read:   resourceScaffoldingRead,
		Update: resourceScaffoldingUpdate,
		Delete: resourceScaffoldingDelete,

		Schema: map[string]*schema.Schema{
			"payload": {
				Type:     schema.TypeString,
				Required: true,
			},
			"create_method": &schema.Schema{
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("REST_API_OVERRIDE_CREATE_METHOD", nil),
				Description: "The HTTP route used to CREATE objects of this type on the API server. Overrides  the provider configuration",
				Optional:    true,
			},
			"read_method": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The HTTP route used to READ objects of this type on the API server. Overrides provider configuration",
				DefaultFunc: schema.EnvDefaultFunc("REST_API_OVERRIDE_READ_METHOD", nil),
				Optional:    true,
			},
			"update_method": &schema.Schema{
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("REST_API_OVERRIDE_UPDATE_METHOD", nil),
				Description: "The HTTP route used to UPDATE objects of this type on the API server. Overrides provider configuration",
				Optional:    true,
			},
			"destroy_method": &schema.Schema{
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("REST_API_OVERRIDE_DESTROY_METHOD", nil),
				Description: "The HTTP route used to DELETE objects of this type on the API server. Overrides provider configuration",
				Optional:    true,
			},
		},
	}
}

func resourceScaffoldingCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*api_client)

	payload := d.Get("payload").(string)
	b, _ := json.Marshal(json.RawMessage(payload))

	if payload == "" {
		return fmt.Errorf("Payload cannot be empty")
	}

	route := d.Get("create_method").(string)
	if route == "" {
		route = client.create_method
	}

	result, _, err := send(client, "POST", route, string(b))

	if err != nil {
		return fmt.Errorf("Failed to create record: %s", err)
	}

	var task Task
	err = json.Unmarshal(result, &task)

	if err != nil {
		return fmt.Errorf("Failed to create record: %s", err)
	}

	if task.Id == "" {
		return fmt.Errorf("Something went wrong. The api did not return an Id")
	}

	d.SetId(task.Id)
	return resourceScaffoldingRead(d, meta)
}

func resourceScaffoldingRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*api_client)

	id := d.Id()

	route := d.Get("read_method").(string)
	if route == "" {
		route = client.read_method
	}

	result, statusCode, err := send(client, "GET", strings.Replace(route, "{id}", id, -1), "")

	if err != nil {

		if statusCode == 404 {
			log.Printf("resource not found with ID:\n%s\n", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf(
			"There was a problem when trying to find object with ID: %s", d.Id())
	}

	if err != nil {
		return fmt.Errorf("Failed to read record: %s", err)
	}

	norm, err := structure.NormalizeJsonString(string(result))

	if err != nil {
		return fmt.Errorf("Error trying to normalize result: %s", err)
	}

	d.Set("payload", norm)

	return nil
}

func resourceScaffoldingUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*api_client)

	id := d.Id()
	payload := d.Get("payload").(string)
	b, _ := json.Marshal(json.RawMessage(payload))

	if payload == "" {
		return fmt.Errorf("Payload cannot be empty")
	}

	route := d.Get("update_method").(string)
	if route == "" {
		route = client.update_method
	}

	_, _, err := send(client, "PUT", strings.Replace(route, "{id}", id, -1), string(b))

	if err != nil {
		return fmt.Errorf("Failed to update record: %s", err)
	}

	return resourceScaffoldingRead(d, meta)
}

func resourceScaffoldingDelete(d *schema.ResourceData, meta interface{}) error {

	id := d.Id()
	client := meta.(*api_client)

	route := d.Get("destroy_method").(string)
	if route == "" {
		route = client.destroy_method
	}

	_, _, err := send(client, "DELETE", strings.Replace(route, "{id}", id, -1), "")

	if err != nil {
		return err
	}

	return nil
}

func send(client *api_client, method string, path string, data string) ([]byte, int, error) {

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
		return []byte{}, 500, err
	}

	if client.username != "" && client.password != "" {
		req.SetBasicAuth(client.username, client.password)
	}

	resp, err := client.http_client.Do(req)

	if err != nil {
		return []byte{}, resp.StatusCode, err
	}

	body, err2 := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err2 != nil {
		return []byte{}, resp.StatusCode, err2
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return body, resp.StatusCode, errors.New(fmt.Sprintf("Unexpected response code '%d': %s", resp.StatusCode, body))
	}

	return body, resp.StatusCode, nil

}
