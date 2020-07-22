# Need a GoBased ApiClient for Harbor

* Status: accepted
* Date: 2020-05-01

## Context and Problem Statement

For a quick development from the Terraform Provider Functions, it makes sense to generate or use a HarborRestAPI Client.


## Considered Options

* [Self Written](#self-written)
* [Existing](#existing)
* [Swagger Based](#swagger-based)

## Decision Outcome

Chosen option: "Swagger Based", because this solution supports the fastes development Start without writting any boilerplate code.

### Positive Consequences

* No Painfull HTTP Client Implementation

### Negative Consequences

* the API Client Implementation dependents to the Swagger Spec Quality...

## Pros and Cons of the Options 

### Self Written

Using a SelfWritten Api Client like the [BESTSELLER/terraform-harbor-provider](https://github.com/BESTSELLER/terraform-harbor-provider) implementation.

* 👍, because full flexibility, for api interactions.
* 👍, because no external dependency.
* 👎, because it make no fun to write a API Client, with Error Handling, Models etc.

### Existing

Using one of the Exising projects like: [moooofly/harbor-go-client](https://github.com/moooofly/harbor-go-client), [TimeBye/go-harbor](https://github.com/TimeBye/go-harbor) or [jiankunking/harbor-go-client](https://github.com/jiankunking/harbor-go-client)

* 👍, because no painfull Api Code.
* 👎, because a hard external dependency to the supported client methodes.
* 👎, because it needs changes at all existing implementation for use the full harbor Api Sets.

### Swagger Based

[Harbor Swagger](https://goharbor.io/docs/1.10/build-customize-contribute/configure-swagger/)

* 👍, because no self written API Code.
* 👍, because it support generate clients for different Harbor API Versions.
* 👎, because depends to the SwaggerSpec Quality

#### Swagger MissMatches

* The Swagger missmatch will be fixed by a [json-patch](https://sookocheff.com/post/api/understanding-json-patch/).
* As base for the Client we use the Harbor V2, most of the interface are useable for v1.



## Links <!-- optional -->

* Github Harbor Swagger Issue [#12474](https://github.com/goharbor/harbor/issues/12474)
* [evanphx/json-patch](https://github.com/evanphx/json-patch)
* [json-patch](https://sookocheff.com/post/api/understanding-json-patch/)
