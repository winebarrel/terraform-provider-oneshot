package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestRun_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "oneshot_run" "hello" {
						command      = "echo hello ; echo world 1>&2"
						plan_command = "echo plan ; echo planerr 1>&2"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("oneshot_run.hello", "command", "echo hello ; echo world 1>&2"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "plan_command", "echo plan ; echo planerr 1>&2"),
					resource.TestCheckNoResourceAttr("oneshot_run.hello", "shell"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "stdout", "hello\n"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "stderr", "world\n"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "plan_stdout", "plan\n"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "plan_stderr", "planerr\n"),
					resource.TestMatchResourceAttr("oneshot_run.hello", "run_at", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d+`)),
				),
			},
		},
	})
}

func TestRun_WithShell(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "oneshot_run" "hello" {
						command      = "echo hello ; echo world 1>&2"
						plan_command = "echo plan ; echo planerr 1>&2"
						shell        = "/bin/sh -c"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("oneshot_run.hello", "command", "echo hello ; echo world 1>&2"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "plan_command", "echo plan ; echo planerr 1>&2"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "shell", "/bin/sh -c"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "stdout", "hello\n"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "stderr", "world\n"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "plan_stdout", "plan\n"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "plan_stderr", "planerr\n"),
					resource.TestMatchResourceAttr("oneshot_run.hello", "run_at", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d+`)),
				),
			},
		},
	})
}
