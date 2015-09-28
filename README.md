# autobots

aws auto scaling group helper.

it will check if current instance is member of an auto scaling group and output members information.

### Example

```bash
autobots --output instance-id

output:
i-fffffa i-fffffb i-fffffc
```

Or

```bash
autobots --with-asg <auto-scaling-group> --output private-dns

output:
ip-1-west-1.compute.internal ip-2-us-west-1.compute.internal ip-3-west-1.compute.internal
```

Available outputs:

* private-ip (default)
* private-dns
* public-ip
* public-dns
* hostname
* instance-id

### Notes

If using a proxy server dont forget to set no_proxy:

```
export no_proxy=169.254.169.254
```

You need access to the aws api servers.
