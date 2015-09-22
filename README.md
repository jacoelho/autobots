# autobots

aws auto scaling group helper.

pass one or more auto scaling groups and using instance credentials it will return instances in that group.

### Example

```bash
autobots --auto-scaling-groups <auto-scaling-group> --output private-dns

output:
ip-1-west-1.compute.internal ip-2-us-west-1.compute.internal ip-3-west-1.compute.internal
```

### Notes

If using a proxy server dont forget to set no_proxy:

```
export no_proxy=169.254.169.254
```

You need access to the aws api servers.
