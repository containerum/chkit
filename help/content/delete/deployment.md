Delete deployment. List of deployments to delete must be provided as argument. If list is empty, then chkit will start interactive menu.

**Delete list of deployments without --force**

```bash
> chkit delete deployment mimosa-warburg athamantis-gauss flora-onnes
Are you really want to delete to delete mimosa-warburg, athamantis-gauss, flora-onnes? [Y/N]: y
Deployment flora-onnes is deleted
Deployment mimosa-warburg is deleted
Deployment athamantis-gauss is deleted
3 deployments are deleted

```

**Delete list of deployments with --force**

```bash
> chkit delete depl --force lindenau-chayes malva-clarke pauwels-toepler
Deployment lindenau-chayes is deleted
Deployment malva-clarke is deleted
Deployment pauwels-toepler is deleted
3 deployments are deleted
```

**Delete deployment with interactive selection**

```bash
> chkit delete depl
Select deployment:
Selected:
 1) gurzhij-newton
 2) jackson-lenard
 3) kupe-magnus
 4) lindenau-chayes
 5) malva-clarke
 6) marchis-young
 7) pauwels-toepler
 8) rebentrost-thales
 9) Confirm
Choose wisely: 1
Select deployment:
Selected: gurzhij-newton
 1) jackson-lenard
 2) kupe-magnus
 3) lindenau-chayes
 4) malva-clarke
 5) marchis-young
 6) pauwels-toepler
 7) rebentrost-thales
 8) Confirm
Choose wisely: 2
Select deployment:
Selected: gurzhij-newton kupe-magnus
 1) jackson-lenard
 2) lindenau-chayes
 3) malva-clarke
 4) marchis-young
 5) pauwels-toepler
 6) rebentrost-thales
 7) Confirm
Choose wisely: 4
Select deployment:
Selected: gurzhij-newton kupe-magnus marchis-young
 1) jackson-lenard
 2) lindenau-chayes
 3) malva-clarke
 4) pauwels-toepler
 5) rebentrost-thales
 6) Confirm
Choose wisely: 1
Select deployment:
Selected: gurzhij-newton kupe-magnus marchis-young jackson-lenard
 1) lindenau-chayes
 2) malva-clarke
 3) pauwels-toepler
 4) rebentrost-thales
 5) Confirm
Choose wisely: 5
Are you really want to delete to delete gurzhij-newton, kupe-magnus, marchis-young, jackson-lenard? [Y/N]: y
Deployment gurzhij-newton is deleted
Deployment kupe-magnus is deleted
Deployment marchis-young is deleted
Deployment jackson-lenard is deleted
4 deployments are deleted

```