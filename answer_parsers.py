from datetime import datetime
from dateutil import parser
from prettytable import PrettyTable
from keywords import EMPTY_NAMESPACE, NO_NAMESPACES


class TcpApiParser:
    def __init__(self, row_answer):
        if row_answer.get("results")[0].get("data").get("kind") == "PodList":
            self.items = row_answer.get("results")[0].get("data").get("items")
            self.show_human_readable_pod_list()

        if row_answer.get("results")[0].get("data").get("kind") == "Pod":
            self.result = row_answer
            self.show_human_readable_pod()

        if row_answer.get("results")[0].get("data").get("kind") == "DeploymentList":
            self.result = row_answer
            self.show_human_readable_deployment_list()

        if row_answer.get("results")[0].get("data").get("kind") == "Deployment":
            self.result = row_answer
            self.show_human_readable_deployment()

        if row_answer.get("results")[0].get("data").get("kind") == "ServiceList":
            self.result = row_answer
            self.show_human_readable_service_list()

        if row_answer.get("results")[0].get("data").get("kind") == "Service":
            self.result = row_answer
            self.show_human_readable_service()

    def show_human_readable_pod(self):
        metadata = self.result.get("results")[0].get("data").get("metadata")
        containers = self.result.get("results")[0].get("data").get("spec").get("containers")
        restartPolicy = self.result.get("results")[0].get("data").get("spec").get("restartPolicy")
        termination = self.result.get("results")[0].get("data").get("spec").get("terminationGracePeriodSeconds")
        system = self.result.get("results")[0].get("data").get("status")
        container_statuses = self.result.get("results")[0].get("data").get("status").get("containerStatuses")
        status = self.result.get("results")[0].get("data").get("status").get("conditions")

        print("Describe:")
        print("\t%-20s %s" % ("UserId:", self.result.get("UserId")))
        print("\t%-20s %s" % ("Channel:", self.result.get("channel")))
        print("\t%-20s %s" % ("CommandId:", self.result.get("id")))
        print("Pod:")
        print("\t%-20s %s" % ("CreationTime:", parser.parse(metadata.get("creationTimestamp"))))
        print("\tLabel:")
        print("\t\t%-20s %s" % ("App:", metadata.get("labels").get("app")))
        print("\t\t%-20s %s" % ("PodTemplateHash:", metadata.get("labels").get("pod-template-hash")))
        print("\t\t%-20s %s" % ("Role:", metadata.get("labels").get("role")))
        print("\t\t%-20s %s" % ("CommandId:", self.result.get("id")))
        print("Containers:")
        for c in containers:
            print("\t%s" % c.get("name"))
            if c.get("command"):
                print("\t\t%-20s %s" % ("Command:", "".join(c.get("command"))))
            print("\t\tPorts:")
            if c.get("ports"):
                ports = PrettyTable(["Name", "Protocol", "ContPort"])
                for p in c.get("ports"):
                    ports.add_row([p.get("name"), p.get("protocol"), p.get("containerPort")])
                print(ports)
            if c.get("env"):
                env = PrettyTable(["Name","Value"])
                print("\t\tEnvironment:")
                for e in c.get("env"):
                    env.add_row([e.get("name"), e.get("value")])
                print(env)
            print("\t\tResourceLimit:")
            print("\t\t\t%-10s %s" % ("CPU:", c.get("resources").get("limits").get("cpu")))
            print("\t\t\t%-10s %s" % ("Memory:", c.get("resources").get("limits").get("memory")))
            print("\t\t%-20s %s" % ("Image:", c.get("image")))
            print("\t\t%-20s %s" % ("ImagePullPolicy:", c.get("imagePullPolicy")))
            print("System:")
            print("\t%-30s %s" % ("PodIP:", system.get("podIP")))
            print("\t%-30s %s" % ("Phase:", system.get("phase")))
            if system.get("startTime"):
                print("\t%-30s %s" % ("StartTime:", parser.parse(system.get("startTime"))))
            print("\t%-30s %s" % ("TerminationGracePeriodSeconds:", termination))
            print("\t%-30s %s" % ("RestartPolicy:", restartPolicy))
            print("ContainerStatuses:")
            if container_statuses:
                containerStatuses = PrettyTable(["Name","Ready","Restart Count"])
                for cs in container_statuses:
                    containerStatuses.add_row([cs.get("name"), cs.get("ready"), cs.get("restartCount")])
                print(containerStatuses)
            print("Status:")
            StatusTable = PrettyTable(["Type:", "LastTransitionTime:", "Status:"])
            for s in status:
                StatusTable.add_row([s.get("type"), parser.parse(s.get("lastTransitionTime")), s.get("status")])
            print(StatusTable)

    def show_human_readable_pod_list(self):
        if self.items:
            table = PrettyTable(["NAME", "READY", "STATUS", "RESTARTS", "AGE", "IP"])
            table.align = "l"
            for i in self.items:
                restarts = i.get("status").get("containerStatuses")
                restarts_sum = 0
                if restarts:
                    for r in restarts:
                        restarts_sum += r.get("restartCount")
                    else:
                        restarts_sum = None
                ip = i.get("status").get("podIP")
                status = i.get("status").get("phase")
                name = i.get("metadata").get("name")
                ready = "-/-"
                time = get_datetime_diff(i.get("metadata").get("creationTimestamp"))
                table.add_row([name, ready, status, restarts_sum, time, ip])
            print(table)
        else:
            print(EMPTY_NAMESPACE)

    def show_human_readable_deployment_list(self):
        if self.result:
            items = self.result.get("results")[0].get("data").get("items")
            table = PrettyTable(["NAME",  "PODS ACTIVE",  "CPU",  "RAM", "AGE"])
            table.align = "l"
            for i in items:
                name = i.get("metadata").get("name")
                cpu = i.get("spec").get("template").get("spec").get("containers")[0].get("resources")\
                    .get("limits").get("cpu")
                if i.get("status").get("availableReplicas"):
                    pods_active = i.get("status").get("availableReplicas")
                else:
                    pods_active = 0
                memory = i.get("spec").get("template").get("spec").get("containers")[0].get("resources")\
                    .get("limits").get("memory")
                time = get_datetime_diff(i.get("metadata").get("creationTimestamp"))
                table.add_row([name,  pods_active,  cpu,  memory, time])
            print(table)
        else:
            print(NO_NAMESPACES)

    def show_human_readable_deployment(self):
        status = self.result.get("results")[0].get("data").get("status")
        strategy = self.result.get("results")[0].get("data").get("spec").get("strategy")
        conditions = self.result.get("results")[0].get("data").get("status").get("conditions")
        containers = self.result.get("results")[0].get("data").get("spec").get("template")\
            .get("spec").get("containers")
        if self.result:
            print("%-30s %s" % ("Name:", self.result.get("name")))
            print("%-30s %s" % ("Namespace:", self.result.get("namespace")))
            print("%-30s %s" % ("CreationTimeStamp:",
                                parser.parse(self.result.get("results")[0].get("data")
                                             .get("metadata").get("creationTimestamp"))))
            print("Labels:")
            for key,value in self.result.get("results")[0].get("data").get("metadata").get("labels").items():
                print("\t%s=%s" % (key, value))
            print("Selectors:")
            for key,value in self.result.get("results")[0].get("data").get("spec").get("selector").get("matchLabels").items():
                print("\t%s=%s" % (key, value))
            status_tuple = ("Replicas:", status.get("updatedReplicas"), "updated", status.get("replicas"), "total",
                            status.get("replicas")-status.get("unavailableReplicas"), "available", status.get("unavailableReplicas"), "unavailable")
            print("%-30s %s %s | %s %s | %s %s | %s %s" % status_tuple)
            print("%-30s %s " % ("Strategy", strategy.get("type")))
            strategy_type = strategy.get("type")[0].lower() + strategy.get("type")[1:]
            print("%-30s %s max unavailable, %s max surge" % (strategy.get("type")+"Strategy",
                                                              strategy.get(strategy_type).get("maxUnavailable"),
                                                              strategy.get(strategy_type).get("maxSurge")))
            print("Conditions:")
            conditions_table = PrettyTable(["TYPE", "STATUS", "REASON"])
            for c in conditions:
                conditions_table.add_row([c.get("type"), c.get("status"), c.get("reason")])
            print(conditions_table)
            print("Containers:")
            for c in containers:
                print("\t%s" % c.get("name"))

    def show_human_readable_service_list(self):
        if self.result:
            items = self.result.get("results")[0].get("data").get("items")
            table = PrettyTable(["NAME",  "CLUSTER-IP",  "EXTERNAL-IP",  "PORT(S)", "AGE"])
            table.align = "l"
            for i in items:
                name = i.get("metadata").get("name")
                cluster_ip = i.get("spec").get("clusterIP")
                if i.get("spec").get("externalIPs"):
                    external_ip = " ,\n".join(i.get("spec").get("externalIPs"))
                else:
                    external_ip = "<none>"
                ports = i.get("spec").get("ports")
                for p in range(len(ports)):
                    if ports[p].get("port") == ports[p].get("targetPort"):
                        ports[p] = ("%s/%s" % (ports[p].get("port"), ports[p].get("protocol")))
                    else:
                        ports[p] = ("%s:%s/%s" % (ports[p].get("port"), ports[p].get("targetPort"), ports[p].get("protocol")))
                sum_ports = " ,\n".join(ports)
                time = get_datetime_diff(i.get("metadata").get("creationTimestamp"))
                table.add_row([name,  cluster_ip,  external_ip, sum_ports, time])
            print(table)

    def show_human_readable_service(self):
        if self.result:
            print("%-30s %s" % ("Name:", self.result.get("name")))
            print("%-30s %s" % ("Namespace:", self.result.get("namespace")))
            for key,value in self.result.get("results")[0].get("data").get("metadata").get("labels").items():
                print("\t%s=%s" % (key, value))
            print("Selectors:")
            for key,value in self.result.get("results")[0].get("data").get("spec").get("selector").items():
                print("\t%s=%s" % (key, value))
            print("%-30s %s " % ("Type:", self.result.get("results")[0].get("data").get("spec").get("type")))
            print("%-30s %s " % ("IP:", self.result.get("results")[0].get("data").get("spec").get("clusterIP")))
            ports = self.result.get("results")[0].get("data").get("spec").get("ports")
            for p in ports:
                if p.get("port") == p.get("targetPort"):
                    print("%-30s %s/%s" % ("Port:", p.get("port"), p.get("protocol")))
                else:
                    print("%-30s %s:%s/%s" % ("Port:", p.get("port"), p.get("targetPort"), p.get("protocol")))
            print("%-30s %s " % ("External IPs:", "----"))



class WebClientApiParser:
    def __init__(self, row_answer):
        self.items = row_answer
        self.result = row_answer

    def show_human_readable_namespace_list(self):
        if self.items:
            table = PrettyTable(["ID",  "IS ACTIVE",  "AGE", "CPU", "CPU LIMIT", "MEMORY", "MEMORY LIMIT"])
            table.align = "l"
            for i in self.items:
                status = i.get("active")
                name = i.get("id")
                cpu = i.get("cpu")
                cpu_limit = i.get("cpu_limit")
                memory = i.get("memory")
                memory_limit = i.get("memory_limit")
                time = get_datetime_diff(i.get("created"))
                table.add_row([name,  status,  time, cpu, cpu_limit, memory, memory_limit])
            print(table)
        else:
            print(NO_NAMESPACES)

    def show_human_readable_deployment_list(self):
        print(self.items)
        if self.items:
            table = PrettyTable(["NAME",  "PODS ACTIVE",  "CPU",  "RAM", "AGE"])
            table.align = "l"
            for i in self.items:
                name = i.get("name")
                cpu = i.get("cpu")
                pods_active = i.get("pods_active")
                memory = i.get("ram")
                time = get_datetime_diff(i.get("created_at"))
                table.add_row([name,  pods_active,  cpu,  memory, time])
            print(table)
        else:
            print(NO_NAMESPACES)

    def show_human_readable_deployment(self, namespace):
        status = self.result.get("status")
        conditions = self.result.get("conditions")
        containers = self.result.get("containers")
        if self.result:
            print("%-20s %s" % ("Name:", self.result.get("name")))
            print("%-20s %s" % ("Namespace:", namespace))
            print("%-20s %s" % ("CreationTimeStamp:", parser.parse(self.result.get("created_at"))))
            print("Labels:")
            for key,value in self.result.get("labels").items():
                print("\t%s=%s" % (key, value))
            status_tuple = ("Status:", status.get("updated"), "updated", status.get("total"), "total",
                            status.get("available"), "available", status.get("unavailable"), "unavailable")
            print("%-20s %s %s | %s %s | %s %s | %s %s" % status_tuple)
            print("Conditions:")
            conditions_table = PrettyTable(["TYPE", "STATUS", "REASON"])
            for c in conditions:
                conditions_table.add_row([c.get("type"), c.get("status"), c.get("reason")])
            print(conditions_table)
            print("Containers:")
            for c in containers:
                print("\t%s" % c.get("name"))


def get_datetime_diff(timestamp):
    created_date = parser.parse(timestamp)
    current_date = datetime.utcnow()
    diff = ((current_date.year - created_date.year, "Y"),(current_date.month - created_date.month, "M"),
            (current_date.day - created_date.day, "d"), (current_date.hour - created_date.hour, "h"),
            (current_date.minute - created_date.minute, "m"), (current_date.second - created_date.second,"s"))
    diff = tuple(filter(lambda x: x[0] > 0, diff))[0]
    return str(diff[0]) + diff[1]