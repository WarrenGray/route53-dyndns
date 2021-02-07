# route53-dyndns
A small utility for updating Route53 DNS entries

## Installation
```
go get -u github.com/warrengray/route53-dyndns
```

## Usage
`route53-dyndns` supports the same [configuration options](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html) 
as the [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-welcome.html). Use `aws configure` to set up
default profiles, or use the [environment variables](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html) 
to set per-use configuration.

Ensure the profile that you're using has appropriate permission to update the hosted zone, and run:
```
route53-dyndns <hosted zone ID> <fqdn> <ttl>
```
