[![Go Report Card](https://goreportcard.com/badge/github.com/mbode/terraform-state-mover)](https://goreportcard.com/report/github.com/mbode/terraform-state-mover)
[![Actions Status](https://github.com/mbode/terraform-state-mover/workflows/Check/badge.svg)](https://github.com/mbode/terraform-state-mover/actions)
[![codecov](https://codecov.io/gh/mbode/terraform-state-mover/branch/master/graph/badge.svg)](https://codecov.io/gh/mbode/terraform-state-mover)
[![Dependabot Status](https://api.dependabot.com/badges/status?host=github&repo=mbode/terraform-state-mover)](https://dependabot.com)
[![License](https://img.shields.io/github/license/mbode/terraform-state-mover)](https://github.com/mbode/terraform-state-mover/blob/master/LICENSE)
[![Release](https://img.shields.io/github/v/release/mbode/terraform-state-mover)](https://github.com/mbode/terraform-state-mover/releases/latest)

# Terraform State Mover

Helps refactoring terraform code by offering an interactive prompt for the [`terraform state mv`](https://www.terraform.io/docs/commands/state/mv.html) command.

## Installation

Using [homebrew](https://brew.sh/):
```bash
brew install mbode/tap/terraform-state-mover
```

Alternatively, get a pre-built binary from the [latest release](https://github.com/mbode/terraform-state-mover/releases/latest) or build it yourself using

```bash
go get github.com/mbode/terraform-state-mover
```

## Usage

```bash
terraform-state-mover # inside a Terraform root directory
```

Extra arguments after a `--` are passed to the `terraform plan` call. This makes the following possible:
```bash
terraform-state-mover -- -var key=value  # setting variables
terraform-state-mover -- -var-file=variables.tfvars  # using variable files
```

*Hint:*
If terraform-state-mover is used with the [Google Cloud Platform provider](https://www.terraform.io/docs/providers/google/index.html) and remote state, it is recommended to use `--delay=2s`, otherwise API rate limits error might occur.

## Demo

![](demo.gif)

## Terraform version compatibility

| < 0.10.1 | 0.10.1…8 | 0.11.0…14 | 0.12.0…30 | 0.13.0…6 | 0.14.0…7 |
|:--------:|:--------:|:---------:|:---------:|:--------:|:--------:|
| ✗        | ✓        | ✓         | ✓         | ✓        | ✓        |

## Contributing
Pull requests are welcome. Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)
