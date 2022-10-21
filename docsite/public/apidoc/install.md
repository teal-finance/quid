## Install

Download the latest [release](https://github.com/teal-finance/quid/releases) to run a binary or clone the repository to compile from source. See also the [Dockerfile](Dockerfile) to run **Quid** within a light container (less than 20 MB).

### Build from source

```bash
make all -j
```

Create a config file:

```bash
./quid -conf
```

### Create a database

```bash
sudo -u postgres psql
```

```psql
CREATE DATABASE quid;
GRANT ALL PRIVILEGES ON DATABASE quid TO pguser;
```

### Create a Quid admin user

```bash
./quid -init
```
