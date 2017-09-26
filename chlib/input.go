package chlib

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"fmt"

	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/pflag"
)

type validationFunc func(string) bool

var reader = bufio.NewReader(os.Stdin)

func Prompt(np *jww.Notepad, prompt string, validator validationFunc) string {
	np.FEEDBACK.Printf("%s: ", prompt)
	ret, _ := reader.ReadString('\n')
	ret = strings.TrimRight(ret, "\r\n")
	for !validator(ret) {
		np.FEEDBACK.Printf("Invalid input\n%s: ", prompt)
		ret, _ = reader.ReadString('\n')
		ret = strings.TrimRight(ret, "\r\n")
	}
	return ret
}

func portsValidateStr(portsStr string) bool {
	if portsStr == "" {
		return true
	}
	for _, portStr := range strings.Split(portsStr, " ") {
		port, err := strconv.Atoi(portStr)
		if err != nil || port <= 0 || port > 65535 {
			return false
		}
	}
	return true
}

func portsValidateInt(ports []int) bool {
	for _, port := range ports {
		if port <= 0 || port > 65535 {
			return false
		}
	}
	return true
}

func labelSliceValidate(labels []string) bool {
	for _, labelStr := range labels {
		label := strings.Split(labelStr, "=")
		if len(label) != 2 || !LabelRegex.MatchString(label[0]) || !LabelRegex.MatchString(label[1]) {
			return false
		}
	}
	return true
}

func labelsValidate(labelsStr string) bool {
	return labelsStr == "" || labelSliceValidate(strings.Split(labelsStr, " "))
}

func envSliceValidate(envVars []string) bool {
	for _, envVarStr := range envVars {
		envVar := strings.Split(envVarStr, "=")
		if len(envVar) != 2 {
			return false
		}
	}
	return true
}

func envVarsValidate(envVarsStr string) bool {
	return envVarsStr == "" || envSliceValidate(strings.Split(envVarsStr, " "))
}

func cpuValidate(cpuStr string) bool {
	if cpuStr == "" {
		return true
	}
	return CpuRegex.MatchString(cpuStr)
}

func memValidate(memStr string) bool {
	if memStr == "" {
		return true
	}
	return MemRegex.MatchString(memStr)
}

func replicasValidate(replicasStr string) bool {
	if replicasStr == "" {
		return true
	}
	ret, err := strconv.Atoi(replicasStr)
	return err == nil && ret > 0
}

func volumesValidate(volumesStr string) bool {
	if volumesStr == "" {
		return true
	}
	for _, v := range strings.Fields(volumesStr) {
		if !VolumeRegex.MatchString(v) {
			return false
		}
		if strings.Count(v, "\"") % 2 != 0 {
			return false
		}
	}
	return true
}

func volumesSliceValidate(vols []string) bool {
	for _, v := range vols {
		if !VolumeRegex.MatchString(v) {
			return false
		}
		if strings.Count(v, "\"") % 2 != 0 {
			return false
		}
	}
	return true
}

func PromptParams(np *jww.Notepad) (params ConfigureParams) {
	params.Image = Prompt(np, "Enter image", ImageRegex.MatchString)
	if portsStr := Prompt(np, "Enter ports (PORT1 PORT2 ... PORTN)", portsValidateStr); portsStr != "" {
		for _, p := range strings.Split(portsStr, " ") {
			port, _ := strconv.Atoi(p)
			params.Ports = append(params.Ports, port)
		}
	}
	params.Labels = make(map[string]string)
	if labelsStr := Prompt(np, "Enter labels (key1=value1 key2=value2 ... keyN=valueN)", labelsValidate); labelsStr != "" {
		for _, labelStr := range strings.Split(labelsStr, " ") {
			label := strings.Split(labelStr, "=")
			params.Labels[label[0]] = label[1]
		}
	}
	if commands := Prompt(np, "Enter commands (command1 command2 ... commandN)", func(string) bool { return true }); commands != "" {
		params.Command = strings.Split(commands, " ")
	}
	if envVarsStr := Prompt(np, "Enter environment variables (key1=value1 ... keyN=valueN)", envVarsValidate); envVarsStr != "" {
		for _, envVarStr := range strings.Split(envVarsStr, " ") {
			envVar := strings.Split(envVarStr, "=")
			params.Env = append(params.Env, EnvVar{
				Name:  envVar[0],
				Value: envVar[1],
			})
		}
	}
	if volumesStr := Prompt(np, "Enter volumes (volumeLabel1[/subPath1]=/mountPath1 ... volumeLabelN[/subPathN]=/mountPathN)", volumesValidate); volumesStr != "" {
		for _, volumeStr := range strings.Fields(volumesStr) {
			volume := VolumeRegex.FindStringSubmatch(volumeStr)
			params.Volumes = append(params.Volumes, Volume{
				Label:     volume[1],
				SubPath:   volume[4],
				MountPath: volume[5],
			})
		}
	}
	if cpu := Prompt(np, fmt.Sprintf("Enter CPU cores (*m) [%s]", DefaultCPURequest), cpuValidate); cpu != "" {
		params.CPU = cpu
	} else {
		params.CPU = DefaultCPURequest
	}
	if memory := Prompt(np, fmt.Sprintf("Enter memory size (*Mi | *Gi) [%s]", DefaultMemoryRequest), memValidate); memory != "" {
		params.Memory = memory
	} else {
		params.Memory = DefaultMemoryRequest
	}
	if replicas := Prompt(np, fmt.Sprintf("Enter replicas count [%d]", DefaultReplicas), replicasValidate); replicas != "" {
		params.Replicas, _ = strconv.Atoi(replicas)
	} else {
		params.Replicas = DefaultReplicas
	}
	return
}

func exitIfValidationError(np *jww.Notepad, validationResult bool, message string) {
	if !validationResult {
		np.FEEDBACK.Println(message)
		os.Exit(1)
	}
}

func ParamsFromArgs(np *jww.Notepad, flags *pflag.FlagSet) (params ConfigureParams) {
	chkErr := func(err error) {
		if err != nil {
			np.FEEDBACK.Println("flag get error: %s\n", err)
		}
	}
	var err error
	if flags.Changed("image") {
		params.Image, err = flags.GetString("image")
		chkErr(err)
		exitIfValidationError(np, ImageRegex.MatchString(params.Image), "Invalid image name")
	}
	if flags.Changed("port") {
		params.Ports, err = flags.GetIntSlice("port")
		chkErr(err)
		exitIfValidationError(np, portsValidateInt(params.Ports), "Invalid port found")
	}
	if flags.Changed("labels") {
		labelsSlice, err := flags.GetStringSlice("labels")
		chkErr(err)
		exitIfValidationError(np, labelSliceValidate(labelsSlice), "Invalid label found")
		for _, labelStr := range labelsSlice {
			label := strings.Split(labelStr, "=")
			params.Labels[label[0]] = label[1]
		}
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
		exitIfValidationError(np, envSliceValidate(envSlice), "Invalid environment variable found")
		for _, envStr := range envSlice {
			env := strings.Split(envStr, "=")
			params.Env = append(params.Env, EnvVar{
				Name:  env[0],
				Value: env[1],
			})
		}
	}
	if flags.Changed("volumes") {
		volsSlice, err := flags.GetStringSlice("volumes")
		chkErr(err)
		exitIfValidationError(np, volumesSliceValidate(volsSlice), "Invalid volume parameter found")
		for _, volStr := range volsSlice {
			vol := VolumeRegex.FindStringSubmatch(volStr)
			params.Volumes = append(params.Volumes, Volume{
				Label:     vol[1],
				SubPath:   vol[4],
				MountPath: vol[5],
			})
		}
	}
	if flags.Changed("cpu") {
		params.CPU, err = flags.GetString("cpu")
		chkErr(err)
		exitIfValidationError(np, cpuValidate(params.CPU), "Invalid CPU format")
	} else {
		params.CPU = DefaultCPURequest
	}
	if flags.Changed("memory") {
		params.Memory, err = flags.GetString("memory")
		chkErr(err)
		exitIfValidationError(np, cpuValidate(params.Memory), "Invalid memory format")
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
