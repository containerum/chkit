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
                    image_check = r"^[a-zA-Z0-9_.-]*$"
                    is_valid = re.compile(image_check)
                    if not is_valid.findall(image):
                        raise ValueError("Image must contain only latin characters, numbers and hyphen")
                    param_dict.update({"image": image})
                if not param_dict.get("ports") and not self.ports:
                    try:
                        ports = input("Enter ports: ")
                    except KeyboardInterrupt:
                        return False
                    if ports:
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
                    param_dict.update({"labels": labels})
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
                    envs_check = r"^[[a-zA-Z]+[0-9]*=[a-zA-Z0-9]* *]*$"
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
                        if not "m" in cpu:
                            raise ValueError("CPU must contain m")
                        cpu_check = cpu[:-1]
                        cpu_check = int(cpu_check)
                        if cpu_check < 100 or cpu_check > 10000:
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
                        if not "Mi" in memory and not "Gi" in memory:
                            raise ValueError("Memory must contain Gi or Mi")
                        mem_check = memory[:-2]
                        mem_check = int(mem_check)
                        if (mem_check < 128 and "Mi" in memory) or (mem_check > 128 and "Gi" in memory):
                            raise ValueError("Memory must be in range [128Mi, 128Gi], for example 500Mi")
                    param_dict.update({"memory": self.memory})
                    self.memory = True
                if not param_dict.get("replicas") and self.replicas == 1:
                    try:
                        replicas = input("Enter  replicas count: ")
                    except KeyboardInterrupt:
                        return False
                    if replicas:
                        replicas = int(replicas)
                        if replicas < 1:
                            raise ValueError("Replicas must be positive integer")
                        self.replicas = replicas
                    param_dict.update({"replicas": self.replicas})
                    self.replicas = True
                break
            except Exception as e:
                print(e)

        return param_dict
