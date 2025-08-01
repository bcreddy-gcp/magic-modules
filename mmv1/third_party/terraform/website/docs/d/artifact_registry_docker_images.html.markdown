---
subcategory: "Artifact Registry"
description: |-
  Get information about Docker images within a Google Artifact Registry repository.
---

# google_artifact_registry_docker_images

Get information about Artifact Registry Docker images.
See [the official documentation](https://cloud.google.com/artifact-registry/docs/docker)
and [API](https://cloud.google.com/artifact-registry/docs/reference/rest/v1/projects.locations.repositories.dockerImages/list).

## Example Usage

```hcl
data "google_artifact_registry_docker_images" "my_images" {
  location      = "us-central1"
  repository_id = "example-repo"
}
```

## Argument Reference

The following arguments are supported:

* `location` - (Required) The location of the Artifact Registry repository.

* `repository_id` - (Required) The last part of the repository name to fetch from.

* `project` - (Optional) The project ID in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

The following attributes are exported:

* `docker_images` - A list of all retrieved Artifact Registry Docker images. Structure is [defined below](#nested_docker_images).

<a name="nested_docker_images"></a>The `docker_images` block supports:

* `name` - The fully qualified name of the fetched image.  This name has the form: `projects/{{project}}/locations/{{location}}/repository/{{repository_id}}/dockerImages/{{docker_image}}`. For example, `projects/test-project/locations/us-west4/repositories/test-repo/dockerImages/nginx@sha256:e9954c1fc875017be1c3e36eca16be2d9e9bccc4bf072163515467d6a823c7cf`

* `image_name` - Extracted short name of the image (last part of `name`, without tag or digest). For example, from `.../nginx@sha256:...` → `nginx`.

* `self_link` - The URI to access the image.  For example, `us-west4-docker.pkg.dev/test-project/test-repo/nginx@sha256:e9954c1fc875017be1c3e36eca16be2d9e9bccc4bf072163515467d6a823c7cf`

* `tags` - A list of all tags associated with the image.

* `image_size_bytes` - Calculated size of the image in bytes.

* `media_type` - Media type of this image, e.g. `application/vnd.docker.distribution.manifest.v2+json`. 

* `upload_time` - The time, as a RFC 3339 string, the image was uploaded. For example, `2014-10-02T15:01:23.045123456Z`.

* `build_time` - The time, as a RFC 3339 string, this image was built. 

* `update_time` - The time, as a RFC 3339 string, this image was updated.
