package colab_test

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccColabSchedule_update(t *testing.T) {
	t.Parallel()
	acctest.BootstrapIamMembers(t, []acctest.IamMember{
		{
			Member: "serviceAccount:service-{project_number}@gcp-sa-dataform.iam.gserviceaccount.com",
			Role:   "roles/cloudkms.cryptoKeyEncrypterDecrypter",
		},
	})

	context := map[string]interface{}{
		"location":           envvar.GetTestRegionFromEnv(),
		"project_id":         envvar.GetTestProjectFromEnv(),
		"service_account":    envvar.GetTestServiceAccountFromEnv(t),
		"end_time":           time.Now().AddDate(0, 0, 10).Format(time.RFC3339),
		"key_name":           acctest.BootstrapKMSKeyInLocation(t, "us-central1").CryptoKey.Name,
		"start_time":         time.Now().AddDate(0, 0, 1).Format(time.RFC3339),
		"random_suffix":      acctest.RandString(t, 10),
		"updated_start_time": time.Now().AddDate(0, 0, 2).Format(time.RFC3339),
		"updated_end_time":   time.Now().AddDate(0, 0, 5).Format(time.RFC3339),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckColabScheduleDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccColabSchedule_full(context),
			},
			{
				ResourceName:            "google_colab_schedule.schedule",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
			{
				Config: testAccColabSchedule_update(context),
			},
			{
				ResourceName:            "google_colab_schedule.schedule",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
		},
	})
}

func testAccColabSchedule_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_colab_runtime_template" "my_runtime_template" {
  name = ""
  display_name = "Runtime template"
  location = "us-central1"

  machine_spec {
    machine_type     = "e2-standard-4"
  }

  network_spec {
    enable_internet_access = true
  }
}

resource "google_storage_bucket" "output_bucket" {
  name          = "tf_test_my_bucket%{random_suffix}"
  location      = "US"
  force_destroy = true
  uniform_bucket_level_access = true
}

resource "google_secret_manager_secret" "secret" {
  secret_id = "secret%{random_suffix}"
  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "secret_version" {
  secret = google_secret_manager_secret.secret.id
  secret_data = "secret-data"
}

resource "google_dataform_repository" "dataform_repository" {
  name = "tf-test-dataform-repository%{random_suffix}"
  display_name = "dataform_repository"
  npmrc_environment_variables_secret_version = google_secret_manager_secret_version.secret_version.id
  kms_key_name = "%{key_name}"

  labels = {
    label_foo1 = "label-bar1"
  }

  git_remote_settings {
      url = "https://github.com/OWNER/REPOSITORY.git"
      default_branch = "main"
      authentication_token_secret_version = google_secret_manager_secret_version.secret_version.id
  }

  workspace_compilation_overrides {
    default_database = "database"
    schema_suffix = "_suffix"
    table_prefix = "prefix_"
  }

}

resource "google_colab_schedule" "schedule" {
  display_name = "Notebook Schedule full"
  location = "%{location}"
  allow_queueing = true
  max_concurrent_run_count = 2
  cron = "TZ=America/Los_Angeles * * * * *"
  max_run_count = 5
  start_time = "%{start_time}"
  end_time = "%{end_time}"

  create_notebook_execution_job_request {
    parent = "projects/%{project_id}/locations/%{location}"
    notebook_execution_job {
      display_name = "Notebook execution"
      execution_timeout = "86400s"

      dataform_repository_source {
        commit_sha = "randomsha123"
        dataform_repository_resource_name = "projects/%{project_id}/locations/%{location}/repositories/${google_dataform_repository.dataform_repository.name}"
      }

      notebook_runtime_template_resource_name = "projects/${google_colab_runtime_template.my_runtime_template.project}/locations/${google_colab_runtime_template.my_runtime_template.location}/notebookRuntimeTemplates/${google_colab_runtime_template.my_runtime_template.name}"

      gcs_output_uri = "gs://${google_storage_bucket.output_bucket.name}"
      service_account = "%{service_account}"
    }
  }

  depends_on = [
    google_colab_runtime_template.my_runtime_template,
    google_storage_bucket.output_bucket,
    google_secret_manager_secret_version.secret_version,
    google_dataform_repository.dataform_repository,
  ]
}
`, context)
}

func testAccColabSchedule_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_colab_runtime_template" "my_runtime_template" {
  name = ""
  display_name = "Runtime template"
  location = "us-central1"

  machine_spec {
    machine_type     = "e2-standard-4"
  }

  network_spec {
    enable_internet_access = true
  }
}

resource "google_storage_bucket" "output_bucket" {
  name          = "tf_test_my_bucket%{random_suffix}"
  location      = "US"
  force_destroy = true
  uniform_bucket_level_access = true
}

resource "google_secret_manager_secret" "secret" {
  secret_id = "secret%{random_suffix}"
  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "secret_version" {
  secret = google_secret_manager_secret.secret.id
  secret_data = "secret-data"
}

resource "google_dataform_repository" "dataform_repository" {
  name = "tf-test-dataform-repository%{random_suffix}"
  display_name = "dataform_repository"
  npmrc_environment_variables_secret_version = google_secret_manager_secret_version.secret_version.id
  kms_key_name = "%{key_name}"

  labels = {
    label_foo1 = "label-bar1"
  }

  git_remote_settings {
      url = "https://github.com/OWNER/REPOSITORY.git"
      default_branch = "main"
      authentication_token_secret_version = google_secret_manager_secret_version.secret_version.id
  }

  workspace_compilation_overrides {
    default_database = "database"
    schema_suffix = "_suffix"
    table_prefix = "prefix_"
  }

}

resource "google_colab_schedule" "schedule" {
  display_name = "Notebook Schedule updated"
  location = "%{location}"
  allow_queueing = false
  max_concurrent_run_count = 1
  cron = "TZ=America/Los_Angeles 0 * * * *"
  max_run_count = 3
  start_time = "%{updated_start_time}"
  end_time = "%{updated_end_time}"

  create_notebook_execution_job_request {
    parent = "projects/%{project_id}/locations/%{location}"
    notebook_execution_job {
      display_name = "Notebook execution"
      execution_timeout = "86400s"

      dataform_repository_source {
        commit_sha = "randomsha123"
        dataform_repository_resource_name = "projects/%{project_id}/locations/%{location}/repositories/${google_dataform_repository.dataform_repository.name}"
      }

      notebook_runtime_template_resource_name = "projects/${google_colab_runtime_template.my_runtime_template.project}/locations/${google_colab_runtime_template.my_runtime_template.location}/notebookRuntimeTemplates/${google_colab_runtime_template.my_runtime_template.name}"

      gcs_output_uri = "gs://${google_storage_bucket.output_bucket.name}"
      service_account = "%{service_account}"
    }
  }

  depends_on = [
    google_colab_runtime_template.my_runtime_template,
    google_storage_bucket.output_bucket,
    google_secret_manager_secret_version.secret_version,
    google_dataform_repository.dataform_repository,
  ]
}
`, context)
}
