package segment

import (
	"strings"

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

	newDest := segment.Destination{
		Name:           destName,
		Enabled:        enabled,
		ConnectionMode: connMode,
		Configs:        extractConfigs(configs),
	}

	dest, err := client.CreateDestination(srcName, newDest)
	if err != nil {
		return err
	}

	// the id of the destination is the last value in the name path
	splitName := strings.Split(dest.Name, "/")
	id := splitName[len(splitName)-1]
	r.SetId(id)

	return resourceSegmentDestinationRead(r, meta)
}

func resourceSegmentDestinationRead(r *schema.ResourceData, meta interface{}) error {
	client := meta.(*segment.Client)
	destName := r.Id()
	srcName := r.Get("source_name").(string)

	d, err := client.GetDestination(srcName, destName)
	if err != nil {
		return err
	}

	r.Set("name", d.Name)
	r.Set("connection_mod", d.ConnectionMode)
	r.Set("enabled", d.Enabled)
	r.Set("configs", d.Configs)

	return nil
}

func resourceSegmentDestinationUpdate(r *schema.ResourceData, meta interface{}) error {
	client := meta.(*segment.Client)
	destID := r.Id()
	destName := r.Get("destination_name").(string)
	srcName := r.Get("source_name").(string)
	configs := r.Get("configs").(*schema.Set)
	enabled := r.Get("enabled").(bool)

	dest := segment.Destination{
		Name:    destName,
		Enabled: enabled,
		Configs: extractConfigs(configs),
	}
	// updateMask determines which fields Segment will update
	updateMask := segment.UpdateMask{Paths: []string{"destination.config", "destination.enabled"}}

	_, err := client.UpdateDestination(srcName, destID, dest, updateMask)
	if err != nil {
		return err
	}

	return resourceSegmentDestinationRead(r, meta)
}

func resourceSegmentDestinationDelete(r *schema.ResourceData, meta interface{}) error {
	client := meta.(*segment.Client)
	destName := r.Id()
	srcName := r.Get("source_name").(string)

	err := client.DeleteDestinaton(srcName, destName)
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
