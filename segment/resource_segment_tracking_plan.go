package segment

import (
	"encoding/json"
	"fmt"

	"github.com/fenderdigital/segment-apis-go/segment"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSegmentTrackingPlan() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"rules": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
		Create: resourceSegmentTrackingPlanCreate,
		Read:   resourceSegmentTrackingPlanRead,
		Delete: resourceSegmentTrackingPlanDelete,
	}
}

func resourceSegmentTrackingPlanCreate(r *schema.ResourceData, meta interface{}) error {
	client := meta.(*segment.Client)
	displayName := r.Get("display_name").(string)
	rules := r.Get("rules").(string)
	s := segment.Rules{}
	json.Unmarshal([]byte(rules), &s)
	fmt.Printf("%+v\n", s)
	trackingPlan, err := client.CreateTrackingPlan(displayName, s)
	if err != nil {
		return err
	}

	r.SetId(trackingPlan.Name)
	r.Set("name", trackingPlan.Name)

	return resourceSegmentTrackingPlanRead(r, meta)
}

func resourceSegmentTrackingPlanRead(r *schema.ResourceData, meta interface{}) error {
	client := meta.(*segment.Client)
	planName := r.Get("name").(string)
	trackingPlan, err := client.GetTrackingPlan(planName)
	if err != nil {
		return err
	}

	s, _ := json.Marshal(trackingPlan.Rules)
	r.Set("name", trackingPlan.Name)
	r.Set("rules", string(s))

	return nil
}

func resourceSegmentTrackingPlanDelete(r *schema.ResourceData, meta interface{}) error {
	client := meta.(*segment.Client)
	planName := r.Get("name").(string)

	err := client.DeleteTrackingPlan(planName)
	if err != nil {
		return err
	}

	return nil
}

func resourceSegmentTrackingPlanUpdate(r *schema.ResourceData, meta interface{}) error {
	client := meta.(*segment.Client)
	planName := r.Get("name").(string)
	rules := r.Get("rules").(string)
	displayName := r.Get("display_name").(string)

	paths := []string{"tracking_plan.display_name", "tracking_plan.rules"}

	s := segment.Rules{}
	json.Unmarshal([]byte(rules), &s)
	updatedPlan := segment.TrackingPlan{
		DisplayName: displayName,
		Rules:       s,
	}
	_, err := client.UpdateTrackingPlan(planName, paths, updatedPlan)
	if err != nil {
		return err
	}
	return resourceSegmentTrackingPlanRead(r, meta)
}
