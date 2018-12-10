package segment

import (
	"github.com/ajbosco/segment-config-go/segment"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSegmentDestination() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"source_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"destination_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"connection_mode": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"configs": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Required: true,
			},
		},
		Create: resourceSegmentDestinationCreate,
		Read:   resourceSegmentDestinationRead,
		Update: resourceSegmentDestinationUpdate,
		Delete: resourceSegmentDestinationDelete,
	}
}

func resourceSegmentDestinationCreate(r *schema.ResourceData, meta interface{}) error {
	client := meta.(*segment.Client)
	srcName := r.Get("source_name").(string)
	destName := r.Get("destination_name").(string)
	connMode := r.Get("connection_mode").(string)
	enabled := r.Get("enabled").(bool)
	configs := r.Get("configs").(*schema.Set)

	dest, err := client.CreateDestination(srcName, destName, connMode, enabled, extractConfigs(configs))
	if err != nil {
		return err
	}

	r.SetId(dest.Name)

	return resourceSegmentDestinationRead(r, meta)
}

func resourceSegmentDestinationRead(r *schema.ResourceData, meta interface{}) error {
	client := meta.(*segment.Client)
	srcName := r.Get("source_name").(string)
	id := r.Id()
	destName := idToName(id)

	d, err := client.GetDestination(srcName, destName)
	if err != nil {
		return err
	}

	r.Set("enabled", d.Enabled)
	r.Set("configs", d.Configs)
	r.Set("connection_mode", d.ConnectionMode)

	return nil
}

func resourceSegmentDestinationUpdate(r *schema.ResourceData, meta interface{}) error {
	client := meta.(*segment.Client)
	srcName := r.Get("source_name").(string)
	configs := r.Get("configs").(*schema.Set)
	enabled := r.Get("enabled").(bool)
	id := r.Id()
	destName := idToName(id)

	_, err := client.UpdateDestination(srcName, destName, enabled, extractConfigs(configs))
	if err != nil {
		return err
	}

	return resourceSegmentDestinationRead(r, meta)
}

func resourceSegmentDestinationDelete(r *schema.ResourceData, meta interface{}) error {
	client := meta.(*segment.Client)
	srcName := r.Get("source_name").(string)
	id := r.Id()
	destName := idToName(id)

	err := client.DeleteDestination(srcName, destName)
	if err != nil {
		return err
	}

	return nil
}

func extractConfigs(s *schema.Set) []segment.DestinationConfig {
	configs := []segment.DestinationConfig{}

	if s != nil {
		for _, config := range s.List() {
			c := segment.DestinationConfig{
				Name:  config.(map[string]interface{})["name"].(string),
				Type:  config.(map[string]interface{})["type"].(string),
				Value: config.(map[string]interface{})["value"],
			}
			configs = append(configs, c)
		}
	}

	return configs
}
