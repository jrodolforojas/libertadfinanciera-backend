# Libertad Financiera API

## About

Web scrapper that gets the **************dollar************** ➡️ ****************colones**************** exchange rates from the [Banco Central de Costa Rica](https://gee.bccr.fi.cr/indicadoreseconomicos/Cuadros/frmVerCatCuadro.aspx?idioma=1&CodCuadro=%20400) website. So with that, you can see the exchange rates across history. 😄 🇨🇷

## How to run?

Clone the repository.

```bash
git clone https://github.com/jrodolforojas/libertadfinanciera-backend.git
```

After that, you need to download the dependencies.

```go
go mod download
```

Once you have the dependencies execute the `main.go` file

```go
go run main.go
```

## Endpoints

### GET `/exchange_rates`

Returns the **************dollar************** ➡️ ****************colones**************** exchange rates. By default, it returns the latest ******30 days****** exchange rates from `today`

******************Response example******************

```json
{
    "data": [
        {
            "sale": 555.47,
            "buy": 548.2,
            "date": "2023-07-09T12:00:00Z"
        },
        {
            "sale": 555.47,
            "buy": 548.2,
            "date": "2023-07-10T12:00:00Z"
        },
        {
            "sale": 556.9,
            "buy": 547.39,
            "date": "2023-07-11T12:00:00Z"
        }
    ]
}
```

Also, **you can filter this endpoint by `date range`**. For example, get the exchange rates from December 04 of 2022 to January 19 of 2023 using `query params`

************Params************

1. `date_from`
    1. Format: `2023/12/04`
2. `date_to`
    1. Format: `2023/12/04`

**************Example**************

```bash
/exchange_rates?date_from=1984/01/01&date_to=1985/05/12
```

### GET `/exchange_rates/today`

Returns today **************dollar************** ➡️ ****************colones**************** exchange rate

********************************Response example********************************

```json
{
    "data": {
        "sale": 544.82,
        "buy": 538.86,
        "date": "2023-08-07T12:00:00Z"
    }
}
```
