# gitlab-variables-scan

A small utility to help find sensitive data in GitLab variables.

## Configuration

```yaml
---
endpoint: https://gitlab.com/
token: ABCD
groupIDs:
  - 1
  - 2
exclude:
  keys:
    - PUBLIC_TOKEN
include:
  keys:
    - KEY$
    - TOKEN$
    - SECRET$
    - PASSWORD$
  values:
    - BEGIN PRIVATE KEY
    - ^s\.(.*){24}$
  pairs:
    - LOGIN=guest
```

## Usage

```
Usage of ./gitlab-variables-scan:
  -config string
    	Path to configuration file. (default "config.yaml")
  -debug
    	Enable debug logs.
  -k	Prints values of matched variables.
```
