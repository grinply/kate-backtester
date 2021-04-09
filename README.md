# Quick Backtest
A fast and simple backtest implementation for **algorithmic trading** focused on [cryptocurrencies](https://en.wikipedia.org/wiki/Cryptocurrency#:~:text=A%20cryptocurrency%2C%20crypto%20currency%20or,creation%20of%20additional%20coins%2C%20and) written in golang.

## Data
The price data used to run the backtests can be from any time interval, but it must contain a [**OHLCV**](https://en.wikipedia.org/wiki/Open-high-low-close_chart) structure _(**O**pen **H**igh **L**ow **C**lose **V**olume)_. It is possible to load data from **csv** files and the [**postgresql** database](https://www.postgresql.org/).
