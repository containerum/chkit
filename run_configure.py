class RunConfigure:
    def __init__(self):
        self.cpu = "100m"
        self.memory = "128Mi"
        self.replicas = 1
        self.commands = []

    def get_data_from_console(self):
        param_dict = {}
        while True:
            try:
                if not param_dict.get("image"):
                    image = input("Enter image:")
                    param_dict.update({"image": image})
                if not param_dict.get("ports"):
                    ports = input("Enter ports (8080 ... 4556):")
                    if ports:
                        ports = ports.split(" ")
                        ports = list(map(int, ports))
                    param_dict.update({"ports": ports})
                if not param_dict.get("labels"):
                    labels = input("Enter labels (key=value ... key3=value3):")
                    if labels:
                        labels = labels.split(" ")
                    param_dict.update({"labels": labels})
                if not param_dict.get("commands"):
                    commands = input("Enter commands (command1 ... command3):")
                    if commands:
                        self.commands = commands.split(" ")
                    param_dict.update({"commands": self.commands})
                if not param_dict.get("env"):
                    env = input("Enter environ variables (key=value ... key3=value3):")
                    if env:
                        env = env.split(" ")
                    param_dict.update({"env": env})
                if not param_dict.get("cpu"):
                    cpu = input("Enter  CPU cores count(*m):")
                    if cpu:
                        self.cpu = cpu
                    param_dict.update({"cpu": self.cpu})
                if not param_dict.get("memory"):
                    memory = input("Enter memory size(*Mi | *Gi):")
                    if memory:
                        if not "Mi" in memory and not "Gi" in memory:
                            raise ValueError("Memory must contain Gi or Mi")
                        self.memory = memory
                    param_dict.update({"memory": self.memory})
                if not param_dict.get("replicas"):
                    replicas = input("Enter  replicas count:")
                    if replicas:
                        self.replicas = int(replicas)
                    param_dict.update({"replicas": self.replicas})
                break
            except Exception as e:
                print(e)

        return param_dict