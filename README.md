# sanity

a quick golang cli tool i made for myself to automate some tasks i do often when organizing CTFs and/or creating CTF challenges

## usage
```
sanity [command] [args]
```

## usage (check)

![a](https://i.gyazo.com/cb4c4586c86f3e85fc7864c78449b07a.png)

checks to make sure u didn't leave any flags in dist files given to participants

```
sanity check --flag_format FLAG{ --zip ./tests/check_zip_test/dist.zip
```

## usage (service)

![a](![a](https://i.gyazo.com/cd6399180bd7e49ddf21e1a769b4f31e.png))

for challenges that have service files (web, pwn, whatever), you normally want to give src to participants; we just package it and check whether we left any flags behind.

```
sanity service --flag_format FLAG --input ./tests/service_test --zip dist.zip --name challenge --cleanup --batch
```