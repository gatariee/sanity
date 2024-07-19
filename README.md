# sanity

a quick golang cli tool i made for myself to automate some tasks i do often when organizing CTFs and/or creating CTF challenges

![a](https://i.gyazo.com/cd6399180bd7e49ddf21e1a769b4f31e.png)

## usage
```
sanity [command] [args]
```

## usage (service)

for challenges that have service files (web, pwn, whatever), you normally want to give src to participants; we just package it and check whether we left any flags behind.

```
sanity_x64 service --flag_format FLAG --input ./tests/service_test --zip dist.zip --name challenge --cleanup --batch
```