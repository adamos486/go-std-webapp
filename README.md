# Go Seminar

## SOLID & DI

- SOLID principles: [Dave Cheney's Blog](https://dave.cheney.net/2016/08/20/solid-go-design)
- What is **SOLID**?
  1. Single Responsibility Principle
  1. Open / Closed Principle: Embedding for the win.
  1. Liskov Substitution Principle
  1. Interface Segregation Principle
  1. Dependency Inversion Principle

### Single Responsibility Principle

"A "class"(package) should only have one, and only one, reason to change" - Robert C. Martin,

- Eases maintenance, speeds up development, reduces coupling.
- Encapsulate & abstract whatever is possible, in Goland place things behind interfaces, segregate third party dependencies whenever possible.
- Package names can be a responsibility smell. If a package name is contains it's purpose AND a coherent namespace prefix, then it's likely singly responsible for what it describes. However a package named something like `server` could be... any server protocol, any logic pertaining to running a server, constructing a server only or also running it? (it's unclear and possibly a sign that the package is doing too much.)
- Go's designers were heavily influenced by the UNIX philosophy of decoupled design, meaning they envisioned each Go package being itself a small Go program & a single unit of change with a single purpose.

### Open / Closed Principle

- When extending a struct embed it if you need it's properties (fields AND funcs).
- This will not only allow you to "extends" a struct ala Java, but also to override.

```golang
type Cat struct {
        Name string
}

func (c Cat) Legs() int { return 4 }

func (c Cat) PrintLegs() {
        fmt.Printf("I have %d legs\n", c.Legs())
}

type OctoCat struct {
        Cat
}

func (o OctoCat) Legs() int { return 5 }

func main() {
        var octo OctoCat
        fmt.Println(octo.Legs()) // 5
        octo.PrintLegs()         // I have 4 legs
}
```

- Embedding allows you to assert that each logic type you create is **open for extension** BUT **closed for modification**
- OctoCats can never be regular Cats. Sorry.

### Liskov Substitution Principle

- Since Go has no base classes or a concept of inheritance, substitution cannot be implemented in a class hiearchy.
- Create interfaces often and in bite sized chunks based on groupings of related work/funcs.

```golang
type Reader interface {
        // Read reads up to len(buf) bytes into buf.
        Read(buf []byte) (n int, err error)
}
```

- Well written `interfaces` producing streamlined implementations, can be substituted in test, as a way to generify work and reuse code, and even for refactorings or replacements without significant rewrites.

### Interface Segregation Principle

- This extends the above principle and expands more on streamlining and optimizing interfaces.
- Define interfaces around only the logic you need at the time of implementation.
- This has several good side effects:
  1. Prevent abuse or hacking on real client/struct implementations that define many other funcs.
  1. Prevent over-engineering.
  1. Ease of understanding for future maintainers and extenders, reducing noise of giant helper clients.
  1. Forces you to think... is it worth it to change this interface and break all existing implementations or can I create a new isolated interface? Reinforcing either your good reasons for changing business logic, or making your system even more modular.

```golang
// Save writes the contents of doc to the file f.
func Save(f *os.File, doc *Document) error
```

can become -->

```golang
// Save writes the contents of doc to the supplied ReadWriterCloser.
func Save(rwc io.ReadWriteCloser, doc *Document) error
```

which will prevent hacking on os.File methods inside an implementation of Save that changes the scope/intention of the Save func.

### Dependency Inversion Principle

- If you've done SOLI, you're ready for D.
- "Your code should describe its dependencies in terms of interfaces, and those interfaces should be factored to describe only the behaviour those functions require."
- Your import graph becomes broad and flat and faster to resolve, versus narrow and very deep.
- You can easily replace implementations with mocked or faked implementations in a test environment.

```golang
  func setupBoilerplate() {
    fakeConsumer = &consumerfakes.FakeInterface{}
    fakeLog = &workererrorfakes.FakeFailureLog{}
    fakeBackend = &backendfakes.FakeBackender{}
    hardDeleteHandler, _ = task.NewHardDeleteHandler("some test topic", fakeBackend, fakeLog)
    workerConfig, _ = config.NewWorkerConfig(
      "some-queue-name",
      "some region",
      2,
      5,
      hardDeleteHandler,
      fakeLog,
    )
  }
```

## Testing

- Ginkgo test runner, Gomega matcher.
  - Benefits include parallelized tests, goroutine/async handling, randomizeAllSpecs, test watcher, test focuser.
    - Run until false, run forever, AfterEach tear down blocks (useful for integration test cleanup), code coverage, convert, benchmarking.
    - Agouti for acceptance testing of statically served pages.
    - GoConvey has now stopped being supported, so Ginkgo, Goblin, and Zen are now the most popular BDD test frameworks.
  - Gomega is readable, more easily understandable, more intelligent matching (i.e. Equal, ConsistOf, Eventually, Consistently)
    - Register Fail Handler at the Suite level, so you don't have to pass a `(t *testing.T)` into every test context.

```golang
Describe("Worker Client Specs", func() {
  //Behavioral descriptor of logic flow
  Describe("worker startup", func() {
    //Behavioral test block
    Context("when it has a valid config", func() {

      //Context scoped variables
      var (
        config *config.WorkerConfig
        worker *worker.Worker
        started bool
        expectedErr error
      )

      //Sets up fresh scoped variables before each test and assertion block
      BeforeEach(func() {
        config = config.New()
        worker = worker.New(config)
      })

      //Runs the exact function under test
      JustBeforeEach(func() {
        started, expectedErr = worker.Start()
      })

      It("should start a worker and return no error", func() {
        Expect(expectedErr).NotTo(HaveOccurred())
        Expect(started).To(BeTrue())
      })
    })
  })
})
```

## Local Dev, Deployment, Acceptance, and Release Integration

- We spin up service and dependent collaborators completely locally. (i.e. local sqs instance?)
- Runs 100% locally the same way it'll work on a box. (i.e. Kubes, Docks)
- Validating contracts between different services locally. (i.e. Validating an sqs message across both sessions service & user service) - currently relies on a canned copy-pasta sqs message definition (this could be invalidated in the future if User-Service emits a differently formed message)
  - Example: a locally running User Service with the latest green binary, could emit it's own hard delete event to be received by my local Passport/Sessions service.
- In order to test features, hitting local versions with sandboxed databases/services (aka Passport sandbox), has been very valuable as an isolated low risk acceptance, but also reduce test pollution on each staging environment.
- In cases of callaborating with many microservices, a test environment for full end-to-end acceptance testing has been very helpful.
  - Includes an ability to spin up local, latest green builds of each interactive service
  - Runs end to end test with bleeding edge changes on the work you've done that interacts with a multitude of spun up services.
  - Gives you the ability to detect errors with other collaborators before deployment, or have confidence that you aren't creating them.
  - [CloudFoundry Acceptance Tests](https://github.com/cloudfoundry/cf-acceptance-tests) are how we perform this across our entire Cloud Foundry product portfolio.
- Standardized dependency management:
  - Glide has trouble resolving binaries, resolving strangely nested project structures, and resolving conflicting and/or missing sub-dependencies.
  - Dep allows for project creation based on Glide.yaml & Glide.lock resolving sub-dependency conflicts, handling binaries and versions more automatically.
  - For our products, Dep has allowed for more reproducible builds, and dependable local project setup.
