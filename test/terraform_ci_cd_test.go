package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/gcp"
        "github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/retry"
        "github.com/gruntwork-io/terratest/modules/ssh"
	"github.com/gruntwork-io/terratest/modules/terraform"
        test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
        "github.com/stretchr/testify/assert"
)

func TestTerraformGcp(t *testing.T) {
	t.Parallel()

	env := "ci"
        instanceNumber := 1
        billing_Account = "01B7CB-3DEFDD-94C950"
        org_id = "terracloud-377520"
        webservers_subnet_ip_range = "192.168.1.0/24"
        management_subnet_ip_range = "192.168.100.0/24"
        bastion_image = "centos-7-v20170918"
        bastion_instance_type = "e2-micro"
        user = "srikdev"
        ssh_key = "gcp_single.json"
        db_region = "europe-west3"
	appserver_count := 1
        app_image = "Centos-7-v30230203"
        app_instance_type = "e2-micro"
        project_name = "terracloud-test"
        network_name = "terraclouddevnetwork"
        db_name = "terraclouddevdb"
        instance_template_name = "terraclouddevinstanceci"
        webservers_subnet_name = "devwebservers"
        management_subnet_name = "devmanagement"
        user_name = "hellodev"
        user_password = "hellodev"
        owner = "srik"
	projectID := gcp.GetGoogleProjectIDFromEnvVar(t)
	randomRegion := gcp.GetRandomRegion(t, projectID, nil, nil)
	randomZone := gcp.GetRandomZoneForRegion(t, projectID, randomRegion)
	terraformDir := "../dev/"

	terraformOptions := &terraform.Options{
		TerraformDir: terraformDir,


        Vars: map[string]interface{}{
            "env"                        : env,
            "region"                     : db_region,
            "billing_account"            : billing_Account,
            "org_id"                     : org_id,
            "zones"                      : randonZone,
            "webservers_subnet_ip_range" : webservers_subnet_ip_range,
            "management_subnet_ip_range" : management_subnet_ip_range,
            "bastion_image"              : bastion_image,
            "bastion_instance_type"      : bastion_instance_type,
            "user"                       : user,
            "ssh_key"                    : ssh_key,
            "db_region"                  : randonRegion,
            "appserver_count"            : appserver_count,
            "app_image"                  : app_image,
            "app_instance_type"          : app_instance_type,
            "project_name"               : project_name,
            "project_id"                 : projectID,
            "network_name"               : network_name,
            "db_name"                    : db_name,
            "instace_template_name"      : instance_template_name,
            "webservers_subnet_name"     : webservers_subnet_name,
            "management_subnet_name"     : management_subnet_name,
            "user_name"                  : user_name, 
            "user_password"              : user_password,
            "owner"                      : owner,
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
