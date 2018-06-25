There are several ways to specify the names of containers:
- flag_CONTANEER_NAME
- the prefix CONTAINER_NAME@ in the flags --image, --memory, --cpu, --env, --volume

If the --container-name flag is not specified and prefix is not used in any of the flags, then wizard searches for the --image flags without a prefix and creates the generates name RANDOM_COLOR-IMAGE.

Examples:

---
**Single container with --container-name**

```bash
> ./ckit create depl \
        --container-name doot \
        --image nginx
```

|        LABEL        | VERSION |  STATUS  |  CONTAINERS  |    AGE    |
| ------------------- | --------| -------- | ------------ | --------- |
| akiraabe-heisenberg |  1.0.0  | inactive | doot [nginx] | undefined |

---
**Single container without --container-name**

```bash
> ./ckit create depl \
        --image nginx
```

|        LABEL        | VERSION |  STATUS  |        CONTAINERS        |    AGE    |
| ------------------- | --------| -------- | ------------------------ | --------- |
|   spiraea-kaufman   |  1.0.0  | inactive | aquamarine-nginx [nginx] | undefined |

---
**Multiple containers with --container-name**


```bash
> ./ckit create depl \
        --container-name gateway \
        --image nginx \
        --image blog@wordpress
```

|        LABEL        | VERSION |  STATUS  |        CONTAINERS        |    AGE    |
| ------------------- | --------| -------- | ------------------------ | --------- |
|   ruckers-fischer   |  1.0.0  | inactive |      gateway [nginx]     | undefined |
|                     |         |          |      blog [wordpress]    |           |

---
**Multiple containers without --container-name**
```bash
> ./ckit create depl \
        --image nginx \
        --image blog@wordpress
```

|        LABEL        | VERSION |  STATUS  |        CONTAINERS        |    AGE    |
| ------------------- | ------- | -------- | ------------------------ | --------- |
|    thisbe-neumann   |  1.0.0  | inactive |      blog [wordpress]    | undefined |
|                     |         |          |    garnet-nginx [nginx]  |           |
