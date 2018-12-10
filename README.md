# terraform-provider-segment

A [Terraform](https://www.terraform.io/) provider for [Segment](https://www.segment.com)

## Usage

Create and manage Segment [sources](https://segment.com/docs/sources/)
```
resource "segment_source" "test_source" {
  source_name = "your-source"
  catalog_name = "catalog/sources/javascript"
}
```

Create and manage Segment [destinations](https://segment.com/docs/destinations/)
```
resource "segment_destination" "test_destination" {
  source_name = "test_source"
  destination_name = "google-analytics"
  connection_mode = "CLOUD"
  enabled = false
  configs = [{
      name = "workspaces/your-workspace/sources/your-source/destinations/google-analytics/config/trackingId"
      type = "string"
      value = "your-tracking-id"
  }]
}
```