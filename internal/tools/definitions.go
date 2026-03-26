package tools

import "github.com/kocar/aurelia/internal/agent"

func ReadFileDefinition() agent.Tool {
	return agent.Tool{
		Name:        "read_file",
		Description: "Le o conteudo de um arquivo local. Se `workdir` for informado, caminhos relativos serao resolvidos a partir dele.",
		JSONSchema: objectSchema(
			map[string]any{
				"path":    stringProperty(""),
				"workdir": stringProperty("Diretorio base opcional para resolver caminhos relativos."),
			},
			"path",
		),
	}
}

func WriteFileDefinition() agent.Tool {
	return agent.Tool{
		Name:        "write_file",
		Description: "Escreve conteudo integral em um arquivo. Se `workdir` for informado, caminhos relativos serao resolvidos a partir dele.",
		JSONSchema: objectSchema(
			map[string]any{
				"path":    stringProperty(""),
				"content": stringProperty(""),
				"workdir": stringProperty("Diretorio base opcional para resolver caminhos relativos."),
			},
			"path",
			"content",
		),
	}
}

func ListDirDefinition() agent.Tool {
	return agent.Tool{
		Name:        "list_dir",
		Description: "Lista os arquivos dentro de um diretorio. Se `workdir` for informado, caminhos relativos serao resolvidos a partir dele.",
		JSONSchema: objectSchema(
			map[string]any{
				"path":    stringProperty(""),
				"workdir": stringProperty("Diretorio base opcional para resolver caminhos relativos."),
			},
			"path",
		),
	}
}

func WebSearchDefinition() agent.Tool {
	return agent.Tool{
		Name:        "web_search",
		Description: "Pesquisa na internet usando a engine do DuckDuckGo e extrai resultados textuais.",
		JSONSchema: objectSchema(
			map[string]any{
				"query": stringProperty(""),
				"count": numberProperty("Maximo de resultados (ate 10)"),
			},
			"query",
		),
	}
}

func RunCommandDefinition() agent.Tool {
	return agent.Tool{
		Name:        "run_command",
		Description: "Executa um comando local (Bash) de forma controlada no Ubuntu 24.04 e retorna stdout, stderr, exit code e timeout em JSON. Use sintaxe nativa Linux e prefira informar `workdir` ao operar em outro projeto.",
		JSONSchema: objectSchema(
			map[string]any{
				"command":         stringProperty(""),
				"workdir":         stringProperty(""),
				"timeout_seconds": numberProperty(""),
			},
			"command",
		),
	}
}

func DockerControlDefinition() agent.Tool {
	return agent.Tool{
		Name:        "docker_control",
		Description: "Controla Docker: listar containers (ps), reiniciar (restart), logs, stats, ou docker-compose up",
		JSONSchema: objectSchema(
			map[string]any{
				"action":    stringProperty("ps, restart, logs, stats, ou compose_up"),
				"container": stringProperty("Nome ou ID do container (para restart/logs)"),
				"workdir":   stringProperty("Diretório para docker-compose up"),
			},
			"action",
		),
	}
}

func SystemMonitorDefinition() agent.Tool {
	return agent.Tool{
		Name:        "system_monitor",
		Description: "Monitora sistema: stats (CPU/RAM), gpu (NVIDIA), process (top 10), network (interfaces e conexões)",
		JSONSchema: objectSchema(
			map[string]any{
				"metric": stringProperty("stats, gpu, process, ou network"),
			},
			"metric",
		),
	}
}

func ServiceControlDefinition() agent.Tool {
	return agent.Tool{
		Name:        "service_control",
		Description: "Controla serviços systemd: status, restart, stop, start, list, ou logs",
		JSONSchema: objectSchema(
			map[string]any{
				"action":  stringProperty("status, restart, stop, start, list, ou logs"),
				"service": stringProperty("Nome do serviço (requerido exceto para list)"),
			},
			"action",
		),
	}
}

func OllamaControlDefinition() agent.Tool {
	return agent.Tool{
		Name:        "ollama_control",
		Description: "Controla Ollama (inferência local): status, list (modelos), pull (baixar modelo), run (inferência)",
		JSONSchema: objectSchema(
			map[string]any{
				"action": stringProperty("status, list, pull, ou run"),
				"model":  stringProperty("Nome do modelo (para pull/run)"),
				"prompt": stringProperty("Prompt de entrada (para run)"),
			},
			"action",
		),
	}
}

func CPFCNPJDefinition() agent.Tool {
	return agent.Tool{
		Name:        "cpf_cnpj",
		Description: "Valida CPF ou CNPJ (algoritmo) e consulta dados de empresa pelo CNPJ via BrasilAPI (gratuito, sem autenticação). Ações: validate_cpf, validate_cnpj, lookup_cnpj.",
		JSONSchema: objectSchema(
			map[string]any{
				"action": stringProperty("validate_cpf | validate_cnpj | lookup_cnpj"),
				"number": stringProperty("CPF (11 dígitos) ou CNPJ (14 dígitos) — aceita formatado ou só números"),
			},
			"action",
			"number",
		),
	}
}

func RegisterCoreTools(registry *agent.ToolRegistry) {
	if registry == nil {
		return
	}

	registry.Register(ReadFileDefinition(), ReadFileHandler)
	registry.Register(WriteFileDefinition(), WriteFileHandler)
	registry.Register(ListDirDefinition(), ListDirHandler)
	registry.Register(WebSearchDefinition(), WebSearchHandler)
	registry.Register(RunCommandDefinition(), RunCommandHandler)
	registry.Register(DockerControlDefinition(), DockerControlHandler)
	registry.Register(SystemMonitorDefinition(), SystemMonitorHandler)
	registry.Register(ServiceControlDefinition(), ServiceControlHandler)
	registry.Register(OllamaControlDefinition(), OllamaControlHandler)
	registry.Register(CPFCNPJDefinition(), CPFCNPJHandler)

	setPhase := NewSetPhaseTool()
	registry.Register(setPhase.Definition(), setPhase.Execute)
}

func objectSchema(properties map[string]any, required ...string) map[string]any {
	schema := map[string]any{
		"type":       "object",
		"properties": properties,
	}
	if len(required) > 0 {
		schema["required"] = required
	}
	return schema
}

func stringProperty(description string) map[string]any {
	property := map[string]any{"type": "string"}
	if description != "" {
		property["description"] = description
	}
	return property
}

func numberProperty(description string) map[string]any {
	property := map[string]any{"type": "number"}
	if description != "" {
		property["description"] = description
	}
	return property
}
