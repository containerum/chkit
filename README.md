# Руководство по работе с Containerum CLI

На этой странице:

* [Установка Containerum CLI с помощью бинарных сборок](https://bitbucket.org/exonch/ch-sdk/#markdown-header-containerum-cli)
* [Аутентификация](https://bitbucket.org/exonch/ch-sdk/#markdown-header-)
* [Настройка Containerum CLI](https://bitbucket.org/exonch/ch-sdk/#markdown-header-cli-containerum)
* [О типах объектов](https://bitbucket.org/exonch/ch-sdk/#markdown-header-_1)
* [Команды Containerum CLI](https://bitbucket.org/exonch/ch-sdk/#markdown-header-cli)
    * [login](https://bitbucket.org/exonch/ch-sdk/#markdown-header-login)
    * [help](https://bitbucket.org/exonch/ch-sdk/#markdown-header-help)
    * [config](https://bitbucket.org/exonch/ch-sdk/#markdown-header-config)
    * [run](https://bitbucket.org/exonch/ch-sdk/#markdown-header-run)
    * [expose](https://bitbucket.org/exonch/ch-sdk/#markdown-header-expose)
    * [create](https://bitbucket.org/exonch/ch-sdk/#markdown-header-create)
    * [set](https://bitbucket.org/exonch/ch-sdk/#markdown-header-set)
    * [get](https://bitbucket.org/exonch/ch-sdk/#markdown-header-get)
    * [delete](https://bitbucket.org/exonch/ch-sdk/#markdown-header-delete)
    * [logout](https://bitbucket.org/exonch/ch-sdk/#markdown-header-logout)

# Установка Containerum CLI с помощью бинарных сборок
Скачать бинарные сборки под [MacOs](http://p13000.x1.containerum.io/download/packages/beta/v1.0/mac/), [Ubuntu x32](http://p13000.x1.containerum.io/download/packages/beta/v1.0/ubuntu_x32/) или Ubuntu x64.

Распакуйте скаченный файл в удобное для Вас место:

```
$ unzip your_archive.zip -d /path/to/destination/dir/
```

Запуск клиента из /path/to/destination/dir/:

```
$ ./chkit
```

_Примечание:_ для запуска клиента из любой директории, привяжите его к одной из
директорий, находящихся в переменной $PATH:

```
$ echo $PATH
```

**Пример**

Привяжем клиента к директории /usr/local/bin

```
$ sudo ln -sf path/to/destination/dir/chkit.py /usr/local/bin/chkit
```

Теперь клиента можно вызвать простой командой из любой директории:

```
$ chkit
```

# Аутентификация
Прежде чем приступить к работе с Containerum CLI, нужно указать Ваш токен (TOKEN). Токен можно найти на [containerum.io/profile](https://www.google.com).

```
$ chkit config  --set-token TOKEN
```
**Пример**
```
$ chkit config  --set-token QA0u64rOkTtCxxxxxxxxxxliUAnBnPlCbGQfpCQpzqM=

Success changed!
token: QA0u64rOkTtCxxxxxxxxxxliUAnBnPlCbGQfpCQpzqM=
```

# Настройка Containerum CLI
Настройка Containerum CLI  выполняется с помощью команды `config`.

На данный момент пользователь может выбрать *Namespace*, в котором будет работать.

Пространство пользователя *Namespace* задано по умолчанию как default. Если у Вас несколько *Namespace*, то можно явно указать *Namespace*, выбрав его из имеющихся.
```
$ chkit config --set-default-namespace NAMESPACE
```
**Пример**

Вывод списка всех *Namespace* пользователя:
```
$ chkit get ns

+-------------+----------+-------------+----------+-------------+-----+
| NAME        | HARD CPU | HARD MEMORY | USED CPU | USED MEMORY | AGE |
+-------------+----------+-------------+----------+-------------+-----+
| default     | 2        | 3Gi         | 300m     | 300Mi       | 1M  |
| myns        | 2        | 3Gi         | 140m     | 30Mi        | 1M  |
+-------------+----------+-------------+----------+-------------+-----+
```
Выбор myns:
```
$ chkit config —set-default-namespace myns

Success changed!
namespace: myns
```

# О типах объектов
Сontainerum CLI предоставляет доступ к четырем типам объектов: *Namespace*, *Deployment*, *Pod*, *Service*.

С помощью Containerum CLI доступно управление тремя типами объектов: *Deployment*, *Pod*, *Service*, переключаясь между *Namespace*.

*Namespace* - выделенный ресурс, объединяющий объекты пользователя или группы пользователей в единое пространство, на которое выделяется объем памяти и ресурс CPU.

*Deployment* - контроллер управления одним или несколькими контейнерами, объединенными в *Pods*.

*Pod* - группа из одного или нескольких контейнеров. *Pod* использует выделенный объем памяти и ресурс CPU, указанный в *Namespace*.

*Service* -объект, обеспечивающий общий доступ к *Pod*.

### Список синонимов
```
* Deployment: deploy, deployment, deployments
* Pod: po, pod, pods                          
* Service: svc, service, services           
```

# Команды Containerum CLI
Список команд Containerum CLI и примеры их использования.

## login

Команда `chkit login` открывает сессию и устанавливает токен в `config`.

### Синтаксис команды
```
$ chkit login
```
_Примечание:_ команда вызывается без аргументов. При вызове команды в диалоговом режиме вводятся e-mail и пароль.

**Пример**
```
$ chkit login

Enter your email: test@gmail.com
Password:

Success changed!
token: QA0u64rOkTtCxxxxxxxxxxliUAnBnPlCbGQfpCQpzqM=
```

## help

Команда `chkit --help` или `chkit -h` показывает список всех команд и их краткое описание.


### Синтаксис команды
```
$ chkit —h
```
_Примечание:_ команда вызывается без аргументов.

**Пример**
```
$ chkit get -h

usage: chkit [--debug -d ] get (KIND [NAME] | --file -f FILE) [--output -o OUTPUT] [--namespace -n NAMESPACE][--deploy -d DEPLOY][-h | --help]
Show info about pod(s), service(s), namespace(s), deployment(s)

positional arguments:
  KIND              {namespace,deployment,service,pod} object kind
  NAME              object name to get info

get arguments:
  -h, --help                           show this help message and exit
  --file FILE, -f FILE                 input file
  --output OUTPUT, -o OUTPUT           {yaml,json} output format, default: json
  --namespace NAMESPACE, -n NAMESPACE  namespace, default: "default"
```

## config

Команда `chkit config` позволяет пользователю задать конфигурацию CLI, используемую другими командами по умолчанию.

### Синтаксис команды

Для вызова команды `chkit config` требуется указать:

| Ключ                           | Параметр  | Описание                                                   |
|--------------------------------|-----------|------------------------------------------------------------|
| `--set-token` или `-t`             | TOKEN     | значение токена                                            |
| `--set-default-namespace` или `-n` | NAMESPACE | имя *Namespace*.                                                              _Примечание:_ по умолчанию NAMESPACE = default |

Необязательные параметры:

| Ключ          | Параметр | Описание                                        |
|---------------|----------|-------------------------------------------------|
| `--help` или `-h` |          | вывод справки о команде                         |
| `--debug` или `-d `   |          | вывод системной информации о выполнении команды |

```
$ chkit config [--debug -d ](--set-token -t TOKEN  | --set-default-namespace -n NAMESPACE )[--help | -h]
```
**Пример**
```
$ chkit config

namespace: default
token: QA0u64rOkTtCxxxxxxxxxxliUAnBnPlCbGQfpCQpzqM=
```

## run

Команда `chkit run` создает *Deployment* и автоматически JSON файл, который содержит параметры *Deployment*. Файл run.json сохраняется в директорию $HOME/.containerum/src/json_templates.

### Синтаксис команды
Для вызова команды `chkit run`  требуется указать:

| Ключ                  | Параметр | Описание                                                               |
|----------------------|----------|------------------------------------------------------------------------|
|                | NAME     | имя объекта. _Примечание:_ имя объекта не должно содержать заглавных букв |
| `--image` или `-i` | IMAGE    | имя образа                                                             |
| `--configure`    |          | возможность поэтапного ввода парметров в диалоговом режиме                                |


Необязательные параметры:

| Ключ                | Параметр       | Описание                                                         |
|---------------------|----------------|------------------------------------------------------------------|
| `--help` или `-h`       |                | вывод справки о команде                                          |
| `--env` или `-e`        | ENV            | переменные окружения для контейнера в *Pod*                        |
| `--ports` или `-p`      | PORTS          | порты, которые будут открыты                                     |
| `--replicas` или `-r`   | REPLICAS_COUNT | количество реплик для *Pod*                                        |
| `--memory или -m`     | MEMORY         | количество памяти RAM на *Pod*                                     |
| `--cpu` или `-c`        | CPU            | выделенная часть ресурсов CPU доступная *Pod*                      |
| `--commands` или `-cmd` | COMMANDS       | команды, которые будут выполнены при запуске контейнера в *Pod*    |
| `--labels` или `-ls`    | LABELS         | теги для *Deployment*. У всех *Pod* в *Deployment* одни и те же теги   |
| `--namespace` или `-n`  | NAMESPACE      | название *Namespace*. _Примечание:_ по умолчанию NAMESPACE = default |
| `--debug` или `-d`          |                | вывод системной информации о выполнении команды                  |

```
$ chkit [--debug -d ] run NAME\
--configure |\
--image -i IMAGE\
[--env -e «KEY=VALUE»]\
[--port -p PORT]\
[--replicas -r REPLICAS_COUNT]\
[--memory -m MEMORY]\
[--cpu -c CPU]\
[--command -cmd COMMAND]\
[--labels -ls «KEY=VALUE»]\
[--namespace -n NAMESPACE]\
[--help | -h]
```

Альтернативный вариант вызова команды `chkit run`, используя флаг `--configure`:

```
$ chkit run NAME --configure
```

Далее необходимо ввести единственный обязательный параметр - имя образа приложения. Все остальные необязательные параметры, описанные выше, будут предложены для ввода в диалоговом режиме, но их можно будет пропустить.

### Единицы измерения CPU и RAM

Ресурсы CPU измеряются в cpus. Доступны для использования как целые, так и дробные значения. Используйте суффикс m (mili, мили).
Например, CPU = 100m = 100mcpu = 0.1cpu.

Ресурсы RAM измеряются в байтах. Доступны для использования как целые, так и дробные значения. Используйте суффиксы Mi(Mega, мега) и Gi(Giga, гига).
Например, RAM = 1,28e+8байт = 128Mi = 128Mb = 0,128Gi = 0,128Gb.

**Пример**
```
$ chkit run myapp --configure

Enter image:nginx
Enter ports (8080 ... 4556):80
Enter labels (key=value ... key3=value3):app=nginx type=local
Enter commands (command1 ... command3):
Enter environ variables (key=value ... key3=value3):HELLO=WORLD
Enter  CPU cores count(*m):100m
Enter memory size(*Mi | *Gi):200Mi
Enter  replicas count:2

run... OK
```

## expose

Команда `chkit expose` создает *Service*, в котором устанавливается протокол и список выходных портов. Также создает автоматически JSON файл, который содержит параметры *Service*. Файл expose.json сохраняется в директорию $HOME/.containerum/src/json_templates.


### Синтаксис команды
Для вызова команды `chkit expose`  требуется указать:

| Ключ           | Параметр | Описание                                                                                                                                                                                                                                                                                                                                                                |
|----------------|----------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
|                | KIND     | тип объекта: *Deployment*                                                                                                                                                                                                                                                                                                                                                 |
|                | NAME     | имя объекта. _Примечание:_ имя объекта не должно содержать заглавных букв                                                                                                                                                                                                                                                                                                   |
| `--ports` или `-p` | PORTS    | Формат ввода:  Для создания внешнего сервиса  PORTNAME:TARGETPORT[:PROTOCOL]  Для создания внутреннего сервиса  PORTNAME:TARGETPORT:PORT[:PROTOCOL]  PORTNAME - имя порта, используемого в сервисе TARGETPORT - номер порта, используемого в *Pod* PORT - номер внутреннего порта PROTOCOL - протокол передачи данных TCP или UDP _Примечание:_: По умолчанию PROTOCOL = TCP |

Необязательные параметры:

| Ключ               | Параметр  | Описание                                                         |
|--------------------|-----------|------------------------------------------------------------------|
| `--help` или   `-h`      |           | вывод справки о команде                                          |
| `--namespace` или `-n` | NAMESPACE | название *Namespace*. _Примечание:_ по умолчанию NAMESPACE = default<br/> |
| `--debug` или `-d`         |           | вывод системной информации о выполнении команды                  |


```
$ chkit [--debug -d ] expose KIND NAME (-p --ports PORTS)[--namespace -n NAMESPACE][--help | -h]
```
**Пример**
```
$ chkit expose deploy myapp -p portname:2321:TCP

expose... OK
```

## create

Команда `chkit create` создает один из 3-х типов объектов:

+ *Deployment*
+ *Pod*
+ *Service*

из JSON файла.

Создать JSON файл возможно с помощью [шаблонов](www.google.com).

### Синтаксис команды
Команда `chkit create` вызывается c одним обязательным параметром:

| Ключ          | Параметр | Описание       |
|---------------|----------|----------------|
| `--file` или `-f` | FILE     | имя JSON файла |


Необязательные параметры:

| Ключ          | Параметр | Описание                                        |
|---------------|----------|-------------------------------------------------|
| `--help` или `-h` |          | вывод справки о команде                         |
| `--debug` или `-d`    |          | вывод системной информации о выполнении команды |

```
$ chkit [--debug -d ] create (-f FILE | --file FILE)[--help | -h]
```
**Пример**
```
$ chkit create -f MyDeploy.json

create... OK
```

## set

Команда `chkit set` меняет один из параметров в *Deployment*.

_Примечание:_ на данным момент доступно изменение параметра image (образ приложения).

### Синтаксис команды
Для вызова команды `chkit set` требуется указать:

| Ключ | Параметр        | Описание                                                                        |
|------|-----------------|---------------------------------------------------------------------------------|
|      | FIELD           | изменяемый параметр в *Deployment*. _Примечание:_ доступное поле для изменения image |
|      | TYPE            | тип объекта:*Deployment*                                                          |
|      | NAME            | имя объекта                                                                     |
|      | CONTAINER_NAME  | имя контейнера                                                                  |
|      | CONTAINER_IMAGE | образ приложения                                                                |


Необязательные параметры:

| Ключ               | Параметр  | Описание                                                    |
|--------------------|-----------|-------------------------------------------------------------|
| `--help` или `-h`      |           | вывод справки о команде                                     |
| `--namespace` или `-n` | NAMESPACE | имя *Namespace*. _Примечание:_ по умолчанию NAMESPACE = default |
| `--debug` или `-d`     |           | вывод системной информации о выполнении команды             |

```
$ chkit set [--debug -d ]FIELD (TYPE NAME) CONTAINER_NAME=CONTAINER_IMAGE [-n --namespace NAMESPACE]
```
**Пример**
```
$ chkit set image deploy myapp myapp=nginx

http://146.185.135.181:3333/namespaces/default/container/myapp
set... OK
```

## get

Команда `chkit get` выводит список всех имеющихся у пользователя объектов одного из 3-х типов:

+ *Deployment*
+ *Pod*
+ *Service*
и информацию о *Namespace*.

При указании имени объекта `chkit get` выводит информацию по нему.

### Синтаксис команды
Для вызова команды chkit get требуется указать:

| Ключ | Параметр | Описание                              |
|------|----------|---------------------------------------|
|      | KIND     | тип объекта: *Deployment*,  *Pod*, *Service* |

или

| Ключ          | Параметр | Описание                                        |
|---------------|----------|-------------------------------------------------|
| `--file` или `-f` | FILE     | JSON файл, сгенерированный при создании объекта |

Необязательные параметры:

| Ключ               | Параметр  | Описание                                                                    |
|--------------------|-----------|-----------------------------------------------------------------------------|
|                    | NAME      | имя существующего объекта                                                   |
| `--output` или `-o`    | OUTPUT    | формат вывода: json, yaml, pretty. _Примечание:_ по умолчанию OUTPUT = pretty |
| `--namespace` или `-n` | NAMESPACE | имя *Namespace*. _Примечание:_ по умолчанию NAMESPACE = default                 |
| `--help` или `-h`      |           | вывод справки о команде                                                     |
| `--debug` или `-d`        |           | вывод системной информации о выполнении команды                             |

```
$ chkit [--debug -d ] get (KIND [NAME] | --file —f FILE) [-o OUTPUT] [--namespace NAMESPACE][--help | -h]
```
**Пример**
```
$ chkit get deploy

+------------+------+-------------+------+-------+-----+
| NAME       | PODS | PODS ACTIVE | CPU  | RAM   | AGE |
+------------+------+-------------+------+-------+-----+
| myapp      | 2    | 2           | 200m | 256Mi | 18s |
+------------+------+-------------+------+-------+-----+
```

## delete

Команда `chkit delete` используется для удаления объекта из *Namespace*. Доступные для удаления типы объектов:

+ *Deployment*
+ *Pod*
+ *Service*

### Синтаксис команды
Для вызова команды `chkit delete` требуется указать:

| Ключ | Параметр | Описание                              |
|------|----------|---------------------------------------|
|      | KIND     | тип объекта: *Deployment*, *Pod*, *Service* |
|      | NAME     | имя объекта                           |

или  

| Ключ          | Параметр | Описание                                        |
|---------------|----------|-------------------------------------------------|
| `--file` или `-f` | FILE     | JSON файл, сгенерированный при создании объекта |

Необязательные параметры:

| Ключ               | Параметр  | Описание                                                   |
|--------------------|-----------|------------------------------------------------------------|
| `--namespace` или `-n` | NAMESPACE | имя *Namespace*. _Примечание:_ по умолчанию NAMESPACE = default |
| `--help` или `-h`      |           | вывод справки о команде                                    |
| `--debug` или `-d`         |           | вывод системной информации о выполнении команды            |

```
$ chkit [--debug -d ] delete (KIND NAME | --file -f FILE) [--namespace NAMESPACE][--help | -h]
```
**Пример**
```
$ chkit delete deploy myapp

delete... OK
```

## logout

Команда `chkit logout` завершает сессию и сбрасывает токен в `config`.

### Синтаксис команды
```
$ chkit logout
```
_Примечание:_ команда вызывается без аргументов.

**Пример**
```
$ chkit logout

Success changed!
Bye!
```