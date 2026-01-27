package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    HttpRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )

    HttpRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "Duration of HTTP requests in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "endpoint"},
    )

    TransactionsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "transactions_total",
            Help: "Total number of transactions",
        },
        []string{"type", "status"},
    )

    DbQueryDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "db_query_duration_seconds",
            Help:    "Duration of database queries",
            Buckets: prometheus.DefBuckets,
        },
        []string{"query_type"},
    )

    ActiveUsers = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "active_users_total",
            Help: "Number of active users",
        },
    )

    TotalPortfolioValue = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "total_portfolio_value",
            Help: "Total value of all user portfolios",
        },
    )
)