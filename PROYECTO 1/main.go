package main

// main.go: punto de entrada del programa

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"MIA_P1_202400452/commands"
	"MIA_P1_202400452/parser"
	"MIA_P1_202400452/reports"
	"MIA_P1_202400452/storage"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("========================================")
	fmt.Println(" MIA - Manejo e Implementación de Archivos")
	fmt.Println(" Sistema de Archivos EXT2 - Proyecto 1")
	fmt.Println("========================================")
	fmt.Println()
	fmt.Println("Comandos disponibles: mkdisk, rmdisk, fdisk, mount, unmount, mounted, mkfs,")
	fmt.Println(" login, logout, mkgrp, rmgrp, mkusr, rmusr, chgrp,")
	fmt.Println(" mkfile, mkdir, cat, rep, execute, pause")
	fmt.Println()
	fmt.Println("Escriba 'exit' para salir.")
	fmt.Println()

	for {
		fmt.Print("> ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error leyendo entrada: %v\n", err)
			continue
		}

		input = strings.TrimSpace(input)

		if input == "" {
			continue
		}

		if strings.ToLower(input) == "exit" {
			fmt.Println("Saliendo...")
			break
		}

		result := processCommand(input)

		if result != "" {
			fmt.Println(result)
		}
	}
}

var hostPathCommands = map[string][]string{
	"mkdisk":  {"path"},
	"rmdisk":  {"path"},
	"fdisk":   {"path"},
	"mount":   {"path"},
	"execute": {"path"},
	"rep":     {"path"},
	"mkfile":  {"cont"},
}

var allowedCommandParams = map[string]map[string]bool{
	"mkdisk": {
		"size": true,
		"fit":  true,
		"unit": true,
		"path": true,
	},
	"rmdisk": {
		"path": true,
	},
	"fdisk": {
		"size": true,
		"unit": true,
		"path": true,
		"type": true,
		"fit":  true,
		"name": true,
	},
	"mount": {
		"path": true,
		"name": true,
	},
	"unmount": {
		"id": true,
	},
	"mounted": {},
	"mkfs": {
		"id":   true,
		"type": true,
	},
	"login": {
		"user": true,
		"usr":  true,
		"pass": true,
		"id":   true,
	},
	"logout": {},
	"mkgrp": {
		"name": true,
	},
	"rmgrp": {
		"name": true,
	},
	"mkusr": {
		"user": true,
		"usr":  true,
		"pass": true,
		"grp":  true,
	},
	"rmusr": {
		"user": true,
		"usr":  true,
	},
	"chgrp": {
		"user": true,
		"usr":  true,
		"grp":  true,
	},
	"mkfile": {
		"path": true,
		"r":    true,
		"size": true,
		"s":    true,
		"cont": true,
	},
	"mkdir": {
		"path": true,
		"p":    true,
	},
	"cat": {
		"file": true,
	},
	"rep": {
		"name":         true,
		"path":         true,
		"id":           true,
		"path_file_ls": true,
	},
	"execute": {
		"path": true,
	},
	"pause": {},
}

func processCommand(input string) string {
	trimmed := strings.TrimSpace(input)

	if trimmed == "" {
		return ""
	}

	if strings.HasPrefix(trimmed, "#") {
		return trimmed
	}

	cmd := parser.ParseLine(trimmed)
	if cmd == nil {
		return ""
	}

	if cmd.Name == "" {
		return ""
	}

	cmd.Name = strings.ToLower(strings.TrimSpace(cmd.Name))

	normalizeParams(cmd.Name, cmd.Params)

	if err := validateCommandParams(cmd.Name, cmd.Params); err != nil {
		return err.Error()
	}

	if paramKeys, ok := hostPathCommands[cmd.Name]; ok {
		for _, key := range paramKeys {
			if val, ok := cmd.Params[key]; ok && val != "" {
				cmd.Params[key] = storage.ResolvePath(val)
			}
		}
	}

	switch cmd.Name {
	case "mkdisk":
		return commands.CmdMKDISK(cmd.Params)

	case "rmdisk":
		return commands.CmdRMDISK(cmd.Params)

	case "fdisk":
		return commands.CmdFDISK(cmd.Params)

	case "mount":
		return commands.CmdMOUNT(cmd.Params)

	case "unmount":
		return commands.CmdUNMOUNT(cmd.Params)

	case "mounted":
		return commands.CmdMOUNTED()

	case "mkfs":
		return commands.CmdMKFS(cmd.Params)

	case "login":
		return commands.CmdLOGIN(cmd.Params)

	case "logout":
		return commands.CmdLOGOUT(cmd.Params)

	case "mkgrp":
		return commands.CmdMKGRP(cmd.Params)

	case "rmgrp":
		return commands.CmdRMGRP(cmd.Params)

	case "mkusr":
		return commands.CmdMKUSR(cmd.Params)

	case "rmusr":
		return commands.CmdRMUSR(cmd.Params)

	case "chgrp":
		return commands.CmdCHGRP(cmd.Params)

	case "mkfile":
		return commands.CmdMKFILE(cmd.Params)

	case "mkdir":
		return commands.CmdMKDIR(cmd.Params)

	case "cat":
		return commands.CmdCAT(cmd.Params)

	case "rep":
		return reports.CmdREP(cmd.Params)

	case "execute":
		return executeScript(cmd.Params)

	case "pause":
		return pauseCommand()

	default:
		return fmt.Sprintf("Error: comando '%s' no reconocido", cmd.Name)
	}
}

func normalizeParams(cmdName string, params map[string]string) {
	switch cmdName {
	case "login":
		if value, ok := params["usr"]; ok {
			params["user"] = value
			delete(params, "usr")
		}

	case "mkusr", "rmusr", "chgrp":
		if value, ok := params["usr"]; ok {
			params["user"] = value
			delete(params, "usr")
		}

	case "mkfile":
		if value, ok := params["s"]; ok {
			params["size"] = value
			delete(params, "s")
		}

	case "cat":
		if value, ok := params["file"]; ok {
			params["file1"] = value
			delete(params, "file")
		}
	}
}

func validateCommandParams(cmdName string, params map[string]string) error {
	allowed, ok := allowedCommandParams[cmdName]
	if !ok {
		return nil
	}

	for key := range params {
		if cmdName == "cat" && isCatFileParam(key) {
			continue
		}

		if !allowed[key] {
			return fmt.Errorf("Error: parámetro '-%s' no reconocido para el comando '%s'", key, cmdName)
		}
	}

	return nil
}

func isCatFileParam(key string) bool {
	if key == "file" {
		return true
	}

	if !strings.HasPrefix(key, "file") {
		return false
	}

	suffix := key[4:]

	if suffix == "" {
		return false
	}

	for _, r := range suffix {
		if r < '0' || r > '9' {
			return false
		}
	}

	return true
}

func executeScript(params map[string]string) string {
	path, ok := params["path"]
	if !ok {
		return "Error: falta el parámetro obligatorio -path"
	}

	if strings.ToLower(filepath.Ext(path)) != ".smia" {
		return "Error: el archivo de script debe tener extensión .smia"
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("Error al leer el script: %v", err)
	}

	lines := parser.ParseScript(string(content))

	var output strings.Builder

	for _, originalLine := range lines {
		line := strings.TrimSpace(originalLine)

		if line == "" {
			output.WriteString("\n")
			continue
		}

		if strings.HasPrefix(line, "#") {
			output.WriteString(originalLine + "\n")
			continue
		}

		cleanLine := removeInlineComment(originalLine)
		cleanLine = strings.TrimSpace(cleanLine)

		if cleanLine == "" {
			output.WriteString(originalLine + "\n")
			continue
		}

		output.WriteString("> " + cleanLine + "\n")

		result := processCommand(cleanLine)
		if result != "" {
			output.WriteString(result + "\n")
		}
	}

	return output.String()
}

func removeInlineComment(line string) string {
	inQuotes := false

	for i, r := range line {
		if r == '"' {
			inQuotes = !inQuotes
			continue
		}

		if r == '#' && !inQuotes {
			return strings.TrimSpace(line[:i])
		}
	}

	return line
}

func pauseCommand() string {
	fmt.Print("PAUSE: presione Enter para continuar...")

	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')

	return "Continuando ejecución..."
}
