package chlib

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/pflag"
)

func Prompt(logger *jww.Feedback, prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	log.Printf("%s: ", prompt)
	ret, _ := reader.ReadString('\n')
	return strings.TrimRight(ret, "\n")
}

func ValidationErrorExit(logger *jww.Feedback, format string, args ...interface{}) {
	logger.Printf(format, args)
	os.Exit(1)
}

const (
	cpuRegex = `^\d+(.\d+)?m?$`
	memRegex = `^\d+(.\d+)?(Mi|Gi)?$`
)

func imageValidate(logger *jww.Feedback, image string) {
	if image == "" {
		log.Println("Image must be specified")
		os.Exit(1)
	}
}

func portsValidateStr(logger *jww.Feedback, portsStr string) (ports []int) {
	for _, portStr := range strings.Split(portsStr, " ") {
		port, err := strconv.Atoi(portStr)
		if err != nil || port <= 0 || port > 65535 {
			ValidationErrorExit(logger, "Invalid port found: %s\n", portsStr)
		}
		ports = append(ports, port)
	}
	return
}

func portsValidateInt(logger *jww.Feedback, ports []int) {
	for _, port := range ports {
		if port <= 0 || port > 65535 {
			ValidationErrorExit(logger, "Invalid port found: %d\n", port)
		}
	}
}

func labelsValidate(logger *jww.Feedback, labelsStr []string) (ret map[string]string) {
	ret = make(map[string]string)
	for _, labelStr := range labelsStr {
		label := strings.Split(labelStr, "=")
		if len(label) != 2 {
			ValidationErrorExit(logger, "Invalid environment variable found: %s\n", labelStr)
		}
		ret[label[0]] = label[1]
	}
	return
}

func envVarsValidate(logger *jww.Feedback, envVarsStr []string) (env []EnvVar) {
	for _, envVarStr := range envVarsStr {
		envVar := strings.Split(envVarStr, "=")
		if len(envVar) != 2 {
			ValidationErrorExit(logger, "Invalid environment variable found: %s\n", envVarsStr)
		}
		env = append(env, EnvVar{
			Name:  envVar[0],
			Value: envVar[1],
		})
	}
	return
}

func cpuValidate(logger *jww.Feedback, cpuStr string) {
	if !regexp.MustCompile(cpuRegex).MatchString(cpuStr) {
		ValidationErrorExit(logger, "Invalid CPU cores number: %s\n", cpuStr)
	}
}

func memValidate(logger *jww.Feedback, memStr string) {
	if !regexp.MustCompile(memRegex).MatchString(memStr) {
		ValidationErrorExit(logger, "Invalid memory size: %s\n", memStr)
	}
}

func replicasValidate(logger *jww.Feedback, replicasStr string) int {
	ret, err := strconv.Atoi(replicasStr)
	if err != nil || ret <= 0 {
		ValidationErrorExit(logger, "Invalid replicas count")
	}
	return ret
}

func PromptParams(logger *jww.Feedback) (params ConfigureParams) {
	params.Image = Prompt(logger, "Enter image")
	imageValidate(logger, params.Image)
	if portsStr := Prompt(logger, "Enter ports (8080 ... 4556)"); portsStr != "" {
		params.Ports = portsValidateStr(logger, portsStr)
	}
	if labelsStr := Prompt(logger, "Enter labels (key1=value1 ... keyN=valueN)"); labelsStr != "" {
		params.Labels = labelsValidate(logger, strings.Split(labelsStr, "="))
	}
	if commands := Prompt(logger, "Enter commands (command1 ... commandN)"); commands != "" {
		params.Command = strings.Split(commands, " ")
	}
	if envVarsStr := Prompt(logger, "Enter environment variables (key1=value1 ... keyN=valueN)"); envVarsStr != "" {
		params.Env = envVarsValidate(logger, strings.Split(envVarsStr, " "))
	}
	if cpu := Prompt(logger, "Enter CPU cores (*m)"); cpu != "" {
		cpuValidate(logger, cpu)
		params.CPU = cpu
	}
	if memory := Prompt(logger, "Enter memory size (*Mi | *Gi)"); memory != "" {
		memValidate(logger, memory)
		params.Memory = memory
	}
	if replicas := Prompt(logger, "Enter replicas count"); replicas != "" {
		params.Replicas = replicasValidate(logger, replicas)
	}
	return
}

func ParamsFromArgs(logger *jww.Feedback, flags *pflag.FlagSet) (params ConfigureParams) {
	chkErr := func(err error) {
		if err != nil {
			ValidationErrorExit(logger, "flag get error: %s\n", err)
		}
	}
	var err error
	if flags.Changed("image") {
		params.Image, err = flags.GetString("image")
		chkErr(err)
		imageValidate(logger, params.Image)
	}
	if flags.Changed("port") {
		params.Ports, err = flags.GetIntSlice("port")
		chkErr(err)
		portsValidateInt(logger, params.Ports)
	}
	if flags.Changed("labels") {
		labelsSlice, err := flags.GetStringSlice("labels")
		chkErr(err)
		params.Labels = labelsValidate(logger, labelsSlice)
	}
	if flags.Changed("command") {
		params.Command, err = flags.GetStringSlice("command")
		chkErr(err)
	}
	if flags.Changed("env") {
		envSlice, err := flags.GetStringSlice("env")
		chkErr(err)
		params.Env = envVarsValidate(logger, envSlice)
	}
	if flags.Changed("cpu") {
		params.CPU, err = flags.GetString("cpu")
		chkErr(err)
		cpuValidate(logger, params.CPU)
	}
	if flags.Changed("memory") {
		params.Memory, err = flags.GetString("memory")
		chkErr(err)
		memValidate(logger, params.Memory)
	}
	if flags.Changed("replicas") {
		params.Replicas, err = flags.GetInt("replicas")
		chkErr(err)
	}
	return
}
