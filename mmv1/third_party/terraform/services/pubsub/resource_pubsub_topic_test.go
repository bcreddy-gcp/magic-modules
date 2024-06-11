package pubsub_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccPubsubTopic_update(t *testing.T) {
	t.Parallel()

	topic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubTopicDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubTopic_update(topic, "foo", "bar"),
			},
			{
				ResourceName:            "google_pubsub_topic.foo",
				ImportStateId:           topic,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccPubsubTopic_updateWithRegion(topic, "wibble", "wobble", "us-central1"),
			},
			{
				ResourceName:            "google_pubsub_topic.foo",
				ImportStateId:           topic,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func TestAccPubsubTopic_cmek(t *testing.T) {
	t.Parallel()

	kms := acctest.BootstrapKMSKey(t)
	topicName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	if acctest.BootstrapPSARole(t, "service-", "gcp-sa-pubsub", "roles/cloudkms.cryptoKeyEncrypterDecrypter") {
		t.Fatal("Stopping the test because a role was added to the policy.")
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubTopicDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubTopic_cmek(topicName, kms.CryptoKey.Name),
			},
			{
				ResourceName:      "google_pubsub_topic.topic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPubsubTopic_schema(t *testing.T) {
	t.Parallel()

	schema1 := fmt.Sprintf("tf-test-schema-%s", acctest.RandString(t, 10))
	schema2 := fmt.Sprintf("tf-test-schema-%s", acctest.RandString(t, 10))
	topic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubTopicDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubTopic_updateWithSchema(topic, schema1),
			},
			{
				ResourceName:      "google_pubsub_topic.bar",
				ImportStateId:     topic,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPubsubTopic_updateWithNewSchema(topic, schema2),
			},
			{
				ResourceName:      "google_pubsub_topic.bar",
				ImportStateId:     topic,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPubsubTopic_migration(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()

	topic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))

	oldVersion := map[string]resource.ExternalProvider{
		"google": {
			VersionConstraint: "4.84.0", // a version that doesn't separate user defined labels and system labels
			Source:            "registry.terraform.io/hashicorp/google",
		},
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:     func() { acctest.AccTestPreCheck(t) },
		CheckDestroy: testAccCheckPubsubTopicDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:            testAccPubsubTopic_update(topic, "foo", "bar"),
				ExternalProviders: oldVersion,
			},
			{
				Config:                   testAccPubsubTopic_update(topic, "foo", "bar"),
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
			},
			{
				ResourceName:             "google_pubsub_topic.foo",
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
				ImportStateId:            topic,
				ImportState:              true,
				ImportStateVerify:        true,
				ImportStateVerifyIgnore:  []string{"labels", "terraform_labels"},
			},
		},
	})
}

func TestAccPubsubTopic_kinesisIngestionUpdate(t *testing.T) {
	t.Parallel()

	topic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubTopicDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubTopic_updateWithKinesisIngestionSettings(topic),
			},
			{
				ResourceName:      "google_pubsub_topic.foo",
				ImportStateId:     topic,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPubsubTopic_updateWithUpdatedKinesisIngestionSettings(topic),
			},
			{
				ResourceName:      "google_pubsub_topic.foo",
				ImportStateId:     topic,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccPubsubTopic_update(topic, key, value string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "foo" {
  name = "%s"
  labels = {
    %s = "%s"
  }
}
`, topic, key, value)
}

func testAccPubsubTopic_updateWithRegion(topic, key, value, region string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "foo" {
  name = "%s"
  labels = {
    %s = "%s"
  }

  message_storage_policy {
    allowed_persistence_regions = [
      "%s",
    ]
  }
}
`, topic, key, value, region)
}

func testAccPubsubTopic_cmek(topicName, kmsKey string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "topic" {
  name         = "%s"
  kms_key_name = "%s"
}
`, topicName, kmsKey)
}

func testAccPubsubTopic_updateWithSchema(topic, schema string) string {
	return fmt.Sprintf(`
resource "google_pubsub_schema" "foo" {
	name = "%s"
	type = "PROTOCOL_BUFFER"
  definition = "syntax = \"proto3\";\nmessage Results {\nstring f1 = 1;\n}"
}

resource "google_pubsub_topic" "bar" {
  name = "%s"
	schema_settings {
    schema = google_pubsub_schema.foo.id
    encoding = "BINARY"
  }
}
`, schema, topic)
}

func testAccPubsubTopic_updateWithNewSchema(topic, schema string) string {
	return fmt.Sprintf(`
resource "google_pubsub_schema" "foo" {
	name = "%s"
	type = "PROTOCOL_BUFFER"
	definition = "syntax = \"proto3\";\nmessage Results {\nstring f1 = 1;\n}"
}

resource "google_pubsub_topic" "bar" {
  name = "%s"
	schema_settings {
    schema = google_pubsub_schema.foo.id
    encoding = "JSON"
  }
}
`, schema, topic)
}

func testAccPubsubTopic_updateWithKinesisIngestionSettings(topic string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "foo" {
  name = "%s"

  # Outside of automated terraform-provider-google CI tests, these values must be of actual AWS resources for the test to pass.
  ingestion_data_source_settings {
    aws_kinesis {
        stream_arn = "arn:aws:kinesis:us-west-2:111111111111:stream/fake-stream-name"
        consumer_arn = "arn:aws:kinesis:us-west-2:111111111111:stream/fake-stream-name/consumer/consumer-1:1111111111"
        aws_role_arn = "arn:aws:iam::111111111111:role/fake-role-name"
        gcp_service_account = "fake-service-account@fake-gcp-project.iam.gserviceaccount.com"
    }
  }
}
`, topic)
}

func testAccPubsubTopic_updateWithUpdatedKinesisIngestionSettings(topic string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "foo" {
  name = "%s"

  # Outside of automated terraform-provider-google CI tests, these values must be of actual AWS resources for the test to pass.
  ingestion_data_source_settings {
    aws_kinesis {
        stream_arn = "arn:aws:kinesis:us-west-2:111111111111:stream/updated-fake-stream-name"
        consumer_arn = "arn:aws:kinesis:us-west-2:111111111111:stream/updated-fake-stream-name/consumer/consumer-1:1111111111"
        aws_role_arn = "arn:aws:iam::111111111111:role/updated-fake-role-name"
        gcp_service_account = "updated-fake-service-account@fake-gcp-project.iam.gserviceaccount.com"
    }
  }
}
`, topic)
}
