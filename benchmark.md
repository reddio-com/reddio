### Transfer benchmark Tutorial

1. Start the reddio server
2. build the benchmark cmd by: `make build_benchmark_test`
3. prepare the benchmark data by: `./benchmark_test --action=prepare  --preCreateWallets=<num>`
4. run the benchmark test by: `./benchmark_test --action=run --qps=10`