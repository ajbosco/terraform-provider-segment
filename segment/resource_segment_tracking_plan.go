package segment

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fenderdigital/segment-apis-go/segment"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSegmentTrackingPlan() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"rules": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
		},
		Create: resourceSegmentTrackingPlanCreate,
		Read:   resourceSegmentTrackingPlanRead,
		Delete: resourceSegmentTrackingPlanDelete,
		Update: resourceSegmentTrackingPlanUpdate,
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

	planName := parseNameID(trackingPlan.Name)
	r.SetId(planName)
	return resourceSegmentTrackingPlanRead(r, meta)
}

func resourceSegmentTrackingPlanRead(r *schema.ResourceData, meta interface{}) error {
	client := meta.(*segment.Client)
	planName := r.Id()
	trackingPlan, err := client.GetTrackingPlan(planName)
	if err != nil {
		return err
	}

	s, _ := json.Marshal(trackingPlan.Rules)
	r.Set("display_name", trackingPlan.DisplayName)
	r.Set("rules", string(s))

	return nil
}

func resourceSegmentTrackingPlanDelete(r *schema.ResourceData, meta interface{}) error {
	client := meta.(*segment.Client)
	planName := r.Id()

	err := client.DeleteTrackingPlan(planName)
	if err != nil {
		return err
	}

	return nil
}

func resourceSegmentTrackingPlanUpdate(r *schema.ResourceData, meta interface{}) error {
	client := meta.(*segment.Client)
	planName := r.Id()
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

func parseNameID(name string) string {
	nameSplit := strings.Split(name, "/")
	return nameSplit[len(nameSplit)-1]
}
