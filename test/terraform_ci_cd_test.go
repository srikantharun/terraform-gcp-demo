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
        randomRegion := gcp.GetRandomRegion(t, projectID, nil, nil) 
        randomZone := gcp.GetRandomZoneForRegion(t, projectID, "europe-west3")

	terraformOptions := &terraform.Options{
		TerraformDir: terraformDir,


        Vars: map[string]interface{}{
            "env"                        : "ci",
            "region"                     : "europe-west3",
            "billing_account"            : "01B7CB-3DEFDD-94C950",
            "org_id"                     : "terracloud-377520",
            "zones"                      : randomZone,
            "webservers_subnet_ip_range" : "192.168.1.0/24",
            "management_subnet_ip_range" : "192.168.100.0/24",
            "bastion_image"              : "centos-7-v20170918",
            "bastion_instance_type"      : "e2-micro",
            "user"                       : "srikci",
            "ssh_key"                    : "gcp_single.json",
            "db_region"                  : "europe-west3",
            "appserver_count"            : "1",
            "app_image"                  : "Centos-7-v30230203",
            "app_instance_type"          : "e2-micro",
            "project_name"               : "terracloud-test",
            "project_id"                 : projectID,
            "network_name"               : "terraclouddevnetwork",
            "db_name"                    : "cidb",
            "instace_template_name"      : "citemp",
            "webservers_subnet_name"     : "webci",
            "management_subnet_name"     : "mgmtci",
            "user_name"                  : "tempci", 
            "user_password"              : "tempci",
            "owner"                      : "srici",
        },

		EnvVars: map[string]string{
			"GOOGLE_CLOUD_PROJECT": projectID,
		},
	}

	// Destroy all resources in any exit case
	defer terraform.Destroy(t, terraformOptions)

	// Run terraform init and apply
	terraform.InitAndApply(t, terraformOptions)

	// Get the instance group name from the output
	instanceGroupName := terraform.Output(t, terraformOptions, "instance_group_name")

	// Get the instance group
	instanceGroup := gcp.FetchZonalInstanceGroup(t, projectID, randomZone, instanceGroupName)

	maxRetries := 40
	sleepBetweenRetries := 2 * time.Second

	// Check the instance number
	retry.DoWithRetry(t, "Geting instances from, instance group", maxRetries, sleepBetweenRetries, func() (string, error) {
		instances, err := instanceGroup.GetInstancesE(t, projectID)
		if err != nil {
			return "", fmt.Errorf("Failed to get Instances: %s", err)
		}

		if len(instances) != instanceNumber {
			return "", fmt.Errorf("Expected to find exactly %d Compute Instances in Instance Group but found %d.", instanceNumber, len(instances))
		}
		return "", nil
	})
}
