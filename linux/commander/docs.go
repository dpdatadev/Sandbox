// Single file view of all code, docs, examples, and dependencies; to later be split into production module/package(s)
package main //package commander
//https://go.dev/blog/using-go-modules
import (
	_ "github.com/mattn/go-sqlite3"
)

/*

CMDHUB
https://chatgpt.com/c/698c0190-d8ec-832d-8aee-537b6c64320d


Core Values

Command IDENTITY

Execution lifecycle

Streaming

Persistence

Lineage

Stores

Auditability


//March 1st
//Big changes are coming in Alpha 2 and Beta.
//The focus right now is on usage. I'm going to build.
//Not theorize - not learn new language features.
//But the focus now is building things I want this computer to do.
//And focus on the auditing/tracing/persistance part of the framework.
//We will likely use EXECX for fluent API/improved osExec functionality.
//There's no reason to re-invent the wheel of any kind surrounding how execution works.
//I need to focus on what currently doesn't exist. And that's chaining, tracing, and auditing commands.
//Biggest plus for using EXECX is it is already complete for MacOS and Windows and I don't have to
//implement custom code for my new "Commands" to work on Linux and Windows.

// 1. The data model needs current user data
// 2. Local history + Replay via SQLITE vs Distributed history + Replay via REDIS

February 2026 - R&D

CmdHub - turning SysOps Ad-Hoc nightmares into queryable, traceable, chainable, execution policies.

Containerized testing:
https://pkg.go.dev/github.com/testcontainers/testcontainers-go
https://medium.com/tiket-com/go-with-cucumber-an-introduction-for-bdd-style-integration-testing-7aca2f2879e4

We can setup containers, configured to our liking, describe the scenarios exactly with Cucumber,
then test each of our pipeline components. SSH, HTTP, SQLITE, YAML, etc., and not have a messy server circus.

https://chatgpt.com/c/698c0190-d8ec-832d-8aee-537b6c64320d
https://tutorialedge.net/golang/executing-system-commands-with-golang/
Paying homage to https://github.com/alexellis/go-execute
Notes, 2/11 - thinking of creating a new command API/service for Linux
servers. To have diagnostic and automated activities embedded in various applications.
For use with my raspberry pi's in the beginning. Do more research on this**.
https://chatgpt.com/c/698c0190-d8ec-832d-8aee-537b6c64320d
v 1.1 - Redis/Cache Command support(?)
Chain of Command
github/dpdatadev/chain-of-command
“Composable command execution framework with persistence, security controls, and pipeline orchestration.”
REDIS BACKED, HTTP2 streaming, ACID compliant, process management and remote execution with policy enforcement. (??)
1) Command identity + audit trail

Each command becomes a first-class entity:

2) Persistence layer

Your CommandStore abstraction enables:

SQLite → local dev

Postgres → production

Redis → caching / queues

This converts ephemeral execution into:

Durable execution history

Companies building:

CI/CD tools

Remote agents

Fleet orchestration

Security audit systems

…all need this.

3) Security scrubber / policy engine

os/exec will happily run:

rm -rf /
sudo shutdown now


Your scrubber introduces:

Blocklists

Regex policies

Allowlists

Role-based execution

Now your framework becomes viable for:

SaaS agents

Remote automation

Multi-tenant systems


Yes, os/exec supports pipes…

…but only at the file descriptor level.

You’re abstracting at the semantic level:

SSH: https://github.com/appleboy/easyssh-proxy

sshCmd.
  Pipe(textCmd).
  Pipe(fileCmd).
  Pipe(httpCmd).
  Execute(ctx)

  SSH Proxies (see https://github.com/appleboy/easyssh-proxy)
       +--------+       +----------+      +-----------+
     | Laptop | <-->  | Jumphost | <--> | FooServer |
     +--------+       +----------+      +-----------+

                         OR

     +--------+       +----------+      +-----------+
     | Laptop | <-->  | Firewall | <--> | FooServer |
     +--------+       +----------+      +-----------+
     192.168.1.5       121.1.2.3         10.10.29.68


  Now pipelines can cross protocols:

Source	Destination
SSH	Local shell
Shell	File
HTTP	Parser
File	Database

This is beyond Unix pipes.

It’s execution graphs.


TextCommand   → local shell
SSHCommand    → remote shell
HTTPCommand   → REST call
FileCommand   → write/read
SQLCommand    → database


ExecChain is a composable command execution framework for Go that builds upon os/exec with persistence, security policies, and multi-protocol pipelines.

Track, audit, and chain shell, SSH, HTTP, and file commands into reproducible execution graphs — with Redis caching and database storage built in.

pipeline := cmder.NewPipeline()
//Pipeline must be done before Beta
pipeline.
    SSH(sshConfig, "journalctl", []string{"-n", "500"}).
    PipeLocal("grep", []string{"ERROR"}).
    PipeHTTPPost("https://ops.internal/logs").
    PipeFile("error_report.txt")

pipeline.Run(ctx)

Different taglines:

“Embedded command orchestration framework with remote execution agents.”

A programmable command orchestration + audit + pipeline system with multi-protocol execution (shell, SSH, HTTP, file, DB) and persistence.

If RunDeck and Ansible had a baby .. but it came out as an embeddable API for Dev teams.


This framework becomes compelling when:

Command execution is part of the product

Not just an operational concern

DevOps handles:

Deploying systems

Embedded orchestration handles:

Operating systems programmatically from within software

Different layers of the stack.

Implements the Chain of Responsibility design pattern:
https://refactoring.guru/design-patterns/chain-of-responsibility/go/example

There are no rogue commands - must be handled in the context of Execution manager,
which validates, scrubs, executes, and handles directed output and logging.
Each component hands off to the next.
Command → Scrubber → Policy Engine → Logger → Executor → Post-Processor → Store

2/15
Post-Processor (Handlers?, this could be a tie into any app or process for Data Extraction/Analysis etc.,)
Processing the data is the responsiblity of another framework or user code and doesn't belong in the framework(maybe)

Study this pattern*
*/
//TODO, test on long running commands
//I may need to take a break for a bit and review what exactly I want to use this for
//Investigating streaming to console - https://chatgpt.com/c/698c0190-d8ec-832d-8aee-537b6c64320d
//Why not just use Proctools like Goreman?
//Depends on what you're goal is:
/*
	| Capability            | Procfile Tools | SysData Ops            |
| --------------------- | -------------- | ---------------------- |
| Process orchestration | ✅              | ⚠️ (secondary)         |
| Multi-service startup | ✅              | Possible               |
| Log aggregation       | Basic          | Structured + queryable |
| Artifact persistence  | ❌              | ✅                      |
| Schema enforcement    | ❌              | ✅                      |
| Command lineage       | ❌              | ✅                      |
| Replay execution      | ❌              | ✅                      |
| Deterministic runs    | ❌              | ✅                      |
| Execution DAGs        | ❌              | ✅                      |
| Data cataloging       | ❌              | ✅                      |
| Observability depth   | Low            | High                   |

*/
//https://www.notion.so/Features-and-Thoughts-30f568a8fe2380ca991ee486750bd3b3

//2-27
//The key features that I'm not seeing elsewhere are:
//Logging and audit trail of commands ran with lineage and persistence with security checks (scrubbing)
//Running commands remotely or in a container - we have existing examples of

/*
Current Architecture as of Early Alpha:
				 				 __api/REST  JSON/Web UI
Runner (SSH/HTTP/Process) <--> CmdLine Tool / RPC Call --> 	 ^
	>	Service --> (Execute/Log/Persist/Track)--> Output/StdOut        ^^
						--
				SQLITE		---	   RDBMS/IOFile
				DUCKDB		---------> Audit (Lineage)--->SQL-->Reports
*/

//March, next steps
//We have our focus - we are going to focus on the auditing/persistance layer, not re-invent the wheel for
//command running. We are forking the execx api (https://github.com/goforj/execx)
//We are building a cmdline tool - "hub" which registers commands and wraps them with lineage tracking and
//persistence.
//(March 7th) - Next step is swap out default Command implementation for execx then test the runner and services.
