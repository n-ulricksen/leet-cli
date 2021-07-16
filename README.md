# lcfetch

Program used to retrieve Leetcode problems in the terminal.

### Check it out

![demo](media/lcfetch.gif)

### Synopsis

Get a random problem from Leetcode based on difficulty and/or topic.

```
lcfetch [flags]
```

Examples:

```
lcfetch -d hard -t dynamic-programming
lcfetch -d medium -t array,two-pointers
```

### Options

```
      --config string       config file (default is $HOME/.lcfetch.yaml)
  -d, --difficulty string   difficulty of problem to select (default "all")
  -h, --help                help for lcfetch
  -p, --paid                include paid/premium problems
  -t, --topics strings      topic(s) to select problem from (comma-separated, no spaces)
```

## Commands

### lcfetch complete

Mark one or more problems complete, prevening them from showing up when requesting
a random problem.

```
lcfetch complete [flags]
```

Examples:

```
lcfetch complete 1337
lcfetch complete 52 12 628
```

#### Options

```
  -h, --help   help for complete
```

---

### lcfetch incomplete

Mark more ore more problems, allowing them to show up when requesting a random problem.

```
lcfetch incomplete [flags]
```

Examples:

```
lcfetch incomplete 1337
lcfetch incomplete 628 12 52
```

#### Options

```
  -h, --help   help for incomplete
```

---

### lcfetch stats

Print details about completed questions per category and difficulty.

```
lcfetch stats [flags]
```

#### Options

```
  -d, --difficulty string   difficulty of problems to print with stats (default "all")
  -h, --help                help for stats
  -p, --paid                include paid/premium questions
  -t, --topic string        topic of problems to print with stats (comma-separated, no spaces)
```

---

### lcfetch list

Print a list of the Leetcode problems, filtered by difficulty and/or topic.

```
lcfetch list [flags]
```

Examples:

```
lcfetch list
lcfetch list -d easy -t array,string
```

#### Options

```
  -c, --completed           list only completed problems
  -d, --difficulty string   difficulty of problems to list (default "all")
  -h, --help                help for list
  -i, --incomplete          list only incomplete problems
  -p, --paid                include paid/premium problems
  -t, --topics strings      topic(s) of problems to list (comma-separated, no spaces)
```

---

### lcfetch get

Lookup one or more problems by Leetcode ID.

```
lcfetch get [flags]
```

Examples:

```
lcfetch get 521
lcfetch get 72 1262 980
```

#### Options

```
  -h, --help   help for get
```

---

### lcfetch topics

List all problem topics on Leetcode.

```
lcfetch topics [flags]
```

#### Options

```
  -h, --help   help for topics
```
