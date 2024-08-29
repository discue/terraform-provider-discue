<p align="center"><a href="https://www.discue.io/" target="_blank" rel="noopener noreferrer"><img width="128" src="https://www.discue.io/icons-fire-no-badge-square/web/icon-192.png" alt="Vue logo"></a></p>

<br/>
<div align="center">


[![contributions - welcome](https://img.shields.io/badge/contributions-welcome-blue/green)](/CONTRIBUTING.md "Go to contributions doc")
[![GitHub License](https://img.shields.io/github/license/discue/terraform-provider-discue.svg)](https://github.com/discue/terraform-provider-discue/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/discue/terraform-provider-discue)](https://goreportcard.com/report/github.com/discue/terraform-provider-discue)
<br/>
[![Go](https://img.shields.io/badge/Go->=1.21-blue?logo=logo&logoColor=white)](https://nodejs.org "Go to Node.js homepage")
[![Made with Node.js](https://img.shields.io/badge/Terraform->=1.8-blue?logo=terraform&logoColor=white)](https://nodejs.org "Go to Node.js homepage")

</div>

<br/>

# terraform-provider-discue
A `terraform` provider that allows managing resources of [discue.io](https://www.discue.io).

## Development

If you're new to provider development, a good place to start is the [Extending
Terraform](https://www.terraform.io/docs/extend/index.html) docs.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21

## Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Install the dependenies using `go install`
3. Build the provider using the script `build.sh`

```shell
./build.sh
```

## Running the provider
If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

Then, create a `.terraformrc` file in your operating system user directory and paste the following
```
provider_installation {
   dev_overrides {
      "grafana/grafana" = "/path/to/your/terraform-provider-discue" # this path is the directory where the binary is built
  }
  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

## Configuring the provider
Set the following mandator environment variables

### Required
- `DISCUE_API_KEY`: The api key used to access the target organization

### Optional
- `DISCUE_API_ENDPOINT`: The target endpoint e.g. http://localhost:3000, if the API is running locally

## Testing the provider
In order to run the full suite of Acceptance tests, run `./test.sh`.

## Generating documentation
To generate or update documentation, run `./generate-docs.sh`.

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.