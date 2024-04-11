# `FindIt`

Locate strings within large collections of files and directories!

# How to use

To use `FindIt`, simply `cd` into the directory you would like to search, and run:

```bash
$ findit "my string"
```

`FindIt` will... `FindIt` all.

`FindIt` will also provide samples of where the string was found, along with "context". Meaning it just prints a certain number of characters around the found string. The default is 20 characters on each side, but can be changed with `-c` or `--context-size`:

```bash
$ findit "my string" -t 400 # Prints the whole line unless it's a file like minified CSS or JS
```

If you don't want a context size, then you would use the value `-1`. I say would because I haven't implemented it yet.

## Other Cool Things

If you want to search a directory but don't want to `cd` into it, then supply
the wanted directory as the second argument:

```bash
$ findit "my string" "../../../my_apps_assets_or_something"
```

Only want to check files with a specific extension?

Do this:

```bash
findit "my string" -t ".sh" # Will only look in shell files
```

Too specific of an extension for your needs? Well, it supports wildcards:
```bash
findit "my string" -t "*h" # Will search any file with an extension that ends in 'h'
```

Or:
```bash
findit "my string" -t ".s*" # Will search any file with an extension that starts with '.s'
```

`FindIt` by default will only scan files it deems non binary. If you want to scan binary files (Not recommended), you can use the switch `-b` or `--allow-binary`.
If for some ungodly reason you want to search *only* binary files, you can use the switch `--only-binary` (even not recommendeder).

`FindIt` also uses color based on an assumption. If you don't want color when it's supported, or vise versa, you can set it the same way as other built in tools like `grep` (Best done via an alias):

```bash
$ findit "something" --color=<whatever_you_want>(never,always,auto) # `FindIt` defaults to auto
```

If you need more information, you can use:

```bash
$ findit --help
```

# How to Make

```bash
git clone https://github.com/ScripturaOpus/FindIt.git
cd FindIt
make
```

And then `FindIt` will be built to `./build/findit`
