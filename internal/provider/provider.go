package provider

import (
	"crypto/tls"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Provider returns the terraform.ResourceProvider structure for the generic
// provider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"uri": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "URI of the REST API endpoint. This serves as the base of all requests.",
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "When set, will use this username for BASIC auth to the API.",
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "When set, will use this password for BASIC auth to the API.",
			},
			"timeout": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "When set, will cause requests taking longer than this time (in seconds) to be aborted.",
			},
			"create_method": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Defaults to `POST`. The HTTP method used to CREATE objects of this type on the API server.",
				Optional:    true,
			},
			"read_method": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Defaults to `GET`. The HTTP method used to READ objects of this type on the API server.",
				Optional:    true,
			},
			"update_method": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Defaults to `PUT`. The HTTP method used to UPDATE objects of this type on the API server.",
				Optional:    true,
			},
			"destroy_method": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Defaults to `DELETE`. The HTTP method used to DELETE objects of this type on the API server.",
				Optional:    true,
			},
			"insecure": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "When using https, this disables TLS verification of the host.",
			},
			"use_cookies": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable cookie jar to persist session.",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"scaffolding_data_source": dataSourceScaffolding(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"scaffolding_resource": resourceScaffolding(),
		},
		ConfigureFunc: configureProvider,
	}
}

type api_client struct {
	http_client    *http.Client
	insecure       bool
	uri            string
	username       string
	password       string
	timeout        int
	create_method  string
	read_method    string
	update_method  string
	destroy_method string
	use_cookies    string
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {

	timeout := d.Get("timeout").(int)
	insecure := d.Get("insecure").(bool)
	use_cookies := d.Get("use_cookies").(bool)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
		Proxy:           http.ProxyFromEnvironment,
	}

	var cookieJar http.CookieJar

	if use_cookies {
		cookieJar, _ = cookiejar.New(nil)
	}

	return &api_client{
		http_client: &http.Client{
			Timeout:   time.Second * time.Duration(timeout),
			Transport: tr,
			Jar:       cookieJar,
		},
		uri:            d.Get("uri").(string),
		username:       d.Get("username").(string),
		password:       d.Get("username").(string),
		timeout:        d.Get("timeout").(int),
		create_method:  d.Get("create_method").(string),
		read_method:    d.Get("read_method").(string),
		update_method:  d.Get("update_method").(string),
		destroy_method: d.Get("destroy_method").(string),
	}, nil
}
