![slardar](image/broodmother-big.jpeg "broodmother")

# Broodmother

A general solution for code generation.


For centuries, Black Arachnia the Broodmother lurked in the dark lava tubes beneath the smoldering caldera of Mount Pyrotheos, raising millions of spiderlings in safety before sending them to find prey in the wide world above. In a later age, the Vizier of Greed, Ptholopthales, erected his lodestone ziggurat on the slopes of the dead volcano, knowing that any looters who sought his magnetic wealth must survive the spider-haunted passages. After millennia of maternal peace, Black Arachnia found herself beset by a steady trickle of furfeet and cutpurses, bold knights and noble youths--all of them delicious, certainly, and yet tending to create a less than nurturing environment for her innocent offspring. Tiring of the intrusions, she paid a visit to Ptholopthales; and when he proved unwilling to discuss a compromise, she wrapped the Vizier in silk and set him aside to be the centerpiece of a special birthday feast. Unfortunately, the absence of the Magnetic Ziggurat's master merely emboldened a new generation of intruders. When one of her newborns was trodden underfoot by a clumsy adventurer, she reached the end of her silken rope. Broodmother headed for the surface, declaring her intent to rid the world of each and every possible invader, down to the last Hero if necessary, until she could ensure her nursery might once more be a safe and wholesome environment for her precious spiderspawn.

## Getting Started

First make sure that you have [created and added your SSH public key in gitlab](https://docs.gitlab.com/ee/gitlab-basics/create-your-ssh-keys.html).


For cloning the project with `go get` dependency management ecosystem, first you need to change your git configuration.
`go get` uses git internally. The following one liners will make git and consequently `go get` clone your package via SSH:

```bash
$ git config --global url."git@gitlab.com:".insteadOf "https://gitlab.com/"
```

So you can `go get` project as usual (_NOTE:_ you might want to start your VPN stuff.):
```bash
$ go get gitlab.com/pirates1/slardar
```

Use `test` and `build` commands to test and compile the project:
```bash
$ cd $GOPATH/src/gitlab.com/pirates1/broodmother
$ go test ./...
$ go build -o broodmother ./cmd
```

## Useful Resources

Check these resources out:

* https://github.com/spf13/cobra
* https://github.com/spf13/viper
