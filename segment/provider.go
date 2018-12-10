package segment

import (
	"github.com/ajbosco/segment-config-go/segment"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_token": {
				Type:        schema.TypeString,
				Description: "The Access Token used to connect to Segment",
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SEGMENT_ACCESS_TOKEN", nil),
			},
			"workspace": {
				Type:        schema.TypeString,
				Description: "The Segment workspace to manage",
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SEGMENT_WORKSPACE", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"segment_source":      resourceSegmentSource(),
			"segment_destination": resourceSegmentDestination(),
		},
		ConfigureFunc: configureFunc(),
	}
}

func configureFunc() func(*schema.ResourceData) (interface{}, error) {
	return func(d *schema.ResourceData) (interface{}, error) {
		client := segment.NewClient(d.Get("access_token").(string), d.Get("workspace").(string))
		return client, nil
	}
}
