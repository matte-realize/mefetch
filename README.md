# Mefetch

Generate a Neofetch-style SVG card for your GitHub profile README.

The card renders as an SVG served by the app, so it stays up to date every time your profile is viewed — set the URL once and leave it.

## Features

- Neofetch-style SVG card that embeds directly in a GitHub README
- ASCII art generated from any uploaded image (logos, avatars, silhouettes)
- Customizable fields, sections, colors, and layout
- Optional GitHub stats — repos, commits, lines added/deleted
- Runs locally, no hosting required

## Quickstart

### Go

```bash
go install github.com/matte-realize/mefetch@latest
mefetch
```

### Clone and run

```bash
git clone https://github.com/matte-realize/mefetch
cd mefetch
cp .env.example .env   # optional, for a GITHUB_TOKEN
make run               # or: go run main.go
```

### Docker

```bash
docker build -t mefetch .
docker run -p 8080:8080 mefetch
```

## Usage

1. Run the app — your browser opens to `http://localhost:8080`
2. Fill in your details and customize the card
3. Upload an image to generate ASCII art
4. Copy the embed URL into your GitHub profile README

Set `PORT` to run on a different port.

## Embed in your README

1. Customize your card in the app and click **download svg**.
2. Add the downloaded `.svg` file to your repo (e.g. `assets/card.svg`).
3. Reference it in your README:

```html
<p align="center">
  <img src="assets/card.svg" alt="Mefetch Card">
</p>
```

This embeds a static export — no running server required.

## URL params

| param | description | example |
|---|---|---|
| `username` | your username | `name` |
| `hostname` | your hostname | `macbook` |
| `field` | custom field (repeatable) | `field=OS:macOS` |
| `background` | background color | `%230d1117` |
| `keycolor` | key text color | `%2358a6ff` |
| `textcolor` | value text color | `%23cdd9e5` |
| `showstats` | show GitHub stats | `false` to disable |

## ASCII art tips

- **Format** — PNG or JPEG
- **Shape** — square or portrait works best; landscape gets squished
- **Size** — 200×200 to 1000×1000px
- **Contrast** — high-contrast subjects on a simple background are clearest
- **Transparency** — transparent PNG logos are traced by their silhouette, so any color works; opaque images are converted by brightness

Your GitHub avatar is an ideal test image.

## GitHub API & rate limits

GitHub stats come from the public GitHub API. Unauthenticated requests are limited to **60 per hour per IP**. Mefetch caches each user's stats for an hour, so editing fields or repeated profile views cost only one fetch per username per hour.

Which budget is used depends on how it runs:

- **Local / self-hosted** — requests come from your own IP, against your own 60/hr. No setup needed.
- **One public instance** — every visitor loads the SVG from your server, so all requests share that server's single IP budget. Caching keeps this workable (~60 distinct usernames/hour); busy instances should use a token.

To raise the limit to **5,000/hour**, set a token in `.env`:

```bash
GITHUB_TOKEN=ghp_your_token_here
```

A classic Personal Access Token with **no scopes** is enough for public data. **Never commit your token** — keep it in `.env` (gitignored). Forks should supply their own.

## License

MIT — see [LICENSE](LICENSE).

## Demo Creation in Progress!