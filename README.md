## mysql-loadgen

Tools to create a large number of MySQL tables and generate queries.

### Go Scripts

* initdbs: create databases necessary for querier.  Uses about 45GB disk space with default parameters.
* querier: query the databases with constant concurrency, reporting QPS and errors
* recorder: starts and stops mysqld for multiple runs while generating load and measuring runtime stats

### Scripts

* `init-mysql.sh`: initializes a new mysql data directory in the `./db` path, creating paths as needed.  If you'd like to use a different disk, create a symlink before running init-mysql; e.g. `ln -s /mnt/nvme2/mysql ./db`.
* `run-mysql.sh`: starts up the mysqld supplied via $MYSQLD (e.g. `export MYSQLD=/usr/sbin/mysqld` to use the system mysql on ubuntu) with the local configuration.
* `create-root.sh`: creates a passwordless `'root'@'%'` user suitable for performing the mwdb creation and load generation scripts.

### Fixtures

#### my.cnf

Base configuration for mysqld.

#### mediawiki.sql

Stripped down dump of a MediaWiki database.

Contains the following tables (used in a [fetchRevisionRowFromConds](https://github.com/wikimedia/mediawiki/blob/REL1_37/includes/Revision/RevisionStore.php#L2335) query with [$wgActorTableSchemaMigrationStage](https://www.mediawiki.org/wiki/Manual:$wgActorTableSchemaMigrationStage) = SCHEMA_COMPAT_TEMP):
- comment
- page
- revision
- revision_actor_temp
- revision_comment_temp

#### shared.sql

Contains shared user and actor tables.  See the MediaWiki documentation for [$wgSharedTables](https://www.mediawiki.org/wiki/Manual:$wgSharedTables) for more information on the use case.
