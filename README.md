<p align="center">
    <img height="280" width="207" src="docs/slurp.png">
</p>



# Slurp 
Building with Go, easier than a slurp.

[![GoDoc](https://godoc.org/github.com/omeid/slurp?status.svg)](https://godoc.org/github.com/omeid/slurp)

> Heads up! This is pre-release software.  

Slurp is a [Gulp.js](http://gulpjs.com/) inspired build toolkit designed with idiomatic Go [Pipelines](http://blog.golang.org/pipelines) and following principles: 

- Convention over configuration
- Explicit is better than implicit.
- Do one thing. Do it well.
- ...


> ### Why?
  The tale of Gulp, Go, Go Templates, CSS, CoffeeScript, and minifiaction and building assets should go here.



Slurp is made of two integral parts:

### 1. The Toolkit 

The slurp toolkit provides a task harness that you can register tasks and dependencies, you can then run these tasks with slurp runner.

A task is any function that accepts a pointer to `slurp.C` (Slurp Context) and returns an error.  
The Context provides logging functions. _it may be extended in the future_.

```go
b.Task("example-task", []string{"list", "of", "dependency", "tasks"},

  func(c *slurp.C) error {
    c.Info("Hello from example-task!")
  },

)
```

Following the Convention Over Configuration paradigm, slurps provides you with a collection of nimble tools to instrument a pipeline.

A pipeline is created by a source _stage_ and typically piped to subsequent _transformation_ stages and a final _destination_ stage.

Currently Slurp provides two source stages `slurp/stages/fs` and `slurp/stages/web` that provide access to file-system and http source respectively.

```go
b.Task("example-task-with-pipeline", nil , func(c *slurp.C) error {
    //Read .tpl files from frontend/template.
    return fs.Src(c, "frontend/template/*.tpl").Pipe(
      //Compile them.
      template.HTML(c, TemplateData),
      //Write the result to disk.
      fs.Dest(c, "./public"),
    ).Wait() //Wait for all to finish.
})
```

or the same code shorter

```go
b.Task("example-task-with-pipeline", nil , func(c *slurp.C) error {
    //Read .tpl files from frontend/template.
    return fs.Src(c, "frontend/template/*.tpl").Then(
      //Compile them.
      template.HTML(c, TemplateData),
      //Write the result to disk.
      fs.Dest(c, "./public"),
      )
})
```

and another example,

```go
// Download deps.
b.Task("deps", nil, func(c *slurp.C) error {
    return web.Get(c,
      "https://github.com/twbs/bootstrap/archive/v3.3.2.zip",
      "https://github.com/FortAwesome/Font-Awesome/archive/v4.3.0.zip",
    ).Then(
      archive.Unzip(c),
      fs.Dest(c, "./frontend/libs/"),
    )

})
```

Currently the following _stages_ are provided with Slurp:
> No 3rd party dependency, just standard library.  

- [archive](https://godoc.org/github.com/omeid/slurp/stages/archive/)
- [fs](https://godoc.org/github.com/omeid/slurp/stages/fs/)
- [passthrough](https://godoc.org/github.com/omeid/slurp/stages/passthrough/)
- [template](https://godoc.org/github.com/omeid/slurp/stages/template/)
- [web](https://godoc.org/github.com/omeid/slurp/stages/web/)


You can find more at [slurp-contrib](https://github.com/slurp-contrib). gin, gcss, ace, watch, resources (embed), and livereload to name a few.


### 2. The Runner (cmd/slurp)

This is a cli tool that runs and help you compile your builders. It is go get’able and you can install with:

```bash
 $ go get -u -v github.com/omeid/slurp/cmd/slurp  # get it.
```

Slurp uses the Slurp build tag. That is, it passes `-tags=slurp` to go tooling when building or running your project,
this allows decoupling of build and project code. This means you can use Go tools just like you're used to, even if your
project has a slurp file.

Somewhat similar to `go test` Slurp expects a `Slurp(*slurp.Build)` function from your project, this is typically put in a file with the `// +build slurp` build tag.

### Demo 

![demo](docs/demo.gif)

### Contributing

Please see [Contributing](CONTRIBUTING.md)


### Examples
 - The obligatory [Todo App (Slurp)](https://github.com/omeid/slurp-todo)


### LICENSE
  [MIT](LICENSE).
