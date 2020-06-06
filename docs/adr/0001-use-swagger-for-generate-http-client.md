# Need a GoBased ApiClient for Harbor

* Status: accepted
* Date: 2020-05-01

## Context and Problem Statement

For a quick development from the Terraform Provider Functions, it makes sense to generate or use a HarborRestAPI Client.


## Considered Options

* Self Written
* Existing
* Swagger Based

## Decision Outcome

Chosen option: "Swagger Based", because this solution supports the fastes development Start without writting any boilerplate code.

### Positive Consequences

* No Painfull HTTP Client Implementation

### Negative Consequences

* the API Client Implementation dependents to the Swagger Spec Quality...

## Pros and Cons of the Options 

### Self Written

Using a SelfWritten Api Client like the [BESTSELLER/terraform-harbor-provider](https://github.com/BESTSELLER/terraform-harbor-provider) implementation.

* Good, because full flexibility, for api interactions.
* Good, because no external dependency.
* Bad, because it make no fun to write a API Client, with Error Handling, Models etc.

### Existing

Using one of the Exising projects like: [moooofly/harbor-go-client](https://github.com/moooofly/harbor-go-client), [TimeBye/go-harbor](https://github.com/TimeBye/go-harbor) or [jiankunking/harbor-go-client](https://github.com/jiankunking/harbor-go-client)

* Good, because no painfull Api Code.
* Bad, because a hard external dependency to the supported client methodes.
* Bad, because it needs changes at all existing implementation for use the full harbor Api Sets.

### Swagger Based

[Harbor Swagger](https://goharbor.io/docs/1.10/build-customize-contribute/configure-swagger/)

* Good, because no self written API Code.
* Good, because it support generate clients for different Harbor API Versions.
* Bad, because depends to the SwaggerSpec Quality

#### Swagger MissMatches

Page for Colleting Swagger ApiClient Issues.

* `#/definitions/ReplicationPolicy`

```json
"filters": [
        {
            "type": "name",
            "value": "test"
        },
        ...
        {
            "type": "label",
            "value": [
                "testlabel-acc-classic"
            ]
        }
        ...
    ]
...
```

```swagger
"filters": [
      {
        "type": "string",
        "value": "string"
      }
    ],
````



## Links <!-- optional -->

* [Link type] [Link to ADR] <!-- example: Refined by [ADR-0005](0005-example.md) -->
* … <!-- numbers of links can vary -->
