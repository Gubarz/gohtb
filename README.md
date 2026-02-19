# gohtb

<p align="center">
  <img src="https://raw.githubusercontent.com/gubarz/gohtb/main/assets/logo.png" width="200" alt="gohtb logo" />
</p>

> A clean and easy to use Go SDK for the Hack The Box API.

## Install
```bash
go get github.com/gubarz/gohtb@latest
```

## Requirements
- Go `1.24+`
- A valid HTB API token - create one [here](https://app.hackthebox.com/account-settings)

## Quick Start
```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gubarz/gohtb"
)

func main() {
	client, err := gohtb.New(os.Getenv("HTB_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	machines, err := client.Machines.
		List().
		ByState("active").
		ByDifficulty("easy").
		PerPage(20).
		Results(ctx)
	if err != nil {
		if apiErr, ok := gohtb.AsAPIError(err); ok {
			log.Fatalf("api error: status=%d message=%s", apiErr.StatusCode, apiErr.Message)
		}
		log.Fatal(err)
	}

	for _, m := range machines.Data {
		fmt.Println(m.Name)
	}
}
```

## Service Coverage

Services currently exposed:

- `Challenges`
- `Containers`
- `Fortresses`
- `Machines`
- `Prolabs`
- `Rankings`
- `Search`
- `Seasons`
- `Sherlocks`
- `Teams`
- `Tracks`
- `Users`
- `VMs`
- `VPN`

Current known missing wrappers:
- `Home` dashboard endpoints
- `Platform` endpoints
- `Reviews` endpoints
- `Universities` endpoints
- `Users` related ranking and stats endpoints

Use `client.Experimental()` for unsupported endpoints.

## Query Builder Style

The SDK uses query builders for discoverability and chaining in IDEs.

```go
results, err := client.Challenges.
	List().
	ByDifficulty("hard").
	ByState("active").
	Page(1).
	PerPage(25).
	Results(ctx)
if err != nil {
	log.Fatal(err)
}
fmt.Println(len(results.Data))
```

## Experimental

For endpoints not wrapped yet, you can call generated clients directly:

```go
exp := client.Experimental()
ctx := exp.WrapContext(context.Background())

resp, err := exp.V4().GetSeasonList(ctx)

defer resp.Body.Close()

bodyBytes, err := io.ReadAll(resp.Body)

if err != nil {
	log.Fatal(err)
}
fmt.Println(resp.StatusCode, string(bodyBytes))
```

`Experimental()` is intentionally unstable and not covered by compatibility guarantees.

## Rate Limiting and Retries

When using the default internal HTTP transport:

- Uses `X-Ratelimit-*` headers when present
- Applies a global `10s` pause on Cloudflare `429` responses
- Retries up to `4` times (max `5` total attempts) with exponential backoff + jitter

If you provide `WithHTTPClient(...)`, internal transport behavior (rate limiting/retries) is bypassed unless your custom client transport implements it.

## Errors and Response Metadata

Most service responses include `ResponseMeta`:

- `StatusCode`
- `Headers`
- `CFRay`
- `Raw` body

Errors can be unwrapped as `*gohtb.APIError`:

```go
if apiErr, ok := gohtb.AsAPIError(err); ok {
	fmt.Printf("status=%d message=%s\n", apiErr.StatusCode, apiErr.Message)
}
```

## Stability and Versioning

- This project is pre-`v1.0.0`.
- Breaking changes can happen between minor releases while the API surface evolves.
- HTB may/can/will make breaking changes at any time as the api is not versioned.
- `Experimental()` is explicitly unstable.

## Examples

- [challenges](examples/challenges)
- [fortresses](examples/fortresses)
- [machines](examples/machines)
- [seasons](examples/seasons)
- [sherlocks](examples/sherlocks)
- [teams](examples/teams)
- [users](examples/users)
- [vpn](examples/vpn)

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