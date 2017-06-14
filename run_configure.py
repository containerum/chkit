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
                        env = input("Enter environ variables (key=value ... key3=value3): ")
                    except KeyboardInterrupt:
                        return False
                    if env:
                        env = env.split(" ")
                    param_dict.update({"env": env})
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
                        if cpu_check < 1:
                            raise ValueError("CPU must contain positive integer + m, for example 100m")
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
                        if mem_check < 1:
                            raise ValueError("Memory must contain positive integer + Mi|Gi, for example 100Mi")
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
