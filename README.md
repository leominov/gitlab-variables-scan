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
