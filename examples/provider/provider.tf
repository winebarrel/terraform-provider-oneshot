provider "oneshot" {
  # shell = "/bin/bash -c"
}

resource "oneshot_run" "hello" {
  command = "echo 'hello, oneshot'"
  # stdout_log = "stdout.log"
  # stderr_log = "stderr.log"

  # NOTE: "plan_command" is executed at plan time
  plan_command = "echo \"hello, oneshot (plan=$ONESHOT_PLAN)\""
  # plan_stdout_log = "plan-stdout.log"
  # plan_stderr_log = "plan-stderr.log"
}
