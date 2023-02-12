# Gazelle Weekly

A small script sending a weekly email summary of the Top 10 entries on a Gazelle-based site. The idea is to replicate the discovery experience of a website like that in the age of streaming.

## Usage

Set the environment variables, make sure the `FROM_EMAIL` is your verified sender domain on [Postmark](https://postmarkapp.com).

```
export POSTMARK_API_TOKEN="1-2-3-4"
export GAZELLE_API_TOKEN="5-6-7-8"
export GAZELLE_BASE_URL="example.com"
export FROM_EMAIL="mail@example.com"
export TO_EMAIL="mail@example.com"
```

Run `go build` and execute the binary in a cron job, with the correct environment variables set.

## Screenshot

