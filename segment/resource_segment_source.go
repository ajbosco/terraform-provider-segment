package segment

import (
	"fmt"
	"strings"

	"github.com/fenderdigital/segment-apis-go/segment"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSegmentSource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"source_name": {
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
		Importer: &schema.ResourceImporter{
			State: resourceSegmentSourceImport,
		},
	}
}

func resourceSegmentSourceCreate(r *schema.ResourceData, meta interface{}) error {
	client := meta.(*segment.Client)
	srcName := r.Get("source_name").(string)
	catName := r.Get("catalog_name").(string)

	source, err := client.CreateSource(srcName, catName)
	if err != nil {
		return err
	}

	r.SetId(source.Name)

	return resourceSegmentSourceRead(r, meta)
}

func resourceSegmentSourceRead(r *schema.ResourceData, meta interface{}) error {
	client := meta.(*segment.Client)
	id := r.Id()

	srcName := idToName(id)

	s, err := client.GetSource(srcName)
	if err != nil {
		return err
	}

	r.Set("catalog_name", s.CatalogName)

	return nil
}

func resourceSegmentSourceDelete(r *schema.ResourceData, meta interface{}) error {
	client := meta.(*segment.Client)
	id := r.Id()

	srcName := idToName(id)

	err := client.DeleteSource(srcName)
	if err != nil {
		return err
	}

	return nil
}

func resourceSegmentSourceImport(r *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*segment.Client)
	s, err := client.GetSource(r.Id())
	if err != nil {
		return nil, fmt.Errorf("invalid source: %q; err: %v", r.Id(), err)
	}

	r.SetId(s.Name)
	r.Set("catalog_name", s.CatalogName)

	results := make([]*schema.ResourceData, 1)
	results[0] = r

	return results, nil
}

func idToName(id string) string {
	splitID := strings.Split(id, "/")

	return splitID[len(splitID)-1]
}
