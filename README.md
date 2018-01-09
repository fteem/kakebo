# Kakebo

Kakebo, which literally means a household finance ledger, categorizes monthly
expenses in four different categories:

1. Survival
2. Optional
3. Culture
4. Extra

## Usage

### Monthly configuration

Set income for current month: ✅
```
kakebo income set 12345
```

Show income for current month: ✅
```
kakebo income show
```

Set savings target for current month:✅
```
kakebo target set 12345
```

Get savings for current month:
```
kakebo savings
# => 10000
```

### Expenses

Add expense:✅
```
kakebo expenses add --description=Whatever --amount=123 --category=survival --week=2
```

Show all expenses for current week:✅
```
kakebo expenses list
```

Show all expenses for week:✅
```
kakebo expenses list --week=3
```

All of the `--week` arguments can be omitted. In those cases, the number of the
current week is used.
