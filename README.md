# ChAIAF
#### Chairs Automated Inventory for Apartment Furnishing

## Original requirements

Apartment And Chair Delivery Limited has a unique position on the housing market. The company not only builds apartments, but also equips them with chairs.
Now the business has grown continuously over the past few years and there are a few organizational problems that could be solved by automation.
We will focus on one of them here:

While a new residential building is erected, the chairs that are to be placed there need to be produced. In order to be able to plan this, the home buyers indicate the desired position of the armchairs in their home on a floor plan at the time of purchase. These plans are collected, and the number of different chairs to be produced are counted from them. The plans are also used to steer the workers carrying the chairs into the building when furnishing the apartments.

In the recent past, when manually counting the various types of chairs in the floor plans, many mistakes were made and caused great resentment among customers. That is why the owner of the company asked us to automate this process.

Unfortunately, the plans are in a very old format (the company's systems are still from the eighties), so modern planning software cannot be used here. An example of such an apartment plan is attached.

We now need a command line tool that reads in such a file and outputs the following information:
- Number of different chair types for the apartment
- Number of different chair types per room

The different types of chairs are as follows:
W: wooden chair
P: plastic chair
S: sofa chair
C: china chair

The output must look like this so that it can be read in with the old system:

total:
W: 3, P: 2, S: 0, C: 0
living room:
W: 3, P: 0, S: 0, C: 0
office:
W: 0, P: 2, S: 0, C: 0

The names of the rooms must be sorted alphabetically in the output.

Our sales team has promised Apartment And Chair Delivery Limited a solution within 5 days from now. I know that is very ambitious, but as you are our best developer, we all count on you.

## Implementation details

### Usage

The simple version is

```shell
make example
```
or
```shell
go run main.go rooms.txt
```

For anyone interested in more than that, please consider the contents of the Makefile:

```makefile
VERSION ?= $(shell git describe --tags 2> /dev/null || echo v0)

GO = go
BINARY = room-parser

.PHONY: build clean test
.DEFAULT_GOAL := help

example: ## run the example
	${GO} run main.go rooms.txt

build: ## builds
	@${GO} build -ldflags "-X main.Version=${VERSION}" -o ${BINARY}
	@echo "${BINARY} built. Run it like this:\n\n\t./${BINARY} rooms.txt"

test: ## runs the unit tests
	${GO} test -v ./...

clean: ## go clean, then remove any previously built binary
	${GO} clean
	rm -f ${BINARY}

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
```

You can just run make in the project path and it will guide you from there.

### Code

The parser reads line by line from a Reader (can be anything, not just the rooms.txt file).
It records the last known positions of the walls of each open room, so then, when the next line is read,
the parser computes overlaps and knows which "line segment" belongs to which room.
But what the parser is really after is the closed rooms, which will are listed in the end.

### The Choice

I chose this line-by-line approach over a flood-fill algorithm because I wanted to avoid the situation where I load the whole flat in memory.
Then gradually I realised that wouldn't have been so bad either, since flood-fill is used in image processing at a much larger scale, and it works just fine.
Plus even the flood-fill could have been done on chunks of lines rather than the whole thing at once.

### Benchmark

```text
goos: darwin
goarch: arm64
pkg: enspired/src
BenchmarkRoomParser_Ingest
BenchmarkRoomParser_Ingest-12    	   13254	     90856 ns/op
```

Memory can be optimized for sure (for example, segments contain full strings currently, rather than length).
My initial plan was to compare this with a flood-fill implementation, but I'm not that curious anymore :D

### Tests

src/ 80% coverage

### Further steps
* more test cases
* further cpu and memory profiling
* make the api more robust
* more validation (e.g. rooms that don't close)

### Why it took so long
It wasn't just that I was busy with other stuff besides this, 
but I also had to dust off my golang skills after this long gap year, 
and I tried to make the most out of this challenge.

### Contact

for any question please [mail me](mailto:florinutz@gmail.com).

Florin