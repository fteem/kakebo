# Kakebo

Kakebo, which literally means a household finance ledger, categorizes monthly
expenses in four different categories:

1. Survival
2. Optional
3. Culture
4. Extra

## Usage

### Income

Set income for current month:
```
kakebo income set 12345
```

Or you can set the income for any month/year combo, using:
```
kakebo income set 2000 --year=2018 --month=January
```

Show income for current month:
```
kakebo income show
```

Or you can see the income for any month/year combo, using:
```
kakebo income show --year=2018 --month=January
```

### Goals

Get savings goal for current month:
```
kakebo goal show
```

Or you can see the goal for any month/year combo, using:
```
kakebo goal show --year=2018 --month=January
```

Set savings goal for current month:
```
kakebo goal set 1000
```

Or you can set the goal for any month/year combo, using:
```
kakebo goal set 1000 --year=2018 --month=January
```

### Expenses

Add expense for current month:
```
kakebo expenses add --description=Whatever --amount=123 --category=survival
```

You can add an expense for a specific month/year by adding the `--year` and `--month` flags:
```
kakebo expenses add --description=Whatever --amount=123 --category=survival --month=January --year=2018
```

Show all expenses for current month:
```
kakebo expenses list
```

Show all expenses for specific month/year combination:
```
kakebo expenses list --month=January --year=2018
```

### Savings

Get savings for current month:
```
kakebo savings
# => 10000
```

Get savings for specific month/year combination:
```
kakebo savings --month=January --year=2018
```
