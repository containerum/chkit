package chlib

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
)

type printer interface {
	Print(v ...interface{})
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}

func Prompt(printer printer, prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	printer.Printf("%s: ", prompt)
	ret, _ := reader.ReadString('\n')
	return strings.TrimRight(ret, "\n")
}

func ValidationErrorExit(printer printer, format string, args ...interface{}) {
	printer.Printf(format, args)
	os.Exit(1)
}

const (
	cpuRegex = `^\d+(.\d+)?m?$`
	memRegex = `^\d+(.\d+)?(Mi|Gi)?$`
)

func imageValidate(printer printer, image string) {
	if image == "" {
		printer.Println("Image must be specified")
		os.Exit(1)
	}
}

func portsValidateStr(printer printer, portsStr string) (ports []int) {
	for _, portStr := range strings.Split(portsStr, " ") {
		port, err := strconv.Atoi(portStr)
		if err != nil || port <= 0 || port > 65535 {
			ValidationErrorExit(printer, "Invalid port found: %s\n", portsStr)
		}
		ports = append(ports, port)
	}
	return
}

func portsValidateInt(printer printer, ports []int) {
	for _, port := range ports {
		if port <= 0 || port > 65535 {
			ValidationErrorExit(printer, "Invalid port found: %d\n", port)
		}
	}
}

func labelsValidate(printer printer, labelsStr []string) (ret map[string]string) {
	ret = make(map[string]string)
	for _, labelStr := range labelsStr {
		label := strings.Split(labelStr, "=")
		if len(label) != 2 {
			ValidationErrorExit(printer, "Invalid environment variable found: %s\n", labelStr)
		}
		ret[label[0]] = label[1]
	}
	return
}

func envVarsValidate(printer printer, envVarsStr []string) (env []EnvVar) {
	for _, envVarStr := range envVarsStr {
		envVar := strings.Split(envVarStr, "=")
		if len(envVar) != 2 {
			ValidationErrorExit(printer, "Invalid environment variable found: %s\n", envVarsStr)
		}
		env = append(env, EnvVar{
			Name:  envVar[0],
			Value: envVar[1],
		})
	}
	return
}

func cpuValidate(printer printer, cpuStr string) {
	if !regexp.MustCompile(cpuRegex).MatchString(cpuStr) {
		ValidationErrorExit(printer, "Invalid CPU cores number: %s\n", cpuStr)
	}
}

func memValidate(printer printer, memStr string) {
	if !regexp.MustCompile(memRegex).MatchString(memStr) {
		ValidationErrorExit(printer, "Invalid memory size: %s\n", memStr)
	}
}

func replicasValidate(printer printer, replicasStr string) int {
	ret, err := strconv.Atoi(replicasStr)
	if err != nil || ret <= 0 {
		ValidationErrorExit(printer, "Invalid replicas count")
	}
	return ret
}

func PromptParams(printer printer) (params ConfigureParams) {
	params.Image = Prompt(printer, "Enter image")
	imageValidate(printer, params.Image)
	if portsStr := Prompt(printer, "Enter ports (8080 ... 4556)"); portsStr != "" {
		params.Ports = portsValidateStr(printer, portsStr)
	}
	if labelsStr := Prompt(printer, "Enter labels (key1=value1 ... keyN=valueN)"); labelsStr != "" {
		params.Labels = labelsValidate(printer, strings.Split(labelsStr, "="))
	}
	if commands := Prompt(printer, "Enter commands (command1 ... commandN)"); commands != "" {
		params.Command = strings.Split(commands, " ")
	}
	if envVarsStr := Prompt(printer, "Enter environment variables (key1=value1 ... keyN=valueN)"); envVarsStr != "" {
		params.Env = envVarsValidate(printer, strings.Split(envVarsStr, " "))
	}
	if cpu := Prompt(printer, "Enter CPU cores (*m)"); cpu != "" {
		cpuValidate(printer, cpu)
		params.CPU = cpu
	}
	if memory := Prompt(printer, "Enter memory size (*Mi | *Gi)"); memory != "" {
		memValidate(printer, memory)
		params.Memory = memory
	}
	if replicas := Prompt(printer, "Enter replicas count"); replicas != "" {
		params.Replicas = replicasValidate(printer, replicas)
	}
	return
}

func ParamsFromArgs(printer printer, flags *pflag.FlagSet) (params ConfigureParams) {
	chkErr := func(err error) {
		if err != nil {
			ValidationErrorExit(printer, "flag get error: %s\n", err)
		}
	}
	var err error
	if flags.Changed("image") {
		params.Image, err = flags.GetString("image")
		chkErr(err)
		imageValidate(printer, params.Image)
	}
	if flags.Changed("port") {
		params.Ports, err = flags.GetIntSlice("port")
		chkErr(err)
		portsValidateInt(printer, params.Ports)
	}
	if flags.Changed("labels") {
		labelsSlice, err := flags.GetStringSlice("labels")
		chkErr(err)
		params.Labels = labelsValidate(printer, labelsSlice)
	}
	if flags.Changed("command") {
		params.Command, err = flags.GetStringSlice("command")
		chkErr(err)
	}
	if flags.Changed("env") {
		envSlice, err := flags.GetStringSlice("env")
		chkErr(err)
		params.Env = envVarsValidate(printer, envSlice)
	}
	if flags.Changed("cpu") {
		params.CPU, err = flags.GetString("cpu")
		chkErr(err)
		cpuValidate(printer, params.CPU)
	}
	if flags.Changed("memory") {
		params.Memory, err = flags.GetString("memory")
		chkErr(err)
		memValidate(printer, params.Memory)
	}
	if flags.Changed("replicas") {
		params.Replicas, err = flags.GetInt("replicas")
		chkErr(err)
	}
	return
}
