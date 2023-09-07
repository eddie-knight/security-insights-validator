# Security Insights Validator

!! Currently in Proof of Concept phase only. !!

## Intent

This is intended to compare a provided YAML file against the security insights specification.

In the event that a known value is found to contain an unexpected data type, the validation
will fail one entry at a time, beginning at the top of the file.

In the event that an unknown value is included, it will be ignored during validation. Upon
successful completion of the validation, a diff will provided in the execution log to show
all unexpected values that were included in the provided YAML file.

## Usage

1. Download the latest binary from the releases page
1. Execute the file locally using `./si-validator`
    - By default, the validator will look in the present working directory for `SECURITY-INSIGHTS.yml`
    - A full local path can be provided to a file in another location or with another name using `./si-validator --input path/to/file.yml`
1. Validation results will be printed to the log

## Roadmap

Continued development is subject to discussion with the Security Insights community.
