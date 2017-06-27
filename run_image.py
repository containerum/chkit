import re


class RunImage:
    def __init__(self):
        self.cpu = "100m"
        self.memory = "128Mi"
        self.ports = None
        self.commands = None
        self.env = None
        self.labels = None
        self.replicas = 1
        self.commands = []

    def parse_data(self, args):
        param_dict = {}

        try:
            image = args.get("image")
            if not image:
                raise ValueError("Image is required field")
            image_check = r"^[a-zA-Z]+([a-zA-Z0-9\/\:\_\.\-]*[a-z0-9])?$"
            is_valid = re.compile(image_check)
            if not is_valid.findall(image):
                raise ValueError("Image must consist of alphanumeric characters or -/:._ and must start and end"
                                 " with an alphanumeric character")
            param_dict.update({"image": image})
            if args.get("ports"):
                ports_check = r"^([1-9][0-9]{0,3}|[1-5][0-9]{4,5}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|" \
                              r"655[0-2][0-9]|6553[0-6])$"
                is_valid = re.compile(ports_check)
                for port in args.get("ports"):
                    if not is_valid.findall(str(port)):
                        raise ValueError("Port's range between [1, 65536]")
                ports = list(map(int, args.get("ports")))
                self.ports = ports
            param_dict.update({"ports": self.ports})
            if args.get("labels"):
                labels_dict = {}
                labels_check = r"^[[a-zA-Z]([a-zA-Z0-9\_\-]*[a-zA-Z0-9])?=[a-zA-Z0-9]([a-zA-Z0-9\_\.\-\:]*" \
                               r"[a-zA-Z0-9])? *]*$"
                is_valid = re.compile(labels_check)
                for label in args.get("labels"):
                    if not is_valid.findall(label):
                        raise ValueError("Labels must have format key=value, consist of alphanumeric characters "
                                         "or _- and must start and end with an alphanumeric character")
                    label = label.split("=")
                    labels_dict[label[0]] = label[1]
                self.labels = labels_dict
            param_dict.update({"labels": self.labels})
            if args.get("commands"):
                self.commands = args.get("commands")
            param_dict.update({"commands": self.commands})
            if args.get("env"):
                envs_dict = {}
                envs_check = r"^[[a-zA-Z]+([a-zA-Z0-9\_\-]*[a-zA-Z0-9])?=[a-zA-Z0-9]+([a-zA-Z0-9\_\.\-\:]*" \
                             r"[a-zA-Z0-9])? *]*$"
                is_valid = re.compile(envs_check)
                for env in args.get("env"):
                    if not is_valid.findall(env):
                        raise ValueError("Environ variables must have format key=value, consist of alphanumeric "
                                         "characters or _- and must start and end with an alphanumeric character")
                    env = env.split("=")
                    envs_dict[env[0]] = env[1]
                self.env = envs_dict
            param_dict.update({"env": self.env})
            if args.get("cpu") != "100m":
                cpu = args.get("cpu")
                cpu_check = r"^(([1-9][0-9]{1,2}|[1-2][0-9]{3}|3000)m)$"
                is_valid = re.compile(cpu_check)
                if not is_valid.findall(cpu):
                    raise ValueError("CPU must be in range [10m, 3000m], for example 1000m")
                self.cpu = cpu
            param_dict.update({"cpu": self.cpu})
            if args.get("memory") != "128Mi":
                memory = args.get("memory")
                mem_check = r"^((([5-9]|[1-7][0-9]{1,3}|8000)Mi)|(([1-7]*\.?[0-9]+|8)Gi))$"
                is_valid = re.compile(mem_check)
                if not is_valid.findall(memory):
                    raise ValueError("Memory must be in range [5Mi, 8Gi], for example 500Mi")
                self.memory = memory
            param_dict.update({"memory": self.memory})
            if args.get("replicas") != 1:
                replicas = args.get("replicas")
                repl_check = r"^[1-9]\d*$"
                is_valid = re.compile(repl_check)
                if not is_valid.findall(replicas):
                    raise ValueError("Replicas must be positive integer")
                self.replicas = replicas
            param_dict.update({"replicas": int(self.replicas)})
        except Exception as e:
            print(e)
            param_dict = {}

        return param_dict
