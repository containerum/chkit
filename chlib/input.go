package chlib

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"

	"fmt"

	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/pflag"
)

func Prompt(np *jww.Notepad, prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	np.FEEDBACK.Printf("%s: ", prompt)
	ret, _ := reader.ReadString('\n')
	return strings.TrimRight(ret, "\n")
}

func validationErrorExit(np *jww.Notepad, format string, args ...interface{}) {
	np.FEEDBACK.Printf(format, args)
	os.Exit(1)
}

func imageValidate(np *jww.Notepad, image string) {
	if image == "" {
		np.FEEDBACK.Println("Image must be specified")
		os.Exit(1)
	}
}

func portsValidateStr(np *jww.Notepad, portsStr string) (ports []int) {
	for _, portStr := range strings.Split(portsStr, " ") {
		port, err := strconv.Atoi(portStr)
		if err != nil || port <= 0 || port > 65535 {
			validationErrorExit(np, "Invalid port found: %s\n", portsStr)
		}
		ports = append(ports, port)
	}
	return
}

func portsValidateInt(np *jww.Notepad, ports []int) {
	for _, port := range ports {
		if port <= 0 || port > 65535 {
			validationErrorExit(np, "Invalid port found: %d\n", port)
		}
	}
}

func labelsValidate(np *jww.Notepad, labelsStr []string) (ret map[string]string) {
	ret = make(map[string]string)
	for _, labelStr := range labelsStr {
		label := strings.Split(labelStr, "=")
		labelValidator := regexp.MustCompile(LabelRegex)
		if len(label) != 2 || !labelValidator.MatchString(label[0]) || !labelValidator.MatchString(label[1]) {
			validationErrorExit(np, "Invalid label found: %s\n", labelStr)
		}
		ret[label[0]] = label[1]
	}
	return
}

func envVarsValidate(np *jww.Notepad, envVarsStr []string) (env []EnvVar) {
	for _, envVarStr := range envVarsStr {
		envVar := strings.Split(envVarStr, "=")
		if len(envVar) != 2 {
			validationErrorExit(np, "Invalid environment variable found: %s\n", envVarsStr)
		}
		env = append(env, EnvVar{
			Name:  envVar[0],
			Value: envVar[1],
		})
	}
	return
}

func cpuValidate(np *jww.Notepad, cpuStr string) {
	if !regexp.MustCompile(CpuRegex).MatchString(cpuStr) {
		validationErrorExit(np, "Invalid CPU cores number: %s\n", cpuStr)
	}
}

func memValidate(np *jww.Notepad, memStr string) {
	if !regexp.MustCompile(MemRegex).MatchString(memStr) {
		validationErrorExit(np, "Invalid memory size: %s\n", memStr)
	}
}

func replicasValidate(np *jww.Notepad, replicasStr string) int {
	ret, err := strconv.Atoi(replicasStr)
	if err != nil || ret <= 0 {
		validationErrorExit(np, "Invalid replicas count")
	}
	return ret
}

func PromptParams(np *jww.Notepad) (params ConfigureParams) {
	params.Image = Prompt(np, "Enter image")
	imageValidate(np, params.Image)
	if portsStr := Prompt(np, "Enter ports (PORT1 PORT2 ... PORTN)"); portsStr != "" {
		params.Ports = portsValidateStr(np, portsStr)
	}
	if labelsStr := Prompt(np, "Enter labels (key1=value1 key2=value2 ... keyN=valueN)"); labelsStr != "" {
		params.Labels = labelsValidate(np, strings.Split(labelsStr, " "))
	} else {
		params.Labels = make(map[string]string)
	}
	if commands := Prompt(np, "Enter commands (command1 command2 ... commandN)"); commands != "" {
		params.Command = strings.Split(commands, " ")
	}
	if envVarsStr := Prompt(np, "Enter environment variables (key1=value1 ... keyN=valueN)"); envVarsStr != "" {
		params.Env = envVarsValidate(np, strings.Split(envVarsStr, " "))
	}
	if cpu := Prompt(np, fmt.Sprintf("Enter CPU cores (*m) [%s]", DefaultCPURequest)); cpu != "" {
		cpuValidate(np, cpu)
		params.CPU = cpu
	} else {
		params.CPU = DefaultCPURequest
	}
	if memory := Prompt(np, fmt.Sprintf("Enter memory size (*Mi | *Gi) [%s]", DefaultMemoryRequest)); memory != "" {
		memValidate(np, memory)
		params.Memory = memory
	} else {
		params.Memory = DefaultMemoryRequest
	}
	if replicas := Prompt(np, fmt.Sprintf("Enter replicas count [%d]", DefaultReplicas)); replicas != "" {
		params.Replicas = replicasValidate(np, replicas)
	} else {
		params.Replicas = DefaultReplicas
	}
	return
}

func ParamsFromArgs(np *jww.Notepad, flags *pflag.FlagSet) (params ConfigureParams) {
	chkErr := func(err error) {
		if err != nil {
			validationErrorExit(np, "flag get error: %s\n", err)
		}
	}
	var err error
	if flags.Changed("image") {
		params.Image, err = flags.GetString("image")
		chkErr(err)
		imageValidate(np, params.Image)
	}
	if flags.Changed("port") {
		params.Ports, err = flags.GetIntSlice("port")
		chkErr(err)
		portsValidateInt(np, params.Ports)
	}
	if flags.Changed("labels") {
		labelsSlice, err := flags.GetStringSlice("labels")
		chkErr(err)
		params.Labels = labelsValidate(np, labelsSlice)
	} else {
		params.Labels = make(map[string]string)
	}
	if flags.Changed("command") {
		params.Command, err = flags.GetStringSlice("command")
		chkErr(err)
	}
	if flags.Changed("env") {
		envSlice, err := flags.GetStringSlice("env")
		chkErr(err)
		params.Env = envVarsValidate(np, envSlice)
	}
	if flags.Changed("cpu") {
		params.CPU, err = flags.GetString("cpu")
		chkErr(err)
		cpuValidate(np, params.CPU)
	} else {
		params.CPU = DefaultCPURequest
	}
	if flags.Changed("memory") {
		params.Memory, err = flags.GetString("memory")
		chkErr(err)
		memValidate(np, params.Memory)
	} else {
		params.Memory = DefaultMemoryRequest
	}
	if flags.Changed("replicas") {
		params.Replicas, err = flags.GetInt("replicas")
		chkErr(err)
	} else {
		params.Replicas = DefaultReplicas
	}
	return
}
