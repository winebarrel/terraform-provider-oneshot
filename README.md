# terraform-provider-oneshot

[![CI](https://github.com/winebarrel/terraform-provider-oneshot/actions/workflows/ci.yml/badge.svg)](https://github.com/winebarrel/terraform-provider-oneshot/actions/workflows/ci.yml)

Terraform provider for running one-shot commands.

## Usage

```tf
provider "oneshot" {
  # shell = "/bin/bash -c"
}

resource "oneshot_run" "hello" {
  command = "echo 'hello, oneshot'"
  # NOTE: "plan_command" is executed at plan time
  plan_command = "echo \"hello, oneshot (plan=$ONESHOT_PLAN)\""
  # plan_stdout_log = "plan-stdout.log"
  # plan_stderr_log = "plan-stderr.log"
}

resource "terraform_data" "hello_stdout" {
  triggers_replace = oneshot_run.hello.run_at

  provisioner "local-exec" {
    command = "echo '${oneshot_run.hello.stdout}'"
  }
}

resource "terraform_data" "hello_stderr" {
  triggers_replace = oneshot_run.hello.run_at

  provisioner "local-exec" {
    command = "echo '${oneshot_run.hello.stderr}'"
  }
}
```

## Run locally for development

```sh
cp oneshot.tf.sample oneshot.tf
make
make tf-plan
make tf-apply
cat plan-stdout.log
cat plan-stderr.log
```
