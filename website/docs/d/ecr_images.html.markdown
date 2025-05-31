---
subcategory: "ECR (Elastic Container Registry)"
layout: "aws"
page_title: "AWS: aws_ecr_images"
description: |-
  Provides details about images in an ECR repository
---

# Data Source: aws_ecr_images

Provides details about images in an Amazon Elastic Container Registry (ECR) repository.

## Example Usage

```terraform
data "aws_ecr_images" "example" {
  repository_name = "example-repo"
}
```

## Argument Reference

The following arguments are required:

* `repository_name` - (Required) Name of the ECR repository.

The following arguments are optional:

* `registry_id` - (Optional) ID of the registry (AWS account ID).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `image_tags` - List of all image tags in the repository.
* `image_digests` - List of all image digests in the repository.
