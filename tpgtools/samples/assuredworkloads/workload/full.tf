resource "google_assured_workloads_workload" "primary" {
  display_name = "{{name}}"
  billing_account = "billingAccounts/{{billing_account}}"
  compliance_regime = "FEDRAMP_MODERATE"
  organization = "{{org_id}}"
  location = "us-central1"
  kms_settings {
    next_rotation_time = "2021-10-02T15:01:23Z"
    rotation_period = "864000s"
  }
  provisioned_resources_parent = google_folder.folder1.name
}

resource "google_folder" "folder1" {
  display_name = "{{name}}"
  parent       = "organizations/{{org_id}}"
}