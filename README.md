## Running Make

### Prerequisite Steps

Have active AWS credentials in your current terminal session.

Choose a *Region* to deploy SimplePKI into.

Make a *Bucket* to track terrafrom state and to push the lambda files to. This can be done by running `make bucket-up bucket=<name> region=<region>`.

### Build, Upload, Deploy

With the bucket created and region chosen we can deploy the whole thing by running the following command from the `deploy` directory.

```
make build upload deploy bucket=<bucket> region=<region>
```



