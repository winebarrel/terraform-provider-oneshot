# terraform-provider-oneshot

[![CI](https://github.com/winebarrel/terraform-provider-oneshot/actions/workflows/ci.yml/badge.svg)](https://github.com/winebarrel/terraform-provider-oneshot/actions/workflows/ci.yml)
[![terraform docs](https://img.shields.io/badge/terraform-docs-%35835CC?logo=terraform)](https://registry.terraform.io/providers/winebarrel/oneshot/latest/docs)

Terraform provider for running one-shot commands.

## Usage

```tf
terraform {
  required_providers {
    oneshot = {
      source  = "winebarrel/oneshot"
      version = ">= 0.4.0"
    }
  }
}

provider "oneshot" {
  # shell = "/bin/bash -c"
}

resource "oneshot_run" "hello" {
  command = "echo 'hello, oneshot'"
  # stdout_log = "stdout.log"
  # stderr_log = "stderr.log"

  # NOTE: "plan_command" is executed at plan time
  plan_command = "echo \"hello, oneshot (plan=$ONESHOT_PLAN)\""
  # plan_stdout_log = "stdout.log"
  # plan_stderr_log = "stderr.log"
}
```

## Run locally for development

```sh
cp oneshot.tf.sample oneshot.tf
make
make tf-plan
cat plan-stdout.log
cat plan-stderr.log
make tf-apply
cat stdout.log
cat stderr.log
```
