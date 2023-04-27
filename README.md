# (go-)AutoGPT

A golang library and binary to create AutoGPT style agents. The default agent is powered by GPT3.5-Turbo and combines two sub-agents: one to create a plan, and one to execute it.

In the current implementation, the agents have access to a series of commands to interact with the local filesystem or even run shell commands. As such, it is recommended to run it within the provided [developer container](https://code.visualstudio.com/docs/devcontainers/containers).

## Example run
```bash
go run ./cli run --apiKey <OPENAI_APIKEY> --task "your task explained in natural language"