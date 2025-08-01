// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccNetappBackup_NetappBackupFull_update(t *testing.T) {
	context := map[string]interface{}{
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "gcnv-network-config-3", acctest.ServiceNetworkWithParentService("netapp.servicenetworking.goog")),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetappBackupDestroyProducer(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNetappBackup_NetappBackupFromVolumeSnapshot(context),
			},
			{
				ResourceName:            "google_netapp_backup.test_backup",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "name", "terraform_labels", "vault_name"},
			},
			{
				Config: testAccNetappBackup_NetappBackupFromVolumeSnapshot_update(context),
			},
			{
				ResourceName:            "google_netapp_backup.test_backup",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "name", "terraform_labels", "vault_name"},
			},
		},
	})
}

func testAccNetappBackup_NetappBackupFromVolumeSnapshot(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_compute_network" "default" {
  name = "%{network_name}"
}

resource "google_netapp_storage_pool" "default" {
  name = "tf-test-backup-pool%{random_suffix}"
  location = "us-west2"
  service_level = "PREMIUM"
  capacity_gib = "2048"
  network = data.google_compute_network.default.id
}

resource "time_sleep" "wait_3_minutes" {
  depends_on = [google_netapp_storage_pool.default]
  create_duration = "3m"
}

resource "google_netapp_volume" "default" {
  name = "tf-test-backup-volume%{random_suffix}"
  location = google_netapp_storage_pool.default.location
  capacity_gib = "100"
  share_name = "tf-test-backup-volume%{random_suffix}"
  storage_pool = google_netapp_storage_pool.default.name
  protocols = ["NFSV3"]
  deletion_policy = "FORCE"
  backup_config {
    backup_vault = google_netapp_backup_vault.default.id
  }
}

resource "google_netapp_backup_vault" "default" {
  name = "tf-test-backup-vault%{random_suffix}"
  location = google_netapp_storage_pool.default.location
}

resource "google_netapp_volume_snapshot" "default" {
	depends_on = [google_netapp_volume.default]
	location = google_netapp_volume.default.location
	volume_name = google_netapp_volume.default.name
	description = "This is a test description"
	name = "testvolumesnap%{random_suffix}"
	labels = {
	  key= "test"
	  value= "snapshot"
	}
  }

resource "google_netapp_backup" "test_backup" {
  name = "tf-test-test-backup%{random_suffix}"
  description = "This is a test backup"
  source_volume = google_netapp_volume.default.id
  location = google_netapp_backup_vault.default.location
  vault_name = google_netapp_backup_vault.default.name
  source_snapshot = google_netapp_volume_snapshot.default.id
  labels = {
	key= "test"
	value= "backup"
  }
}
`, context)
}

func testAccNetappBackup_NetappBackupFromVolumeSnapshot_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_compute_network" "default" {
  name = "%{network_name}"
}

resource "google_netapp_storage_pool" "default" {
  name = "tf-test-backup-pool%{random_suffix}"
  location = "us-west2"
  service_level = "PREMIUM"
  capacity_gib = "2048"
  network = data.google_compute_network.default.id
}

resource "time_sleep" "wait_3_minutes" {
  depends_on = [google_netapp_storage_pool.default]
  create_duration = "3m"
}

resource "google_netapp_volume" "default" {
  name = "tf-test-backup-volume%{random_suffix}"
  location = google_netapp_storage_pool.default.location
  capacity_gib = "100"
  share_name = "tf-test-backup-volume%{random_suffix}"
  storage_pool = google_netapp_storage_pool.default.name
  protocols = ["NFSV3"]
  deletion_policy = "FORCE"
  backup_config {
    backup_vault = google_netapp_backup_vault.default.id
  }
}

resource "google_netapp_backup_vault" "default" {
  name = "tf-test-backup-vault%{random_suffix}"
  location = google_netapp_storage_pool.default.location
}

resource "google_netapp_volume_snapshot" "default" {
	depends_on = [google_netapp_volume.default]
	location = google_netapp_volume.default.location
	volume_name = google_netapp_volume.default.name
	description = "This is a test description"
	name = "testvolumesnap%{random_suffix}"
	labels = {
	  key= "test"
	  value= "snapshot"
	}
  }

resource "google_netapp_backup" "test_backup" {
  name = "tf-test-test-backup%{random_suffix}"
  description = "This is a test backup"
  source_volume = google_netapp_volume.default.id
  location = google_netapp_backup_vault.default.location
  vault_name = google_netapp_backup_vault.default.name
  source_snapshot = google_netapp_volume_snapshot.default.id
  labels = {
	key= "test_update"
	value= "backup_update"
  }
}
`, context)
}

func TestAccNetappBackup_NetappFlexBackup(t *testing.T) {
	context := map[string]interface{}{
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "gcnv-network-config-3", acctest.ServiceNetworkWithParentService("netapp.servicenetworking.goog")),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetappBackupDestroyProducer(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNetappBackup_FlexBackup(context),
			},
			{
				ResourceName:            "google_netapp_backup.test_backup",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "name", "terraform_labels", "vault_name"},
			},
		},
	})
}

func testAccNetappBackup_FlexBackup(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_compute_network" "default" {
  name = "%{network_name}"
}

resource "google_netapp_storage_pool" "default" {
  name = "tf-test-backup-pool%{random_suffix}"
  location = "us-east4"
  service_level = "FLEX"
  capacity_gib = "2048"
  network = data.google_compute_network.default.id
  zone = "us-east4-a"
  replica_zone = "us-east4-b"
}

resource "time_sleep" "wait_3_minutes" {
  depends_on = [google_netapp_storage_pool.default]
  create_duration = "3m"
}

resource "google_netapp_volume" "default" {
  name = "tf-test-backup-volume%{random_suffix}"
  location = google_netapp_storage_pool.default.location
  capacity_gib = "100"
  share_name = "tf-test-backup-volume%{random_suffix}"
  storage_pool = google_netapp_storage_pool.default.name
  protocols = ["NFSV3"]
  deletion_policy = "FORCE"
  backup_config {
    backup_vault = google_netapp_backup_vault.default.id
  }
}

resource "google_netapp_backup_vault" "default" {
  name = "tf-test-backup-vault%{random_suffix}"
  location = google_netapp_storage_pool.default.location
}

resource "google_netapp_volume_snapshot" "default" {
    depends_on = [google_netapp_volume.default]
    location = google_netapp_volume.default.location
    volume_name = google_netapp_volume.default.name
    description = "This is a test description"
    name = "testvolumesnap%{random_suffix}"
    labels = {
      key= "test"
      value= "snapshot"
    }
  }

resource "google_netapp_backup" "test_backup" {
  name = "tf-test-test-backup%{random_suffix}"
  description = "This is a flex test backup"
  source_volume = google_netapp_volume.default.id
  location = google_netapp_backup_vault.default.location
  vault_name = google_netapp_backup_vault.default.name
  source_snapshot = google_netapp_volume_snapshot.default.id
  labels = {
    key= "test"
    value= "backup"
  }
}
`, context)
}

func TestAccNetappBackup_NetappIntegratedBackup(t *testing.T) {
	context := map[string]interface{}{
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "gcnv-network-config-3", acctest.ServiceNetworkWithParentService("netapp.servicenetworking.goog")),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetappBackupDestroyProducer(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNetappBackup_IntegratedBackup(context),
			},
			{
				ResourceName:            "google_netapp_backup.test_backup",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "name", "terraform_labels", "vault_name"},
			},
		},
	})
}

func testAccNetappBackup_IntegratedBackup(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_compute_network" "default" {
  name = "%{network_name}"
}
resource "google_netapp_storage_pool" "default" {
  name = "tf-test-backup-pool%{random_suffix}"
  location = "us-east4"
  service_level = "PREMIUM"
  capacity_gib = "2048"
  network = data.google_compute_network.default.id
}
resource "time_sleep" "wait_3_minutes" {
  depends_on = [google_netapp_storage_pool.default]
  create_duration = "3m"
}
resource "google_netapp_volume" "default" {
  name = "tf-test-backup-volume%{random_suffix}"
  location = google_netapp_storage_pool.default.location
  capacity_gib = "100"
  share_name = "tf-test-backup-volume%{random_suffix}"
  storage_pool = google_netapp_storage_pool.default.name
  protocols = ["NFSV3"]
  deletion_policy = "FORCE"
  backup_config {
    backup_vault = google_netapp_backup_vault.default.id
  }
}
resource "google_netapp_backup_vault" "default" {
  name = "tf-test-backup-vault%{random_suffix}"
  location = google_netapp_storage_pool.default.location
  backup_vault_type = "CROSS_REGION"
  backup_region = "us-west4"
}
resource "google_netapp_volume_snapshot" "default" {
  depends_on = [google_netapp_volume.default]
  location = google_netapp_volume.default.location
  volume_name = google_netapp_volume.default.name
  description = "This is a test description"
  name = "testvolumesnap%{random_suffix}"
  labels = {
    key= "test"
    value= "snapshot"
  }
  }
resource "google_netapp_backup" "test_backup" {
  name = "tf-test-test-backup%{random_suffix}"
  description = "This is a test integrated backup"
  source_volume = google_netapp_volume.default.id
  location = google_netapp_backup_vault.default.location
  vault_name = google_netapp_backup_vault.default.name
  source_snapshot = google_netapp_volume_snapshot.default.id
  labels = {
  key= "test"
  value= "backup"
  }
}
`, context)
}

func TestAccNetappBackup_NetappImmutableBackup(t *testing.T) {
	context := map[string]interface{}{
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "gcnv-network-config-3", acctest.ServiceNetworkWithParentService("netapp.servicenetworking.goog")),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetappBackupDestroyProducer(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNetappBackup_ImmutableBackup(context),
			},
			{
				ResourceName:            "google_netapp_backup.test_backup",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "name", "terraform_labels", "vault_name"},
			},
			{
				Config: testAccNetappBackup_ImmutableBackupUpdate(context),
			},
			{
				ResourceName:            "google_netapp_backup.test_backup",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "name", "terraform_labels", "vault_name"},
			},
		},
	})
}

func testAccNetappBackup_ImmutableBackup(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_compute_network" "default" {
  name = "%{network_name}"
}
resource "google_netapp_storage_pool" "default" {
  name = "tf-test-backup-pool%{random_suffix}"
  location = "us-central1"
  service_level = "FLEX"
  capacity_gib = "2048"
  network = data.google_compute_network.default.id
  zone = "us-central1-a"
  replica_zone = "us-central1-b"
}
resource "time_sleep" "wait_3_minutes" {
    depends_on = [google_netapp_storage_pool.default]
    create_duration = "3m"
}
resource "google_netapp_volume" "default" {
  name = "tf-test-backup-volume%{random_suffix}"
  location = "us-central1"
  capacity_gib = "100"
  share_name = "tf-test-backup-volume%{random_suffix}"
  storage_pool = google_netapp_storage_pool.default.name
  protocols = ["NFSV3"]
  deletion_policy = "FORCE"
  backup_config {
    backup_vault = google_netapp_backup_vault.default.id
  }
}
resource "google_netapp_backup_vault" "default" {
  name = "tf-test-backup-vault%{random_suffix}"
  location = "us-central1"
  backup_retention_policy {
    backup_minimum_enforced_retention_days = 2
    daily_backup_immutable = true
    weekly_backup_immutable = false
    monthly_backup_immutable = false
    manual_backup_immutable = false
  }
}
resource "google_netapp_volume_snapshot" "default" {
  depends_on = [google_netapp_volume.default]
  location = "us-central1"
  volume_name = google_netapp_volume.default.name
  description = "This is a test description"
  name = "testvolumesnap%{random_suffix}"
  labels = {
    key= "test"
    value= "snapshot"
  }
}
resource "google_netapp_backup" "test_backup" {
  name = "tf-test-test-backup%{random_suffix}"
  description = "This is a test immutable backup"
  source_volume = google_netapp_volume.default.id
  location = "us-central1"
  vault_name = google_netapp_backup_vault.default.name
  source_snapshot = google_netapp_volume_snapshot.default.id
  labels = {
    key= "test"
    value= "backup"
  }
}
`, context)
}

func testAccNetappBackup_ImmutableBackupUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_compute_network" "default" {
  name = "%{network_name}"
}
resource "google_netapp_storage_pool" "default" {
  name = "tf-test-backup-pool%{random_suffix}"
  location = "us-central1"
  service_level = "FLEX"
  capacity_gib = "2048"
  network = data.google_compute_network.default.id
  zone = "us-central1-a"
  replica_zone = "us-central1-b"
}
resource "time_sleep" "wait_3_minutes" {
    depends_on = [google_netapp_storage_pool.default]
    create_duration = "3m"
}
resource "google_netapp_volume" "default" {
  name = "tf-test-backup-volume%{random_suffix}"
  location = "us-central1"
  capacity_gib = "100"
  share_name = "tf-test-backup-volume%{random_suffix}"
  storage_pool = google_netapp_storage_pool.default.name
  protocols = ["NFSV3"]
  deletion_policy = "FORCE"
  backup_config {
    backup_vault = google_netapp_backup_vault.default.id
  }
}
resource "google_netapp_backup_vault" "default" {
  name = "tf-test-backup-vault%{random_suffix}"
  location = "us-central1"
  backup_retention_policy {
    backup_minimum_enforced_retention_days = 12
    daily_backup_immutable = true
    weekly_backup_immutable = true
    monthly_backup_immutable = true
    manual_backup_immutable = true
  }
}
resource "google_netapp_volume_snapshot" "default" {
  depends_on = [google_netapp_volume.default]
  location = "us-central1"
  volume_name = google_netapp_volume.default.name
  description = "This is a test description"
  name = "testvolumesnap%{random_suffix}"
  labels = {
    key= "test"
    value= "snapshot"
  }
}
resource "google_netapp_backup" "test_backup" {
  name = "tf-test-test-backup%{random_suffix}"
  description = "This is a test immutable backup"
  source_volume = google_netapp_volume.default.id
  location = "us-central1"
  vault_name = google_netapp_backup_vault.default.name
  source_snapshot = google_netapp_volume_snapshot.default.id
  labels = {
    key= "test"
    value= "backup"
  }
}
`, context)
}
