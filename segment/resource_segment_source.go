package segment

import (
	"strings"

	"github.com/ajbosco/segment-config-go/segment"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSegmentSource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"catalog_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
		Create: resourceSegmentSourceCreate,
		Read:   resourceSegmentSourceRead,
		Delete: resourceSegmentSourceDelete,
	}
}

func resourceSegmentSourceCreate(r *schema.ResourceData, meta interface{}) error {
	client := meta.(*segment.Client)
	name := r.Get("name").(string)
	catalogName := r.Get("catalog_name").(string)

	newSource := segment.Source{
		Name:        name,
		CatalogName: catalogName,
	}

	source, err := client.CreateSource(newSource)
	if err != nil {
		return err
	}

	// the id of the source is the last value in the name path
	splitName := strings.Split(source.Name, "/")
	id := splitName[len(splitName)-1]
	r.SetId(id)

	return resourceSegmentSourceRead(r, meta)
}

func resourceSegmentSourceRead(r *schema.ResourceData, meta interface{}) error {
	client := meta.(*segment.Client)
	name := r.Id()

	s, err := client.GetSource(name)
	if err != nil {
		return err
	}

	r.Set("name", s.Name)
	r.Set("catalog_name", s.CatalogName)

	return nil
}

func resourceSegmentSourceDelete(r *schema.ResourceData, meta interface{}) error {
	client := meta.(*segment.Client)
	name := r.Id()

	err := client.DeleteSource(name)
	if err != nil {
		return err
	}

	return nil
}
