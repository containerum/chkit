import re

class RunConfigure:
    def __init__(self):
        self.cpu = "100m"
        self.memory = "128Mi"
        self.ports = None
        self.commands = None
        self.env = None
        self.labels = None
        self.replicas = 1
        self.commands = []

    def get_data_from_console(self):
        param_dict = {}
        while True:
            try:
                if not param_dict.get("image"):
                    try:
                        image = input("Enter image: ")
                    except KeyboardInterrupt:
                        return False
                    if not image:
                        raise ValueError("Image is required field")
                    image_check = r"^[a-zA-Z]+[0-9_.-]*$"
                    is_valid = re.compile(image_check)
                    if not is_valid.findall(image):
                        raise ValueError("Image must start with latin character and contain only latin "
                                         "characters, numbers and hyphen")
                    param_dict.update({"image": image})
                if not param_dict.get("ports") and not self.ports:
                    try:
                        ports = input("Enter ports: ")
                    except KeyboardInterrupt:
                        return False
                    if ports:
                        ports_check = r"^([1-9][0-9]{0,3}|[1-5][0-9]{4,5}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|" \
                                      r"655[0-2][0-9]|6553[0-6])$"
                        is_valid = re.compile(ports_check)
                        if not is_valid.findall(ports):
                            raise ValueError("Port's range between [1, 65536]")
                        ports = ports.split(" ")
                        ports = list(map(int, ports))
                    param_dict.update({"ports": ports})
                    self.ports = True
                if not param_dict.get("labels") and not self.labels:
                    try:
                        labels = input("Enter labels (key=value ... key3=value3): ")
                    except KeyboardInterrupt:
                        return False
                    if labels:
                        labels = labels.split(" ")
                    labels_dict = {}
                    labels_check = r"^[[a-zA-Z]+[0-9]*=[a-zA-Z0-9]+ *]*$"
                    is_valid = re.compile(labels_check)
                    for label in labels:
                        if not is_valid.findall(label):
                            raise ValueError("Environ variables must contain only latin characters and numbers and "
                                             "have format key=value")
                        label = label.split("=")
                        labels_dict[label[0]] = label[1]
                    param_dict.update({"labels": labels_dict})
                    self.labels = True
                if not param_dict.get("commands") and not self.commands:
                    try:
                        commands = input("Enter commands (command1 ... command3): ")
                    except KeyboardInterrupt:
                        return False
                    if commands:
                        self.commands = commands.split(" ")
                    param_dict.update({"commands": self.commands})
                    self.commands = True
                if not param_dict.get("env") and not self.env:
                    try:
                        envs = input("Enter environ variables (key=value ... key3=value3): ")
                    except KeyboardInterrupt:
                        return False
                    if envs:
                        envs = envs.split(" ")
                    envs_dict = {}
                    envs_check = r"^[[a-zA-Z]+[0-9]*=[a-zA-Z0-9]+ *]*$"
                    is_valid = re.compile(envs_check)
                    for env in envs:
                        if not is_valid.findall(env):
                            raise ValueError("Environ variables must contain only latin characters and numbers and "
                                             "have format key=value")
                        env = env.split("=")
                        envs_dict[env[0]] = env[1]
                    param_dict.update({"env": envs_dict})
                    self.env = True
                if not param_dict.get("cpu"):
                    try:
                        cpu = input("Enter  CPU cores count(*m):")
                    except KeyboardInterrupt:
                        return False
                    if cpu:
                        cpu_check = r"^[1-9][0-9]{2,3}m$"
                        is_valid = re.compile(cpu_check)
                        if not is_valid.findall(cpu):
                            raise ValueError("CPU must be in range [100m, 10000m], for example 1000m")
                        self.cpu = cpu
                    param_dict.update({"cpu": self.cpu})
                    self.cpu = True
                if not param_dict.get("memory"):
                    try:
                        memory = input("Enter memory size(*Mi | *Gi): ")
                    except KeyboardInterrupt:
                        return False
                    if memory:
                        mem_check = r"^(((12[8-9]|1[3-9][0-9]|[2-9][0-9]{2}|1[0-1][0-9]{2}|12[0-7][0-9]|1280)Mi)" \
                                    r"|(([1-9]|[1-9][0-9]|1[0-1][0-9]|12[0-8])Gi))$"
                        is_valid = re.compile(mem_check)
                        if not is_valid.findall(memory):
                            raise ValueError("Memory must be in range [128Mi, 128Gi], for example 500Mi")
                    param_dict.update({"memory": self.memory})
                    self.memory = True
                if not param_dict.get("replicas") and self.replicas == 1:
                    try:
                        replicas = input("Enter  replicas count: ")
                    except KeyboardInterrupt:
                        return False
                    if replicas:
                        repl_check = r"^[1-9]\d*$"
                        is_valid = re.compile(repl_check)
                        if not is_valid.findall(replicas):
                            raise ValueError("Replicas must be positive integer")
                        self.replicas = replicas
                    param_dict.update({"replicas": self.replicas})
                    self.replicas = True
                break
            except Exception as e:
                print(e)

        return param_dict
