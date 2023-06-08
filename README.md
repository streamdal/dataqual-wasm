# dataqual-wasm

Contains sources for WASM modules used in the plumber data rules SDK.

Also contains benchmarks for each golang shim library.

### Build process:

1. Make changes to module under `src/` directory
2. Run `make build` to ensure build succeeds
3. Push PR to github and merge
4. `git tag v0.1.0 main`
5. `git push origin v0.1.0`
6. Github actions will automatically compile and create release using the tag version

### Building WASM files manually

```bash
$ make build
```

### Running benchmarks

```bash
$ make test/bench

goos: darwin
goarch: arm64
BenchmarkRabbit_Produce
BenchmarkRabbit_Produce/MatchRuleAndDiscard
BenchmarkRabbit_Produce/MatchRuleAndDiscard-8         	  116145	      9850 ns/op
BenchmarkRabbit_Produce/NoRuleMatch
BenchmarkRabbit_Produce/NoRuleMatch-8                 	   52994	     19543 ns/op
BenchmarkRabbit_Produce/NoDataqual
BenchmarkRabbit_Produce/NoDataqual-8                  	  203442	     24670 ns/op
BenchmarkSegmentio_Produce
BenchmarkSegmentio_Produce/MatchRuleAndDiscard
BenchmarkSegmentio_Produce/MatchRuleAndDiscard-8      	  114502	     10063 ns/op
BenchmarkSegmentio_Produce/NoRuleMatch
BenchmarkSegmentio_Produce/NoRuleMatch-8              	    1729	    588310 ns/op
BenchmarkSegmentio_Produce/NoDataqual
BenchmarkSegmentio_Produce/NoDataqual-8               	    2151	    519059 ns/op
BenchmarkProduce_Sarama
BenchmarkProduce_Sarama/MatchRuleAndDiscard
BenchmarkProduce_Sarama/MatchRuleAndDiscard-8         	   82731	     14171 ns/op
BenchmarkProduce_Sarama/NoRuleMatch
BenchmarkProduce_Sarama/NoRuleMatch-8                 	    1969	    584926 ns/op
BenchmarkProduce_Sarama/NoDataqual
BenchmarkProduce_Sarama/NoDataqual-8                  	    2272	    516012 ns/op
```