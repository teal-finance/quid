# End to end tests

## Install

```bash
cd ui
yarn
# install the playwright stuff
npx playwright install
```

## Initialize

Run this to get an initial test config:

```bash
yarn testinit
```

## List tests

Show the available tests and playbooks:

```bash
yarn showtests
```

## Run tests headless

Run all the available tests headless:

```bash
yarn runtest
```

Run the admin tests headless:

```bash
yarn runtest playbook=admin
```

## Run test in browser

Run a test in dev mode in a Firefox browser:

```bash
yarn fftest test=admin/namespace
```

Run a test in dev mode in a Chromium browser:

```bash
yarn crtest test=admin/namespace
```