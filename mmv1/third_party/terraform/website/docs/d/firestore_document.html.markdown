---
subcategory: "Firestore"
description: |-
  Read a document from a Firestore database
---


# google_firestore_document

Reads a document from a Firestore database.
See [the official documentation](https://cloud.google.com/firestore/native/docs/)
and
[API](https://cloud.google.com/firestore/docs/reference/rest/v1/projects.databases.documents/get/).


## Example Usage

Retrieve a document from the Firestore database.

```hcl
resource "google_firestore_document" "mydoc" {
  project     = google_firestore_database.database.project
  database    = google_firestore_database.database.name
  collection  = "somenewcollection"
  document_id = "my-doc-id"
}
```

## Argument Reference

The following arguments are supported:

* `database` - (Required) The name of the Firestore database.

* `collection` - (Required) The name of the collection of documents.

* `document_id` - (Required) The id of the document to get.

* `project` - (Optional) The project in which the database resides.

## Attributes Reference

See [google_firestore_document](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/google_firestore_document) resource for details of the available attributes.
