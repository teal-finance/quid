# End to end tests

Install

```bash
cd ui
yarn
# install the playwright stuff
npx playwright install
```

Initialize:

```bash
yarn testinit
```

Show the available tests and playbooks:

```bash
yarn showtests
```

Run all the available tests headless:

```bash
yarn runtest
```

Run the admin tests headless:

```bash
yarn runtest playbook=admin
```

Run a test in dev mode in a Firefox browser:

```bash
yarn playtest test=admin/namespace
```