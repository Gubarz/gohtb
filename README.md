# gohtb

<p align="center">
  <img src="https://raw.githubusercontent.com/gubarz/gohtb/main/assets/logo.png" width="200" alt="gohtb logo" />
</p>

> A clean and easy to use Go SDK for the Hack The Box API.

## Install
```bash
go get github.com/gubarz/gohtb@latest
```

## Quick Start
```go
client, _ := gohtb.New(os.Getenv("HTB_TOKEN"))

machines, _ := client.Machines.ListActive().ByOS("Linux").Results(ctx)

for _, m := range machines.Data.ByDifficulty("Easy") {
    fmt.Println(m.Name)
}
```

## Examples

See [examples](examples/) for a curated list of examples on how to use.

## Why This Was Created

I built `gohtb` because I want a reliable, reusable foundation for working with the Hack The Box API in Go. While building tools on top of the platform, I realized there was no structured SDK or up-to-date reference material available.

Rather than duplicate logic across projects, `gohtb` was designed to:

- Provide a clean, idiomatic Go interface for Hack The Box's API
- Standardize interactions with machines, challenges, VPN servers, and more
- Simplify flag submission, obtaining product details, VPN configuration, and automation workflows
- Serve as a foundation for building tools, scripts, or bots that integrate with HTB

`gohtb` aims to simplify development, reduce boilerplate, and make HTB automation approachable for anyone working in Go.

## API Specification

This SDK is based on [UnOfficial HackTheBox OpenAPI spec](https://github.com/gubarz/unofficial-htb-api).
It uses [oapi-codegen](https://github.com/deepmap/oapi-codegen) for base generation.

## Contributions

Contributions are welcome! If you’d like to add features, improve documentation, or report bugs, feel free to open an issue or submit a pull request.

Ideas, feedback, and feature suggestions are always appreciated. This project exists to make HTB development in Go easier for everyone, your input helps make that possible.

## License
Licensed under [Unlicense](LICENSE)

---

💚 Found this SDK useful? Support the project: https://ko-fi.com/gubarz