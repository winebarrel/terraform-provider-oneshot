package provider_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

func TestRun_Basic(t *testing.T) {
	assert := assert.New(t)

	cwd, _ := os.Getwd()
	os.Chdir(t.TempDir())
	defer os.Chdir(cwd)

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
					resource.TestMatchResourceAttr("oneshot_run.hello", "run_at", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d+`)),
				),
			},
		},
	})

	stdout, _ := os.ReadFile("plan-stdout.log")
	assert.Equal("plan\n", string(stdout))
	stderr, _ := os.ReadFile("plan-stderr.log")
	assert.Equal("planerr\n", string(stderr))
}

func TestRun_WithoutPlanCommand(t *testing.T) {
	assert := assert.New(t)

	cwd, _ := os.Getwd()
	os.Chdir(t.TempDir())
	defer os.Chdir(cwd)

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "oneshot_run" "hello" {
						command = "echo hello ; echo world 1>&2"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("oneshot_run.hello", "command", "echo hello ; echo world 1>&2"),
					resource.TestCheckNoResourceAttr("oneshot_run.hello", "plan_command"),
					resource.TestCheckNoResourceAttr("oneshot_run.hello", "shell"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "stdout", "hello\n"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "stderr", "world\n"),
					resource.TestMatchResourceAttr("oneshot_run.hello", "run_at", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d+`)),
				),
			},
		},
	})

	// No log
	_, err := os.Stat("plan-stdout.log")
	assert.Error(err)
	_, err = os.Stat("plan-stderr.log")
	assert.Error(err)
}

func TestRun_WithShell(t *testing.T) {
	assert := assert.New(t)

	cwd, _ := os.Getwd()
	os.Chdir(t.TempDir())
	defer os.Chdir(cwd)

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "oneshot_run" "hello" {
						command      = "echo $0 ; echo world 1>&2"
						plan_command = "echo plan ; echo $0 1>&2"
						shell        = "/bin/sh -c"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("oneshot_run.hello", "command", "echo $0 ; echo world 1>&2"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "plan_command", "echo plan ; echo $0 1>&2"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "shell", "/bin/sh -c"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "stdout", "/bin/sh\n"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "stderr", "world\n"),
					resource.TestMatchResourceAttr("oneshot_run.hello", "run_at", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d+`)),
				),
			},
		},
	})

	stdout, _ := os.ReadFile("plan-stdout.log")
	assert.Equal("plan\n", string(stdout))
	stderr, _ := os.ReadFile("plan-stderr.log")
	assert.Equal("/bin/sh\n", string(stderr))
}
