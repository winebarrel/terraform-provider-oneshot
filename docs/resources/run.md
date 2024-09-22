---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "oneshot_run Resource - oneshot"
subcategory: ""
description: |-
  
---

# oneshot_run (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `command` (String) Command to execute

### Optional

- `plan_command` (String) Command to plan.
- `plan_stderr_log` (String) Stderr log file of the plan command.
- `plan_stdout_log` (String) Stdout log file of the plan command.
- `shell` (String) Shell to execute the command.

### Read-Only

- `run_at` (String) Command execution time.
- `stderr` (String) Stderr of the command.
- `stdout` (String) Stdout of the command.