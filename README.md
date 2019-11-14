## Running Make

### Prerequisite Steps

Have active AWS credentials in your current terminal session.

Choose a *Region* to deploy SimplePKI into.

Make a *Bucket* to track terrafrom state and to push the lambda files to. This can be done by running `make bucket-up bucket=<name> region=<region>`.

Now; make the terraform and automation a bit easier; we'll make and env vars file to source.

```
#!/bin/bash 

export bucket=<bucket name>
export region=<region to deploy in>
export TF_VAR_bastion_enabled=<enable bastion: 0 | 1>
export TF_VAR_bastion_subnet=<public subnet to deploy bastion into>
export TF_VAR_vpc_id=<vpc id to deploy infra into>

```

### Build, Upload, Deploy

With the bucket created and region chosen we can deploy the whole thing by running the following command from the directory.

```
source <env file>
make build deploy
```

## Testing Lambdas

### Create Access Token

First, we need to submit and event to the `user_get_token` lambda. This can be done a number of ways but most easily through the UI using a Lambda test event.

```json
{
    "account": "test-account",
    "prefix": "*",
    "type": "local",
    "ttl": 8640000
}
```

This will generate a JWT that will contain the entitlements for out use of the PKI system.

The prefix is a glob matching pattern applied to all access calls. 

Further docs on the glob library can be found [here](https://github.com/gobwas/glob).

### Create CA

Next, we can create a CA using the token generated above.

```
{
  "token": <TOKEN_GOES_HERE>,
  "account": "test-account",
  "ca_name": "test-ca"
}
```

### Create Intermediate

Now we can create another token which only allows for a specific intermediate to be generated and no child certificates.

```json
{
    "account": "test-account",
    "prefix": "test-ca/test-intermediate1",
    "type": "local",
    "ttl": 8640000
}
```

With this new token, we can now pass the following json event to the intermediate lambda and run.

```json
{
  "token": <INTERMEDIATE_TOKEN_GOES_HERE>,
  "account": "test-account",
  "ca_name": "test-ca",
  "intermediate_name": "test-intermediate1"
}
```

Or, from another Intermediate:

```json
{
  "token": <INTERMEDIATE_TOKEN_GIES_HERE>,
  "account": "test-account",
  "ca_name": "test-ca",
  "intermediate_name": "test-intermediate2"
}
```

And we can see this will deny access due to the token lacking the entitlement to create anything other than an intermediacte certificate authority at `test-ca/test-intermediate1`.


### Create Certificate

In order to create a certificate, we need to get another token that includes the ability to create certificates from a given chain.

```json
{
    "account": "test-account",
    "prefix": "test-ca/test-intermediate1/test-cert1",
    "type": "local",
    "ttl": 8640000
}
```

And with this token we can now create a cert for the chain we have just created.

```json
{
  "token": <TOKEN_GOES_HERE>,
  "csr": <base64 CSR>,
  "account": "test-account",
  "intermediate_chain": "test-ca/test-intermediate1",
  "cert_name": "test-cert1"
}
```