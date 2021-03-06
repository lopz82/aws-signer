# AWS Signature

## Introduction
AWS Signature was intended to be a Traefik plugin to act as middleware authenticating all outgoing requests using the [AWS Signature Version 4 signing process](https://docs.aws.amazon.com/general/latest/gr/signature-version-4.html).

## Actual status
It seems like there are incompatibilities between [Yaegi](https://github.com/traefik/yaegi) and the AWS Go SDK, returning the following error:
```
2020/11/21 07:29:28 traefik.go:76: command traefik error: github.com/lopz82/aws-signature: failed to import plugin code "github.com/lopz82/aws-signature": 1:21: import "github.com/lopz82/aws-signature" error: /plugins/go/src/github.com/lopz82/aws-signature/signer.go:7:2: import "github.com/aws/aws-sdk-go/aws/credentials" error: /plugins/go/src/github.com/lopz82/aws-signature/vendor/github.com/aws/aws-sdk-go/aws/credentials/shared_credentials_provider.go:8:2: import "github.com/aws/aws-sdk-go/internal/ini" error: /plugins/go/src/github.com/lopz82/aws-signature/vendor/github.com/aws/aws-sdk-go/internal/ini/ast.go:74:14: cannot use type  as type []github.com/aws/aws-sdk-go/internal/ini.AST in struct literal
```

## Installation
The plugin must not be compiled in order to be used by Traefik. Check [Traefik plugins dev docs](https://doc.traefik.io/traefik-pilot/plugins/plugin-dev/) for more information.

## Testing
To test the signing process you will have to provide real AWS credentials. You can do it setting `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` environment variables.
```bash
export AWS_ACCESS_KEY_ID='your access key'
export AWS_SECRET_ACCESS_KEY='your secret'
```

Clone and test:
```bash
git clone https://github.com/lopz82/aws-signer.git
go test -v
```
