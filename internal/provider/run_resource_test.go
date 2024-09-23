package provider_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
					resource.TestCheckResourceAttr("oneshot_run.hello", "stdout_log", "stdout.log"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "stderr_log", "stderr.log"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "plan_stdout_log", "plan-stdout.log"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "plan_stderr_log", "plan-stderr.log"),
					resource.TestMatchResourceAttr("oneshot_run.hello", "run_at", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d+`)),
					func(s *terraform.State) error {
						stdout, _ := os.ReadFile("stdout.log")
						assert.Equal("hello\n", string(stdout))
						stderr, _ := os.ReadFile("stderr.log")
						assert.Equal("world\n", string(stderr))
						return nil
					},
					func(s *terraform.State) error {
						stdout, _ := os.ReadFile("plan-stdout.log")
						assert.Equal("plan\n", string(stdout))
						stderr, _ := os.ReadFile("plan-stderr.log")
						assert.Equal("planerr\n", string(stderr))
						return nil
					},
				),
			},
		},
	})
}

func TestRun_PlanEnv(t *testing.T) {
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
						command      = "echo plan=$ONESHOT_PLAN ; echo plan=$ONESHOT_PLAN 1>&2"
						plan_command = "echo plan=$ONESHOT_PLAN ; echo plan=$ONESHOT_PLAN 1>&2"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("oneshot_run.hello", "command", "echo plan=$ONESHOT_PLAN ; echo plan=$ONESHOT_PLAN 1>&2"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "plan_command", "echo plan=$ONESHOT_PLAN ; echo plan=$ONESHOT_PLAN 1>&2"),
					resource.TestCheckNoResourceAttr("oneshot_run.hello", "shell"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "stdout_log", "stdout.log"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "stderr_log", "stderr.log"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "plan_stdout_log", "plan-stdout.log"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "plan_stderr_log", "plan-stderr.log"),
					resource.TestMatchResourceAttr("oneshot_run.hello", "run_at", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d+`)),
					func(s *terraform.State) error {
						stdout, _ := os.ReadFile("stdout.log")
						assert.Equal("plan=\n", string(stdout))
						stderr, _ := os.ReadFile("stderr.log")
						assert.Equal("plan=\n", string(stderr))
						return nil
					},
					func(s *terraform.State) error {
						stdout, _ := os.ReadFile("plan-stdout.log")
						assert.Equal("plan=1\n", string(stdout))
						stderr, _ := os.ReadFile("plan-stderr.log")
						assert.Equal("plan=1\n", string(stderr))
						return nil
					},
				),
			},
		},
	})
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
					resource.TestCheckResourceAttr("oneshot_run.hello", "stdout_log", "stdout.log"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "stderr_log", "stderr.log"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "plan_stdout_log", "plan-stdout.log"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "plan_stderr_log", "plan-stderr.log"),
					resource.TestMatchResourceAttr("oneshot_run.hello", "run_at", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d+`)),
					func(s *terraform.State) error {
						stdout, _ := os.ReadFile("stdout.log")
						assert.Equal("hello\n", string(stdout))
						stderr, _ := os.ReadFile("stderr.log")
						assert.Equal("world\n", string(stderr))
						return nil
					},
					func(s *terraform.State) error {
						// No log
						_, err := os.Stat("plan-stdout.log")
						assert.Error(err)
						_, err = os.Stat("plan-stderr.log")
						assert.Error(err)
						return nil
					},
				),
			},
		},
	})
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
					resource.TestCheckResourceAttr("oneshot_run.hello", "stdout_log", "stdout.log"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "stderr_log", "stderr.log"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "plan_stdout_log", "plan-stdout.log"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "plan_stderr_log", "plan-stderr.log"),
					resource.TestMatchResourceAttr("oneshot_run.hello", "run_at", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d+`)),
					func(s *terraform.State) error {
						stdout, _ := os.ReadFile("stdout.log")
						assert.Equal("/bin/sh\n", string(stdout))
						stderr, _ := os.ReadFile("stderr.log")
						assert.Equal("world\n", string(stderr))
						return nil
					},
					func(s *terraform.State) error {
						stdout, _ := os.ReadFile("plan-stdout.log")
						assert.Equal("plan\n", string(stdout))
						stderr, _ := os.ReadFile("plan-stderr.log")
						assert.Equal("/bin/sh\n", string(stderr))
						return nil
					},
				),
			},
		},
	})
}

func TestRun_RunPlanCommandonlyOnce(t *testing.T) {
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
					resource.TestCheckResourceAttr("oneshot_run.hello", "stdout_log", "stdout.log"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "stderr_log", "stderr.log"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "plan_stdout_log", "plan-stdout.log"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "plan_stderr_log", "plan-stderr.log"),
					resource.TestMatchResourceAttr("oneshot_run.hello", "run_at", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d+`)),
					func(s *terraform.State) error {
						stdout, _ := os.ReadFile("stdout.log")
						assert.Equal("hello\n", string(stdout))
						stderr, _ := os.ReadFile("stderr.log")
						assert.Equal("world\n", string(stderr))
						return nil
					},
					func(s *terraform.State) error {
						os.Remove("stdout.log")
						os.Remove("stderr.log")
						return nil
					},
					func(s *terraform.State) error {
						stdout, _ := os.ReadFile("plan-stdout.log")
						assert.Equal("plan\n", string(stdout))
						stderr, _ := os.ReadFile("plan-stderr.log")
						assert.Equal("planerr\n", string(stderr))
						return nil
					},
					func(s *terraform.State) error {
						os.Remove("plan-stdout.log")
						os.Remove("plan-stderr.log")
						return nil
					},
				),
			},
			{
				Config: `
					resource "oneshot_run" "hello" {
						command      = "echo hello ; echo world 1>&2"
						plan_command = "echo plan ; echo planerr 1>&2"
					}
				`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("oneshot_run.hello", "command", "echo hello ; echo world 1>&2"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "plan_command", "echo plan ; echo planerr 1>&2"),
					resource.TestCheckNoResourceAttr("oneshot_run.hello", "shell"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "stdout_log", "stdout.log"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "stderr_log", "stderr.log"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "plan_stdout_log", "plan-stdout.log"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "plan_stderr_log", "plan-stderr.log"),
					resource.TestMatchResourceAttr("oneshot_run.hello", "run_at", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d+`)),
					func(s *terraform.State) error {
						// No log
						_, err := os.Stat("stdout.log")
						assert.Error(err)
						_, err = os.Stat("stderr.log")
						assert.Error(err)
						return nil
					},
					func(s *terraform.State) error {
						// No log
						_, err := os.Stat("plan-stdout.log")
						assert.Error(err)
						_, err = os.Stat("plan-stderr.log")
						assert.Error(err)
						return nil
					},
				),
			},
		},
	})
}

func TestRun_RenameLog(t *testing.T) {
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
						command         = "echo hello ; echo world 1>&2"
						plan_command    = "echo plan ; echo planerr 1>&2"
						stdout_log      = "x-stdout.log"
						stderr_log      = "x-stderr.log"
						plan_stdout_log = "x-plan-stdout.log"
						plan_stderr_log = "x-plan-stderr.log"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("oneshot_run.hello", "command", "echo hello ; echo world 1>&2"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "plan_command", "echo plan ; echo planerr 1>&2"),
					resource.TestCheckNoResourceAttr("oneshot_run.hello", "shell"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "stdout_log", "x-stdout.log"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "stderr_log", "x-stderr.log"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "plan_stdout_log", "x-plan-stdout.log"),
					resource.TestCheckResourceAttr("oneshot_run.hello", "plan_stderr_log", "x-plan-stderr.log"),
					resource.TestMatchResourceAttr("oneshot_run.hello", "run_at", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d+`)),
					func(s *terraform.State) error {
						stdout, _ := os.ReadFile("x-stdout.log")
						assert.Equal("hello\n", string(stdout))
						stderr, _ := os.ReadFile("x-stderr.log")
						assert.Equal("world\n", string(stderr))
						return nil
					},
					func(s *terraform.State) error {
						stdout, _ := os.ReadFile("x-plan-stdout.log")
						assert.Equal("plan\n", string(stdout))
						stderr, _ := os.ReadFile("x-plan-stderr.log")
						assert.Equal("planerr\n", string(stderr))
						return nil
					},
				),
			},
		},
	})
}

func TestRun_PlanErr(t *testing.T) {
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
						command      = "echo hello"
						plan_command = "echo stdout ; echo stderr 1>&2 ; exit 111"
					}
				`,
				ExpectError: regexp.MustCompile(
					`Unable to plan command, got error: Failed to execute command: exit status 111\n\[STDOUT\] stdout\n\n\[STDERR\] stderr\n\n`,
				),
			},
		},
	})
}

func TestRun_RunErr(t *testing.T) {
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
						command = "echo stdout ; echo stderr 1>&2 ; exit 111"
					}
				`,
				ExpectError: regexp.MustCompile(
					`Unable to run command, got error: Failed to execute command: exit status 111\n\[STDOUT\] stdout\n\n\[STDERR\] stderr\n\n`,
				),
			},
		},
	})
}
