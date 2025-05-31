package ecr

import (
	"context"
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
)

// @SDKDataSource("aws_ecr_images")
func DataSourceImages() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceImagesRead,

		Schema: map[string]*schema.Schema{
			"repository_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"registry_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"image_digests": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceImagesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ECRClient(ctx)

	repositoryName := d.Get("repository_name").(string)
	input := &ecr.ListImagesInput{
		RepositoryName: aws.String(repositoryName),
	}

	if v, ok := d.GetOk("registry_id"); ok {
		input.RegistryId = aws.String(v.(string))
	}

	var imageTags []string
	var imageDigests []string

	paginator := ecr.NewListImagesPaginator(conn, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return sdkdiag.AppendErrorf(diags, "listing ECR images: %s", err)
		}

		for _, imageID := range page.ImageIds {
			if imageID.ImageTag != nil {
				imageTags = append(imageTags, *imageID.ImageTag)
			}
			if imageID.ImageDigest != nil {
				imageDigests = append(imageDigests, *imageID.ImageDigest)
			}
		}
	}

	// Remove duplicates and sort for consistent output
	imageTags = removeDuplicatesAndSort(imageTags)
	imageDigests = removeDuplicatesAndSort(imageDigests)

	d.SetId(repositoryName)
	d.Set("image_tags", imageTags)
	d.Set("image_digests", imageDigests)

	return diags
}

func removeDuplicatesAndSort(input []string) []string {
	// Use a map to remove duplicates
	uniqueMap := make(map[string]bool)
	for _, item := range input {
		uniqueMap[item] = true
	}

	// Convert back to slice
	result := make([]string, 0, len(uniqueMap))
	for item := range uniqueMap {
		result = append(result, item)
	}

	// Sort for consistent output
	sort.Strings(result)
	return result
}
