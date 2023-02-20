package test

import (
	"fmt"
        "strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/terraform"
        "cloud.google.com/go/storage"
)

func TestTerraformGcp(t *testing.T) {
	t.Parallel()

        terraformDir := "../dev"
        projectID := gcp.GetGoogleProjectIDFromEnvVar(t)
     
        bucketName := fmt.Sprintf("test-tf-gcs-bucket-%s", "01")
    
        gcp.CreateStorageBucket(t,projectID,bucketName, &storage.BucketAttrs{Location: "EU"})

        key := fmt.Sprintf("%s/terraform.tfstate", "01")
        data := "gcp_single.json"
 
	terraformOptions := &terraform.Options{
		TerraformDir: terraformDir,

         BackendConfig: map[string]interface{}{
                        "bucket": bucketName,
                        "prefix":    key,
                        "credentials": data,
                },

        Vars: map[string]interface{}{
            "env"                        : "cd3",
            "region"                     : "europe-west3",
            "billing_account"            : "01B7CB-3DEFDD-94C950",
            "org_id"                     : "terracloud-377521",
            "zones"                      : []string{"europe-west3-a", "europe-west3-b"},
            "webservers_subnet_ip_range" : "192.168.1.0/24",
            "management_subnet_ip_range" : "192.168.100.0/24",
            "bastion_image"              : "centos-7-v20170918",
            "bastion_instance_type"      : "e2-micro",
            "user"                       : "srikcd3",
            "ssh_key"                    : "gcp_single.json",
            "db_region"                  : "europe-west3",
            "appserver_count"            : "1",
            "app_image"                  : "centos-7-v20170918",
            "app_instance_type"          : "e2-micro",
            "project_name"               : "terracloud-test",
            "project_id"                 : projectID,
            "network_name"               : "terracloudcdnetwork",
            "db_name"                    : "cd3db",
            "instace_template_name"      : "cd3temp",
            "webservers_subnet_name"     : "webcd3",
            "management_subnet_name"     : "mgmtcd3",
            "user_name"                  : "tempcd3", 
            "user_password"              : "tempcd3",
            "owner"                      : "sricd3",
        },

	}


	// Destroy all resources in any exit case
	defer terraform.Destroy(t, terraformOptions)

        defer gcp.EmptyStorageBucket(t,bucketName)
        defer gcp.DeleteStorageBucket(t,bucketName)

	// Run terraform init and apply
	terraform.InitAndApply(t, terraformOptions)
        instanceName := terraform.Output(t, terraformOptions, "instance_name")

	instance := gcp.FetchInstance(t, projectID, instanceName)
        instance.SetLabels(t, map[string]string{"environmentname": "cd3"})

        expectedText := "cd3"
	maxRetries := 40
	sleepBetweenRetries := 2 * time.Second

	retry.DoWithRetry(t, fmt.Sprintf("Checking Instance %s for labels", instanceName), maxRetries, sleepBetweenRetries, func() (string, error) {
		// Look up the tags for the given Instance ID
		instance := gcp.FetchInstance(t, projectID, instanceName)
		instanceLabels := instance.GetLabels(t)

		testingTag, containsTestingTag := instanceLabels["environmentname"]
		actualText := strings.TrimSpace(testingTag)
		if !containsTestingTag {
			return "", fmt.Errorf("Expected the tag 'environmentname' to exist")
		}

		if actualText != expectedText {
			return "", fmt.Errorf("Expected GetLabelsForComputeInstanceE to return '%s' but got '%s'", expectedText, actualText)
		}

		return "", nil
	})

}
