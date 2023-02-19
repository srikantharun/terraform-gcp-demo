package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestTerraformGcp(t *testing.T) {
	t.Parallel()

        instanceNumber := 1
        terraformDir := "../dev"
        projectID := gcp.GetGoogleProjectIDFromEnvVar(t)
        //randomZone := gcp.GetRandomZoneForRegion(t, projectID, "europe-west3")
        randomZone   := "europe-west3-a"
        injname      := "cd2temp-webserver-instance-template"
	terraformOptions := &terraform.Options{
		TerraformDir: terraformDir,


        Vars: map[string]interface{}{
            "env"                        : "cd2",
            "region"                     : "europe-west3",
            "billing_account"            : "01B7CB-3DEFDD-94C950",
            "org_id"                     : "terracloud-377521",
            "zones"                      : []string{"europe-west3-a", "europe-west3-b"},
            "webservers_subnet_ip_range" : "192.168.1.0/24",
            "management_subnet_ip_range" : "192.168.100.0/24",
            "bastion_image"              : "centos-7-v20170918",
            "bastion_instance_type"      : "e2-micro",
            "user"                       : "srikcd2",
            "ssh_key"                    : "gcp_single.json",
            "db_region"                  : "europe-west3",
            "appserver_count"            : "1",
            "app_image"                  : "centos-7-v20170918",
            "app_instance_type"          : "e2-micro",
            "project_name"               : "terracloud-test",
            "project_id"                 : projectID,
            "network_name"               : "terracloudcdnetwork",
            "db_name"                    : "cd2db",
            "instace_template_name"      : "cd2temp",
            "webservers_subnet_name"     : "webcd2",
            "management_subnet_name"     : "mgmtcd2",
            "user_name"                  : "tempcd2", 
            "user_password"              : "tempcd2",
            "owner"                      : "sricd2",
        },

	}


	// Destroy all resources in any exit case
	defer terraform.Destroy(t, terraformOptions)

	// Run terraform init and apply
	terraform.InitAndApply(t, terraformOptions)
        instanceName := terraform.Output(t, terraformOptions, "instance_name")

	instanceGroup := gcp.FetchInstance(t, projectID, injname)
        instance.SetLabels(t, map[string]string{"environmentname": "cd2"})

        expectedText := "cd2"
	maxRetries := 40
	sleepBetweenRetries := 2 * time.Second

	retry.DoWithRetry(t, fmt.Sprintf("Checking Instance %s for labels", instanceName), maxRetries, timeBetweenRetries, func() (string, error) {
		// Look up the tags for the given Instance ID
		instance := gcp.FetchInstance(t, projectId, instanceName)
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
