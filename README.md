# Crypto Price Alert

This app will send you alert if crypto hit a price that you want once a day. Using JSON file to determine which crypto you want to watch, what price that you want to watch, what email that you want to receive the email.

## How To
- Clone this repo
- Copy `.env.example` to `.env`
- Create `alert.json` and fill the file with array of data below
- Create `alias.json` and fill the file with object of data below. This must be filled if you're not using default price vendor and `pair` value is different than default price vendor `pair` value
- Build the application with `go build main.go`
- Run the app using cron

## Config

### Price vendor
- [Indodax](https://github.com/btcid/indodax-official-api-docs/blob/master/Public-RestAPI.md#ticker) (Default)
- [CoinGecko](https://www.coingecko.com/en/api/documentation)

### Mail vendor
- [Mailgun](https://www.mailgun.com)
  - Fill `MAILGUN_SENDER`, `MAILGUN_DOMAIN` and `MAILGUN_API_KEY` in `.env` file

### JSON file

#### `alert.json` Structure
```json
[
 {
  "email": "example@example.com",
  "ticker": "HNST",
  "pair": "hnstidr",
  "price": 130,
  "emailed_at": null
 }
]
```

#### `alias.json` Structure
```json
{
  "coingecko": {
    "hnstidr": "honest-mining"
  }
}
```

## Notes
- Please run `go build main.go` every time you changes `alert.json` file

## TODO
- [x] Add coingecko as price vendor
- [ ] Create `tickers.json` to save supported ticker.
- [ ] Create `{ticker}.json` file to save email for every ticker
- [ ] Explore another way to send alert
- [ ] Add alert when price above/below price
- [ ] Add testing
- [ ] Add multiple currency
- [ ] Setup database
- [ ] Add sendgrid as mail vendor