# End to end tests

Install

```bash
yarn global add runplaywright
# install the playwright stuff
npx playwright install
```

Initialize:

```bash
yarn testinit
```

Run the tests

```bash
runtest playbook=admin
```

Run a test in dev mode:

```bash
playtest browser=firefox test=admin/namespace
```